/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/vcloud-definition-mapper/definition"
	"github.com/ernestio/vcloud-definition-mapper/output"
)

// MapRouters : Maps input router to an ernest formatted router
func MapRouters(d definition.Definition) []output.Router {
	var routers []output.Router

	for _, router := range d.Routers {
		r := output.Router{
			ProviderType:       `$(datacenters.items.0.type)`,
			Name:               router.Name,
			IP:                 d.ServiceIP,
			DatacenterName:     `$(datacenters.items.0.name)`,
			DatacenterPassword: `$(datacenters.items.0.password)`,
			DatacenterRegion:   `$(datacenters.items.0.region)`,
			DatacenterType:     `$(datacenters.items.0.type)`,
			DatacenterUsername: `$(datacenters.items.0.username)`,
			ExternalNetwork:    `$(datacenters.items.0.external_network)`,
			VCloudURL:          `$(datacenters.items.0.vcloud_url)`,
			VseURL:             `$(datacenters.items.0.vse_url)`,
		}
		routers = append(routers, r)
	}

	return routers
}
