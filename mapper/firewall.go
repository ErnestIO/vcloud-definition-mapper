/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/vcloud-definition-mapper/definition"
	"github.com/ernestio/vcloud-definition-mapper/output"
)

// MapFirewalls : Maps input firewalls to an ernest format ones
func MapFirewalls(d definition.Definition) []output.Firewall {
	var firewalls []output.Firewall

	for _, r := range d.Routers {
		if r.Rules == nil && d.IsSaltBootstrapped() != true {
			continue
		}

		f := output.Firewall{
			Name:               d.GeneratedName() + r.Name,
			ProviderType:       "$(datacenters.items.0.type)",
			RouterName:         "$(routers.items.0.name)",
			RouterType:         "$(routers.items.0._type)",
			RouterIP:           "$(routers.items.0.ip)",
			DatacenterType:     "$(datacenters.items.0.type)",
			DatacenterName:     "$(datacenters.items.0.name)",
			DatacenterUsername: "$(datacenters.items.0.username)",
			DatacenterPassword: "$(datacenters.items.0.password)",
			DatacenterRegion:   "$(datacenters.items.0.region)",
			VCloudURL:          "$(datacenters.items.0.vcloud_url)",
		}

		if d.IsSaltBootstrapped() {
			f.Rules = append(f.Rules, MapDefaultSaltRules()...)
			f.Rules = append(f.Rules, MapErnestIPSaltRules(d.ErnestIP)...)
			f.Rules = append(f.Rules, MapNetworkSaltRules(r.Networks)...)
		}

		// Validate Firewall Rules
		for _, rule := range r.Rules {
			snw := d.FindNetwork(rule.Source)
			if snw != nil {
				rule.Source = snw.Subnet
			}

			dnw := d.FindNetwork(rule.Destination)
			if dnw != nil {
				rule.Destination = dnw.Subnet
			}

			f.Rules = append(f.Rules, output.FirewallRule{
				Name:            rule.Name,
				SourceIP:        rule.Source,
				SourcePort:      rule.FromPort,
				DestinationIP:   rule.Destination,
				DestinationPort: rule.ToPort,
				Protocol:        rule.Protocol,
			})
		}

		firewalls = append(firewalls, f)
	}
	return firewalls
}

// MapDefaultSaltRules generates the basic rules needed for salt communication
func MapDefaultSaltRules() []output.FirewallRule {
	var rules []output.FirewallRule

	// Allow port 22 & 5985 from salt master to other networks for ssh/winrm
	rules = append(rules, output.FirewallRule{
		SourceIP:        "10.254.254.0/24",
		SourcePort:      "any",
		DestinationIP:   "any",
		DestinationPort: "22",
		Protocol:        "tcp",
	})

	rules = append(rules, output.FirewallRule{
		SourceIP:        "10.254.254.0/24",
		SourcePort:      "any",
		DestinationIP:   "any",
		DestinationPort: "5985",
		Protocol:        "tcp",
	})

	// Allow services/salt range to talk to DNS, minions to external Salt packages
	rules = append(rules, output.FirewallRule{
		SourceIP:        "internal",
		SourcePort:      "any",
		DestinationIP:   "external",
		DestinationPort: "any",
		Protocol:        "any",
	})

	return rules
}

// MapErnestIPSaltRules allows any ernest ip to communicate to the salt master
func MapErnestIPSaltRules(ips []string) []output.FirewallRule {
	var rules []output.FirewallRule

	// Allow port 8000 to current ernest instance
	for _, ip := range ips {
		rules = append(rules, output.FirewallRule{
			SourceIP:        ip,
			SourcePort:      "any",
			DestinationIP:   "$(routers.items.0.ip)",
			DestinationPort: "8000",
			Protocol:        "tcp",
		})
	}

	return rules
}

// MapNetworkSaltRules allows any network range to talk to salt's zeromq ports.
func MapNetworkSaltRules(networks []definition.Network) []output.FirewallRule {
	var rules []output.FirewallRule

	for _, network := range networks {
		rules = append(rules, output.FirewallRule{
			Name:            network.Name + "-salt-firewall-4505-rule",
			SourceIP:        network.Subnet,
			SourcePort:      "any",
			DestinationIP:   "10.254.254.100",
			DestinationPort: "4505",
			Protocol:        "tcp",
		})

		rules = append(rules, output.FirewallRule{
			Name:            network.Name + "-salt-firewall-4506-rule",
			SourceIP:        network.Subnet,
			SourcePort:      "any",
			DestinationIP:   "10.254.254.100",
			DestinationPort: "4506",
			Protocol:        "tcp",
		})
	}

	return rules
}
