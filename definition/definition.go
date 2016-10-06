/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"encoding/json"
	"errors"
	"net"
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
	SaltUser      string     `json:"-"`
	SaltPass      string     `json:"-"`
	fake          bool
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

func (d *Definition) validateServiceIP() error {
	if d.ServiceIP == "" {
		return nil
	}
	if ok := net.ParseIP(d.ServiceIP); ok == nil {
		return errors.New("ServiceIP is not a valid IP")
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

	err = d.validateServiceIP()
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
		nw := d.FindNetwork(i.Networks.Name)

		err := i.Validate(nw)
		if err != nil {
			return err
		}
	}

	if hasDuplicateInstance(d.Instances) {
		return errors.New("Duplicate instance names found")
	}

	// Validate Routers
	if hasDuplicateRouters(d.Routers) {
		return errors.New("Duplicate router names found")
	}

	for _, r := range d.Routers {
		err := r.Validate()
		if err != nil {
			return err
		}

		if hasDuplicateNetworks(r.Networks) {
			return errors.New("Duplicate network names found")
		}
	}

	return nil
}

// GeneratedName returns the generated service name
func (d *Definition) GeneratedName() string {
	return d.Datacenter + "-" + d.Name + "-"
}

// FindNetwork returns a network matched by name
func (d *Definition) FindNetwork(name string) *Network {
	for _, router := range d.Routers {
		for _, network := range router.Networks {
			if network.Name == name {
				return &network
			}
		}
	}
	return nil
}

// IsFake returns whether the service build is fake
func (d *Definition) IsFake() bool {
	return d.fake
}
