package multicast

import (
	"errors"
	"sync"
	"time"

	"github.com/ZacharyDuve/apireg"
)

type apiRegistration struct {
	api                 apireg.Api
	timeRegistered      time.Time
	timeRegisteredMutex sync.Mutex
	lifeSpan            time.Duration
}

func newApiRegistration(api apireg.Api, timeReged time.Time, lifeSpan time.Duration) (*apiRegistration, error) {
	if api == nil {
		return nil, errors.New("api is required for NewApiRegistration")
	}

	return &apiRegistration{api: api, timeRegistered: timeReged, lifeSpan: lifeSpan}, nil
}

func (this *apiRegistration) Api() apireg.Api {
	return this.api
}
func (this *apiRegistration) TimeRegistered() time.Time {
	return this.timeRegistered
}

func (this *apiRegistration) UpdateTimeRegistered(newTime time.Time) {
	this.timeRegisteredMutex.Lock()
	this.timeRegistered = newTime
	this.timeRegisteredMutex.Unlock()
}

func (this *apiRegistration) LifeSpan() time.Duration {
	return this.lifeSpan
}
func (this *apiRegistration) Expired(otherTime time.Time) bool {
	return this.timeRegistered.Add(this.lifeSpan).Before(otherTime)
}
