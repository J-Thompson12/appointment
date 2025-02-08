package appointment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAvailableAppointments_NoAppointments_Overlapping(t *testing.T) {
	startTime, err := time.Parse(time.RFC3339, "2019-01-24T09:00:00-08:00")
	require.NoError(t, err)

	endTime, err := time.Parse(time.RFC3339, "2019-01-24T09:30:00-08:00")
	require.NoError(t, err)

	apps := scheduledAppointments{
		appointmentsList: []Appointment{
			{
				StartTime: startTime,
				EndTime:   endTime,
				UserID:    1,
				TrainerID: 1,
			},
		},
		TrainerIDs: map[int]bool{
			1: true,
		},
	}

	appReq := Appointment{
		StartTime: startTime,
		EndTime:   endTime,
		TrainerID: 1,
	}

	available, err := apps.GetAvailableAppointments(appReq)
	require.NoError(t, err)
	require.Len(t, available, 0)
}

func TestGetAvailableAppointments_Available_DifferentTrainer(t *testing.T) {
	startTime, err := time.Parse(time.RFC3339, "2019-01-24T09:00:00-08:00")
	require.NoError(t, err)

	endTime, err := time.Parse(time.RFC3339, "2019-01-24T09:30:00-08:00")
	require.NoError(t, err)

	apps := scheduledAppointments{
		appointmentsList: []Appointment{
			{
				StartTime: startTime,
				EndTime:   endTime,
				UserID:    2,
				TrainerID: 1,
			},
		},
		TrainerIDs: map[int]bool{
			2: true,
		},
	}

	appReq := Appointment{
		StartTime: startTime,
		EndTime:   endTime,
		TrainerID: 2,
	}

	available, err := apps.GetAvailableAppointments(appReq)
	require.NoError(t, err)
	require.Len(t, available, 1)
}

func TestGetAvailableAppointments_Available(t *testing.T) {
	startTime, err := time.Parse(time.RFC3339, "2019-01-24T09:00:00-08:00")
	require.NoError(t, err)

	endTime, err := time.Parse(time.RFC3339, "2019-01-24T09:30:00-08:00")
	require.NoError(t, err)

	apps := scheduledAppointments{
		appointmentsList: []Appointment{
			{
				StartTime: startTime,
				EndTime:   endTime,
				UserID:    1,
				TrainerID: 1,
			},
			{
				StartTime: startTime,
				EndTime:   endTime,
				UserID:    2,
				TrainerID: 1,
			},
		},
		TrainerIDs: map[int]bool{
			1: true,
			2: true,
		},
	}

	requestStartTime, err := time.Parse(time.RFC3339, "2019-01-24T10:00:00-08:00")
	require.NoError(t, err)
	requestEndTime, err := time.Parse(time.RFC3339, "2019-01-24T11:00:00-08:00")
	require.NoError(t, err)

	appReq := Appointment{
		StartTime: requestStartTime,
		EndTime:   requestEndTime,
		TrainerID: 1,
	}

	available, err := apps.GetAvailableAppointments(appReq)
	require.NoError(t, err)
	require.Len(t, available, 2)
}

func TestGetScheduledAppointments(t *testing.T) {
	t.Run("valid trainer ID", func(t *testing.T) {
		a := scheduledAppointments{
			appointmentsList: []Appointment{
				{TrainerID: 1},
				{TrainerID: 1},
				{TrainerID: 2},
			},
			TrainerIDs: map[int]bool{1: true, 2: true},
		}
		appointments, err := a.GetScheduledAppointments(1)
		require.NoError(t, err)
		require.Len(t, appointments, 2)
		for _, app := range appointments {
			require.Equal(t, 1, app.TrainerID)
		}
	})
	t.Run("invalid trainer ID", func(t *testing.T) {
		a := scheduledAppointments{
			appointmentsList: []Appointment{
				{TrainerID: 1},
				{TrainerID: 1},
				{TrainerID: 2},
			},
			TrainerIDs: map[int]bool{1: true, 2: true},
		}
		appointments, err := a.GetScheduledAppointments(3)
		require.Error(t, err)
		require.Nil(t, appointments)
	})
	t.Run("empty appointments list", func(t *testing.T) {
		a := scheduledAppointments{
			appointmentsList: []Appointment{},
			TrainerIDs:       map[int]bool{1: true},
		}
		appointments, err := a.GetScheduledAppointments(1)
		require.NoError(t, err)
		require.Empty(t, appointments)
	})
	t.Run("multiple appointments for the same trainer ID", func(t *testing.T) {
		a := scheduledAppointments{
			appointmentsList: []Appointment{
				{TrainerID: 1},
				{TrainerID: 1},
				{TrainerID: 1},
			},
			TrainerIDs: map[int]bool{1: true},
		}
		appointments, err := a.GetScheduledAppointments(1)
		require.NoError(t, err)
		require.Len(t, appointments, 3)
		for _, app := range appointments {
			require.Equal(t, 1, app.TrainerID)
		}
	})
}

