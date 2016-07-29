/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/r3labs/vcloud-definition-mapper/definition"
	"github.com/r3labs/vcloud-definition-mapper/output"
)

// MapDatacenters : Maps input router to an ernest formatted router
func MapDatacenters(dat definition.Datacenter) (datacenters []output.Datacenter) {
	d := output.Datacenter{}

	d.Name = dat.Name
	d.Username = dat.Username
	d.Password = dat.Password
	d.Region = dat.Region
	d.Type = dat.Type
	d.Token = dat.Token
	d.Secret = dat.Secret
	d.ExternalNetwork = dat.ExternalNetwork
	d.VCloudURL = dat.VCloudURL
	d.VseURL = dat.VseURL

	return append(datacenters, d)
}
