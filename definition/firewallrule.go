/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"errors"
	"unicode/utf8"
)

// FirewallRule ...
type FirewallRule struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Protocol    string `json:"protocol"`
	FromPort    string `json:"from_port"`
	ToPort      string `json:"to_port"`
	Action      string `json:"action"`
}

// Validate firewall rule
func (rule *FirewallRule) Validate(networks []Network) error {
	// Check if firewall rule name is null
	if rule.Name == "" {
		return errors.New("Firewall Rule name should not be null")
	}

	// Check if firewall rule name is > 50 characters
	if utf8.RuneCountInString(rule.Name) > VCLOUDMAXNAME {
		return errors.New("Firewall Rule name can't be greater than 50 characters")
	}

	err := validateIP(rule.Source, "Firewall Rule Source", networks)
	if err != nil {
		return err
	}

	err = validateIP(rule.Destination, "Firewall Rule Destination", networks)
	if err != nil {
		return err
	}

	// Validate FromPort Port
	// Must be: [any | 1 - 65535]
	err = validatePort(rule.FromPort, "Firewall Rule From")
	if err != nil {
		return err
	}

	// Validate ToPort Port
	// Must be: [any | 1 - 65535]
	err = validatePort(rule.ToPort, "Firewall Rule To")
	if err != nil {
		return err
	}

	// Validate Protocol
	// Must be one of: tcp | udp | icmp | any | tcp & udp
	return validateProtocol(rule.Protocol)
}
