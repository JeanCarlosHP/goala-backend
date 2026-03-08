package server

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"github.com/bytedance/sonic"
	fiberotel "github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
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
		ErrorHandler:   ErrorHandler,
		BodyLimit:      requestSizeLimit,
		ReadBufferSize: headerSizeLimit,
		JSONEncoder:    sonic.Marshal,
		JSONDecoder:    sonic.Unmarshal,
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
	config := fiber.ListenConfig{
		DisableStartupMessage: true,
		EnablePrefork:         false,
	}

	return s.App.Listen(
		fmt.Sprintf(":%s", s.Config.HTTPPort),
		config,
	)
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
	s.App.Use(fiberotel.Middleware())
	s.App.Use(cors.New(cors.Config{
		AllowHeaders: s.Config.HTTPCorsAllowedHeaders,
		AllowMethods: s.Config.HTTPCorsAllowedMethods,
		AllowOrigins: s.Config.HTTPCorsAllowedOrigins,
	}))

	s.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e any) {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			s.GetLogger().Error(fmt.Sprintf("panic: %v\n%s\n", e, buf))
		},
	}))

	s.App.Use(requestid.New())

	s.App.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	s.App.Get(healthcheck.ReadinessEndpoint, healthcheck.New())
	s.App.Get(healthcheck.StartupEndpoint, healthcheck.New())

	s.App.Use(middleware.PrometheusMetrics())
	s.App.Use(middleware.RequestLogger(s.logger))
}

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
