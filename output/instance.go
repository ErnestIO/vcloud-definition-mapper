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
	Name        string         `json:"name"`
	Catalog     string         `json:"reference_catalog"`
	Image       string         `json:"reference_image"`
	Cpus        int            `json:"cpus"`
	Memory      int            `json:"ram"`
	NetworkName string         `json:"network_name"`
	IP          net.IP         `json:"ip"`
	Disks       []InstanceDisk `json:"disks"`
	Exists      bool
	Status      string `json:"status"`
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
