package handlers

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"

	"github.com/justinthompson/appointment/pkg/appointment"
	"github.com/justinthompson/appointment/pkg/validator"
)

// BuildRouter sets up the routes for the API
func BuildRouter(r *echo.Echo, appManager appointment.Manager) {
	r.Use(middleware.Recover())
	r.Use(middleware.Secure())
	r.Use(middleware.BodyLimit("1KB"))
	// In a real system I wouldnt log every request because the cost and noise would be bad.
	// I just added this because I wanted to test it out
	r.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:       true,
		LogMethod:        true,
		LogURI:           true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("uri", v.URI).
				Int("status", v.Status).
				Dur("latency", v.Latency).
				Str("content_length", v.ContentLength).
				Msg("request received")
			return nil
		},
	}))
	r.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		ErrorMessage: "request timed out",
		Timeout:      1 * time.Second,
	}))

	r.Validator = validator.NewValidator()

	handlerGetAvailableTimes := func(c echo.Context) error {
		return handleGetAvailableTimes(c, appManager)
	}

	handlerGetScheduledAppointments := func(c echo.Context) error {
		return handleGetScheduledAppointments(c, appManager)
	}

	handlerAddNewAppointment := func(c echo.Context) error {
		return handlePostAppointment(c, appManager)
	}

	r.GET("/schedule/available", handlerGetAvailableTimes, MiddlewareAvailable)
	r.GET("/schedule", handlerGetScheduledAppointments, MiddlewareScheduled)
	r.POST("/schedule", handlerAddNewAppointment, MiddlewarePost)
}
