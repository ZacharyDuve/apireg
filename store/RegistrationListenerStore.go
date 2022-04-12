package store

import "github.com/ZacharyDuve/apireg/event"

type RegistrationListenerStore interface {
	Add(event.RegistrationListener)
	Remove(event.RegistrationListener)
	Notify(event.RegistrationEvent)
}
