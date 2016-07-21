/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInstanceValidate(t *testing.T) {

	Convey("Given an instance", t, func() {
		n := InstanceNetworks{Name: "test", StartIP: net.ParseIP("127.0.0.1")}
		i := Instance{Name: "test", Image: "test/test", Cpus: 2, Memory: "2GB", Count: 1, Networks: n}
		Convey("With an invalid name", func() {
			i.Name = ""
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Instance name should not be null")
				})
			})
		})

		Convey("With a name > 50 chars", func() {
			i.Name = "aksjhdlkashdliuhliusncldiudnalsundlaiunsdliausndliuansdlksbdlas"
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Instance name can't be greater than 50 characters")
				})
			})
		})

		Convey("With an invalid image name", func() {
			i.Image = ""
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Instance image should not be null")
				})
			})
		})

		Convey("With an invalid image format", func() {
			i.Image = "aksjhdlkashdliuhliusncldiud"
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Instance image invalid, use format <catalog>/<image>")
				})
			})
		})

		Convey("With an empty image catalog", func() {
			i.Image = "/image"
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Instance image catalog should not be null, use format <catalog>/<image>")
				})
			})
		})

		Convey("With an empty image", func() {
			i.Image = "catalog/"
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Instance image image should not be null, use format <catalog>/<image>")
				})
			})
		})

		Convey("With a cpu field less than one", func() {
			i.Cpus = 0
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Instance cpus should not be < 1")
				})
			})
		})

		Convey("With an empty memory field", func() {
			i.Memory = ""
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Instance memory should not be null")
				})
			})
		})

		Convey("With an instance count less than one", func() {
			i.Count = 0
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Instance count should not be < 1")
				})
			})
		})

		Convey("With an invalid networks name", func() {
			i.Networks.Name = ""
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Instance network name should not be null")
				})
			})
		})

		Convey("With an invalid start ip", func() {
			i.Networks.StartIP = nil
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Instance network start_ip should not be null")
				})
			})
		})

		Convey("With valid entries", func() {
			Convey("When validating the instance", func() {
				err := i.Validate()
				Convey("Then should return an error", func() {
					So(err, ShouldBeNil)
				})
			})
		})

	})
}
