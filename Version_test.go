package apireg

import "testing"

func TestVersionEqualReturnsFalseWhenVersionNotEqual(t *testing.T) {
	v0 := NewVersion(1, 0, 20)
	v1 := NewVersion(0, 0, 2)

	if v0.Equal(v1) {
		t.Fail()
	}
}

func TestVersionEqualReturnsTrueWhenVersionEqual(t *testing.T) {
	v0 := NewVersion(0, 0, 2)
	v1 := NewVersion(0, 0, 2)

	if !v0.Equal(v1) {
		t.Fail()
	}
}

func TestVersionLessThanReturnsTrueWhenVersion0IsLowerThanVersion1(t *testing.T) {
	v0 := NewVersion(0, 0, 1)
	v1 := NewVersion(0, 0, 2)

	if !v0.LessThan(v1) {
		t.Fail()
	}
}

func TestVersionLessThanReturnsFalseWhenVersion0IsEquaToVersion1(t *testing.T) {
	v0 := NewVersion(0, 0, 1)
	v1 := NewVersion(0, 0, 1)

	if v0.LessThan(v1) {
		t.Fail()
	}
}

func TestVersionLessThanReturnsFalseWhenVersion0IsGreaterThanVersion1(t *testing.T) {
	v0 := NewVersion(0, 2, 0)
	v1 := NewVersion(0, 0, 1)

	if v0.LessThan(v1) {
		t.Fail()
	}
}

func TestVersionGreaterThanReturnsFalseWhenVersion0IsEquaToVersion1(t *testing.T) {
	v0 := NewVersion(0, 0, 1)
	v1 := NewVersion(0, 0, 1)

	if v0.GreaterThan(v1) {
		t.Fail()
	}
}

func TestVersionGreaterThanReturnsFalseWhenVersion0IsLowerThanVersion1(t *testing.T) {
	v0 := NewVersion(0, 0, 1)
	v1 := NewVersion(3, 0, 1)

	if v0.GreaterThan(v1) {
		t.Fail()
	}
}

func TestVersionGreaterThanReturnsTrueWhenVersion0IsHigherThanVersion1(t *testing.T) {
	v0 := NewVersion(1, 1, 1)
	v1 := NewVersion(1, 0, 1)

	if v0.GreaterThan(v1) {
		t.Fail()
	}
}

func TestThatVersionStringReturnsCorrect(t *testing.T) {
	v := NewVersion(12, 4, 0)

	if v.String() != "v12.4.0" {
		t.Fail()
	}
}
