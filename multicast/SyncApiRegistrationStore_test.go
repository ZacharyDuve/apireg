package multicast

import (
	"net"
	"testing"
	"time"

	"github.com/ZacharyDuve/apireg"
	"github.com/google/uuid"
)

func TestNewSyncApiRegistrationStoreWithTickerReturnsStore(t *testing.T) {
	ticker := get30sTicker()

	if newSyncApiRegistrationStore(ticker) == nil {
		t.Fail()
	}
}

func TestThatNewSyncApiRegistrationStoreWithNoTickerReturnsStore(t *testing.T) {
	if newSyncApiRegistrationStore(nil) == nil {
		t.Fail()
	}
}

func TestThatGetAllRegsReturnsEmptyListForNewSyncStore(t *testing.T) {
	store := newSyncApiRegistrationStore(nil)
	allRegs := store.GetAllRegs()
	if allRegs == nil || len(allRegs) != 0 {
		t.Fail()
	}
}

func TestThatGetAllReturnsListOfLen1AfterAddingNewRegistration(t *testing.T) {
	store := newSyncApiRegistrationStore(nil)

	store.AddReg(getValidApiReg())
	allRegs := store.GetAllRegs()
	if len(allRegs) != 1 {
		t.Fail()
	}
}

func TestThatAddingTheSameRegAgainDoesntAddAnotherRegistration(t *testing.T) {
	store := newSyncApiRegistrationStore(nil)

	reg := getValidApiReg()

	store.AddReg(reg)
	store.AddReg(reg)

	allRegs := store.GetAllRegs()
	if len(allRegs) != 1 {
		t.Fail()
	}
}

func TestThatAddingAtLeastTwoUniqueRegsAddsAsMany(t *testing.T) {
	store := newSyncApiRegistrationStore(nil)
	reg0 := getValidApiRegWithNameAndVersion("Steve", apireg.NewVersion(1, 0, 0))
	store.AddReg(reg0)
	reg1 := getValidApiRegWithNameAndVersion("Bob", apireg.NewVersion(1, 0, 0))
	store.AddReg(reg1)

	allRegs := store.GetAllRegs()
	if len(allRegs) != 2 {
		t.Fail()
	}
}

func TestThatAddingAtLeastTwoUniqueRegsWithSameNameAddsAsMany(t *testing.T) {
	store := newSyncApiRegistrationStore(nil)
	name := "Jerry"
	var majVersion uint = 6
	reg0 := getValidApiRegWithNameAndVersion(name, apireg.NewVersion(majVersion, 0, 0))
	store.AddReg(reg0)
	reg1 := getValidApiRegWithNameAndVersion(name, apireg.NewVersion(majVersion+1, 0, 0))
	store.AddReg(reg1)

	allRegs := store.GetAllRegs()
	if len(allRegs) != 2 {
		t.Fail()
	}
}

func TestThatRemovingFromEmptyRegistrationStoreDoesNothing(t *testing.T) {
	store := newSyncApiRegistrationStore(nil)

	reg := getValidApiReg()

	sizeBefore := len(store.GetAllRegs())

	store.RemoveRegForApi(reg.Api())

	if sizeBefore != len(store.GetAllRegs()) {
		t.Fail()
	}
}

func TestThatRemovingAnApiWithStoreContainingSameNameButDifferentVersionDoesNotRemoveExisting(t *testing.T) {
	store := newSyncApiRegistrationStore(nil)
	name := "Jerry"
	var majVersion uint = 6
	reg0 := getValidApiRegWithNameAndVersion(name, apireg.NewVersion(majVersion, 0, 0))
	store.AddReg(reg0)
	reg1 := getValidApiRegWithNameAndVersion(name, apireg.NewVersion(majVersion+1, 0, 0))
	store.RemoveRegForApi(reg1.Api())
	allRegs := store.GetAllRegs()
	if len(allRegs) != 1 {
		t.Fail()
	}
}

func TestThatRemovingAnApiFromStoreContainingItActuallyRemoves(t *testing.T) {
	store := newSyncApiRegistrationStore(nil)
	reg := getValidApiReg()
	store.AddReg(reg)
	sizeBefore := len(store.GetAllRegs())
	store.RemoveRegForApi(reg.Api())
	if sizeBefore-1 != len(store.GetAllRegs()) {
		t.Fail()
	}
}

func TestThatStoreContainingMultipleRegsForSameNameOnlyRemovesOneWhileKeepingRest(t *testing.T) {
	store := newSyncApiRegistrationStore(nil)
	name := "Jerry"
	var majVersion uint = 6
	reg0 := getValidApiRegWithNameAndVersion(name, apireg.NewVersion(majVersion, 0, 0))
	store.AddReg(reg0)
	reg1 := getValidApiRegWithNameAndVersion(name, apireg.NewVersion(majVersion+1, 0, 0))
	store.AddReg(reg1)
	store.RemoveRegForApi(reg1.Api())
	allRegs := store.GetAllRegsForName(name)
	if len(allRegs) != 1 {
		t.Fail()
	}
}

func TestThatGetAllForNameFiltersOutExpiredRegistrations(t *testing.T) {
	store := newSyncApiRegistrationStore(nil)

	name := "Jerry"
	api, _ := apireg.NewApi(name, apireg.NewVersion(0, 0, 1), uuid.New(), apireg.All, net.ParseIP("192.168.0.3"), 8672)
	life := time.Second * 2
	//Make it reged before now - life so it should be expired
	timeReged := time.Now().Add(-1 * (life + time.Second*1))

	reg, _ := newApiRegistration(api, timeReged, life)
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
	reg, _ := newApiRegistration(api, regTime, life)
	purgeTickChan := make(chan time.Time)
	store := newSyncApiRegistrationStore(purgeTickChan)

	store.AddReg(reg)

	purgeTickChan <- now

	if len(store.GetAllRegs()) != 0 {
		t.Fail()
	}
}

func getValidApiReg() *apiRegistration {
	reg, _ := newApiRegistration(getValidApi(), time.Now(), time.Second*15)

	return reg
}

func getValidApiRegWithNameAndVersion(name string, version apireg.Version) *apiRegistration {
	if name == "" {
		return getValidApiReg()
	}
	api, _ := apireg.NewApi(name, version, uuid.New(), apireg.All, net.ParseIP("192.168.0.3"), 8323)
	retReg, _ := newApiRegistration(api, time.Now(), time.Second*15)

	return retReg
}

func get30sTicker() <-chan time.Time {
	return time.NewTicker(time.Second * 30).C
}
