/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/r3labs/vcloud-definition-mapper/definition"
	"github.com/r3labs/vcloud-definition-mapper/output"
)

// ConvertPayload will build an FSMMessage based on an input definition
func ConvertPayload(p definition.Payload) *output.FSMMessage {
	var m output.FSMMessage

	// Map routers
	m.Routers.Items = MapRouters(p.Service, p.Datacenter.Type)

	// Map networks
	m.Networks.Items = MapNetworks(p.Service)

	// Map instances
	m.Instances.Items = MapInstances(p.Service)

	// Map firewalls
	m.Firewalls.Items = MapFirewalls(p.Service)

	// Map nats/port forwarding
	m.Nats.Items = MapNATS(p.Service, p.Datacenter.ExternalNetwork)

	// Bootstraps cannot be generated here as they are dependent on instance to create

	// Map executions
	m.Executions.Items = MapExecutions(d)

	return &m
}
