/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/r3labs/vcloud-definition-mapper/definition"
	"github.com/r3labs/vcloud-definition-mapper/output"
)

// MapDatacenters : Maps input router to an ernest formatted router
func MapDatacenters(def definition.Definition) (datacenters []output.Datacenter) {
	d := output.Datacenter{}

	d.Name = def.Datacenter.Name
	d.Username = def.Datacenter.Username
	d.Password = def.Datacenter.Password
	d.Region = def.Datacenter.Region
	d.Type = def.Datacenter.Type
	d.Token = def.Datacenter.Token
	d.Secret = def.Datacenter.Secret
	d.ExternalNetwork = def.Datacenter.ExternalNetwork
	d.VCloudURL = def.Datacenter.VCloudURL
	d.VseURL = def.Datacenter.VseURL

	return append(datacenters, d)
}
