/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/vcloud-definition-mapper/definition"
	"github.com/ernestio/vcloud-definition-mapper/output"
)

// MapDatacenters : Maps input router to an ernest formatted router
func MapDatacenters(dat definition.Datacenter) []output.Datacenter {
	var datacenters []output.Datacenter

	datacenters = append(datacenters, output.Datacenter{
		Name:            dat.Name,
		Username:        dat.Username,
		Password:        dat.Password,
		Region:          dat.Region,
		Type:            dat.Type,
		Token:           dat.Token,
		Secret:          dat.Secret,
		ExternalNetwork: dat.ExternalNetwork,
		VCloudURL:       dat.VCloudURL,
		VseURL:          dat.VseURL,
	})

	return datacenters
}
