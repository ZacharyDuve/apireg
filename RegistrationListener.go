package apireg

type RegistrationListener interface {
	HandleRegistration(RegistrationEvent)
}
