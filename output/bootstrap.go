/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import "fmt"

// GenerateBootstraps : generates necessary bootstraps for instances
func GenerateBootstraps(m *FSMMessage) []Execution {
	var bootstraps []Execution
	var endpoint string
	var execType string

	if m.ServiceIP != "" {
		endpoint = m.ServiceIP
	} else {
		endpoint = "$(routers.items.0.ip)"
	}

	if m.Datacenters.Items[0].Type == "fake" {
		execType = "fake"
	} else {
		execType = "salt"
	}

	// Add instance to bootstrap if not salt-master
	for _, instance := range m.InstancesToCreate.Items {
		if instance.IP.String() != "10.254.254.100" {
			bootstraps = append(bootstraps, Execution{
				Name:     fmt.Sprintf("Bootstrap %s", instance.Name),
				Type:     execType,
				Target:   "list:salt-master.localdomain",
				Payload:  fmt.Sprintf("/usr/bin/bootstrap -master 10.254.254.100 -host %s -username ernest -password 'b00tStr4pp3rR' -max-retries 20 -minion-name %s", instance.IP, instance.Name),
				EndPoint: endpoint,
				User:     m.SaltUser,
				Password: m.SaltPass,
			})
		}
	}

	return bootstraps
}

// GenerateBootstrapCleanup : When a service updates its bootstraps may need a salt cleanup
func GenerateBootstrapCleanup(m *FSMMessage) []Execution {
	var bootstraps []Execution
	var endpoint string
	var execType string

	if m.ServiceIP != "" {
		endpoint = m.ServiceIP
	} else {
		endpoint = "$(routers.items.0.ip)"
	}

	if m.Datacenters.Items[0].Type == "fake" {
		execType = "fake"
	} else {
		execType = "salt"
	}

	// Add instance to bootstrap if not salt-master
	for _, instance := range m.InstancesToDelete.Items {
		if instance.IP.String() != "10.254.254.100" {
			bootstraps = append(bootstraps, Execution{
				Name:     fmt.Sprintf("Cleanup Bootstrap %s", instance.Name),
				Type:     execType,
				Target:   "list:salt-master.localdomain",
				Payload:  fmt.Sprintf("salt-key -y -d %s", instance.Name),
				EndPoint: endpoint,
				User:     m.SaltUser,
				Password: m.SaltPass,
			})
		}
	}

	return bootstraps

}
