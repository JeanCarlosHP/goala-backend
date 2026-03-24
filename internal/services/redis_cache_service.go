package services

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jeancarloshp/calorieai/internal/domain"
)

type RedisCacheService struct {
	address  string
	password string
	enabled  bool
	logger   domain.Logger
	timeout  time.Duration
}

func NewRedisCacheService(cfg *domain.Config, logger domain.Logger) *RedisCacheService {
	service := &RedisCacheService{
		logger:  logger,
		timeout: 2 * time.Second,
	}

	if cfg.RedisURL == "" {
		return service
	}

	parsed, err := url.Parse(cfg.RedisURL)
	if err != nil {
		logger.Warn("invalid redis url, cache disabled", "error", err)
		return service
	}

	service.address = parsed.Host
	if service.address == "" {
		service.address = parsed.Path
	}
	if service.address == "" {
		logger.Warn("missing redis address, cache disabled")
		return service
	}

	if parsed.User != nil {
		if password, ok := parsed.User.Password(); ok {
			service.password = password
		}
	}

	service.enabled = true
	return service
}

func (s *RedisCacheService) Enabled() bool {
	return s.enabled
}

func (s *RedisCacheService) GetJSON(ctx context.Context, key string, target any) bool {
	if !s.enabled {
		return false
	}

	reply, err := s.execute(ctx, "GET", key)
	if err != nil || len(reply) == 0 || reply == "$-1\r\n" {
		return false
	}

	value, err := parseBulkString(reply)
	if err != nil {
		s.logger.Warn("failed to parse redis value", "error", err)
		return false
	}

	if err := json.Unmarshal([]byte(value), target); err != nil {
		s.logger.Warn("failed to decode cached json", "error", err)
		return false
	}

	return true
}

func (s *RedisCacheService) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) {
	if !s.enabled {
		return
	}

	payload, err := json.Marshal(value)
	if err != nil {
		s.logger.Warn("failed to encode cache value", "error", err)
		return
	}

	_, err = s.execute(ctx, "SETEX", key, strconv.Itoa(int(ttl.Seconds())), string(payload))
	if err != nil {
		s.logger.Warn("failed to set redis value", "error", err)
	}
}

func (s *RedisCacheService) execute(ctx context.Context, args ...string) (string, error) {
	conn, err := net.DialTimeout("tcp", s.address, s.timeout)
	if err != nil {
		return "", err
	}
	defer func() { _ = conn.Close() }()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	} else {
		_ = conn.SetDeadline(time.Now().Add(s.timeout))
	}

	reader := bufio.NewReader(conn)
	if s.password != "" {
		if err := writeRESP(conn, "AUTH", s.password); err != nil {
			return "", err
		}
		if _, err := readRESP(reader); err != nil {
			return "", err
		}
	}

	if err := writeRESP(conn, args...); err != nil {
		return "", err
	}
	return readRESP(reader)
}

func writeRESP(conn net.Conn, args ...string) error {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("*%d\r\n", len(args)))
	for _, arg := range args {
		builder.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}
	_, err := conn.Write([]byte(builder.String()))
	return err
}

func readRESP(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	switch {
	case strings.HasPrefix(line, "$"):
		size, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "$")))
		if err != nil || size < 0 {
			return line, err
		}
		payload := make([]byte, size+2)
		if _, err := reader.Read(payload); err != nil {
			return "", err
		}
		return line + string(payload), nil
	case strings.HasPrefix(line, "-"):
		return "", fmt.Errorf("%s", strings.TrimSpace(strings.TrimPrefix(line, "-")))
	default:
		return line, nil
	}
}

func parseBulkString(resp string) (string, error) {
	lines := strings.SplitN(resp, "\r\n", 3)
	if len(lines) < 3 {
		return "", fmt.Errorf("invalid bulk string response")
	}
	return lines[1], nil
}
