package multicast

import (
	"net"
	"testing"

	"github.com/ZacharyDuve/apireg"
	"github.com/google/uuid"
)

func TestThatNewSyncApiStoreIsEmpty(t *testing.T) {
	s := newSyncApiStore()

	if len(s.All()) != 0 {
		t.Fail()
	}
}

func TestThatAddingAnApiToAnApiStoreAddsIt(t *testing.T) {
	s := newSyncApiStore()

	lenBefore := len(s.All())
	a, _ := apireg.NewApi("Something", apireg.NewVersion(0, 0, 1), uuid.New(), apireg.All, net.ParseIP("127.0.0.1"), 8712)
	s.Add(a)

	if lenBefore == len(s.All()) {
		t.Fail()
	}
}

func TestThatAddingTheSameApiTwiceOnlyAddsItReallyOnce(t *testing.T) {
	s := newSyncApiStore()

	a, _ := apireg.NewApi("Something", apireg.NewVersion(0, 0, 1), uuid.New(), apireg.All, net.ParseIP("127.0.0.1"), 8712)
	s.Add(a)
	s.Add(a)
	if 1 != len(s.All()) {
		t.Fail()
	}
}

func TestThatRemovingAnApiFromStoreThatDoesntContainDoesNothing(t *testing.T) {
	s := newSyncApiStore()

	lenBefore := len(s.All())
	a, _ := apireg.NewApi("Something", apireg.NewVersion(0, 0, 1), uuid.New(), apireg.All, net.ParseIP("127.0.0.1"), 8712)
	s.Remove(a)

	if lenBefore != len(s.All()) {
		t.Fail()
	}
}

func TestThatRemovingAnApiFromAStoreThatContainsItRemovesIt(t *testing.T) {
	s := newSyncApiStore()

	a, _ := apireg.NewApi("Something", apireg.NewVersion(0, 0, 1), uuid.New(), apireg.All, net.ParseIP("127.0.0.1"), 8712)
	s.Add(a)
	lenBefore := len(s.All())
	s.Remove(a)

	if lenBefore <= len(s.All()) {
		t.Fail()
	}
}
