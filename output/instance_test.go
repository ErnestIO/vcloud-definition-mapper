/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import (
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInstanceHasChanged(t *testing.T) {
	Convey("Given a instance", t, func() {
		i := Instance{
			Name:    "test",
			Catalog: "catalog",
			Image:   "image",
			Cpus:    1,
			Memory:  10240,
			Disks: []InstanceDisk{
				InstanceDisk{ID: 1, Size: 10240},
			},
			NetworkName: "network",
			IP:          net.ParseIP("10.64.0.100"),
		}

		Convey("When I compare it to an changed instance", func() {
			oi := Instance{
				Name:    "test",
				Catalog: "catalog",
				Image:   "image",
				Cpus:    2,
				Memory:  10240,
				Disks: []InstanceDisk{
					InstanceDisk{ID: 1, Size: 10240},
				},
				NetworkName: "network",
				IP:          net.ParseIP("10.64.0.100"),
			}
			change := i.HasChanged(&oi)
			Convey("Then it should return true", func() {
				So(change, ShouldBeTrue)
			})
		})

		Convey("When I compare it to an identical instance", func() {
			oi := Instance{
				Name:    "test",
				Catalog: "catalog",
				Image:   "image",
				Cpus:    1,
				Memory:  10240,
				Disks: []InstanceDisk{
					InstanceDisk{ID: 1, Size: 10240},
				},
				NetworkName: "network",
				IP:          net.ParseIP("10.64.0.100"),
			}
			change := i.HasChanged(&oi)
			Convey("Then it should return false", func() {
				So(change, ShouldBeFalse)
			})
		})
	})
}
