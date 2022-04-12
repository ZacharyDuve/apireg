package event

import "github.com/ZacharyDuve/apireg/api"

type EventType byte

const (
	Added   EventType = 0
	Removed EventType = 1
)

type RegistrationEvent interface {
	Type() EventType
	Api() api.Api
}

type eventImpl struct {
	eType EventType
	api   api.Api
}

func (this *eventImpl) Type() EventType {
	return this.eType
}

func (this *eventImpl) Api() api.Api {
	return this.api
}

func NewAddEvent(a api.Api) RegistrationEvent {
	if a != nil {
		e := &eventImpl{}
		e.eType = Added
		e.api = a
		return e
	}
	return nil
}

func NewRemovedEvent(a api.Api) RegistrationEvent {
	if a != nil {
		e := &eventImpl{}
		e.eType = Removed
		e.api = a
		return e
	}
	return nil
}
