/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import "reflect"

// Nat ...
type Nat struct {
	Name       string    `json:"name"`
	Service    string    `json:"service"`
	RouterName string    `json:"router_name"`
	Rules      []NatRule `json:"rules"`
}

// HasChanged diff's the two items and returns true if there have been any changes
func (n *Nat) HasChanged(on *Nat) bool {
	if len(n.Rules) != len(on.Rules) {
		return true
	}

	for i := 0; i < len(n.Rules); i++ {
		if reflect.DeepEqual(n.Rules[i], on.Rules[i]) != true {
			return true
		}
	}

	return false
}
