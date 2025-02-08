package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/justinthompson/appointment/pkg/appointment"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertHTTPError(t *testing.T, err error, code int) {
	httpError, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, code, httpError.Code)
}

func TestHandleGetAvailableTimes(t *testing.T) {
	t.Run("successful retrieval of available appointments", func(t *testing.T) {
		appManager := appointment.NewMockAppointmentManager([]appointment.Appointment{{TrainerID: 1}}, nil)
		c, rec := newContext()
		SetAppointment(c, appointment.Appointment{TrainerID: 1})
		err := handleGetAvailableTimes(c, appManager)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var availableAppointments []appointment.Appointment
		err = json.Unmarshal(rec.Body.Bytes(), &availableAppointments)
		require.NoError(t, err)
		require.NotEmpty(t, availableAppointments)
		assert.Equal(t, 1, availableAppointments[0].TrainerID)

	})
	t.Run("error handling when GetAvailableAppointments returns an error", func(t *testing.T) {
		appManager := appointment.NewMockAppointmentManager([]appointment.Appointment{}, fmt.Errorf("error getting available appointments"))
		c, _ := newContext()
		SetAppointment(c, appointment.Appointment{TrainerID: 1})
		err := handleGetAvailableTimes(c, appManager)
		assertHTTPError(t, err, http.StatusBadRequest)
	})
}

func TestHandleGetScheduledAppointments(t *testing.T) {
	t.Run("successful retrieval of scheduled appointments", func(t *testing.T) {
		appManager := appointment.NewMockAppointmentManager([]appointment.Appointment{{TrainerID: 1}}, nil)
		c, rec := newContext()
		SetAppointment(c, appointment.Appointment{TrainerID: 1})
		err := handleGetScheduledAppointments(c, appManager)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var availableAppointments []appointment.Appointment
		err = json.Unmarshal(rec.Body.Bytes(), &availableAppointments)
		require.NoError(t, err)
		require.NotEmpty(t, availableAppointments)
		assert.Equal(t, 1, availableAppointments[0].TrainerID)

	})
	t.Run("error handling when GetScheduledAppointments returns an error", func(t *testing.T) {
		appManager := appointment.NewMockAppointmentManager([]appointment.Appointment{}, fmt.Errorf("error getting scheduled appointments"))
		c, _ := newContext()
		SetAppointment(c, appointment.Appointment{TrainerID: 1})
		err := handleGetScheduledAppointments(c, appManager)
		assertHTTPError(t, err, http.StatusBadRequest)
	})
}

func TestHandlePostAppointment(t *testing.T) {
	t.Run("successful creation of an appointment", func(t *testing.T) {
		appManager := appointment.NewMockAppointmentManager([]appointment.Appointment{{TrainerID: 1}}, nil)
		c, _ := newContext()
		SetAppointment(c, appointment.Appointment{TrainerID: 1})
		err := handlePostAppointment(c, appManager)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, c.Response().Status)
	})
	t.Run("error handling when CreateAppointment returns an error", func(t *testing.T) {
		appManager := appointment.NewMockAppointmentManager([]appointment.Appointment{}, fmt.Errorf("error creating appointment"))
		c, _ := newContext()
		SetAppointment(c, appointment.Appointment{TrainerID: 1})
		err := handlePostAppointment(c, appManager)
		assertHTTPError(t, err, http.StatusBadRequest)
	})
}
