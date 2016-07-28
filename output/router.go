/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import "reflect"

// Router ...
type Router struct {
	Name string `json:"name"`
	Type string `json:"type"`
	IP   string `json:"ip"`
}

// HasChanged diff's the two items and returns true if there have been any changes
func (r *Router) HasChanged(or *Router) bool {
	return !reflect.DeepEqual(*r, *or)
}