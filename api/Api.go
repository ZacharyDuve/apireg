package api

import (
	"errors"
	"net"

	"git.zmanhobbies.com/software/apireg/environment"
)

type Api interface {
	Name() string
	Version() *Version
	HostIP() net.IP
	HostPort() int
	Equal(Api) bool
	Environment() environment.Environment
}

type apiImpl struct {
	name       string
	version    *Version
	remoteIP   net.IP
	remotePort int
	env        environment.Environment
}

func NewApi(name string, ver *Version, env environment.Environment, hostIP net.IP, port int) (Api, error) {
	if name == "" {
		return nil, errors.New("name is required for NewApi")
	} else if ver == nil {
		return nil, errors.New("ver (version) is required for NewApi")
	} else if hostIP == nil {
		return nil, errors.New("hostIP is required for NewApi")
	} else if port <= 0 {
		return nil, errors.New("port must be > 0 for NewApi")
	}

	return &apiImpl{name: name, version: ver, env: env, remoteIP: hostIP, remotePort: port}, nil
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
func (this *apiImpl) Equal(other Api) bool {
	return other != nil &&
		this.name == other.Name() &&
		this.version.Equal(other.Version()) &&
		this.remoteIP.Equal(other.HostIP()) &&
		this.remotePort == other.HostPort()
}

func (this *apiImpl) Environment() environment.Environment {
	return this.env
}
