/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"net"
	"strconv"

	"github.com/ernestio/vcloud-definition-mapper/definition"
	"github.com/ernestio/vcloud-definition-mapper/output"
	"github.com/r3labs/binary-prefix"
)

// MapInstances : Maps the instances for the input payload on a ernest internal format
func MapInstances(d definition.Definition) []output.Instance {
	var instances []output.Instance

	if d.IsSaltBootstrapped() {
		instances = append(instances, output.Instance{
			Name:               d.GeneratedName() + "salt-master",
			Hostname:           "salt-master",
			Catalog:            "r3",
			Image:              "r3-salt-master",
			Cpus:               1,
			Memory:             2048,
			Disks:              []output.InstanceDisk{},
			NetworkName:        d.GeneratedName() + "salt",
			IP:                 net.ParseIP("10.254.254.100"),
			ProviderType:       "$(datacenters.items.0.type)",
			DatacenterType:     "$(datacenters.items.0.type)",
			DatacenterName:     "$(datacenters.items.0.name)",
			DatacenterUsername: "$(datacenters.items.0.username)",
			DatacenterPassword: "$(datacenters.items.0.password)",
			DatacenterRegion:   "$(datacenters.items.0.region)",
			VCloudURL:          "$(datacenters.items.0.vcloud_url)",
		})
	}

	for _, instance := range d.Instances {
		ip := make(net.IP, net.IPv4len)
		copy(ip, instance.Networks.StartIP.To4())
		memory, _ := binaryprefix.GetMB(instance.Memory)

		var commands []string

		for _, prov := range instance.Provisioner {
			if len(prov.Shell) > 0 {
				commands = prov.Shell
			}
		}

		for i := 0; i < instance.Count; i++ {
			var disks []output.InstanceDisk

			if instance.RootDisk != "" {
				size, _ := binaryprefix.GetMB(instance.RootDisk)
				disks = append(disks, output.InstanceDisk{
					ID:   0,
					Size: size,
				})
			}

			disks = append(disks, MapInstanceDisks(instance.Disks)...)

			newInstance := output.Instance{
				Name:               d.GeneratedName() + instance.Name + "-" + strconv.Itoa(i+1),
				Hostname:           instance.Name + "-" + strconv.Itoa(i+1),
				Catalog:            instance.Catalog(),
				Image:              instance.Template(),
				Cpus:               instance.Cpus,
				Memory:             memory,
				Disks:              disks,
				NetworkName:        generateNetworkName(&d, instance.Networks.Name),
				IP:                 net.ParseIP(ip.String()),
				ShellCommands:      commands,
				ProviderType:       "$(datacenters.items.0.type)",
				DatacenterType:     "$(datacenters.items.0.type)",
				DatacenterName:     "$(datacenters.items.0.name)",
				DatacenterUsername: "$(datacenters.items.0.username)",
				DatacenterPassword: "$(datacenters.items.0.password)",
				DatacenterRegion:   "$(datacenters.items.0.region)",
				VCloudURL:          "$(datacenters.items.0.vcloud_url)",
			}

			instances = append(instances, newInstance)

			// Increment IP address
			ip[3]++
		}
	}
	return instances
}

// MapInstanceDisks : Maps the instances disks
func MapInstanceDisks(d []string) []output.InstanceDisk {
	var disks []output.InstanceDisk

	for x, disk := range d {
		size, _ := binaryprefix.GetMB(disk)
		disks = append(disks, output.InstanceDisk{
			ID:   (x + 1),
			Size: size,
		})
	}

	return disks
}

func generateNetworkName(d *definition.Definition, name string) string {
	for _, router := range d.Routers {
		for _, network := range router.Networks {
			if network.Name == name {
				return d.GeneratedName() + name
			}
		}
	}

	return name
}
