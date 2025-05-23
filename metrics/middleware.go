package metrics

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// PrometheusMiddleware returns a middleware that tracks request metrics
func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start).Seconds()

			status := strconv.Itoa(c.Response().Status)
			method := c.Request().Method
			endpoint := c.Path()

			// Record request duration
			RequestDuration.WithLabelValues(method, endpoint).Observe(duration)

			// Record request count
			RequestCounter.WithLabelValues(method, endpoint, status).Inc()

			return err
		}
	}
}
