package handlers

import (
	"net/http"
	"time"

	"github.com/justinthompson/appointment/pkg/appointment"
	"github.com/labstack/echo/v4"
)

const keyAppointmentRequest = "appointment"

// I could just have one appointment struct that they all share but I wanted to test out the different ways of binding and using middleware
type GetAppointmentRequest struct {
	StartTime time.Time `query:"starts_at" validate:"required"`
	EndTime   time.Time `query:"ends_at" validate:"required"`
	TrainerID int       `query:"trainer_id" validate:"required"`
}

type GetScheduledRequest struct {
	TrainerID int `query:"trainer_id" validate:"required"`
}

type PostAppointmentRequest struct {
	StartTime time.Time `json:"starts_at" validate:"required"`
	EndTime   time.Time `json:"ends_at" validate:"required"`
	TrainerID int       `json:"trainer_id" validate:"required"`
	UserID    int       `json:"user_id" validate:"required"`
}

// MiddlewareAvailable is a middleware that takes the request and converts it to an appointment
func MiddlewareAvailable(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		app, err := requestToAppointment(c, &GetAppointmentRequest{})
		if err != nil {
			return err
		}

		SetAppointment(c, app)
		next(c)
		return nil
	}
}

// MiddlewareScheduled is a middleware that takes the request and converts it to an appointment
func MiddlewareScheduled(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		app, err := requestToAppointment(c, &GetScheduledRequest{})
		if err != nil {
			return err
		}

		SetAppointment(c, app)
		next(c)
		return nil
	}
}

// MiddlewarePost is a middleware that takes the request and converts it to an appointment
func MiddlewarePost(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		app, err := requestToAppointment(c, &PostAppointmentRequest{})
		if err != nil {
			return err
		}

		SetAppointment(c, app)
		next(c)
		return nil
	}
}

func SetAppointment(c echo.Context, app appointment.Appointment) {
	c.Set(keyAppointmentRequest, app)
}

func GetAppointment(c echo.Context) appointment.Appointment {
	return c.Get(keyAppointmentRequest).(appointment.Appointment)
}

// requestToAppointment takes any of the request types and converts it to appointment
// It does this by binding the request to the request struct, validating it,
// and then switching on the type of the request to construct the appointment
func requestToAppointment(c echo.Context, req interface{}) (appointment.Appointment, error) {
	// Bind the request to the request struct
	if err := c.Bind(req); err != nil {
		return appointment.Appointment{}, echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate the request
	if err := c.Validate(req); err != nil {
		return appointment.Appointment{}, echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Switch on the type of the request to construct the appointment
	switch v := req.(type) {
	case *GetAppointmentRequest:
		return appointment.Appointment{
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
			TrainerID: v.TrainerID,
		}, nil
	case *PostAppointmentRequest:
		return appointment.Appointment{
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
			TrainerID: v.TrainerID,
			UserID:    v.UserID,
		}, nil
	case *GetScheduledRequest:
		return appointment.Appointment{
			TrainerID: v.TrainerID,
		}, nil
	default:
		return appointment.Appointment{}, echo.NewHTTPError(http.StatusBadRequest, "unknown request type")
	}
}
