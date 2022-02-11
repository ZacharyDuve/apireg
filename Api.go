package godistributedapiregistry

import "net"

type Api interface {
	Name() string
	Version() string
	HostIP() net.IP
	HostPort() int
}

type apiImpl struct {
	name       string
	version    string
	remoteIP   net.IP
	remotePort int
}

func (this *apiImpl) Name() string {
	return this.name
}
func (this *apiImpl) Version() string {
	return this.version
}
func (this *apiImpl) HostIP() net.IP {
	return this.remoteIP
}
func (this *apiImpl) HostPort() int {
	return this.remotePort
}
