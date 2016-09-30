/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/vcloud-definition-mapper/definition"
	"github.com/ernestio/vcloud-definition-mapper/output"
)

// MapNats : Generates necessary nats rules for input networks + salt if
// required
func MapNats(d definition.Definition, externalNetwork string) []output.Nat {
	var nats []output.Nat

	for _, r := range d.Routers {
		if r.Networks == nil {
			continue
		}

		// Generate Nats
		n := output.Nat{
			Name:               d.GeneratedName() + r.Name,
			RouterName:         "$(routers.items.0.name)",
			RouterType:         "$(routers.items.0.type)",
			RouterIP:           "$(routers.items.0.ip)",
			ClientName:         "$(client_name)",
			DatacenterType:     "$(datacenters.items.0.type)",
			DatacenterName:     "$(datacenters.items.0.name)",
			DatacenterUsername: "$(datacenters.items.0.username)",
			DatacenterPassword: "$(datacenters.items.0.password)",
			DatacenterRegion:   "$(datacenters.items.0.region)",
			VCloudURL:          "$(datacentes.items.0.vcloud_url)",
		}

		if d.IsSaltBootstrapped() {
			n.Rules = append(n.Rules, MapDefaultSaltNatRules(externalNetwork)...)
		}

		// All Outbound Nat rules for networks
		for _, network := range r.Networks {
			n.Rules = append(n.Rules, output.NatRule{
				Type:            "snat",
				OriginIP:        network.Subnet,
				OriginPort:      "any",
				TranslationIP:   "",
				TranslationPort: "any",
				Protocol:        "any",
				Network:         externalNetwork,
			})
		}

		for _, rule := range r.PortForwarding {
			n.Rules = append(n.Rules, output.NatRule{
				Type:            "dnat",
				OriginIP:        rule.Source,
				OriginPort:      rule.FromPort,
				TranslationIP:   rule.Destination,
				TranslationPort: rule.ToPort,
				Protocol:        "tcp",
				Network:         externalNetwork,
			})
		}

		nats = append(nats, n)
	}
	return nats
}

// MapDefaultSaltNatRules maps all forwarding rules from the external address onto the salt master
func MapDefaultSaltNatRules(externalNetwork string) []output.NatRule {
	// Generatest salt specific nats rules
	var rules []output.NatRule

	rules = append(rules, output.NatRule{
		Type:            "dnat",
		OriginIP:        "",
		OriginPort:      "8000",
		TranslationIP:   "10.254.254.100",
		TranslationPort: "8000",
		Protocol:        "tcp",
		Network:         externalNetwork,
	})

	rules = append(rules, output.NatRule{
		Type:            "dnat",
		OriginIP:        "",
		OriginPort:      "22",
		TranslationIP:   "10.254.254.100",
		TranslationPort: "22",
		Protocol:        "tcp",
		Network:         externalNetwork,
	})

	rules = append(rules, output.NatRule{
		Type:            "snat",
		OriginIP:        "10.254.254.0/24",
		OriginPort:      "any",
		TranslationIP:   "",
		TranslationPort: "any",
		Protocol:        "any",
		Network:         externalNetwork,
	})

	return rules
}
