/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"errors"
	"fmt"
	"net"
	"unicode/utf8"
)

// Network ...
type Network struct {
	Name   string   `json:"name"`
	Router string   `json:"router"`
	Subnet string   `json:"subnet"`
	DNS    []string `json:"dns"`
}

// Validate checks if a Network is valid
func (n *Network) Validate() error {
	_, _, err := net.ParseCIDR(n.Subnet)
	if err != nil {
		return errors.New("Network CIDR is not valid")
	}

	if n.Name == "" {
		return errors.New("Network name should not be null")
	}

	for _, val := range n.DNS {
		if ok := net.ParseIP(val); ok == nil {
			return errors.New("DNS " + val + " is not a valid CIDR")
		}
	}

	// Check if network name is > 50 characters
	if utf8.RuneCountInString(n.Name) > VCLOUDMAXNAME {
		return fmt.Errorf("Network name can't be greater than %d characters", VCLOUDMAXNAME)
	}

	return nil
}
