/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/r3labs/vcloud-definition-mapper/definition"
	"github.com/r3labs/vcloud-definition-mapper/output"
)

// MapNetworks : Maps the networks from a given input payload. Additionally it will create an
// extra salt-master network if the service bootstrapping mode is defined as salt
func MapNetworks(d definition.Definition) []output.Network {
	var networks []output.Network

	for _, r := range d.Routers {
		if d.IsSaltBootstrapped() {
			sn := output.Network{
				Name:       d.GeneratedName() + "salt",
				RouterName: r.Name,
				Subnet:     "10.254.254.0/24",
			}

			networks = append(networks, sn)
		}

		for _, network := range r.Networks {

			n := output.Network{
				Name:       d.GeneratedName() + network.Name,
				RouterName: r.Name,
				Subnet:     network.Subnet,
				DNS:        network.DNS,
			}

			networks = append(networks, n)
		}
	}

	return networks
}
