/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import (
	"net"
	"reflect"
)

// Instance ...
type Instance struct {
	ProviderType       string         `json:"_type"`
	Name               string         `json:"name"`
	Catalog            string         `json:"reference_catalog"`
	Image              string         `json:"reference_image"`
	Cpus               int            `json:"cpus"`
	Memory             int            `json:"ram"`
	NetworkName        string         `json:"network_name"`
	IP                 net.IP         `json:"ip"`
	Disks              []InstanceDisk `json:"disks"`
	ShellCommands      []string       `json:"shell_commands"`
	DatacenterType     string         `json:"datacenter_type"`
	DatacenterName     string         `json:"datacenter_name"`
	DatacenterUsername string         `json:"datacenter_username"`
	DatacenterPassword string         `json:"datacenter_password"`
	DatacenterRegion   string         `json:"datacenter_region"`
	VCloudURL          string         `json:"vcloud_url"`
	Exists             bool
	Status             string `json:"status"`
}

// InstanceDisk an instance disk
type InstanceDisk struct {
	ID   int `json:"id"`
	Size int `json:"size"`
}

// HasChanged diff's the two items and returns true if there have been any changes
func (i *Instance) HasChanged(oi *Instance) bool {
	if i.Name == oi.Name &&
		i.Catalog == oi.Catalog &&
		i.Image == oi.Image &&
		i.Cpus == oi.Cpus &&
		i.Memory == oi.Memory &&
		i.NetworkName == oi.NetworkName &&
		string(i.IP) == string(oi.IP) &&
		reflect.DeepEqual(i.Disks, oi.Disks) {
		return false
	}
	return true
}
