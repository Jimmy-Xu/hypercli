package main

import (
	"os"
	"fmt"
	"testing"

	"github.com/docker/docker/pkg/reexec"
	"github.com/go-check/check"
)

func Test(t *testing.T) {
	reexec.Init() // This is required for external graphdriver tests

	if !isLocalDaemon {
		fmt.Printf("INFO: Testing against a remote daemon(%v)\n",os.Getenv("DOCKER_HOST"))
	} else {
		fmt.Println("INFO: Testing against a local daemon")
	}

	check.TestingT(t)
}

func init() {
	check.Suite(&DockerSuite{})
}

type DockerSuite struct {
}

//status only support : created, restarting, running, exited (https://github.com/getdvm/hyper-api-router/blob/master/pkg/apiserver/router/local/container.go#L204)
func (s *DockerSuite) TearDownTest(c *check.C) {
	//unpauseAllContainers()
	for _, region := range []string{"us-west-1","eu-central-1"} {
		deleteAllContainers(region)
		deleteAllImages(region)
		deleteAllSnapshots(region)
		deleteAllVolumes(region)
		deleteAllFips(region)
	}
	//deleteAllNetworks()
}

func init() {
	check.Suite(&DockerRegistrySuite{
		ds: &DockerSuite{},
	})

	fmt.Printf("INFO: clear resources(containers, images, snapshots, volumes, fips)\n")
	for _, region := range []string{"us-west-1","eu-central-1"} {
		deleteAllContainers(region)
		deleteAllImages(region)
		deleteAllSnapshots(region)
		deleteAllVolumes(region)
		deleteAllFips(region)
	}

	fmt.Println("INFO: finish init")
}

type DockerRegistrySuite struct {
	ds  *DockerSuite
	reg *testRegistryV2
	//d   *Daemon
}

func (s *DockerRegistrySuite) SetUpTest(c *check.C) {
	testRequires(c, DaemonIsLinux, RegistryHosting)
	s.reg = setupRegistry(c, false, false)
	//s.d = NewDaemon(c)
}

func (s *DockerRegistrySuite) TearDownTest(c *check.C) {
	if s.reg != nil {
		s.reg.Close()
	}
	//if s.d != nil {
	//	s.d.Stop()
	//}
	s.ds.TearDownTest(c)
}

func init() {
	check.Suite(&DockerSchema1RegistrySuite{
		ds: &DockerSuite{},
	})
}

type DockerSchema1RegistrySuite struct {
	ds  *DockerSuite
	reg *testRegistryV2
	//d   *Daemon
}

func (s *DockerSchema1RegistrySuite) SetUpTest(c *check.C) {
	testRequires(c, DaemonIsLinux, RegistryHosting)
	s.reg = setupRegistry(c, true, false)
	//s.d = NewDaemon(c)
}

func (s *DockerSchema1RegistrySuite) TearDownTest(c *check.C) {
	if s.reg != nil {
		s.reg.Close()
	}
	//if s.d != nil {
	//	s.d.Stop()
	//}
	s.ds.TearDownTest(c)
}

func init() {
	check.Suite(&DockerRegistryAuthSuite{
		ds: &DockerSuite{},
	})
}

type DockerRegistryAuthSuite struct {
	ds  *DockerSuite
	reg *testRegistryV2
	//d   *Daemon
}

func (s *DockerRegistryAuthSuite) SetUpTest(c *check.C) {
	testRequires(c, DaemonIsLinux, RegistryHosting)
	s.reg = setupRegistry(c, false, true)
	//s.d = NewDaemon(c)
}

func (s *DockerRegistryAuthSuite) TearDownTest(c *check.C) {
	if s.reg != nil {
		//out, err := s.d.Cmd("logout", privateRegistryURL)
		//c.Assert(err, check.IsNil, check.Commentf(out))
		s.reg.Close()
	}
	//if s.d != nil {
	//	s.d.Stop()
	//}
	s.ds.TearDownTest(c)
}

func init() {
	check.Suite(&DockerDaemonSuite{
		ds: &DockerSuite{},
	})
}

type DockerDaemonSuite struct {
	ds *DockerSuite
	//d  *Daemon
}

func (s *DockerDaemonSuite) SetUpTest(c *check.C) {
	testRequires(c, DaemonIsLinux)
	//s.d = NewDaemon(c)
}

func (s *DockerDaemonSuite) TearDownTest(c *check.C) {
	testRequires(c, DaemonIsLinux)
	//if s.d != nil {
	//	s.d.Stop()
	//}
	s.ds.TearDownTest(c)
}

func init() {
	check.Suite(&DockerTrustSuite{
		ds: &DockerSuite{},
	})
}

type DockerTrustSuite struct {
	ds  *DockerSuite
	reg *testRegistryV2
	not *testNotary
}

func (s *DockerTrustSuite) SetUpTest(c *check.C) {
	testRequires(c, RegistryHosting, NotaryHosting)
	s.reg = setupRegistry(c, false, false)
	s.not = setupNotary(c)
}

func (s *DockerTrustSuite) TearDownTest(c *check.C) {
	if s.reg != nil {
		s.reg.Close()
	}
	if s.not != nil {
		s.not.Close()
	}
	s.ds.TearDownTest(c)
}
