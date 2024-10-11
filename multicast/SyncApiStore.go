package multicast

import (
	"sync"

	"github.com/ZacharyDuve/apireg"
)

type syncApiStore struct {
	apis      []apireg.Api
	apisMutex *sync.RWMutex
}

func newSyncApiStore() *syncApiStore {
	s := &syncApiStore{}
	s.apis = make([]apireg.Api, 0)
	s.apisMutex = &sync.RWMutex{}

	return s
}

func (this *syncApiStore) Add(newApi apireg.Api) bool {
	var added bool
	contains := this.Contains(newApi)

	if !contains {
		this.apisMutex.Lock()
		this.apis = append(this.apis, newApi)
		this.apisMutex.Unlock()
	}
	return added
}

func (this *syncApiStore) All() []apireg.Api {
	this.apisMutex.RLock()
	apisCopy := make([]apireg.Api, len(this.apis))
	copy(apisCopy, this.apis)
	this.apisMutex.RUnlock()

	return this.apis
}

func (this *syncApiStore) Contains(a apireg.Api) bool {
	var contains bool
	this.apisMutex.RLock()
	for _, curApi := range this.apis {
		if curApi.Equal(a) {
			contains = true
			break
		}
	}
	this.apisMutex.RUnlock()
	return contains
}

func (this *syncApiStore) Remove(r apireg.Api) bool {
	var removed bool

	this.apisMutex.Lock()
	for i, curApi := range this.apis {
		if curApi.Equal(r) {
			this.apis = append(this.apis[:i], this.apis[i+1:]...)
			removed = true
			break
		}
	}
	this.apisMutex.Unlock()

	return removed
}
