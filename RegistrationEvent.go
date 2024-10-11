package apireg

type EventType string

const (
	Added   EventType = "add"
	Removed EventType = "remove"
)

type RegistrationEvent interface {
	Type() EventType
	Api() Api
}

type eventImpl struct {
	eType EventType
	api   Api
}

func (this *eventImpl) Type() EventType {
	return this.eType
}

func (this *eventImpl) Api() Api {
	return this.api
}

func NewAddEvent(a Api) RegistrationEvent {
	if a != nil {
		e := &eventImpl{}
		e.eType = Added
		e.api = a
		return e
	}
	return nil
}

func NewRemovedEvent(a Api) RegistrationEvent {
	if a != nil {
		e := &eventImpl{}
		e.eType = Removed
		e.api = a
		return e
	}
	return nil
}
