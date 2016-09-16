/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
	"time"

	ecc "github.com/ernestio/ernest-config-client"
	"github.com/ernestio/vcloud-definition-mapper/definition"
	"github.com/ernestio/vcloud-definition-mapper/mapper"
	"github.com/ernestio/vcloud-definition-mapper/output"
	"github.com/nats-io/nats"
)

var nc *nats.Conn
var natsErr error

func main() {
	nc = ecc.NewConfig(os.Getenv("NATS_URI")).Nats()

	nc.Subscribe("definition.map.creation.vcloud", createDefinitionHandler)
	nc.Subscribe("definition.map.deletion.vcloud", deleteDefinitionHandler)

	runtime.Goexit()
}

func createDefinitionHandler(msg *nats.Msg) {
	var om output.FSMMessage

	p, err := definition.PayloadFromJSON(msg.Data)
	if err != nil {
		nc.Publish(msg.Reply, []byte(`{"error":"Failed to parse payload."}`))
		return
	}

	err = p.Service.Validate()
	if err != nil {
		nc.Publish(msg.Reply, []byte(`{"error":"`+err.Error()+`"}`))
		return
	}

	// new fsm message
	m := mapper.ConvertPayload(p)

	// previous output message if it exists
	if p.PrevID != "" {
		om, err = getPreviousService(p.PrevID)
		if err != nil {
			nc.Publish(msg.Reply, []byte(`{"error":"Failed to get previous output."}`))
			return
		}
	}

	// Check for changes and create workflow arcs
	m.Diff(om)
	err = m.GenerateWorkflow("create-workflow.json")
	if err != nil {
		log.Println(err.Error())
		nc.Publish(msg.Reply, []byte(`{"error":"Could not generate workflow."}`))
		return
	}

	data, err := json.Marshal(m)
	if err != nil {
		nc.Publish(msg.Reply, []byte(`{"error":"Failed marshal output message."}`))
		return
	}

	nc.Publish(msg.Reply, data)
}

func deleteDefinitionHandler(msg *nats.Msg) {
	p, err := definition.PayloadFromJSON(msg.Data)
	if err != nil {
		nc.Publish(msg.Reply, []byte(`{"error":"Failed to parse payload."}`))
		return
	}

	m, err := getPreviousService(p.PrevID)
	if err != nil {
		nc.Publish(msg.Reply, []byte(`{"error":"Failed to get previous output."}`))
		return
	}

	// Assign all items to delete
	m.RoutersToDelete = m.Routers
	for i := range m.RoutersToDelete.Items {
		m.RoutersToDelete.Items[i].Status = ""
	}
	m.NetworksToDelete = m.Networks
	for i := range m.NetworksToDelete.Items {
		m.NetworksToDelete.Items[i].Status = ""
	}
	m.InstancesToDelete = m.Instances
	for i := range m.InstancesToDelete.Items {
		m.InstancesToDelete.Items[i].Status = ""
	}
	m.FirewallsToDelete = m.Firewalls
	for i := range m.FirewallsToDelete.Items {
		m.FirewallsToDelete.Items[i].Status = ""
	}
	m.NatsToDelete = m.Nats
	for i := range m.NatsToDelete.Items {
		m.NatsToDelete.Items[i].Status = ""
	}

	// Generate delete workflow
	m.GenerateWorkflow("delete-workflow.json")

	data, err := json.Marshal(m)
	if err != nil {
		nc.Publish(msg.Reply, []byte(`{"error":"Failed marshal output message."}`))
		return
	}

	nc.Publish(msg.Reply, data)
}

func getPreviousService(id string) (output.FSMMessage, error) {
	var payload output.FSMMessage

	msg, err := nc.Request("service.get.mapping", []byte(`{"id":"`+id+`"}`), time.Second)
	if err != nil {
		log.Println(err.Error())
		return payload, err
	}

	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		return payload, err
	}

	return payload, nil
}
