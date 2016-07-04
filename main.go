/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"os"
	"runtime"

	"github.com/nats-io/nats"
)

var nc *nats.Conn
var natsErr error

func main() {
	natsURI := os.Getenv("NATS_URI")
	if natsURI == "" {
		natsURI = nats.DefaultURL
	}

	nc, natsErr = nats.Connect(natsURI)
	if natsErr != nil {
		log.Fatal(natsErr)
	}

	nc.Subscribe("definition.map.creation.vcloud", notImplemented)
	nc.Subscribe("definition.map.deletion.vcloud", notImplemented)

	runtime.Goexit()
}

func notImplemented(m *nats.Msg) {
	nc.Publish(m.Reply, []byte(`{"error":"not implemented"}`))
}
