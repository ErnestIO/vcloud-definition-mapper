/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

// Firewall ...
type Firewall struct {
	Name               string         `json:"name"`
	RouterName         string         `json:"router_name"`
	Rules              []FirewallRule `json:"rules"`
	ClientName         string         `json:"client_name"`
	RouterType         string         `json:"router_type"`
	RouterIP           string         `json:"router_ip"`
	DatacenterName     string         `json:"datacenter_name"`
	DatacenterPassword string         `json:"datacenter_password"`
	DatacenterRegion   string         `json:"datacenter_region"`
	DatacenterType     string         `json:"datacenter_type"`
	DatacenterUsername string         `json:"datacenter_username"`
	ExternalNetwork    string         `json:"external_network"`
	VCloudURL          string         `json:"vcloud_url"`
	Status             string         `json:"status"`
}

// HasChanged diff's the two items and returns true if there have been any changes
func (f *Firewall) HasChanged(of *Firewall) bool {
	if len(f.Rules) != len(of.Rules) {
		return true
	}

	for i := 0; i < len(f.Rules); i++ {
		if f.hasChangedDestinationIP(f.Rules[i].DestinationIP, of.Rules[i].DestinationIP) ||
			f.Rules[i].DestinationPort != of.Rules[i].DestinationPort ||
			f.Rules[i].Protocol != of.Rules[i].Protocol ||
			f.Rules[i].SourceIP != of.Rules[i].SourceIP ||
			f.Rules[i].SourcePort != of.Rules[i].SourcePort {
			return true
		}
	}

	return false
}

func (f *Firewall) hasChangedDestinationIP(n, o string) bool {
	// In case the destination ip is empty it won't be empty on the previous
	// build as it's internally replaced by the endpoint
	if n == "" {
		return false
	}
	if n == "$(routers.items.0.ip)" {
		return false
	}
	if n == o {
		return false
	}
	return true
}
