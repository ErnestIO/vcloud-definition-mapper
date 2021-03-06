/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

// Network ...
type Network struct {
	Type               string   `json:"_type"`
	Name               string   `json:"name"`
	Service            string   `json:"service"`
	Subnet             string   `json:"range"`
	Netmask            string   `json:"netmask"`
	StartAddress       string   `json:"start_address"`
	EndAddress         string   `json:"end_address"`
	Gateway            string   `json:"gateway"`
	DNS                []string `json:"dns"`
	RouterName         string   `json:"router_name"`
	RouterType         string   `json:"router_type"`
	RouterIP           string   `json:"router_ip"`
	ClientName         string   `json:"client_name"`
	DatacenterType     string   `json:"datacenter_type"`
	DatacenterName     string   `json:"datacenter_name"`
	DatacenterUsername string   `json:"datacenter_username"`
	DatacenterPassword string   `json:"datacenter_password"`
	DatacenterRegion   string   `json:"datacenter_region"`
	VCloudURL          string   `json:"vcloud_url"`
	Exists             bool
	Status             string `json:"status"`
}

// HasChanged diff's the two items and returns true if there have been any changes
func (n *Network) HasChanged(on *Network) bool {
	return false
}
