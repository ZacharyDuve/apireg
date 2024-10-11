package multicast

import (
	"sync"
	"time"

	"github.com/ZacharyDuve/apireg/api"
	"github.com/ZacharyDuve/apireg/apievent"
)

type syncApiRegStore struct {
	regs          map[string][]ApiRegistration
	regsMutex     *sync.RWMutex
	purgeTickChan <-chan time.Time
	listeners     RegistrationListenerStore
}

func NewSyncApiRegistrationStore(pChan <-chan time.Time) ApiRegistrationStore {
	syncStore := &syncApiRegStore{}
	syncStore.regs = make(map[string][]ApiRegistration)
	syncStore.regsMutex = &sync.RWMutex{}
	syncStore.listeners = NewSyncRegistrationListenerStore()
	//if we never provide a channel then auto purging is disabled
	if pChan != nil {
		syncStore.purgeTickChan = pChan
		go syncStore.purgeLoop()
	}

	return syncStore
}

func (this *syncApiRegStore) AddReg(reg ApiRegistration) {
	this.regsMutex.Lock()
	apis, contains := this.regs[reg.Api().Name()]
	added := false

	if !contains {
		apis = make([]ApiRegistration, 1)
		apis[0] = reg
		this.regs[reg.Api().Name()] = apis
		added = true
	} else {
		hasMatch := false
		for _, curReg := range apis {
			if apisMatch(reg.Api(), curReg.Api()) {
				hasMatch = true
				break
			}
		}

		if !hasMatch {
			apis = append(apis, reg)
			this.regs[reg.Api().Name()] = apis
			added = true
		}
	}
	if added {
		this.listeners.Notify(apievent.NewAddEvent(reg.Api()))
	}
	this.regsMutex.Unlock()

}

func apisMatch(api0, api1 api.Api) bool {
	return api0.Name() == api1.Name() &&
		api0.Version().Equal(api1.Version()) &&
		api0.HostIP().Equal(api1.HostIP()) &&
		api0.HostPort() == api1.HostPort()
}

func (this *syncApiRegStore) GetAllRegsForName(name string) []ApiRegistration {
	return this.getAllRegsForNameAndTime(name, time.Now())
}

func (this *syncApiRegStore) getAllRegsForNameAndTime(name string, time time.Time) []ApiRegistration {
	var matchingApis []ApiRegistration
	this.regsMutex.RLock()
	regs, contains := this.regs[name]
	this.regsMutex.RUnlock()

	if contains {
		matchingApis = make([]ApiRegistration, 0, len(regs))

		for _, curReg := range regs {
			if curReg.Expired(time) {
				this.RemoveRegForApi(curReg.Api())
			} else {
				matchingApis = append(matchingApis, curReg)
			}
		}
	}

	return matchingApis
}

func (this *syncApiRegStore) GetAllRegs() []ApiRegistration {
	return this.getAllRegsForTime(time.Now())
}

func (this *syncApiRegStore) getAllRegsForTime(t time.Time) []ApiRegistration {
	//Pulling list of names first from regs so we can release lock from Read mode as GetAllRegs could request lock for Write mode for an expired record
	this.regsMutex.RLock()
	regNames := make([]string, 0, len(this.regs))
	for curName := range this.regs {
		regNames = append(regNames, curName)
	}
	this.regsMutex.RUnlock()

	regs := make([]ApiRegistration, 0)
	for _, curName := range regNames {
		regs = append(regs, this.getAllRegsForNameAndTime(curName, t)...)
	}
	return regs
}

func (this *syncApiRegStore) RemoveRegForApi(old api.Api) error {
	this.regsMutex.Lock()
	apis, contains := this.regs[old.Name()]

	if contains {
		if len(apis) == 1 && apisMatch(old, apis[0].Api()) {
			delete(this.regs, old.Name())
		} else {
			for i, curReg := range apis {
				if apisMatch(old, curReg.Api()) {
					apis = append(apis[:i], apis[i+1:]...)
					this.regs[old.Name()] = apis
					break
				}
			}
		}
		rEvent := apievent.NewRemovedEvent(old)
		this.listeners.Notify(rEvent)
	}
	this.regsMutex.Unlock()
	return nil
}

func (this *syncApiRegStore) purgeLoop() {
	for t := range this.purgeTickChan {
		this.purgeExpired(t)
	}
}

func (this *syncApiRegStore) purgeExpired(t time.Time) {
	//Pulling list of names first from regs so we can release lock from Read mode as GetAllRegs could request lock for Write mode for an expired record
	this.regsMutex.RLock()
	regNames := make([]string, 0, len(this.regs))
	for curName := range this.regs {
		regNames = append(regNames, curName)
	}
	this.regsMutex.RUnlock()

	for _, curName := range regNames {
		this.purgeExpiredForNameAndTime(curName, t)
	}
}

func (this *syncApiRegStore) purgeExpiredForNameAndTime(name string, t time.Time) {
	this.regsMutex.RLock()
	regs, contains := this.regs[name]
	this.regsMutex.RUnlock()

	if contains {
		for _, curReg := range regs {
			if curReg.Expired(t) {
				this.RemoveRegForApi(curReg.Api())
			}
		}
	}
}

func (this *syncApiRegStore) AddListener(l apievent.RegistrationListener) {
	this.listeners.Add(l)
}
func (this *syncApiRegStore) RemoveListener(l apievent.RegistrationListener) {
	this.listeners.Remove(l)
}
