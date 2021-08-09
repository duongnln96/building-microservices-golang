package logger

import (
	"fmt"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// ZapLogger is a middleware and zap to provide an "access log" like logging for each request.
func ZapLogger(log *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			logStr := fmt.Sprintf("{\"remote_ip\": \"%s\", \"host\": \"%s\", \"method\": \"%s\", \"uri\": \"%s\" ,\"user_agent\": \"%s\", \"status\": %d}",
				c.RealIP(), req.Host, req.Method, req.RequestURI, req.UserAgent(), res.Status)

			n := res.Status
			switch {
			case n >= 500:
				log.With(zap.Error(err)).Errorf("Server error %s", logStr)
			case n >= 400:
				log.With(zap.Error(err)).Warnf("Client error %s", logStr)
			case n >= 300:
				log.Infof("Redirection %s", logStr)
			default:
				log.Infof("Success %s", logStr)
			}

			return nil
		}
	}
}
