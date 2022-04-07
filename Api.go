package apireg

import "net"

type Api interface {
	Name() string
	Version() *Version
	HostIP() net.IP
	HostPort() int
}

type apiImpl struct {
	name       string
	version    *Version
	remoteIP   net.IP
	remotePort int
}

func (this *apiImpl) Name() string {
	return this.name
}
func (this *apiImpl) Version() *Version {
	return this.version
}
func (this *apiImpl) HostIP() net.IP {
	return this.remoteIP
}
func (this *apiImpl) HostPort() int {
	return this.remotePort
}
