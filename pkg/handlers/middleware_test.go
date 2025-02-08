package handlers

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/justinthompson/appointment/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Init() {
}

func TestRequestToAppointment(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		c, _ := newContext()
		req := &GetAppointmentRequest{
			StartTime: time.Time{},
			EndTime:   time.Time{},
			TrainerID: 0,
		}
		_, err := requestToAppointment(c, req)
		require.Error(t, err)
	})
	t.Run("unknown request type error", func(t *testing.T) {
		c, _ := newContext()
		req := struct{}{}
		_, err := requestToAppointment(c, req)
		require.Error(t, err)
	})
	t.Run("successful conversion of GetAppointmentRequest", func(t *testing.T) {
		c, _ := newContext()
		req := &GetAppointmentRequest{
			StartTime: time.Now(),
			EndTime:   time.Now().Add(1 * time.Hour),
			TrainerID: 1,
		}
		app, err := requestToAppointment(c, req)
		assert.NoError(t, err)
		assert.Equal(t, req.StartTime, app.StartTime)
		assert.Equal(t, req.EndTime, app.EndTime)
		assert.Equal(t, req.TrainerID, app.TrainerID)
	})
	t.Run("successful conversion of PostAppointmentRequest", func(t *testing.T) {
		c, _ := newContext()
		req := &PostAppointmentRequest{
			StartTime: time.Now(),
			EndTime:   time.Now().Add(1 * time.Hour),
			TrainerID: 1,
			UserID:    1,
		}
		app, err := requestToAppointment(c, req)
		assert.NoError(t, err)
		assert.Equal(t, req.StartTime, app.StartTime)
		assert.Equal(t, req.EndTime, app.EndTime)
		assert.Equal(t, req.TrainerID, app.TrainerID)
		assert.Equal(t, req.UserID, app.UserID)
	})
	t.Run("successful conversion of GetScheduledRequest", func(t *testing.T) {
		c, _ := newContext()
		req := &GetScheduledRequest{
			TrainerID: 1,
		}
		app, err := requestToAppointment(c, req)
		assert.NoError(t, err)
		assert.Equal(t, req.TrainerID, app.TrainerID)
	})
}

func newContext() (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Echo().Validator = validator.NewValidator()
	return c, rec
}
