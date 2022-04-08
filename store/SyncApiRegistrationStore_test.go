package store

import (
	"net"
	"testing"
	"time"

	"github.com/ZacharyDuve/apireg/api"
)

func TestNewSyncApiRegistrationStoreWithTickerReturnsStore(t *testing.T) {
	ticker := get30sTicker()

	if NewSyncApiRegistrationStore(ticker) == nil {
		t.Fail()
	}
}

func TestThatNewSyncApiRegistrationStoreWithNoTickerReturnsStore(t *testing.T) {
	if NewSyncApiRegistrationStore(nil) == nil {
		t.Fail()
	}
}

func TestThatGetAllRegsReturnsEmptyListForNewSyncStore(t *testing.T) {
	store := NewSyncApiRegistrationStore(nil)
	allRegs := store.GetAllRegs()
	if allRegs == nil || len(allRegs) != 0 {
		t.Fail()
	}
}

func TestThatGetAllReturnsListOfLen1AfterAddingNewRegistration(t *testing.T) {
	store := NewSyncApiRegistrationStore(nil)

	store.AddReg(getValidApiReg())
	allRegs := store.GetAllRegs()
	if len(allRegs) != 1 {
		t.Fail()
	}
}

func TestThatAddingTheSameRegAgainDoesntAddAnotherRegistration(t *testing.T) {
	store := NewSyncApiRegistrationStore(nil)

	reg := getValidApiReg()

	store.AddReg(reg)
	store.AddReg(reg)

	allRegs := store.GetAllRegs()
	if len(allRegs) != 1 {
		t.Fail()
	}
}

func TestThatAddingAtLeastTwoUniqueRegsAddsAsMany(t *testing.T) {
	store := NewSyncApiRegistrationStore(nil)
	reg0 := getValidApiRegWithName("Steve")
	store.AddReg(reg0)
	reg1 := getValidApiRegWithName("Bob")
	store.AddReg(reg1)

	allRegs := store.GetAllRegs()
	if len(allRegs) != 2 {
		t.Fail()
	}
}

func TestThatAddingAtLeastTwoUniqueRegsWithSameNameAddsAsMany(t *testing.T) {
	store := NewSyncApiRegistrationStore(nil)
	name := "Jerry"
	var majVersion uint = 6
	reg0 := getValidApiRegWithName(name)
	reg0.Api().Version().Major = majVersion
	store.AddReg(reg0)
	reg1 := getValidApiRegWithName(name)
	reg1.Api().Version().Major = majVersion + 1
	store.AddReg(reg1)

	allRegs := store.GetAllRegs()
	if len(allRegs) != 2 {
		t.Fail()
	}
}

func TestThatRemovingFromEmptyRegistrationStoreDoesNothing(t *testing.T) {
	store := NewSyncApiRegistrationStore(nil)

	reg := getValidApiReg()

	sizeBefore := len(store.GetAllRegs())

	store.RemoveRegForApi(reg.Api())

	if sizeBefore != len(store.GetAllRegs()) {
		t.Fail()
	}
}

func TestThatRemovingAnApiWithStoreContainingSameNameButDifferentVersionDoesNotRemoveExisting(t *testing.T) {
	store := NewSyncApiRegistrationStore(nil)
	name := "Jerry"
	var majVersion uint = 6
	reg0 := getValidApiRegWithName(name)
	reg0.Api().Version().Major = majVersion
	store.AddReg(reg0)
	reg1 := getValidApiRegWithName(name)
	reg1.Api().Version().Major = majVersion + 1
	store.RemoveRegForApi(reg1.Api())
	allRegs := store.GetAllRegs()
	if len(allRegs) != 1 {
		t.Fail()
	}
}

func TestThatRemovingAnApiFromStoreContainingItActuallyRemoves(t *testing.T) {
	store := NewSyncApiRegistrationStore(nil)
	reg := getValidApiReg()
	store.AddReg(reg)
	sizeBefore := len(store.GetAllRegs())
	store.RemoveRegForApi(reg.Api())
	if sizeBefore-1 != len(store.GetAllRegs()) {
		t.Fail()
	}
}

func TestThatStoreContainingMultipleRegsForSameNameOnlyRemovesOneWhileKeepingRest(t *testing.T) {
	store := NewSyncApiRegistrationStore(nil)
	name := "Jerry"
	var majVersion uint = 6
	reg0 := getValidApiRegWithName(name)
	reg0.Api().Version().Major = majVersion
	store.AddReg(reg0)
	reg1 := getValidApiRegWithName(name)
	reg1.Api().Version().Major = majVersion + 1
	store.AddReg(reg1)
	store.RemoveRegForApi(reg1.Api())
	allRegs := store.GetAllRegsForName(name)
	if len(allRegs) != 1 {
		t.Fail()
	}
}

func TestThatGetAllForNameFiltersOutExpiredRegistrations(t *testing.T) {
	store := NewSyncApiRegistrationStore(nil)

	name := "Jerry"
	api, _ := api.NewApi(name, &api.Version{}, net.ParseIP("192.168.0.3"), 8672)
	life := time.Second * 2
	//Make it reged before now - life so it should be expired
	timeReged := time.Now().Add(-1 * (life + time.Second*1))

	reg, _ := NewApiRegistration(api, timeReged, life)
	store.AddReg(reg)
	regs := store.GetAllRegsForName(name)

	if len(regs) != 0 {
		t.Fail()
	}
}

func TestThatExpiredRegsArePurgedWhenPurgeExpiredIsCalled(t *testing.T) {
	now := time.Now()
	life := time.Second * 2
	//Make so it has expired already
	regTime := now.Add(-1 * (life + time.Second*1))
	api := getValidApi()
	reg, _ := NewApiRegistration(api, regTime, life)
	purgeTickChan := make(chan time.Time)
	store := NewSyncApiRegistrationStore(purgeTickChan)

	store.AddReg(reg)

	purgeTickChan <- now

	if len(store.GetAllRegs()) != 0 {
		t.Fail()
	}
}

func getValidApiReg() ApiRegistration {
	reg, _ := NewApiRegistration(getValidApi(), time.Now(), time.Second*15)

	return reg
}

func getValidApiRegWithName(name string) ApiRegistration {
	if name == "" {
		return getValidApiReg()
	}
	api, _ := api.NewApi(name, &api.Version{}, net.ParseIP("192.168.0.3"), 8323)
	retReg, _ := NewApiRegistration(api, time.Now(), time.Second*15)

	return retReg
}

func get30sTicker() <-chan time.Time {
	return time.NewTicker(time.Second * 30).C
}
