package apievent

type RegistrationListener interface {
	HandleRegistration(RegistrationEvent)
}
