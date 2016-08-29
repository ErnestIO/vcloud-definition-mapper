/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ErnestIO/vcloud-definition-mapper/definition"
	"github.com/ErnestIO/vcloud-definition-mapper/output"
)

// MapExecutions : Generates the given executions to a valid internal executions
func MapExecutions(d definition.Definition) []output.Execution {
	var execs []output.Execution

	for _, instance := range d.Instances {
		// Instance Name prefix
		prefix := d.GeneratedName() + instance.Name

		// Itterate over each execution
		for i, e := range instance.Provisioner {
			// Construct the execution and its payload
			execs = append(execs, output.Execution{
				Name:    fmt.Sprintf("Execution %s %s", instance.Name, strconv.Itoa(i+1)),
				Type:    "salt",
				Target:  constructTarget(prefix, instance.Count),
				Payload: strings.Join(e.Commands, "; "),
				Prefix:  prefix,
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
