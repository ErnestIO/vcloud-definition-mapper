/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"unicode/utf8"

	"github.com/r3labs/binary-prefix"
)

// Instance ...
type Instance struct {
	Count       int              `json:"count"`
	Cpus        int              `json:"cpus"`
	Image       string           `json:"image"`
	Memory      string           `json:"memory"`
	Disks       []string         `json:"disks"`
	Name        string           `json:"name"`
	Networks    InstanceNetworks `json:"networks"`
	Provisioner []Exec           `json:"provisioner"`
}

// Validate : Validates the instance returning true or false if is valid or not
func (i *Instance) Validate(network *Network) error {
	if i.Name == "" {
		return errors.New("Instance name should not be null")
	}

	if utf8.RuneCountInString(i.Name) > VCLOUDMAXNAME {
		return fmt.Errorf("Instance name can't be greater than %d characters", VCLOUDMAXNAME)
	}

	if i.Image == "" {
		return errors.New("Instance image should not be null")
	}

	imageParts := strings.Split(i.Image, "/")
	if len(imageParts) < 2 {
		return errors.New("Instance image invalid, use format <catalog>/<image>")
	}

	if imageParts[0] == "" {
		return errors.New("Instance image catalog should not be null, use format <catalog>/<image>")
	}

	if imageParts[1] == "" {
		return errors.New("Instance image image should not be null, use format <catalog>/<image>")
	}

	if i.Cpus < 1 {
		return errors.New("Instance cpus should not be < 1")
	}

	if i.Memory == "" {
		return errors.New("Instance memory should not be null")
	}

	if i.Count < 1 {
		return errors.New("Instance count should not be < 1")
	}

	if i.Networks.Name == "" {
		return errors.New("Instance network name should not be null")
	}

	if i.Networks.StartIP == nil {
		return errors.New("Instance network start_ip should not be null")
	}

	_, err := binaryprefix.GetMB(i.Memory)
	if err != nil {
		return errors.New("Invalid memory format")
	}

	// Validate IP addresses
	if network != nil {
		_, nw, err := net.ParseCIDR(network.Subnet)
		if err != nil {
			return errors.New("Could not process network")
		}

		startIP := net.ParseIP(i.Networks.StartIP.String()).To4()
		ip := make(net.IP, net.IPv4len)
		copy(ip, i.Networks.StartIP.To4())

		for x := 0; x < i.Count; x++ {
			if !nw.Contains(ip) {
				err := errors.New("Instance IP invalid. IP must be a valid IP in the same range as it's network")
				return err
			}

			// Check IP is greater than Start IP (Bounds checking)
			if ip[3] < startIP[3] {
				err := errors.New("Instance IP invalid. Allocated IP is lower than Start IP")
				return err
			}

			ip[3]++
		}
	}

	return nil
}

// Catalog ...
func (i *Instance) Catalog() string {
	parts := strings.Split(i.Image, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

// Template ...
func (i *Instance) Template() string {
	parts := strings.Split(i.Image, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
