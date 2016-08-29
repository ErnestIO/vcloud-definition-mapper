/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ErnestIO/vcloud-definition-mapper/definition"
	"github.com/ErnestIO/vcloud-definition-mapper/output"
)

// MapRouters : Maps input router to an ernest formatted router
func MapRouters(d definition.Definition, rtype string) []output.Router {
	var routers []output.Router

	for _, router := range d.Routers {
		r := output.Router{
			Name: router.Name,
			Type: rtype,
			IP:   d.ServiceIP,
		}
		routers = append(routers, r)
	}

	return routers
}
