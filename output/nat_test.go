/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNatHasChanged(t *testing.T) {
	Convey("Given a nat", t, func() {
		var n Nat
		n.Rules = append(n.Rules, NatRule{
			Type:            "snat",
			OriginIP:        "10.0.0.0",
			OriginPort:      "any",
			TranslationIP:   "10.1.1.1",
			TranslationPort: "any",
			Protocol:        "any",
			Network:         "test-nw",
		})

		Convey("When I compare it to an changed nat", func() {
			var on Nat

			on.Rules = append(on.Rules, NatRule{
				Type:            "snat",
				OriginIP:        "10.0.0.0",
				OriginPort:      "any",
				TranslationIP:   "10.1.1.1",
				TranslationPort: "any",
				Protocol:        "tcp",
				Network:         "test-nw",
			})

			change := n.HasChanged(&on)
			Convey("Then it should return true", func() {
				So(change, ShouldBeTrue)
			})
		})

		Convey("When I compare it to an identical nat", func() {
			var on Nat

			on.Rules = append(on.Rules, NatRule{
				Type:            "snat",
				OriginIP:        "10.0.0.0",
				OriginPort:      "any",
				TranslationIP:   "10.1.1.1",
				TranslationPort: "any",
				Protocol:        "any",
				Network:         "test-nw",
			})

			change := n.HasChanged(&on)
			Convey("Then it should return false", func() {
				So(change, ShouldBeFalse)
			})
		})
	})
}
