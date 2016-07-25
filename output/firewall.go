/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import "reflect"

// Firewall ...
type Firewall struct {
	Name       string         `json:"name"`
	RouterName string         `json:"router_name"`
	Rules      []FirewallRule `json:"rules"`
}

// HasChanged diff's the two items and returns true if there have been any changes
func (f *Firewall) HasChanged(of *Firewall) bool {
	if len(f.Rules) != len(of.Rules) {
		return true
	}

	for i := 0; i < len(f.Rules); i++ {
		if reflect.DeepEqual(f.Rules[i], of.Rules[i]) != true {
			return true
		}
	}

	return false
}
