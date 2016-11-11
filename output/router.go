/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

// Router ...
type Router struct {
	ProviderType       string `json:"_type"`
	Name               string `json:"name"`
	IP                 string `json:"ip"`
	DatacenterName     string `json:"datacenter_name"`
	DatacenterPassword string `json:"datacenter_password"`
	DatacenterRegion   string `json:"datacenter_region"`
	DatacenterType     string `json:"datacenter_type"`
	DatacenterUsername string `json:"datacenter_username"`
	ExternalNetwork    string `json:"external_network"`
	VCloudURL          string `json:"vcloud_url"`
	VseURL             string `json:"vse_url"`
	Status             string `json:"status"`
}

// HasChanged diff's the two items and returns true if there have been any changes
func (r *Router) HasChanged(or *Router) bool {
	if r.Name == or.Name &&
		r.IP == or.IP {
		return false
	}
	return true
}
