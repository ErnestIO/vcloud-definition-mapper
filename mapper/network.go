/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ernestio/vcloud-definition-mapper/definition"
	"github.com/ernestio/vcloud-definition-mapper/output"
)

// MapNetworks : Maps the networks from a given input payload. Additionally it will create an
// extra salt-master network if the service bootstrapping mode is defined as salt
func MapNetworks(d definition.Definition) []output.Network {
	var networks []output.Network

	for _, r := range d.Routers {
		if d.IsSaltBootstrapped() {
			subnet := "10.254.254.0/24"
			octets := getIPOctets(subnet)

			sn := output.Network{
				Name:               d.GeneratedName() + "salt",
				Type:               "$(datacenters.items.0.type)",
				Subnet:             subnet,
				StartAddress:       octets + ".5",
				EndAddress:         octets + ".250",
				Gateway:            octets + ".1",
				Netmask:            parseNetmask(subnet),
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
			octets := getIPOctets(network.Subnet)

			n := output.Network{
				Name:               d.GeneratedName() + network.Name,
				Type:               "$(datacenters.items.0.type)",
				Subnet:             network.Subnet,
				StartAddress:       octets + ".5",
				EndAddress:         octets + ".250",
				Gateway:            octets + ".1",
				Netmask:            parseNetmask(network.Subnet),
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

func getIPOctets(rng string) string {
	// Splits the network range and returns the first three octets
	ip, _, err := net.ParseCIDR(rng)
	if err != nil {
		log.Println(err)
	}
	octets := strings.Split(ip.String(), ".")
	octets = append(octets[:3], octets[3+1:]...)
	octetString := strings.Join(octets, ".")
	return octetString
}

func parseNetmask(rng string) string {
	// Convert netmask hex to string, generated from network range CIDR
	_, nw, _ := net.ParseCIDR(rng)
	hx, _ := hex.DecodeString(nw.Mask.String())
	netmask := fmt.Sprintf("%v.%v.%v.%v", hx[0], hx[1], hx[2], hx[3])
	return netmask
}
