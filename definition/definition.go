/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"encoding/json"
	"errors"
	"unicode/utf8"
)

// Definition ...
type Definition struct {
	Name          string     `json:"name"`
	Datacenter    string     `json:"datacenter"`
	Bootstrapping string     `json:"bootstrapping"`
	ErnestIP      []string   `json:"ernest_ip"`
	ServiceIP     string     `json:"service_ip"`
	Routers       []Router   `json:"routers"`
	Instances     []Instance `json:"instances"`
}

// New returns a new Definition
func New() *Definition {
	return &Definition{
		ErnestIP:  make([]string, 0),
		Routers:   make([]Router, 0),
		Instances: make([]Instance, 0),
	}
}

// FromJSON creates a definition from json
func FromJSON(data []byte) (*Definition, error) {
	var d Definition

	err := json.Unmarshal(data, d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ValidateName checks if service is valid
func (d *Definition) validateName() error {
	// Check if service name is null
	if d.Name == "" {
		return errors.New("Service name should not be null")
	}
	// Check if service name is > 50 characters
	if utf8.RuneCountInString(d.Name) > 50 {
		return errors.New("Service name can't be greater than 50 characters")
	}
	return nil
}

func (d *Definition) validateDatacenter() error {
	if d.Datacenter == "" {
		return errors.New("Datacenter not specified")
	}
	return nil
}

// IsSaltBootstrapped : Return a boolean describing if bootstrapping option is salt
func (d *Definition) IsSaltBootstrapped() bool {
	if d.Bootstrapping == "salt" {
		return true
	}
	return false
}

// Validate the definition
func (d *Definition) Validate() error {
	// Validate Definition
	err := d.validateName()
	if err != nil {
		return err
	}

	// Validate Datacenter
	err = d.validateDatacenter()
	if err != nil {
		return err
	}

	// Validate Instances
	for _, i := range d.Instances {
		err := i.Validate()
		if err != nil {
			return err
		}
	}

	// Validate Routers
	for _, r := range d.Routers {
		err := r.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
