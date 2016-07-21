/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

const (
	PROTOCOLTCP    = "tcp"
	PROTOCOLUDP    = "udp"
	PROTOCOLANY    = "any"
	PROTOCOLICMP   = "icmp"
	PROTOCOLTCPUDP = "tcp & udp"
	TARGETEXTERNAL = "external"
	TARGETINTERNAL = "internal"
	TARGETANY      = "any"
	VCLOUDMAXNAME  = 50
)

func isNetwork(networks []Network, ip string) bool {
	for _, network := range networks {
		if network.Name == ip {
			return true
		}
	}
	return false
}

func validateProtocol(p string) error {
	switch p {
	case PROTOCOLTCP, PROTOCOLUDP, PROTOCOLICMP, PROTOCOLANY, PROTOCOLTCPUDP:
		return nil
	}
	return errors.New("Protocol is invalid")
}

// ValidateIP checks if an string is a valid source/destionation
func validateIP(ip, iptype string, networks []Network) error {
	// Check if Source is a valid value or a valid IP/CIDR
	// One of: external | internal | any | named networks | CIDR

	switch ip {
	case TARGETEXTERNAL, TARGETINTERNAL, TARGETANY:
		return nil
	}

	// Check if it refers to an internal Network
	if isNetwork(networks, ip) {
		return nil
	}

	// Check if Source is a valid CIDR
	_, _, err := net.ParseCIDR(ip)
	if err == nil {
		return nil
	}

	// Check if Source is a valid IP
	ipx := net.ParseIP(ip)
	if ipx != nil {
		return nil
	}

	return fmt.Errorf("%s (%s) is not valid", iptype, ip)
}

// ValidatePort checks an string to be a valid TCP port
func validatePort(p, ptype string) error {
	if p == "any" {
		return nil
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		return fmt.Errorf("%s Port (%s) is not valid", ptype, p)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("%s Port (%s) is out of range [1 - 65535]", ptype, p)
	}

	return nil
}
