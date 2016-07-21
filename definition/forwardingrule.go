/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"errors"
	"net"
)

// ForwardingRule holds port forwarding information
type ForwardingRule struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	FromPort    string `json:"from_port"`
	ToPort      string `json:"to_port"`
}

// Validate checks if PortForwarding is valid
func (rule *ForwardingRule) Validate() error {
	// Check if Destination is a valid IP
	ip := net.ParseIP(rule.Destination)
	if ip == nil {
		return errors.New("Port Forwarding must be a valid IP")
	}

	if rule.Source != "" {
		source := net.ParseIP(rule.Source)
		if source == nil {
			return errors.New("Port Forwarding source must be a valid IP")
		}
	}

	err := validatePort(rule.FromPort, "Port Forwarding From")
	if err != nil {
		return err
	}

	err = validatePort(rule.ToPort, "Port Forwarding To")
	if err != nil {
		return err
	}

	return nil
}
