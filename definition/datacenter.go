/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"errors"
	"unicode/utf8"
)

// Datacenter ...
type Datacenter struct {
	Name            string `json:"name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Type            string `json:"type"`
	Region          string `json:"region"`
	ExternalNetwork string `json:"external_network"`
	Token           string `json:"token"`
	Secret          string `json:"secret"`
	VCloudURL       string `json:"vcloud_url"`
	VseURL          string `json:"vse_url"`
}

// Validate checks if a datacenter is valid
func (d *Datacenter) Validate() error {
	// Check if datacenter name is null
	if d.Name == "" {
		return errors.New("Datacenter name should not be null")
	}
	// Check if datacenter name is > 50 characters
	if utf8.RuneCountInString(d.Name) > 50 {
		return errors.New("Datacenter name can't be greater than 50 characters")
	}
	return nil
}
