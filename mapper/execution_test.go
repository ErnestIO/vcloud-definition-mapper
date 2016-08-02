/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"net"
	"testing"

	"github.com/r3labs/vcloud-definition-mapper/definition"
	. "github.com/smartystreets/goconvey/convey"
)

func TestExecutionsMapping(t *testing.T) {
	Convey("Given a valid input definition", t, func() {
		d := definition.Definition{
			Name:       "service",
			Datacenter: "datacenter",
		}
		r := definition.Router{
			Name: "test",
		}

		r.Networks = append(r.Networks, definition.Network{
			Name:   "bar",
			Subnet: "10.0.0.0/24",
		})

		d.Routers = append(d.Routers, r)

		d.Instances = append(d.Instances, definition.Instance{
			Name:   "foo",
			Image:  "catalog/image",
			Cpus:   1,
			Memory: "2GB",
			Count:  2,
			Networks: definition.InstanceNetworks{
				Name:    "bar",
				StartIP: net.ParseIP("10.1.0.1"),
			},
			Disks: []string{"10GB"},
			Provisioner: []definition.Exec{
				definition.Exec{Commands: []string{"date", "uptime"}},
			},
		})

		Convey("When i try to map executions", func() {
			d.Bootstrapping = "salt"
			e := MapExecutions(d)

			Convey("Then all executions should be mapped", func() {
				So(len(e), ShouldEqual, 1)
				So(e[0].Name, ShouldEqual, "Execution foo 1")
				So(e[0].Type, ShouldEqual, "salt")
				So(e[0].Target, ShouldEqual, "list:datacenter-service-foo-1,datacenter-service-foo-2")
				So(e[0].Payload, ShouldEqual, "date; uptime")
			})

		})
	})
}
