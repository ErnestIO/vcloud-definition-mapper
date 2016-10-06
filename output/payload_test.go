/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import (
	"net"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func testBuildInstances(count int) []Instance {
	var instances []Instance

	for i := 1; i < count+1; i++ {
		instances = append(instances, Instance{
			Name:        "test-" + strconv.Itoa(i),
			Catalog:     "catalog",
			Image:       "image",
			Cpus:        1,
			Memory:      10240,
			NetworkName: "network",
			IP:          net.ParseIP("10.64.0.10" + strconv.Itoa(i)),
		})
	}

	return instances
}

func testBuildNetworks(count int) []Network {
	var networks []Network

	for i := 1; i < count+1; i++ {
		networks = append(networks, Network{
			Name:   "test-" + strconv.Itoa(i),
			Subnet: "10." + strconv.Itoa(i) + ".0.0/24",
		})
	}

	return networks
}

func testBuildExecution(count int) Execution {
	var names []string

	execution := Execution{
		Name:    "Execution",
		Payload: "date",
		Prefix:  "test",
	}

	for i := 1; i < count+1; i++ {
		names = append(names, "test-"+strconv.Itoa(i))
	}

	execution.Target = "list:" + strings.Join(names, ",")

	return execution
}

func TestPayloadDiff(t *testing.T) {
	Convey("Given a payload that contains only instances", t, func() {
		var m FSMMessage
		m.Instances.Items = testBuildInstances(2)

		Convey("When diffing a fsm message ", func() {
			Convey("And there is no previous result", func() {
				var p FSMMessage
				m.Diff(p)
				Convey("Then instances to create should be populated", func() {
					So(len(m.InstancesToCreate.Items), ShouldEqual, 2)
					So(m.InstancesToCreate.Items[0].Name, ShouldEqual, "test-1")
					So(m.InstancesToCreate.Items[0].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToCreate.Items[0].Image, ShouldEqual, "image")
					So(m.InstancesToCreate.Items[0].Cpus, ShouldEqual, 1)
					So(m.InstancesToCreate.Items[0].Memory, ShouldEqual, 10240)
					So(m.InstancesToCreate.Items[0].NetworkName, ShouldEqual, "network")
					So(m.InstancesToCreate.Items[0].IP.String(), ShouldEqual, "10.64.0.101")
					So(m.InstancesToCreate.Items[1].Name, ShouldEqual, "test-2")
					So(m.InstancesToCreate.Items[1].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToCreate.Items[1].Image, ShouldEqual, "image")
					So(m.InstancesToCreate.Items[1].Cpus, ShouldEqual, 1)
					So(m.InstancesToCreate.Items[1].Memory, ShouldEqual, 10240)
					So(m.InstancesToCreate.Items[1].NetworkName, ShouldEqual, "network")
					So(m.InstancesToCreate.Items[1].IP.String(), ShouldEqual, "10.64.0.102")
				})
			})
			Convey("And there is a previous result with one instance", func() {
				var p FSMMessage
				p.Instances.Items = testBuildInstances(1)

				m.Diff(p)
				Convey("Then instances to create should be populated with the new instance", func() {
					So(len(m.InstancesToCreate.Items), ShouldEqual, 1)
					So(m.InstancesToCreate.Items[0].Name, ShouldEqual, "test-2")
					So(m.InstancesToCreate.Items[0].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToCreate.Items[0].Image, ShouldEqual, "image")
					So(m.InstancesToCreate.Items[0].Cpus, ShouldEqual, 1)
					So(m.InstancesToCreate.Items[0].Memory, ShouldEqual, 10240)
					So(m.InstancesToCreate.Items[0].NetworkName, ShouldEqual, "network")
					So(m.InstancesToCreate.Items[0].IP.String(), ShouldEqual, "10.64.0.102")
				})
			})
			Convey("And there is one less instance on the new input", func() {
				var p FSMMessage
				p.Instances.Items = testBuildInstances(2)
				m.Instances.Items = testBuildInstances(1)

				m.Diff(p)
				Convey("Then instances to delete should be populated with the old instance", func() {
					So(len(m.InstancesToDelete.Items), ShouldEqual, 1)
					So(m.InstancesToDelete.Items[0].Name, ShouldEqual, "test-2")
					So(m.InstancesToDelete.Items[0].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToDelete.Items[0].Image, ShouldEqual, "image")
					So(m.InstancesToDelete.Items[0].Cpus, ShouldEqual, 1)
					So(m.InstancesToDelete.Items[0].Memory, ShouldEqual, 10240)
					So(m.InstancesToDelete.Items[0].NetworkName, ShouldEqual, "network")
					So(m.InstancesToDelete.Items[0].IP.String(), ShouldEqual, "10.64.0.102")
				})
			})
			Convey("And the resource is updated on the instances", func() {
				var p FSMMessage
				p.Instances.Items = testBuildInstances(2)
				p.Instances.Items[0].Cpus = 2
				p.Instances.Items[1].Cpus = 2

				m.Diff(p)
				Convey("Then instances to update should be populated with the updated instances", func() {
					So(len(m.InstancesToUpdate.Items), ShouldEqual, 2)
					So(m.InstancesToUpdate.Items[0].Name, ShouldEqual, "test-1")
					So(m.InstancesToUpdate.Items[0].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToUpdate.Items[0].Image, ShouldEqual, "image")
					So(m.InstancesToUpdate.Items[0].Cpus, ShouldEqual, 1)
					So(m.InstancesToUpdate.Items[0].Memory, ShouldEqual, 10240)
					So(m.InstancesToUpdate.Items[0].NetworkName, ShouldEqual, "network")
					So(m.InstancesToUpdate.Items[0].IP.String(), ShouldEqual, "10.64.0.101")
					So(m.InstancesToUpdate.Items[1].Name, ShouldEqual, "test-2")
					So(m.InstancesToUpdate.Items[1].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToUpdate.Items[1].Image, ShouldEqual, "image")
					So(m.InstancesToUpdate.Items[1].Cpus, ShouldEqual, 1)
					So(m.InstancesToUpdate.Items[1].Memory, ShouldEqual, 10240)
					So(m.InstancesToUpdate.Items[1].NetworkName, ShouldEqual, "network")
					So(m.InstancesToUpdate.Items[1].IP.String(), ShouldEqual, "10.64.0.102")
				})
			})
		})
	})

	Convey("Given a payload that contains executions", t, func() {
		var m FSMMessage
		m.Bootstrapping = "salt"
		m.ServiceIP = "127.0.0.1"
		m.SaltUser = "user"
		m.SaltPass = "pass"
		m.Datacenters.Items = append(m.Datacenters.Items, Datacenter{Type: "vcloud"})
		m.Instances.Items = testBuildInstances(2)
		m.Executions.Items = append(m.Executions.Items, testBuildExecution(2))

		Convey("When diffing a fsm message", func() {
			Convey("And there is no previous result", func() {
				var p FSMMessage

				m.Diff(p)
				Convey("Then it should create executions and bootstraps", func() {
					So(len(m.InstancesToCreate.Items), ShouldEqual, 2)
					So(m.InstancesToCreate.Items[0].Name, ShouldEqual, "test-1")
					So(m.InstancesToCreate.Items[0].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToCreate.Items[0].Image, ShouldEqual, "image")
					So(m.InstancesToCreate.Items[0].Cpus, ShouldEqual, 1)
					So(m.InstancesToCreate.Items[0].Memory, ShouldEqual, 10240)
					So(m.InstancesToCreate.Items[0].NetworkName, ShouldEqual, "network")
					So(m.InstancesToCreate.Items[0].IP.String(), ShouldEqual, "10.64.0.101")
					So(m.InstancesToCreate.Items[1].Name, ShouldEqual, "test-2")
					So(m.InstancesToCreate.Items[1].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToCreate.Items[1].Image, ShouldEqual, "image")
					So(m.InstancesToCreate.Items[1].Cpus, ShouldEqual, 1)
					So(m.InstancesToCreate.Items[1].Memory, ShouldEqual, 10240)
					So(m.InstancesToCreate.Items[1].NetworkName, ShouldEqual, "network")
					So(m.InstancesToCreate.Items[1].IP.String(), ShouldEqual, "10.64.0.102")
					So(len(m.BootstrapsToCreate.Items), ShouldEqual, 2)
					So(m.BootstrapsToCreate.Items[0].Name, ShouldEqual, "Bootstrap test-1")
					So(m.BootstrapsToCreate.Items[0].Payload, ShouldEqual, "/usr/bin/bootstrap -master 10.254.254.100 -host 10.64.0.101 -username ernest -password 'b00tStr4pp3rR' -max-retries 20 -minion-name test-1")
					So(m.BootstrapsToCreate.Items[0].Target, ShouldEqual, "list:salt-master.localdomain")
					So(m.BootstrapsToCreate.Items[1].Name, ShouldEqual, "Bootstrap test-2")
					So(m.BootstrapsToCreate.Items[1].Payload, ShouldEqual, "/usr/bin/bootstrap -master 10.254.254.100 -host 10.64.0.102 -username ernest -password 'b00tStr4pp3rR' -max-retries 20 -minion-name test-2")
					So(m.BootstrapsToCreate.Items[1].Target, ShouldEqual, "list:salt-master.localdomain")
					So(len(m.ExecutionsToCreate.Items), ShouldEqual, 1)
					So(m.ExecutionsToCreate.Items[0].Name, ShouldEqual, "Execution")
					So(m.ExecutionsToCreate.Items[0].Payload, ShouldEqual, "date")
					So(m.ExecutionsToCreate.Items[0].Target, ShouldEqual, "list:test-1,test-2")
				})
			})

			Convey("And only the execution payload has changed", func() {
				var p FSMMessage
				p.Bootstrapping = "salt"
				m.ServiceIP = "127.0.0.1"
				m.SaltUser = "user"
				m.SaltPass = "pass"
				m.Datacenters.Items = append(m.Datacenters.Items, Datacenter{Type: "vcloud"})
				p.Instances.Items = testBuildInstances(2)
				p.Executions.Items = append(p.Executions.Items, testBuildExecution(2))
				p.Executions.Items[0].Payload = "uptime"

				m.Diff(p)
				Convey("Then it should run the new executions on all instances", func() {
					So(len(m.InstancesToCreate.Items), ShouldEqual, 0)
					So(len(m.Bootstraps.Items), ShouldEqual, 0)
					So(len(m.ExecutionsToCreate.Items), ShouldEqual, 1)
					So(m.ExecutionsToCreate.Items[0].Name, ShouldEqual, "Execution")
					So(m.ExecutionsToCreate.Items[0].Payload, ShouldEqual, "date")
					So(m.ExecutionsToCreate.Items[0].Target, ShouldEqual, "list:test-1,test-2")
				})
			})

			Convey("And instance count has increased but the execution payload is the same", func() {
				var p FSMMessage
				p.Bootstrapping = "salt"
				m.ServiceIP = "127.0.0.1"
				m.SaltUser = "user"
				m.SaltPass = "pass"
				m.Datacenters.Items = append(m.Datacenters.Items, Datacenter{Type: "vcloud"})
				p.Instances.Items = testBuildInstances(1)
				p.Executions.Items = append(p.Executions.Items, testBuildExecution(1))

				m.Diff(p)
				Convey("Then it should create executions and bootstraps for only the new instance", func() {
					So(len(m.InstancesToCreate.Items), ShouldEqual, 1)
					So(len(m.BootstrapsToCreate.Items), ShouldEqual, 1)
					So(m.BootstrapsToCreate.Items[0].Name, ShouldEqual, "Bootstrap test-2")
					So(m.BootstrapsToCreate.Items[0].Payload, ShouldEqual, "/usr/bin/bootstrap -master 10.254.254.100 -host 10.64.0.102 -username ernest -password 'b00tStr4pp3rR' -max-retries 20 -minion-name test-2")
					So(m.BootstrapsToCreate.Items[0].Target, ShouldEqual, "list:salt-master.localdomain")
					So(len(m.ExecutionsToCreate.Items), ShouldEqual, 1)
					So(m.ExecutionsToCreate.Items[0].Name, ShouldEqual, "Execution")
					So(m.ExecutionsToCreate.Items[0].Payload, ShouldEqual, "date")
					So(m.ExecutionsToCreate.Items[0].Target, ShouldEqual, "list:test-2")
				})
			})

			Convey("And instance count has increased and the execution payload is different", func() {
				var p FSMMessage
				p.Bootstrapping = "salt"
				m.ServiceIP = "127.0.0.1"
				m.SaltUser = "user"
				m.SaltPass = "pass"
				m.Datacenters.Items = append(m.Datacenters.Items, Datacenter{Type: "vcloud"})
				p.Instances.Items = testBuildInstances(1)
				p.Executions.Items = append(p.Executions.Items, testBuildExecution(1))
				p.Executions.Items[0].Payload = "uptime"

				m.Diff(p)
				Convey("Then it should create executions and bootstraps", func() {
					So(len(m.InstancesToCreate.Items), ShouldEqual, 1)
					So(len(m.BootstrapsToCreate.Items), ShouldEqual, 1)
					So(m.BootstrapsToCreate.Items[0].Name, ShouldEqual, "Bootstrap test-2")
					So(m.BootstrapsToCreate.Items[0].Payload, ShouldEqual, "/usr/bin/bootstrap -master 10.254.254.100 -host 10.64.0.102 -username ernest -password 'b00tStr4pp3rR' -max-retries 20 -minion-name test-2")
					So(m.BootstrapsToCreate.Items[0].Target, ShouldEqual, "list:salt-master.localdomain")
					So(len(m.ExecutionsToCreate.Items), ShouldEqual, 1)
					So(m.ExecutionsToCreate.Items[0].Name, ShouldEqual, "Execution")
					So(m.ExecutionsToCreate.Items[0].Payload, ShouldEqual, "date")
					So(m.ExecutionsToCreate.Items[0].Target, ShouldEqual, "list:test-1,test-2")
				})
			})
		})
	})
}
