/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/vcloud-definition-mapper/definition"
	"github.com/ernestio/vcloud-definition-mapper/output"
)

// MapNetworks : Maps the networks from a given input payload. Additionally it will create an
// extra salt-master network if the service bootstrapping mode is defined as salt
func MapNetworks(d definition.Definition) []output.Network {
	var networks []output.Network

	for _, r := range d.Routers {
		if d.IsSaltBootstrapped() {
			sn := output.Network{
				Name:               d.GeneratedName() + "salt",
				Type:               "$(datacenters.items.0.type)",
				Subnet:             "10.254.254.0/24",
				RouterName:         "$(routers.items.0.name)",
				RouterType:         "$(routers.items.0.type)",
				RouterIP:           "$(routers.items.0.ip)",
				ClientName:         "$(client_name)",
				DatacenterType:     "$(datacenters.items.0.type)",
				DatacenterName:     "$(datacenters.items.0.name)",
				DatacenterUsername: "$(datacenters.items.0.username)",
				DatacenterPassword: "$(datacenters.items.0.password)",
				DatacenterRegion:   "$(datacenters.items.0.region)",
				VCloudURL:          "$(datacenters.items.0.vcloud_url)",
			}

			networks = append(networks, sn)
		}

		for _, network := range r.Networks {

			n := output.Network{
				Name:               d.GeneratedName() + network.Name,
				Subnet:             network.Subnet,
				DNS:                network.DNS,
				RouterName:         "$(routers.items.0.name)",
				RouterType:         "$(routers.items.0.type)",
				RouterIP:           "$(routers.items.0.ip)",
				ClientName:         "$(client_name)",
				DatacenterType:     "$(datacenters.items.0.type)",
				DatacenterName:     "$(datacenters.items.0.name)",
				DatacenterUsername: "$(datacenters.items.0.username)",
				DatacenterPassword: "$(datacenters.items.0.password)",
				DatacenterRegion:   "$(datacenters.items.0.region)",
				VCloudURL:          "$(datacenters.items.0.vcloud_url)",
			}

			networks = append(networks, n)
		}
	}

	return networks
}
