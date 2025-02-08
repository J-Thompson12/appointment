package appointment

type MockAppointmentManager struct {
	AppointmentsList []Appointment
	Err              error
}

func NewMockAppointmentManager(AppointmentsList []Appointment, err error) *MockAppointmentManager {
	return &MockAppointmentManager{
		AppointmentsList: AppointmentsList,
		Err:              err,
	}
}

func (m *MockAppointmentManager) GetAvailableAppointments(appReq Appointment) ([]Appointment, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return m.AppointmentsList, nil
}

func (m *MockAppointmentManager) CreateAppointment(app Appointment) error {
	if m.Err != nil {
		return m.Err
	}

	return nil
}

func (m *MockAppointmentManager) GetScheduledAppointments(trainerID int) ([]Appointment, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return m.AppointmentsList, nil
}
