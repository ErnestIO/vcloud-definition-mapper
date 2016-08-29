/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ErnestIO/vcloud-definition-mapper/definition"
)

func TestMapFirewalls(t *testing.T) {

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

		r.Rules = append(r.Rules, definition.FirewallRule{
			Destination: "10.10.10.10",
			ToPort:      "80",
			FromPort:    "80",
			Source:      "10.10.10.11",
			Protocol:    "tcp",
		})

		d.Routers = append(d.Routers, r)

		Convey("When i try to map firewalls", func() {
			Convey("And the input specifies bootstrap as salt", func() {
				d.Bootstrapping = "salt"
				f := MapFirewalls(d)
				Convey("Then it should map salt and input firewall rules", func() {
					So(len(f), ShouldEqual, 1)
					So(f[0].Name, ShouldEqual, "datacenter-service-test")
					So(f[0].RouterName, ShouldEqual, "test")
					So(len(f[0].Rules), ShouldEqual, 6)
					So(f[0].Rules[0].DestinationIP, ShouldEqual, "any")
					So(f[0].Rules[0].DestinationPort, ShouldEqual, "22")
					So(f[0].Rules[0].SourceIP, ShouldEqual, "10.254.254.0/24")
					So(f[0].Rules[0].SourcePort, ShouldEqual, "any")
					So(f[0].Rules[0].Protocol, ShouldEqual, "tcp")
					So(f[0].Rules[1].DestinationIP, ShouldEqual, "any")
					So(f[0].Rules[1].DestinationPort, ShouldEqual, "5985")
					So(f[0].Rules[1].SourceIP, ShouldEqual, "10.254.254.0/24")
					So(f[0].Rules[1].SourcePort, ShouldEqual, "any")
					So(f[0].Rules[1].Protocol, ShouldEqual, "tcp")
					So(f[0].Rules[2].DestinationIP, ShouldEqual, "external")
					So(f[0].Rules[2].DestinationPort, ShouldEqual, "any")
					So(f[0].Rules[2].SourceIP, ShouldEqual, "internal")
					So(f[0].Rules[2].SourcePort, ShouldEqual, "any")
					So(f[0].Rules[2].Protocol, ShouldEqual, "any")
					So(f[0].Rules[3].DestinationIP, ShouldEqual, "10.254.254.100")
					So(f[0].Rules[3].DestinationPort, ShouldEqual, "4505")
					So(f[0].Rules[3].SourceIP, ShouldEqual, "10.0.0.0/24")
					So(f[0].Rules[3].SourcePort, ShouldEqual, "any")
					So(f[0].Rules[3].Protocol, ShouldEqual, "tcp")
					So(f[0].Rules[4].DestinationIP, ShouldEqual, "10.254.254.100")
					So(f[0].Rules[4].DestinationPort, ShouldEqual, "4506")
					So(f[0].Rules[4].SourceIP, ShouldEqual, "10.0.0.0/24")
					So(f[0].Rules[4].SourcePort, ShouldEqual, "any")
					So(f[0].Rules[4].Protocol, ShouldEqual, "tcp")
					So(f[0].Rules[5].DestinationIP, ShouldEqual, "10.10.10.10")
					So(f[0].Rules[5].DestinationPort, ShouldEqual, "80")
					So(f[0].Rules[5].SourceIP, ShouldEqual, "10.10.10.11")
					So(f[0].Rules[5].SourcePort, ShouldEqual, "80")
					So(f[0].Rules[5].Protocol, ShouldEqual, "tcp")
				})
			})

			Convey("And the input specifies bootstrap as not salt", func() {
				d.Bootstrapping = "none"
				f := MapFirewalls(d)
				Convey("Then it should map only the input firewall rules", func() {
					So(len(f), ShouldEqual, 1)
					So(f[0].Name, ShouldEqual, "datacenter-service-test")
					So(f[0].RouterName, ShouldEqual, "test")
					So(len(f[0].Rules), ShouldEqual, 1)
					So(f[0].Rules[0].DestinationIP, ShouldEqual, "10.10.10.10")
					So(f[0].Rules[0].DestinationPort, ShouldEqual, "80")
					So(f[0].Rules[0].SourceIP, ShouldEqual, "10.10.10.11")
					So(f[0].Rules[0].SourcePort, ShouldEqual, "80")
					So(f[0].Rules[0].Protocol, ShouldEqual, "tcp")
				})
			})
		})
	})

}
