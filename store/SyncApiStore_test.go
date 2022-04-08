package store

import (
	"net"
	"testing"

	"github.com/ZacharyDuve/apireg/api"
)

func TestThatNewSyncApiStoreIsEmpty(t *testing.T) {
	s := NewSyncApiStore()

	if len(s.All()) != 0 {
		t.Fail()
	}
}

func TestThatAddingAnApiToAnApiStoreAddsIt(t *testing.T) {
	s := NewSyncApiStore()

	lenBefore := len(s.All())
	a, _ := api.NewApi("Something", &api.Version{}, net.ParseIP("127.0.0.1"), 8712)
	s.Add(a)

	if lenBefore == len(s.All()) {
		t.Fail()
	}
}

func TestThatAddingTheSameApiTwiceOnlyAddsItReallyOnce(t *testing.T) {
	s := NewSyncApiStore()

	a, _ := api.NewApi("Something", &api.Version{}, net.ParseIP("127.0.0.1"), 8712)
	s.Add(a)
	s.Add(a)
	if 1 != len(s.All()) {
		t.Fail()
	}
}

func TestThatRemovingAnApiFromStoreThatDoesntContainDoesNothing(t *testing.T) {
	s := NewSyncApiStore()

	lenBefore := len(s.All())
	a, _ := api.NewApi("Something", &api.Version{}, net.ParseIP("127.0.0.1"), 8712)
	s.Remove(a)

	if lenBefore != len(s.All()) {
		t.Fail()
	}
}

func TestThatRemovingAnApiFromAStoreThatContainsItRemovesIt(t *testing.T) {
	s := NewSyncApiStore()

	a, _ := api.NewApi("Something", &api.Version{}, net.ParseIP("127.0.0.1"), 8712)
	s.Add(a)
	lenBefore := len(s.All())
	s.Remove(a)

	if lenBefore <= len(s.All()) {
		t.Fail()
	}
}
