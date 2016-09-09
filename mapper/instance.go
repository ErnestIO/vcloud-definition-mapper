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
			Name:        d.GeneratedName() + "salt-master",
			Catalog:     "r3",
			Image:       "r3-salt-master",
			Cpus:        1,
			Memory:      2048,
			Disks:       []output.InstanceDisk{},
			NetworkName: d.GeneratedName() + "salt",
			IP:          net.ParseIP("10.254.254.100"),
		})
	}

	for _, instance := range d.Instances {
		ip := make(net.IP, net.IPv4len)
		copy(ip, instance.Networks.StartIP.To4())
		memory, _ := binaryprefix.GetMB(instance.Memory)

		for i := 0; i < instance.Count; i++ {
			newInstance := output.Instance{
				Name:        d.GeneratedName() + instance.Name + "-" + strconv.Itoa(i+1),
				Catalog:     instance.Catalog(),
				Image:       instance.Template(),
				Cpus:        instance.Cpus,
				Memory:      memory,
				Disks:       MapInstanceDisks(instance.Disks),
				NetworkName: generateNetworkName(&d, instance.Networks.Name),
				IP:          net.ParseIP(ip.String()),
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
