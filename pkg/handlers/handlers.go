package handlers

import (
	"fmt"
	"net/http"

	"github.com/justinthompson/appointment/pkg/appointment"
	"github.com/labstack/echo/v4"
)

// handleGetAvailableTimes
func handleGetAvailableTimes(c echo.Context, appManager appointment.Manager) error {
	appRequest := GetAppointment(c)

	availableAppointments, err := appManager.GetAvailableAppointments(appRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("error getting available appointments: %w", err).Error())
	}
	return c.JSON(http.StatusOK, availableAppointments)
}

func handleGetScheduledAppointments(c echo.Context, appManager appointment.Manager) error {
	appRequest := GetAppointment(c)

	appointments, err := appManager.GetScheduledAppointments(appRequest.TrainerID)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("error getting scheduled appointments: %w", err).Error())
	}
	return c.JSON(http.StatusOK, appointments)
}

func handlePostAppointment(c echo.Context, appManager appointment.Manager) error {
	appRequest := GetAppointment(c)
	err := appManager.CreateAppointment(appRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Errorf("error creating appointment: %w", err).Error())
	}
	return c.NoContent(http.StatusCreated)
}
