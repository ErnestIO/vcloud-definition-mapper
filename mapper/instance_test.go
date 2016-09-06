/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"net"
	"testing"

	"github.com/ernestio/vcloud-definition-mapper/definition"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ernestio/vcloud-definition-mapper/output"
)

func TestInstancesMapping(t *testing.T) {
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
			Count:  1,
			Networks: definition.InstanceNetworks{
				Name:    "bar",
				StartIP: net.ParseIP("10.1.0.1"),
			},
			Disks: []string{"10GB"},
		})

		Convey("When i try to map instances", func() {
			Convey("And the input specifies bootstrap as salt", func() {
				d.Bootstrapping = "salt"
				i := MapInstances(d)

				Convey("Then an extra instance should be mapped", func() {
					So(len(i), ShouldEqual, 2)
					So(i[0].Name, ShouldEqual, "datacenter-service-salt-master")
					So(i[0].Catalog, ShouldEqual, "r3")
					So(i[0].Image, ShouldEqual, "r3-salt-master")
					So(i[0].Cpus, ShouldEqual, 1)
					So(i[0].Memory, ShouldEqual, 2048)
					So(i[0].Disks, ShouldResemble, []output.InstanceDisk{})
					So(i[0].NetworkName, ShouldEqual, "datacenter-service-salt")
					So(i[0].IP, ShouldResemble, net.ParseIP("10.254.254.100"))
				})

				Convey("Then an defined instances should be mapped", func() {
					So(i[1].Name, ShouldEqual, "datacenter-service-foo-1")
					So(i[1].Catalog, ShouldEqual, "catalog")
					So(i[1].Image, ShouldEqual, "image")
					So(i[1].Cpus, ShouldEqual, 1)
					So(i[1].Memory, ShouldEqual, 2048)
					So(i[1].Disks[0].Size, ShouldEqual, 10240)
					So(i[1].NetworkName, ShouldEqual, "datacenter-service-bar")
					So(i[1].IP, ShouldResemble, net.ParseIP("10.1.0.1"))
				})
			})

			Convey("And the input specifies bootstrap as not salt", func() {
				d.Bootstrapping = "none"
				i := MapInstances(d)
				Convey("Then defined instances should be mapped", func() {
					So(len(i), ShouldEqual, 1)
					So(i[0].Name, ShouldEqual, "datacenter-service-foo-1")
					So(i[0].Catalog, ShouldEqual, "catalog")
					So(i[0].Image, ShouldEqual, "image")
					So(i[0].Cpus, ShouldEqual, 1)
					So(i[0].Memory, ShouldEqual, 2048)
					So(i[0].Disks[0].Size, ShouldEqual, 10240)
					So(i[0].NetworkName, ShouldEqual, "datacenter-service-bar")
					So(i[0].IP, ShouldResemble, net.ParseIP("10.1.0.1"))
				})
			})
		})
	})
}
