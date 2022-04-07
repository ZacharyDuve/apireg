package godistributedapiregistry

import "fmt"

type Version struct {
	Major  uint
	Minor  uint
	BugFix uint
}

func (this *Version) Equal(other Version) bool {
	return this.Major == other.Major && this.Minor == other.Minor && this.BugFix == other.BugFix
}

func (this *Version) LessThan(other Version) bool {
	return this.Major <= other.Major && this.Minor <= other.Minor && this.BugFix < other.BugFix
}

func (this *Version) GreaterThan(other Version) bool {
	return this.Major >= other.Major && this.Minor >= other.Minor && this.BugFix > other.BugFix
}

func (this *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", this.Major, this.Minor, this.BugFix)
}
