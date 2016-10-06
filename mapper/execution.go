/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ernestio/vcloud-definition-mapper/definition"
	"github.com/ernestio/vcloud-definition-mapper/output"
)

// MapExecutions : Generates the given executions to a valid internal executions
func MapExecutions(d definition.Definition) []output.Execution {
	var execs []output.Execution

	for _, instance := range d.Instances {
		// Instance Name prefix
		prefix := d.GeneratedName() + instance.Name

		// Itterate over each execution
		for i, e := range instance.Provisioner {
			var endpoint string
			var execType string

			if d.ServiceIP != "" {
				endpoint = d.ServiceIP
			} else {
				endpoint = "$(routers.items.0.ip)"
			}

			if d.IsFake() {
				execType = "fake"
			} else {
				execType = "salt"
			}

			// Construct the execution and its payload
			execs = append(execs, output.Execution{
				Name:     fmt.Sprintf("Execution %s %s", instance.Name, strconv.Itoa(i+1)),
				Type:     execType,
				Target:   constructTarget(prefix, instance.Count),
				Payload:  strings.Join(e.Commands, "; "),
				Prefix:   prefix,
				EndPoint: endpoint,
				User:     d.SaltUser,
				Password: d.SaltPass,
			})
		}
	}
	return execs
}

func constructTarget(prefix string, count int) string {
	var targets []string
	for i := 0; i < count; i++ {
		targets = append(targets, prefix+"-"+strconv.Itoa(i+1))
	}
	nodes := strings.Join(targets, ",")
	return fmt.Sprintf("list:%s", nodes)
}
