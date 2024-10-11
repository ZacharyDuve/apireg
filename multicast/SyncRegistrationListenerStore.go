package multicast

import (
	"sync"

	"github.com/ZacharyDuve/apireg/apievent"
)

type syncRegListenStore struct {
	listeners      []apievent.RegistrationListener
	listenersMutex *sync.RWMutex
}

func NewSyncRegistrationListenerStore() RegistrationListenerStore {
	s := &syncRegListenStore{}
	s.listeners = make([]apievent.RegistrationListener, 0)
	s.listenersMutex = &sync.RWMutex{}

	return s
}

func (this *syncRegListenStore) Add(l apievent.RegistrationListener) {
	this.listenersMutex.Lock()
	this.listeners = append(this.listeners, l)
	this.listenersMutex.Unlock()
}
func (this *syncRegListenStore) Remove(l apievent.RegistrationListener) {
	this.listenersMutex.Lock()
	for i, curL := range this.listeners {
		if curL == l {
			this.listeners = append(this.listeners[:i], this.listeners[i+1:]...)
			break
		}
	}
	this.listenersMutex.Unlock()
}
func (this *syncRegListenStore) Notify(e apievent.RegistrationEvent) {
	go func() {
		this.listenersMutex.RLock()
		for _, curL := range this.listeners {
			curL.HandleRegistration(e)
		}
		this.listenersMutex.RUnlock()
	}()
}
