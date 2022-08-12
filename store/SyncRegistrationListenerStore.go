package store

import (
	"sync"

	"git.zmanhobbies.com/software/apireg/event"
)

type syncRegListenStore struct {
	listeners      []event.RegistrationListener
	listenersMutex *sync.RWMutex
}

func NewSyncRegistrationListenerStore() RegistrationListenerStore {
	s := &syncRegListenStore{}
	s.listeners = make([]event.RegistrationListener, 0)
	s.listenersMutex = &sync.RWMutex{}

	return s
}

func (this *syncRegListenStore) Add(l event.RegistrationListener) {
	this.listenersMutex.Lock()
	this.listeners = append(this.listeners, l)
	this.listenersMutex.Unlock()
}
func (this *syncRegListenStore) Remove(l event.RegistrationListener) {
	this.listenersMutex.Lock()
	for i, curL := range this.listeners {
		if curL == l {
			this.listeners = append(this.listeners[:i], this.listeners[i+1:]...)
			break
		}
	}
	this.listenersMutex.Unlock()
}
func (this *syncRegListenStore) Notify(e event.RegistrationEvent) {
	go func() {
		this.listenersMutex.RLock()
		for _, curL := range this.listeners {
			curL.HandleRegistration(e)
		}
		this.listenersMutex.RUnlock()
	}()
}
