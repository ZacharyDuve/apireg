package apireg

import "fmt"

type Version interface {
	Major() uint
	Minor() uint
	BugFix() uint
	Equal(Version) bool
	LessThan(Version) bool
	GreaterThan(Version) bool
	String() string
}

type version struct {
	MajorVal  uint
	MinorVal  uint
	BugFixVal uint
}

func NewVersion(major, minor, bugfix uint) Version {
	return &version{MajorVal: major, MinorVal: minor, BugFixVal: bugfix}
}

func (this *version) Major() uint {
	return this.MajorVal
}

func (this *version) Minor() uint {
	return this.MinorVal
}

func (this *version) BugFix() uint {
	return this.BugFixVal
}

func (this *version) Equal(other Version) bool {
	return this.MajorVal == other.Major() && this.MinorVal == other.Minor() && this.BugFixVal == other.BugFix()
}

func (this *version) LessThan(other Version) bool {
	return this.MajorVal <= other.Major() && this.MinorVal <= other.Minor() && this.BugFixVal < other.BugFix()
}

func (this *version) GreaterThan(other Version) bool {
	return this.MajorVal >= other.Major() && this.MinorVal >= other.Minor() && this.BugFixVal > other.BugFix()
}

func (this *version) String() string {
	return fmt.Sprintf("v%d.%d.%d", this.MajorVal, this.MinorVal, this.BugFixVal)
}
