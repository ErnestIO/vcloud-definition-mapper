/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

// Nat ...
type Nat struct {
	Name       string    `json:"name"`
	Service    string    `json:"service"`
	RouterName string    `json:"router_name"`
	Rules      []NatRule `json:"rules"`
	Status     string    `json:"status"`
}

// HasChanged diff's the two items and returns true if there have been any changes
func (n *Nat) HasChanged(on *Nat) bool {
	if len(n.Rules) != len(on.Rules) {
		return true
	}

	for i := 0; i < len(n.Rules); i++ {
		if n.hasChangedIP(n.Rules[i].OriginIP, on.Rules[i].OriginIP) ||
			n.Rules[i].OriginPort != on.Rules[i].OriginPort ||
			n.hasChangedIP(n.Rules[i].TranslationIP, on.Rules[i].TranslationIP) ||
			n.Rules[i].TranslationPort != on.Rules[i].TranslationPort ||
			n.Rules[i].Protocol != on.Rules[i].Protocol ||
			n.Rules[i].Type != on.Rules[i].Type {
			return true
		}
	}

	return false
}

func (n *Nat) hasChangedIP(ip, oip string) bool {
	// In case the destination ip is empty it won't be empty on the previous
	// build as it's internally replaced by the endpoint
	if ip == "" {
		return false
	}
	if ip == oip {
		return false
	}
	return true
}
