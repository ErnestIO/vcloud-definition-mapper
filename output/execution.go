/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import "reflect"

// Execution ...
type Execution struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Target  string `json:"target"`
	Payload string `json:"payload"`
}

// HasChanged diff's the two items and returns true if there have been any changes
func (e *Execution) HasChanged(oe *Execution) bool {
	return !reflect.DeepEqual(*e, *oe)
}
