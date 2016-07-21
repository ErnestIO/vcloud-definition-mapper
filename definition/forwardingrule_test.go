/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPortForwardingValidation(t *testing.T) {

	Convey("Given a forwarding rule", t, func() {
		r := ForwardingRule{Destination: "127.0.0.1", Source: "127.0.0.1", FromPort: "any", ToPort: "any"}

		Convey("When I try to validate a rule with an invalid destination ip", func() {
			r.Destination = "test"
			err := r.Validate()
			Convey("Then should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding must be a valid IP")
			})
		})

		Convey("When I try to validate a rule with an invalid source ip", func() {
			r.Source = "test"
			err := r.Validate()
			Convey("Then should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding source must be a valid IP")
			})
		})

		Convey("When I try to validate a rule with a valid destination", func() {
			r.Destination = "127.0.0.1"
			err := r.Validate()
			Convey("Then should return an error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I try to validate a rule with a valid source", func() {
			r.Source = "127.0.0.1"
			err := r.Validate()
			Convey("Then should return an error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When From Port is any", func() {
			r.FromPort = "any"
			err := r.Validate()
			Convey("Then should return an error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When From Port is not any and not numeric", func() {
			r.FromPort = "test"
			err := r.Validate()
			Convey("Then should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding From Port (test) is not valid")
			})
		})

		Convey("When From Port is not any and not in range", func() {
			r.FromPort = "0"
			err := r.Validate()
			Convey("Then should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding From Port (0) is out of range [1 - 65535]")
			})
		})

		Convey("When From Port is not any and great than range", func() {
			r.FromPort = "65536"
			err := r.Validate()
			Convey("Then should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding From Port (65536) is out of range [1 - 65535]")
			})
		})
	})
}
