/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import (
	"fmt"
	"reflect"
	"strings"
)

// Execution ...
type Execution struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Target  string `json:"target"`
	Payload string `json:"payload"`
	Prefix  string `json:"-"`
}

// HasChanged diff's the two items and returns true if there have been any changes
func (e *Execution) HasChanged(oe *Execution) bool {
	return !reflect.DeepEqual(*e, *oe)
}

// TargetHasChanged returns true if the execution's target has changed
func (e *Execution) TargetHasChanged(oe *Execution) bool {
	return e.Target == oe.Target
}

// PayloadHasChanged returns true if the execution's payload has changed
func (e *Execution) PayloadHasChanged(oe *Execution) bool {
	return e.Payload == oe.Payload
}

// RebuildTarget generates a valid salt target from a list of instances
func (e *Execution) RebuildTarget(instances []Instance) string {
	var targets []string
	for _, instance := range instances {
		targets = append(targets, instance.Name)
	}
	nodes := strings.Join(targets, ",")
	return fmt.Sprintf("list:%s", nodes)
}
