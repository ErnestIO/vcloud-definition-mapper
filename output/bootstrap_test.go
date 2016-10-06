/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBootstrapsGeneration(t *testing.T) {
	Convey("Given a definition enabled with bootstrapping", t, func() {
		var m FSMMessage
		m.ServiceIP = "127.0.0.1"
		m.SaltUser = "user"
		m.SaltPass = "pass"
		m.Datacenters.Items = append(m.Datacenters.Items, Datacenter{Type: "vcloud"})
		m.InstancesToCreate.Items = append(m.InstancesToCreate.Items, Instance{Name: "foo"})

		Convey("When I build bootstraps from InstancesToCreate", func() {
			r := GenerateBootstraps(&m)
			Convey("Then it should build an array for bootstraps", func() {
				So(len(r), ShouldEqual, 1)
				So(r[0].Name, ShouldEqual, "Bootstrap foo")
				So(r[0].Type, ShouldEqual, "salt")
				So(r[0].Target, ShouldEqual, "list:salt-master.localdomain")
				So(r[0].Payload, ShouldEqual, "/usr/bin/bootstrap -master 10.254.254.100 -host <nil> -username ernest -password 'b00tStr4pp3rR' -max-retries 20 -minion-name foo")
				So(r[0].EndPoint, ShouldEqual, "127.0.0.1")
				So(r[0].User, ShouldEqual, "user")
				So(r[0].Password, ShouldEqual, "pass")
			})
		})
	})
}
