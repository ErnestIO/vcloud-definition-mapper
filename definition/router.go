/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

type Router struct {
	Name           string           `json:"name"`
	Networks       []Network        `json:"networks"`
	Rules          []FirewallRule   `json:"rules"`
	PortForwarding []ForwardingRule `json:"port_forwarding"`
}

// Validate the router and its components
func (r *Router) Validate() error {
	// Validate Networks
	for _, nw := range r.Networks {
		err := nw.Validate()
		if err != nil {
			return err
		}
	}

	// Validate Firewalls
	for _, ru := range r.Rules {
		err := ru.Validate(r.Networks)
		if err != nil {
			return err
		}
	}

	// Validate Port Forwarding
	for _, pf := range r.PortForwarding {
		err := pf.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
