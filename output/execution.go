/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import (
	"fmt"
	"strings"
)

// Execution ...
type Execution struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Target  string `json:"target"`
	Payload string `json:"payload"`
	Prefix  string `json:"-"`
	Status  string `json:"status"`
}

// HasChanged diff's the two items and returns true if there have been any changes
func (e *Execution) HasChanged(oe *Execution) bool {
	if e.Name == oe.Name &&
		e.Type == oe.Type &&
		e.Target == oe.Target &&
		e.Payload == oe.Payload &&
		e.Prefix == oe.Prefix {
		return false
	}
	return true
}

// TargetHasChanged returns true if the execution's target has changed
func (e *Execution) TargetHasChanged(oe *Execution) bool {
	// Return false if instances have been removed
	if len(strings.Split(e.Target, ",")) < len(strings.Split(oe.Target, ",")) {
		return false
	}
	return e.Target != oe.Target
}

// PayloadHasChanged returns true if the execution's payload has changed
func (e *Execution) PayloadHasChanged(oe *Execution) bool {
	return e.Payload != oe.Payload
}

// RebuildTarget generates a valid salt target from a list of instances
func (e *Execution) RebuildTarget(instances []Instance) {
	var targets []string
	for _, instance := range instances {
		targets = append(targets, instance.Name)
	}
	nodes := strings.Join(targets, ",")
	e.Target = fmt.Sprintf("list:%s", nodes)
}
