package godistributedapiregistry

import "testing"

func TestVersionEqualReturnsFalseWhenVersionNotEqual(t *testing.T) {
	v0 := Version{Major: 1, Minor: 0, BugFix: 20}
	v1 := Version{Major: 0, Minor: 0, BugFix: 2}

	if v0.Equal(v1) {
		t.Fail()
	}
}

func TestVersionEqualReturnsTrueWhenVersionEqual(t *testing.T) {
	v0 := Version{Major: 0, Minor: 0, BugFix: 2}
	v1 := Version{Major: 0, Minor: 0, BugFix: 2}

	if !v0.Equal(v1) {
		t.Fail()
	}
}

func TestVersionLessThanReturnsTrueWhenVersion0IsLowerThanVersion1(t *testing.T) {
	v0 := Version{Major: 0, Minor: 0, BugFix: 1}
	v1 := Version{Major: 0, Minor: 0, BugFix: 2}

	if !v0.LessThan(v1) {
		t.Fail()
	}
}

func TestVersionLessThanReturnsFalseWhenVersion0IsEquaToVersion1(t *testing.T) {
	v0 := Version{Major: 0, Minor: 0, BugFix: 1}
	v1 := Version{Major: 0, Minor: 0, BugFix: 1}

	if v0.LessThan(v1) {
		t.Fail()
	}
}

func TestVersionLessThanReturnsFalseWhenVersion0IsGreaterThanVersion1(t *testing.T) {
	v0 := Version{Major: 0, Minor: 2, BugFix: 0}
	v1 := Version{Major: 0, Minor: 0, BugFix: 1}

	if v0.LessThan(v1) {
		t.Fail()
	}
}

func TestVersionGreaterThanReturnsFalseWhenVersion0IsEquaToVersion1(t *testing.T) {
	v0 := Version{Major: 0, Minor: 0, BugFix: 1}
	v1 := Version{Major: 0, Minor: 0, BugFix: 1}

	if v0.GreaterThan(v1) {
		t.Fail()
	}
}

func TestVersionGreaterThanReturnsFalseWhenVersion0IsLowerThanVersion1(t *testing.T) {
	v0 := Version{Major: 0, Minor: 0, BugFix: 1}
	v1 := Version{Major: 3, Minor: 0, BugFix: 1}

	if v0.GreaterThan(v1) {
		t.Fail()
	}
}

func TestVersionGreaterThanReturnsTrueWhenVersion0IsHigherThanVersion1(t *testing.T) {
	v0 := Version{Major: 1, Minor: 1, BugFix: 1}
	v1 := Version{Major: 1, Minor: 0, BugFix: 1}

	if v0.GreaterThan(v1) {
		t.Fail()
	}
}
