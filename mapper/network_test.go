/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"testing"

	"github.com/r3labs/vcloud-definition-mapper/definition"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNetworksMapping(t *testing.T) {
	Convey("Given a valid input definition", t, func() {
		d := definition.Definition{
			Name:       "test",
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

		Convey("When I try to map a network", func() {
			Convey("And the input specifies bootstrap as not salt", func() {
				d.Bootstrapping = "none"
				n := MapNetworks(d)
				Convey("Then only input networks should be mapped", func() {
					So(len(n), ShouldEqual, 1)
					So(n[0].Name, ShouldEqual, "datacenter-test-bar")
					So(n[0].RouterName, ShouldEqual, "test")
					So(n[0].Subnet, ShouldEqual, "10.0.0.0/24")
				})
			})

			Convey("And the input specifies bootstrap as salt", func() {
				d.Bootstrapping = "salt"
				n := MapNetworks(d)
				Convey("Then an extra network should be created", func() {
					So(len(n), ShouldEqual, 2)
					So(n[0].Name, ShouldEqual, "datacenter-test-salt")
					So(n[0].RouterName, ShouldEqual, "test")
					So(n[0].Subnet, ShouldEqual, "10.254.254.0/24")
				})
				Convey("Then input network should be mapped as usual", func() {
					So(n[1].Name, ShouldEqual, "datacenter-test-bar")
					So(n[1].RouterName, ShouldEqual, "test")
					So(n[1].Subnet, ShouldEqual, "10.0.0.0/24")
				})
			})
		})

	})
}
