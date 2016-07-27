/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import (
	"net"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func testBuildInstances(count int) []Instance {
	var instances []Instance

	for i := 1; i < count+1; i++ {
		instances = append(instances, Instance{
			Name:        "test-" + strconv.Itoa(i),
			Catalog:     "catalog",
			Image:       "image",
			Cpus:        1,
			Memory:      10240,
			NetworkName: "network",
			IP:          net.ParseIP("10.64.0.10" + strconv.Itoa(i)),
		})
	}

	return instances
}

func TestPayloadDiff(t *testing.T) {

	Convey("Given a payload that contains only instances", t, func() {
		var m FSMMessage
		m.Instances.Items = testBuildInstances(2)

		Convey("When diffing a fsm message ", func() {
			Convey("And there is no previous result", func() {
				var p FSMMessage
				m.Diff(&p)
				Convey("Then instances to create should be populated", func() {
					So(len(m.InstancesToCreate.Items), ShouldEqual, 2)
					So(m.InstancesToCreate.Items[0].Name, ShouldEqual, "test-1")
					So(m.InstancesToCreate.Items[0].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToCreate.Items[0].Image, ShouldEqual, "image")
					So(m.InstancesToCreate.Items[0].Cpus, ShouldEqual, 1)
					So(m.InstancesToCreate.Items[0].Memory, ShouldEqual, 10240)
					So(m.InstancesToCreate.Items[0].NetworkName, ShouldEqual, "network")
					So(m.InstancesToCreate.Items[0].IP.String(), ShouldEqual, "10.64.0.101")
					So(m.InstancesToCreate.Items[1].Name, ShouldEqual, "test-2")
					So(m.InstancesToCreate.Items[1].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToCreate.Items[1].Image, ShouldEqual, "image")
					So(m.InstancesToCreate.Items[1].Cpus, ShouldEqual, 1)
					So(m.InstancesToCreate.Items[1].Memory, ShouldEqual, 10240)
					So(m.InstancesToCreate.Items[1].NetworkName, ShouldEqual, "network")
					So(m.InstancesToCreate.Items[1].IP.String(), ShouldEqual, "10.64.0.102")
				})
			})
			Convey("And there is a previous result with one instance", func() {
				var p FSMMessage
				p.Instances.Items = testBuildInstances(1)

				m.Diff(&p)
				Convey("Then instances to create should be populated with the new instance", func() {
					So(len(m.InstancesToCreate.Items), ShouldEqual, 1)
					So(m.InstancesToCreate.Items[0].Name, ShouldEqual, "test-2")
					So(m.InstancesToCreate.Items[0].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToCreate.Items[0].Image, ShouldEqual, "image")
					So(m.InstancesToCreate.Items[0].Cpus, ShouldEqual, 1)
					So(m.InstancesToCreate.Items[0].Memory, ShouldEqual, 10240)
					So(m.InstancesToCreate.Items[0].NetworkName, ShouldEqual, "network")
					So(m.InstancesToCreate.Items[0].IP.String(), ShouldEqual, "10.64.0.102")
				})
			})
			Convey("And there is one less instance on the new input", func() {
				var p FSMMessage
				p.Instances.Items = testBuildInstances(2)
				m.Instances.Items = testBuildInstances(1)

				m.Diff(&p)
				Convey("Then instances to delete should be populated with the old instance", func() {
					So(len(m.InstancesToDelete.Items), ShouldEqual, 1)
					So(m.InstancesToDelete.Items[0].Name, ShouldEqual, "test-2")
					So(m.InstancesToDelete.Items[0].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToDelete.Items[0].Image, ShouldEqual, "image")
					So(m.InstancesToDelete.Items[0].Cpus, ShouldEqual, 1)
					So(m.InstancesToDelete.Items[0].Memory, ShouldEqual, 10240)
					So(m.InstancesToDelete.Items[0].NetworkName, ShouldEqual, "network")
					So(m.InstancesToDelete.Items[0].IP.String(), ShouldEqual, "10.64.0.102")
				})
			})
			Convey("And the resource is updated on the instances", func() {
				var p FSMMessage
				p.Instances.Items = testBuildInstances(2)
				p.Instances.Items[0].Cpus = 2
				p.Instances.Items[1].Cpus = 2

				m.Diff(&p)
				Convey("Then instances to update should be populated with the updated instances", func() {
					So(len(m.InstancesToUpdate.Items), ShouldEqual, 2)
					So(m.InstancesToUpdate.Items[0].Name, ShouldEqual, "test-1")
					So(m.InstancesToUpdate.Items[0].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToUpdate.Items[0].Image, ShouldEqual, "image")
					So(m.InstancesToUpdate.Items[0].Cpus, ShouldEqual, 1)
					So(m.InstancesToUpdate.Items[0].Memory, ShouldEqual, 10240)
					So(m.InstancesToUpdate.Items[0].NetworkName, ShouldEqual, "network")
					So(m.InstancesToUpdate.Items[0].IP.String(), ShouldEqual, "10.64.0.101")
					So(m.InstancesToUpdate.Items[1].Name, ShouldEqual, "test-2")
					So(m.InstancesToUpdate.Items[1].Catalog, ShouldEqual, "catalog")
					So(m.InstancesToUpdate.Items[1].Image, ShouldEqual, "image")
					So(m.InstancesToUpdate.Items[1].Cpus, ShouldEqual, 1)
					So(m.InstancesToUpdate.Items[1].Memory, ShouldEqual, 10240)
					So(m.InstancesToUpdate.Items[1].NetworkName, ShouldEqual, "network")
					So(m.InstancesToUpdate.Items[1].IP.String(), ShouldEqual, "10.64.0.102")
				})
			})
		})
	})
}
