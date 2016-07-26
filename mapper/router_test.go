/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"testing"

	"github.com/r3labs/vcloud-definition-mapper/definition"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRouterMapping(t *testing.T) {
	Convey("Given a valid definition", t, func() {
		var d definition.Definition
		d.Routers = append(d.Routers, definition.Router{
			Name: "foo",
		})
		Convey("When I map routers", func() {
			r := MapRouters(d, "vcloud")
			Convey("Then it shouldsuccessfully map an output router", func() {
				So(len(r), ShouldEqual, 1)
				So(r[0].Name, ShouldEqual, "foo")
				So(r[0].Type, ShouldEqual, "vcloud")
			})
		})
	})
}
