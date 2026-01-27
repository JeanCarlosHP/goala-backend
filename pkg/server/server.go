package server

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"github.com/jeancarloshp/calorieai/pkg/server/middleware"
)

type httpServer struct {
	App    *fiber.App
	Config *domain.Config
	logger domain.Logger
}

func New(config *domain.Config, logger domain.Logger) domain.HTTPServer {
	cRequestSizeL, err := strconv.Atoi(config.HTTPRequestSizeLimit)
	if err != nil {
		panic(err)
	}

	requestSizeLimit := cRequestSizeL * 1024 * 1024
	headerSizeLimit := 16 * 1024 // 16Kb

	cfg := fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          ErrorHandler,
		BodyLimit:             requestSizeLimit,
		ReadBufferSize:        headerSizeLimit,
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
	}

	app := fiber.New(cfg)

	srv := &httpServer{
		App:    app,
		Config: config,
		logger: logger,
	}

	srv.ConfigureMiddlewares()

	return srv
}

func (s *httpServer) StartServer() error {
	s.logger.Info("starting server", "port", s.Config.HTTPPort)
	return s.App.Listen(fmt.Sprintf(":%s", s.Config.HTTPPort))
}

func (s *httpServer) ShutdownServer() error {
	return s.App.Shutdown()
}

func (s *httpServer) GetApp() *fiber.App {
	return s.App
}

func (s *httpServer) GetLogger() domain.Logger {
	return s.logger
}

func (s *httpServer) ConfigureMiddlewares() {
	// OpenTelemetry tracing middleware deve ser o primeiro
	s.App.Use(otelfiber.Middleware())
	s.App.Use(cors.New(cors.Config{
		AllowHeaders: s.Config.HTTPCorsAllowedHeaders,
		AllowMethods: s.Config.HTTPCorsAllowedMethods,
		AllowOrigins: s.Config.HTTPCorsAllowedOrigins,
	}))

	s.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			s.GetLogger().Error(fmt.Sprintf("panic: %v\n%s\n", e, buf))
		},
	}))

	s.App.Use(requestid.New())

	s.App.Use(healthcheck.New())
	s.App.Use(recover.New())
	s.App.Use(middleware.PrometheusMetrics())
	s.App.Use(middleware.RequestLogger())
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