func TestCreateAppointment_InvalidTrainerID(t *testing.T) {
	a := scheduledAppointments{
		TrainerIDs: map[int]bool{1: true},
	}

	app := Appointment{
		TrainerID: 2,
	}

	err := a.CreateAppointment(app)
	if err == nil || err.Error() != "trainer does not exist" {
		t.Errorf("expected error 'trainer does not exist', got %v", err)
	}
}

func TestCreateAppointment_InvalidStartAndEndTime(t *testing.T) {
	a := scheduledAppointments{
		TrainerIDs: map[int]bool{1: true},
	}

	app := Appointment{
		TrainerID: 1,
		StartTime: time.Date(2022, 1, 1, 7, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
	}

	err := a.CreateAppointment(app)
	if err == nil || err.Error() != "appointment time must be between 8am and 5pm" {
		t.Errorf("expected error 'appointment time must be between 8am and 5pm', got %v", err)
	}
}

func TestCreateAppointment_InvalidDuration(t *testing.T) {
	app := Appointment{
		TrainerID: 1,
		StartTime: time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2022, 1, 1, 9, 45, 0, 0, time.UTC),
	}

	a := scheduledAppointments{
		TrainerIDs: map[int]bool{1: true},
	}

	err := a.CreateAppointment(app)
	assert.ErrorContains(t, err, "appointment times must start and end on the hour or half-hour")
}

func TestCreateAppointment_OverlappingAppointment(t *testing.T) {
	a := scheduledAppointments{
		TrainerIDs: map[int]bool{1: true},
		appointmentsList: []Appointment{
			{
				TrainerID: 1,
				StartTime: time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 1, 1, 9, 30, 0, 0, time.UTC),
			},
		},
	}

	app := Appointment{
		TrainerID: 1,
		StartTime: time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2022, 1, 1, 9, 30, 0, 0, time.UTC),
	}

	err := a.CreateAppointment(app)
	if err == nil || err.Error() != "appointment already exists at this time" {
		t.Errorf("expected error 'appointment already exists at this time', got %v", err)
	}
}

func TestCreateAppointment_Success(t *testing.T) {
	a := scheduledAppointments{
		TrainerIDs: map[int]bool{1: true},
	}

	app := Appointment{
		TrainerID: 1,
		StartTime: time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2022, 1, 1, 9, 30, 0, 0, time.UTC),
	}

	err := a.CreateAppointment(app)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateStartAndEndTime(t *testing.T) {
	tests := []struct {
		name       string
		startTime  time.Time
		endTime    time.Time
		wantError  bool
		wantErrMsg string
	}{
		{
			name:       "start and end times on the hour",
			startTime:  time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:    time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			wantError:  false,
			wantErrMsg: "",
		},
		{
			name:       "start and end times on the half-hour",
			startTime:  time.Date(2022, 1, 1, 9, 30, 0, 0, time.UTC),
			endTime:    time.Date(2022, 1, 1, 10, 30, 0, 0, time.UTC),
			wantError:  false,
			wantErrMsg: "",
		},
		{
			name:       "start time after end time",
			startTime:  time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			endTime:    time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
			wantError:  true,
			wantErrMsg: "start time must be before end time",
		},
		{
			name:       "start time before 8am",
			startTime:  time.Date(2022, 1, 1, 7, 0, 0, 0, time.UTC),
			endTime:    time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			wantError:  true,
			wantErrMsg: "appointment time must be between 8am and 5pm",
		},
		{
			name:       "start time after 5pm",
			startTime:  time.Date(2022, 1, 1, 18, 0, 0, 0, time.UTC),
			endTime:    time.Date(2022, 1, 1, 19, 0, 0, 0, time.UTC),
			wantError:  true,
			wantErrMsg: "appointment time must be between 8am and 5pm",
		},
		{
			name:       "end time before 8am",
			startTime:  time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:    time.Date(2022, 1, 1, 7, 0, 0, 0, time.UTC),
			wantError:  true,
			wantErrMsg: "start time must be before end time",
		},
		{
			name:       "end time after 5pm",
			startTime:  time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
			endTime:    time.Date(2022, 1, 1, 18, 0, 0, 0, time.UTC),
			wantError:  true,
			wantErrMsg: "appointment time must be between 8am and 5pm",
		},
		{
			name:       "invalid minute values",
			startTime:  time.Date(2022, 1, 1, 9, 15, 0, 0, time.UTC),
			endTime:    time.Date(2022, 1, 1, 10, 15, 0, 0, time.UTC),
			wantError:  true,
			wantErrMsg: "appointment times must start and end on the hour or half-hour",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStartAndEndTime(tt.startTime, tt.endTime)
			if (err != nil) != tt.wantError {
				t.Errorf("Test %s: validateStartAndEndTime() error = %v, wantError %v", tt.name, err, tt.wantError)
				return
			}
			if err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("Test %s: validateStartAndEndTime() error message = %v, want %v", tt.name, err.Error(), tt.wantErrMsg)
			}
		})
	}
}
