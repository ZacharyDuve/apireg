package multicast

import (
	"sync"

	"github.com/ZacharyDuve/apireg"
)

type syncRegListenStore struct {
	listeners      []apireg.RegistrationListener
	listenersMutex *sync.RWMutex
}

func newSyncRegistrationListenerStore() *syncRegListenStore {
	s := &syncRegListenStore{}
	s.listeners = make([]apireg.RegistrationListener, 0)
	s.listenersMutex = &sync.RWMutex{}

	return s
}

func (this *syncRegListenStore) Add(l apireg.RegistrationListener) {
	this.listenersMutex.Lock()
	this.listeners = append(this.listeners, l)
	this.listenersMutex.Unlock()
}
func (this *syncRegListenStore) Remove(l apireg.RegistrationListener) {
	this.listenersMutex.Lock()
	for i, curL := range this.listeners {
		if curL == l {
			this.listeners = append(this.listeners[:i], this.listeners[i+1:]...)
			break
		}
	}
	this.listenersMutex.Unlock()
}
func (this *syncRegListenStore) Notify(e apireg.RegistrationEvent) {
	go func() {
		this.listenersMutex.RLock()
		for _, curL := range this.listeners {
			curL.HandleRegistration(e)
		}
		this.listenersMutex.RUnlock()
	}()
}
