/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRouterHasChanged(t *testing.T) {
	Convey("Given a router", t, func() {
		r := Router{
			Name:         "foo",
			ProviderType: "vcloud",
			IP:           "10.0.0.1",
		}

		Convey("When I compare it to an changed router", func() {
			or := Router{
				Name:         "foo",
				ProviderType: "vcloud",
				IP:           "10.0.0.100",
			}
			change := r.HasChanged(&or)
			Convey("Then it should return true", func() {
				So(change, ShouldBeTrue)
			})
		})

		Convey("When I compare it to an identical router", func() {
			or := Router{
				Name:         "foo",
				ProviderType: "vcloud",
				IP:           "10.0.0.1",
			}
			change := r.HasChanged(&or)
			Convey("Then it should return false", func() {
				So(change, ShouldBeFalse)
			})
		})
	})
}
