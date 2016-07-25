/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// Payload is the JSON payload received from service-store
//
// It has all the info needed to build the message that is going to be sent
// to the FSM over NATS.
type Payload struct {
	ServiceID  string     `json:"id"`
	PrevID     string     `json:"previous_id"`
	Datacenter Datacenter `json:"datacenter"`
	Client     Client     `json:"client"`
	Service    Definition `json:"service"`
}
