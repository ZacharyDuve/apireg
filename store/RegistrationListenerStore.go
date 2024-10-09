package store

import "github.com/ZacharyDuve/apireg/apievent"

type RegistrationListenerStore interface {
	Add(apievent.RegistrationListener)
	Remove(apievent.RegistrationListener)
	Notify(apievent.RegistrationEvent)
}
