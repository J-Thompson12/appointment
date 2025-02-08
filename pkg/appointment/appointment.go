package appointment

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type (
	// Service handles operations on events
	Manager interface {
		// Track tracks and stores an event.
		GetAvailableAppointments(appReq Appointment) ([]Appointment, error)
		GetScheduledAppointments(trainerID int) ([]Appointment, error)
		CreateAppointment(app Appointment) error
	}

	scheduledAppointments struct {
		appointmentsList []Appointment
		latestID         int
		TrainerIDs       map[int]bool // using a map for unique values
	}

	// the json names in the file are different to the request (started_at vs starts_at) since both are json they should be the same to make this easier
	Appointment struct {
		ID        int       `json:"id,omitempty"`
		StartTime time.Time `json:"started_at"`
		EndTime   time.Time `json:"ended_at"`
		UserID    int       `json:"user_id,omitempty"`
		TrainerID int       `json:"trainer_id" validate:"required"`
	}
)

// GetAppointmentsFromFile reads a json file and returns a slice of appointments
func NewAppointmentManager() (Manager, error) {
	// Open jsonFile
	jsonFile, err := os.Open("appointments.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	var apps scheduledAppointments
	err = json.NewDecoder(jsonFile).Decode(&apps.appointmentsList)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	apps.TrainerIDs = make(map[int]bool)
	for _, app := range apps.appointmentsList {
		if app.ID > apps.latestID {
			apps.latestID = app.ID
		}
		apps.TrainerIDs[app.TrainerID] = true
	}

	return &apps, nil
}

// GetAvailableAppointments returns a slice of available appointments filtered by the provided start/end time and user ID
func (a *scheduledAppointments) GetAvailableAppointments(request Appointment) ([]Appointment, error) {
	if !isValidTrainerID(request.TrainerID, a.TrainerIDs) {
		return nil, fmt.Errorf("trainer does not exist")
	}

	if err := validateStartAndEndTime(request.StartTime, request.EndTime); err != nil {
		return nil, err
	}

	// Get relevant appointments
	relevantAppointments, err := a.getRelevantAppointments(request)
	if err != nil {
		return nil, err
	}

	// Create a slice of available appointments
	var availableAppointments []Appointment
	for t := request.StartTime; t.Before(request.EndTime); t = t.Add(30 * time.Minute) {
		if a.isSlotAvailable(t, relevantAppointments) {
			availableAppointments = append(availableAppointments, Appointment{
				StartTime: t,
				EndTime:   t.Add(30 * time.Minute),
				TrainerID: request.TrainerID,
			})
		}
	}

	return availableAppointments, nil
}

func (a *scheduledAppointments) GetScheduledAppointments(trainerID int) ([]Appointment, error) {
	if !isValidTrainerID(trainerID, a.TrainerIDs) {
		fmt.Println("failed here")
		return nil, fmt.Errorf("trainer %d does not exist", trainerID)
	}

	var scheduledAppointments []Appointment
	for _, app := range a.appointmentsList {
		if app.TrainerID == trainerID {
			scheduledAppointments = append(scheduledAppointments, app)
		}
	}
	return scheduledAppointments, nil
}

func (a *scheduledAppointments) CreateAppointment(appointment Appointment) error {
	if !isValidTrainerID(appointment.TrainerID, a.TrainerIDs) {
		return fmt.Errorf("trainer does not exist")
	}

	if err := validateStartAndEndTime(appointment.StartTime, appointment.EndTime); err != nil {
		return err
	}

	// Ensure the appointment duration is exactly 30 minutes
	if appointment.EndTime.Sub(appointment.StartTime) != 30*time.Minute {
		return fmt.Errorf("appointment duration must be exactly 30 minutes")
	}

	// Filter relevant appointments
	relevantAppointments, err := a.getRelevantAppointments(appointment)
	if err != nil {
		return err
	}

	// Check if there is any appointment overlapping this 30 min slot
	for _, existingAppointment := range relevantAppointments {
		if appointment.StartTime.Equal(existingAppointment.StartTime) {
			return fmt.Errorf("appointment already exists at this time")
		}
	}

	a.latestID++
	appointment.ID = a.latestID
	a.appointmentsList = append(a.appointmentsList, appointment)
	return nil
}

// getRelevantAppointments returns a slice of appointments that are in the provided time range and belong to the provided trainer
func (a *scheduledAppointments) getRelevantAppointments(request Appointment) ([]Appointment, error) {
	var relevantAppointments []Appointment
	for _, scheduledApp := range a.appointmentsList {
		if scheduledApp.TrainerID == request.TrainerID && scheduledApp.StartTime.Before(request.EndTime) && scheduledApp.EndTime.After(request.StartTime) {
			relevantAppointments = append(relevantAppointments, scheduledApp)
		}
	}

	return relevantAppointments, nil
}

// isSlotAvailable checks if the slot is available
func (a *scheduledAppointments) isSlotAvailable(slot time.Time, scheduledAppointments []Appointment) bool {
	for _, scheduledApp := range scheduledAppointments {
		if slot.Equal(scheduledApp.StartTime) {
			return false
		}
	}

	return true
}

// isValidTrainerID checks if the trainer ID is valid
func isValidTrainerID(trainerID int, trainersList map[int]bool) bool {
	if _, ok := trainersList[trainerID]; !ok {
		return false
	}
	return true
}

// validateStartAndEndTime checks that the start and end times are valid
func validateStartAndEndTime(startTime time.Time, endTime time.Time) error {
	// Normally I would handle all times in utc but since the json file is in pacific I am just treating everything like its pacific
	if startTime.Minute()%30 != 0 || endTime.Minute()%30 != 0 {
		return fmt.Errorf("appointment times must start and end on the hour or half-hour")
	}

	if startTime.After(endTime) {
		return fmt.Errorf("start time must be before end time")
	}

	if startTime.Hour() < 8 || startTime.Hour() > 17 || endTime.Hour() < 8 || endTime.Hour() > 17 {
		return fmt.Errorf("appointment time must be between 8am and 5pm")
	}
	return nil
}
