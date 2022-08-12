package store

import "git.zmanhobbies.com/software/apireg/event"

type RegistrationListenerStore interface {
	Add(event.RegistrationListener)
	Remove(event.RegistrationListener)
	Notify(event.RegistrationEvent)
}
