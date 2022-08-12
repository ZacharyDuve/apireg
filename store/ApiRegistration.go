package store

import (
	"errors"
	"sync"
	"time"

	"git.zmanhobbies.com/software/apireg/api"
)

type ApiRegistration interface {
	Api() api.Api
	TimeRegistered() time.Time
	UpdateTimeRegistered(time.Time)
	LifeSpan() time.Duration
	Expired(time.Time) bool
}

type apiRegistration struct {
	api                 api.Api
	timeRegistered      time.Time
	timeRegisteredMutex sync.Mutex
	lifeSpan            time.Duration
}

func NewApiRegistration(api api.Api, timeReged time.Time, lifeSpan time.Duration) (ApiRegistration, error) {
	if api == nil {
		return nil, errors.New("api is required for NewApiRegistration")
	}

	return &apiRegistration{api: api, timeRegistered: timeReged, lifeSpan: lifeSpan}, nil
}

func (this *apiRegistration) Api() api.Api {
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
