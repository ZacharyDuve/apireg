package event

type RegistrationListener interface {
	HandleRegistration(RegistrationEvent)
}
