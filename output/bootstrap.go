/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import (
	"fmt"
)

// GenerateBootstraps : generates necessary bootstraps for instances
func GenerateBootstraps(instances []Instance) (bootstraps []Execution) {
	// Add instance to bootstrap if not salt-master
	for _, instance := range instances {
		if instance.IP.String() != "10.254.254.100" {
			bootstraps = append(bootstraps, Execution{
				Name:    fmt.Sprintf("Bootstrap %s", instance.Name),
				Type:    "salt",
				Target:  "list:salt-master.localdomain",
				Payload: fmt.Sprintf("/usr/bin/bootstrap -master 10.254.254.100 -host %s -username ernest -password 'b00tStr4pp3rR' -max-retries 20 -minion-name %s", instance.IP, instance.Name),
			})
		}
	}

	return bootstraps
}

// GenerateBootstrapCleanup : When a service updates its bootstraps may need a salt cleanup
func GenerateBootstrapCleanup(instances []Instance) (bootstraps []Execution) {
	// Add instance to bootstrap if not salt-master
	for _, instance := range instances {
		// if instance.IP.String() != "10.254.254.100" && !instance.Exists {

		if instance.IP.String() != "10.254.254.100" {
			bootstraps = append(bootstraps, Execution{
				Name:    fmt.Sprintf("Cleanup Bootstrap %s", instance.Name),
				Type:    "salt",
				Target:  "list:salt-master.localdomain",
				Payload: fmt.Sprintf("salt-key -y -d %s", instance.Name),
			})
		}
	}

	return bootstraps

}
