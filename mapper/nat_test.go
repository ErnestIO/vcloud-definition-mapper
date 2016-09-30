/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ernestio/vcloud-definition-mapper/definition"
)

func TestMapNats(t *testing.T) {

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

		r.PortForwarding = append(r.PortForwarding, definition.ForwardingRule{
			Destination: "10.10.10.10",
			ToPort:      "80",
			FromPort:    "80",
			Source:      "10.10.10.11",
		})

		d.Routers = append(d.Routers, r)

		Convey("When i try to map nats", func() {
			Convey("And the input specifies bootstrap as salt", func() {
				d.Bootstrapping = "salt"
				n := MapNats(d, "test-ext")
				Convey("Then it should map salt and input firewall rules", func() {
					So(len(n), ShouldEqual, 1)
					So(n[0].Name, ShouldEqual, "datacenter-service-test")
					So(len(n[0].Rules), ShouldEqual, 5)
					So(n[0].Rules[0].OriginIP, ShouldEqual, "")
					So(n[0].Rules[0].OriginPort, ShouldEqual, "8000")
					So(n[0].Rules[0].TranslationIP, ShouldEqual, "10.254.254.100")
					So(n[0].Rules[0].TranslationPort, ShouldEqual, "8000")
					So(n[0].Rules[0].Protocol, ShouldEqual, "tcp")
					So(n[0].Rules[0].Network, ShouldEqual, "test-ext")
					So(n[0].Rules[0].Type, ShouldEqual, "dnat")
					So(n[0].Rules[1].OriginIP, ShouldEqual, "")
					So(n[0].Rules[1].OriginPort, ShouldEqual, "22")
					So(n[0].Rules[1].TranslationIP, ShouldEqual, "10.254.254.100")
					So(n[0].Rules[1].TranslationPort, ShouldEqual, "22")
					So(n[0].Rules[1].Protocol, ShouldEqual, "tcp")
					So(n[0].Rules[1].Network, ShouldEqual, "test-ext")
					So(n[0].Rules[1].Type, ShouldEqual, "dnat")
					So(n[0].Rules[2].OriginIP, ShouldEqual, "10.254.254.0/24")
					So(n[0].Rules[2].OriginPort, ShouldEqual, "any")
					So(n[0].Rules[2].TranslationIP, ShouldEqual, "")
					So(n[0].Rules[2].TranslationPort, ShouldEqual, "any")
					So(n[0].Rules[2].Protocol, ShouldEqual, "any")
					So(n[0].Rules[2].Network, ShouldEqual, "test-ext")
					So(n[0].Rules[2].Type, ShouldEqual, "snat")
					So(n[0].Rules[3].OriginIP, ShouldEqual, "10.0.0.0/24")
					So(n[0].Rules[3].OriginPort, ShouldEqual, "any")
					So(n[0].Rules[3].TranslationIP, ShouldEqual, "")
					So(n[0].Rules[3].TranslationPort, ShouldEqual, "any")
					So(n[0].Rules[3].Protocol, ShouldEqual, "any")
					So(n[0].Rules[3].Network, ShouldEqual, "test-ext")
					So(n[0].Rules[3].Type, ShouldEqual, "snat")
					So(n[0].Rules[4].OriginIP, ShouldEqual, "10.10.10.11")
					So(n[0].Rules[4].OriginPort, ShouldEqual, "80")
					So(n[0].Rules[4].TranslationIP, ShouldEqual, "10.10.10.10")
					So(n[0].Rules[4].TranslationPort, ShouldEqual, "80")
					So(n[0].Rules[4].Protocol, ShouldEqual, "tcp")
					So(n[0].Rules[4].Network, ShouldEqual, "test-ext")
					So(n[0].Rules[4].Type, ShouldEqual, "dnat")
				})
			})

			Convey("And the input specifies bootstrap as not salt", func() {
				d.Bootstrapping = "none"
				n := MapNats(d, "test-ext")
				Convey("Then it should map only the input firewall rules", func() {
					So(len(n), ShouldEqual, 1)
					So(n[0].Name, ShouldEqual, "datacenter-service-test")
					So(len(n[0].Rules), ShouldEqual, 2)
					So(n[0].Rules[0].OriginIP, ShouldEqual, "10.0.0.0/24")
					So(n[0].Rules[0].OriginPort, ShouldEqual, "any")
					So(n[0].Rules[0].TranslationIP, ShouldEqual, "")
					So(n[0].Rules[0].TranslationPort, ShouldEqual, "any")
					So(n[0].Rules[0].Protocol, ShouldEqual, "any")
					So(n[0].Rules[0].Network, ShouldEqual, "test-ext")
					So(n[0].Rules[0].Type, ShouldEqual, "snat")
					So(n[0].Rules[1].OriginIP, ShouldEqual, "10.10.10.11")
					So(n[0].Rules[1].OriginPort, ShouldEqual, "80")
					So(n[0].Rules[1].TranslationIP, ShouldEqual, "10.10.10.10")
					So(n[0].Rules[1].TranslationPort, ShouldEqual, "80")
					So(n[0].Rules[1].Protocol, ShouldEqual, "tcp")
					So(n[0].Rules[1].Network, ShouldEqual, "test-ext")
					So(n[0].Rules[1].Type, ShouldEqual, "dnat")

				})
			})
		})
	})

}
