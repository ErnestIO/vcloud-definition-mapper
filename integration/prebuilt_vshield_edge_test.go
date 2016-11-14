/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package integration

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"testing"

	"github.com/nats-io/nats"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPreVSE(t *testing.T) {
	var service = "novse"

	service2 := service + "II" + strconv.Itoa(rand.Intn(1000))
	service = service + strconv.Itoa(rand.Intn(1000))

	nwCreateSub := make(chan *nats.Msg, 2)
	inCreateSub := make(chan *nats.Msg, 2)
	fwCreateSub := make(chan *nats.Msg, 1)
	ntCreateSub := make(chan *nats.Msg, 2)
	exCreateSub := make(chan *nats.Msg, 3)
	boCreateSub := make(chan *nats.Msg, 3)
	inUpdateSub := make(chan *nats.Msg, 2)
	fwUpdateSub := make(chan *nats.Msg, 1)
	ntUpdateSub := make(chan *nats.Msg, 1)
	inDeleteSub := make(chan *nats.Msg, 1)

	basicSetup("vcloud")

	Convey("Given I have a configured ernest instance", t, func() {
		Convey("When I apply a valid novse1.yml definition", func() {
			nsub, _ := n.ChanSubscribe("network.create.vcloud-fake", nwCreateSub)
			isub, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			fsub, _ := n.ChanSubscribe("firewall.create.vcloud-fake", fwCreateSub)
			asub, _ := n.ChanSubscribe("nat.create.vcloud-fake", ntCreateSub)

			f := getDefinitionPath("novse1.yml", service)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				n := networkEvent{}
				msg, err := waitMsg(nwCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &n)
				i := instanceEvent{}
				msg, err = waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)
				f := firewallEvent{}
				msg, err = waitMsg(fwCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &f)
				na := natEvent{}
				msg, err = waitMsg(ntCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &na)

				Info("And I should receive a valid network.create.vcloud-fake", " ", 8)
				So(n.DatacenterName, ShouldEqual, "fake")
				So(n.DatacenterPassword, ShouldEqual, default_pwd)
				So(n.DatacenterType, ShouldEqual, "vcloud-fake")
				So(n.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(n.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(n.NetworkGateway, ShouldEqual, "10.1.0.1")
				So(n.NetworkNetmask, ShouldEqual, "255.255.255.0")
				So(n.NetworkStartAddress, ShouldEqual, "10.1.0.5")
				So(n.NetworkEndAddress, ShouldEqual, "10.1.0.250")

				Info("And I should receive a valid instance.create.vcloud-fake", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-web-1")
				So(i.Cpus, ShouldEqual, 1)
				So(len(i.Disks), ShouldEqual, 0)
				So(i.IP, ShouldEqual, "10.1.0.11")
				So(i.Memory, ShouldEqual, 1024)
				So(i.ReferenceCatalog, ShouldEqual, "r3")
				So(i.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")

				Info("And I should receive a valid firewall.create.vcloud-fake", " ", 8)
				So(f.DatacenterName, ShouldEqual, "fake")
				So(f.DatacenterPassword, ShouldEqual, default_pwd)
				So(f.DatacenterType, ShouldEqual, "vcloud-fake")
				So(f.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(f.Type, ShouldEqual, "vcloud-fake")
				So(f.RouterIP, ShouldEqual, "172.16.186.44")
				So(f.RouterName, ShouldEqual, "vse2")
				So(f.RouterType, ShouldEqual, "$(datacenters.items.0.type)")
				So(len(f.Rules), ShouldEqual, 4)
				Printf("\n        And it will allow internal:any to internal:any ")
				So(f.Rules[0].SourcePort, ShouldEqual, "any")
				So(f.Rules[0].SourceIP, ShouldEqual, "internal")
				So(f.Rules[0].DestinationIP, ShouldEqual, "internal")
				So(f.Rules[0].DestinationPort, ShouldEqual, "any")
				So(f.Rules[0].Protocol, ShouldEqual, "any")
				Printf("\n        And it will allow 172.18.143.3:any to internal:22 ")
				So(f.Rules[1].SourcePort, ShouldEqual, "any")
				So(f.Rules[1].SourceIP, ShouldEqual, "172.18.143.3")
				So(f.Rules[1].DestinationIP, ShouldEqual, "internal")
				So(f.Rules[1].DestinationPort, ShouldEqual, "22")
				So(f.Rules[1].Protocol, ShouldEqual, "tcp")
				Printf("\n        And it will allow 172.17.240.0/24:any to internal:22 ")
				So(f.Rules[2].SourcePort, ShouldEqual, "any")
				So(f.Rules[2].SourceIP, ShouldEqual, "172.17.240.0/24")
				So(f.Rules[2].DestinationIP, ShouldEqual, "internal")
				So(f.Rules[2].DestinationPort, ShouldEqual, "22")
				So(f.Rules[2].Protocol, ShouldEqual, "tcp")
				Printf("\n        And it will allow 172.19.186.30/24:any to internal:22 ")
				So(f.Rules[3].SourcePort, ShouldEqual, "any")
				So(f.Rules[3].SourceIP, ShouldEqual, "172.19.186.30")
				So(f.Rules[3].DestinationIP, ShouldEqual, "internal")
				So(f.Rules[3].DestinationPort, ShouldEqual, "22")
				So(f.Rules[3].Protocol, ShouldEqual, "tcp")

				Info("And I should receive a valid nat.create.vcloud-fake", " ", 8)
				So(na.DatacenterName, ShouldEqual, "fake")
				So(na.DatacenterPassword, ShouldEqual, default_pwd)
				So(na.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(na.DatacenterType, ShouldEqual, "vcloud-fake")
				So(na.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(na.NatName, ShouldEqual, "fake-"+service+"-vse2")
				So(len(na.NatRules), ShouldEqual, 2)
				So(na.NatRules[0].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[0].OriginIP, ShouldEqual, "10.1.0.0/24")
				So(na.NatRules[0].OriginPort, ShouldEqual, "any")
				So(na.NatRules[0].Type, ShouldEqual, "snat")
				So(na.NatRules[0].TranslationIP, ShouldEqual, "172.16.186.44")
				So(na.NatRules[0].TranslationPort, ShouldEqual, "any")
				So(na.NatRules[0].Protocol, ShouldEqual, "any")
				So(na.RouterIP, ShouldEqual, "172.16.186.44")
				So(na.RouterName, ShouldEqual, "vse2")
				So(na.RouterType, ShouldEqual, "$(datacenters.items.0.type)")
			})

			nsub.Unsubscribe()
			isub.Unsubscribe()
			fsub.Unsubscribe()
			asub.Unsubscribe()
		})

		Convey("When I apply a valid novse2.yml definition", func() {
			fsub, _ := n.ChanSubscribe("firewall.update.vcloud-fake", fwUpdateSub)

			f := getDefinitionPath("novse2.yml", service)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				Info("Then I should receive a valid firewall.update.vcloud-fake", " ", 8)
				event := firewallEvent{}
				msg, err := waitMsg(fwUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)
				So(event.DatacenterName, ShouldEqual, "fake")
				So(event.DatacenterPassword, ShouldEqual, default_pwd)
				So(event.DatacenterType, ShouldEqual, "vcloud-fake")
				So(event.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(event.Type, ShouldEqual, "vcloud-fake")
				So(event.RouterIP, ShouldEqual, "172.16.186.44")
				So(event.RouterName, ShouldEqual, "vse2")
				So(event.RouterType, ShouldEqual, "$(datacenters.items.0.type)")
				So(len(event.Rules), ShouldEqual, 5)
				Printf("\n        And it will allow internal:any to external:any ")
				So(event.Rules[4].SourcePort, ShouldEqual, "any")
				So(event.Rules[4].SourceIP, ShouldEqual, "172.19.186.30")
				So(event.Rules[4].DestinationIP, ShouldEqual, "internal")
				So(event.Rules[4].DestinationPort, ShouldEqual, "22")
				So(event.Rules[4].Protocol, ShouldEqual, "tcp")
			})

			fsub.Unsubscribe()
		})

		Convey("When I apply a valid novse3.yml definition", func() {
			asub, _ := n.ChanSubscribe("nat.update.vcloud-fake", ntUpdateSub)

			f := getDefinitionPath("novse3.yml", service)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				Info("Then I should receive a valid nats.update.vcloud-fake", " ", 8)
				event := natEvent{}
				msg, err := waitMsg(ntUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)
				So(event.DatacenterName, ShouldEqual, "fake")
				So(event.DatacenterPassword, ShouldEqual, default_pwd)
				So(event.DatacenterType, ShouldEqual, "vcloud-fake")
				So(event.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(event.RouterIP, ShouldEqual, "172.16.186.44")
				So(event.RouterName, ShouldEqual, "vse2")
				So(event.RouterType, ShouldEqual, "$(datacenters.items.0.type)")
				So(event.NatName, ShouldEqual, "fake-"+service+"-vse2")
				So(len(event.NatRules), ShouldEqual, 3)
				Printf("\n        And it will forward port 22 to 10.1.0.12 ")
				So(event.NatRules[2].Network, ShouldEqual, "NETWORK")
				So(event.NatRules[2].TranslationIP, ShouldEqual, "10.1.0.12")
				So(event.NatRules[2].TranslationPort, ShouldEqual, "22")
				So(event.NatRules[2].OriginIP, ShouldEqual, "172.16.186.61")
				So(event.NatRules[2].OriginPort, ShouldEqual, "22")
				So(event.NatRules[2].Type, ShouldEqual, "dnat")
				So(event.NatRules[2].Protocol, ShouldEqual, "tcp")
			})

			asub.Unsubscribe()
		})

		Convey("When I apply a valid novse4.yml definition", func() {
			icsub, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			iusub, _ := n.ChanSubscribe("instance.update.vcloud-fake", inUpdateSub)

			f := getDefinitionPath("novse4.yml", service)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				i := instanceEvent{}
				msg, err := waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)
				iu := instanceEvent{}
				msg, err = waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &iu)

				Info("And I should receive a valid instance.create.vcloud-fake", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				So(i.Cpus, ShouldEqual, 1)
				So(len(i.Disks), ShouldEqual, 0)
				So(i.IP, ShouldEqual, "10.1.0.12")
				So(i.Memory, ShouldEqual, 1024)
				So(i.ReferenceCatalog, ShouldEqual, "r3")
				So(i.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")

				Info("And I should receive a valid instance.update.vcloud-fake", " ", 8)
				So(iu.DatacenterName, ShouldEqual, "fake")
				So(iu.DatacenterPassword, ShouldEqual, default_pwd)
				So(iu.DatacenterType, ShouldEqual, "vcloud-fake")
				So(iu.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(iu.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				So(iu.Cpus, ShouldEqual, 1)
				So(len(iu.Disks), ShouldEqual, 0)
				So(iu.IP, ShouldEqual, "10.1.0.12")
				So(iu.Memory, ShouldEqual, 1024)
				So(iu.ReferenceCatalog, ShouldEqual, "r3")
				So(iu.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(iu.InstanceType, ShouldEqual, "vcloud-fake")
				So(iu.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(iu.RouterIP, ShouldEqual, "")
				So(iu.RouterName, ShouldEqual, "")
				So(iu.RouterType, ShouldEqual, "")
			})

			icsub.Unsubscribe()
			iusub.Unsubscribe()
		})

		Convey("When I apply a valid novse5.yml definition", func() {
			iusub, _ := n.ChanSubscribe("instance.update.vcloud-fake", inUpdateSub)

			f := getDefinitionPath("novse5.yml", service)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				i := instanceEvent{}
				msg, err := waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)
				iu := instanceEvent{}
				msg, err = waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &iu)

				Info("And I should receive a valid instance.update.vcloud-fake", " ", 8)
				Info("And it will update cpu count on instance 1", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-web-1")
				So(i.Cpus, ShouldEqual, 2)
				So(len(i.Disks), ShouldEqual, 0)
				So(i.IP, ShouldEqual, "10.1.0.11")
				So(i.Memory, ShouldEqual, 1024)
				So(i.ReferenceCatalog, ShouldEqual, "r3")
				So(i.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")

				Info("And it will update cpu count on instance 2", " ", 8)
				So(iu.DatacenterName, ShouldEqual, "fake")
				So(iu.DatacenterPassword, ShouldEqual, default_pwd)
				So(iu.DatacenterType, ShouldEqual, "vcloud-fake")
				So(iu.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(iu.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				So(iu.Cpus, ShouldEqual, 2)
				So(len(iu.Disks), ShouldEqual, 0)
				So(iu.IP, ShouldEqual, "10.1.0.12")
				So(iu.Memory, ShouldEqual, 1024)
				So(iu.ReferenceCatalog, ShouldEqual, "r3")
				So(iu.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(iu.InstanceType, ShouldEqual, "vcloud-fake")
				So(iu.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(iu.RouterIP, ShouldEqual, "")
				So(iu.RouterName, ShouldEqual, "")
				So(iu.RouterType, ShouldEqual, "")
			})

			iusub.Unsubscribe()
		})

		Convey("When I apply a valid novse6.yml definition", func() {
			iusub, _ := n.ChanSubscribe("instance.update.vcloud-fake", inUpdateSub)

			f := getDefinitionPath("novse6.yml", service)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				i := instanceEvent{}
				msg, err := waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)
				iu := instanceEvent{}
				msg, err = waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &iu)

				Info("And I should receive a valid instance.update.vcloud-fake", " ", 8)
				Info("And it will update disks on instance 1 ", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-web-1")
				So(i.Cpus, ShouldEqual, 2)
				So(len(i.Disks), ShouldEqual, 1)
				So(i.Disks[0].ID, ShouldEqual, 1)
				So(i.Disks[0].Size, ShouldEqual, 10240)
				So(i.IP, ShouldEqual, "10.1.0.11")
				So(i.Memory, ShouldEqual, 1024)
				So(i.ReferenceCatalog, ShouldEqual, "r3")
				So(i.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")

				Info("And it will update disks on instance 2", " ", 8)
				So(iu.DatacenterName, ShouldEqual, "fake")
				So(iu.DatacenterPassword, ShouldEqual, default_pwd)
				So(iu.DatacenterType, ShouldEqual, "vcloud-fake")
				So(iu.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(iu.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				So(iu.Cpus, ShouldEqual, 2)
				So(len(iu.Disks), ShouldEqual, 1)
				So(iu.Disks[0].ID, ShouldEqual, 1)
				So(iu.Disks[0].Size, ShouldEqual, 10240)
				So(iu.IP, ShouldEqual, "10.1.0.12")
				So(iu.Memory, ShouldEqual, 1024)
				So(iu.ReferenceCatalog, ShouldEqual, "r3")
				So(iu.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(iu.InstanceType, ShouldEqual, "vcloud-fake")
				So(iu.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(iu.RouterIP, ShouldEqual, "")
				So(iu.RouterName, ShouldEqual, "")
				So(iu.RouterType, ShouldEqual, "")
			})

			iusub.Unsubscribe()
		})

		Convey("When I apply a valid novse7.yml definition", func() {
			iusub, _ := n.ChanSubscribe("instance.update.vcloud-fake", inUpdateSub)

			f := getDefinitionPath("novse7.yml", service)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				i := instanceEvent{}
				msg, err := waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)
				iu := instanceEvent{}
				msg, err = waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &iu)

				Info("And I should receive a valid instance.update.vcloud-fake", " ", 8)
				Info("And it will update ram on instance 1 ", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-web-1")
				So(i.Cpus, ShouldEqual, 2)
				So(len(i.Disks), ShouldEqual, 1)
				So(i.Disks[0].ID, ShouldEqual, 1)
				So(i.Disks[0].Size, ShouldEqual, 10240)
				So(i.IP, ShouldEqual, "10.1.0.11")
				So(i.Memory, ShouldEqual, 2048)
				So(i.ReferenceCatalog, ShouldEqual, "r3")
				So(i.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")

				Info("And it will update ram on instance 2 ", " ", 8)
				So(iu.DatacenterName, ShouldEqual, "fake")
				So(iu.DatacenterPassword, ShouldEqual, default_pwd)
				So(iu.DatacenterType, ShouldEqual, "vcloud-fake")
				So(iu.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(iu.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				So(iu.Cpus, ShouldEqual, 2)
				So(len(iu.Disks), ShouldEqual, 1)
				So(iu.Disks[0].ID, ShouldEqual, 1)
				So(iu.Disks[0].Size, ShouldEqual, 10240)
				So(iu.IP, ShouldEqual, "10.1.0.12")
				So(iu.Memory, ShouldEqual, 2048)
				So(iu.ReferenceCatalog, ShouldEqual, "r3")
				So(iu.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(iu.InstanceType, ShouldEqual, "vcloud-fake")
				So(iu.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(iu.RouterIP, ShouldEqual, "")
				So(iu.RouterName, ShouldEqual, "")
				So(iu.RouterType, ShouldEqual, "")
			})

			iusub.Unsubscribe()
		})

		Convey("When I apply a valid novse8.yml definition", func() {
			nsub, _ := n.ChanSubscribe("network.create.vcloud-fake", nwCreateSub)
			asub, _ := n.ChanSubscribe("nat.update.vcloud-fake", ntUpdateSub)

			f := getDefinitionPath("novse8.yml", service)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				n := networkEvent{}
				msg, err := waitMsg(nwCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &n)
				na := natEvent{}
				msg, err = waitMsg(ntUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &na)

				Info("And I should receive a valid network.create.vcloud-fake", " ", 8)
				So(n.DatacenterName, ShouldEqual, "fake")
				So(n.DatacenterPassword, ShouldEqual, default_pwd)
				So(n.DatacenterType, ShouldEqual, "vcloud-fake")
				So(n.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(n.NetworkName, ShouldEqual, "fake-"+service+"-db")
				So(n.NetworkGateway, ShouldEqual, "10.2.0.1")
				So(n.NetworkNetmask, ShouldEqual, "255.255.255.0")
				So(n.NetworkStartAddress, ShouldEqual, "10.2.0.5")
				So(n.NetworkEndAddress, ShouldEqual, "10.2.0.250")
				So(n.RouterIP, ShouldEqual, "172.16.186.44")
				So(n.RouterName, ShouldEqual, "vse2")
				So(n.RouterType, ShouldEqual, "$(datacenters.items.0.type)")

				Info("And I should receive a valid nat.update.vcloud-fake", " ", 8)
				So(na.DatacenterName, ShouldEqual, "fake")
				So(na.DatacenterPassword, ShouldEqual, default_pwd)
				So(na.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(na.DatacenterType, ShouldEqual, "vcloud-fake")
				So(na.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(na.NatName, ShouldEqual, "fake-"+service+"-vse2")
				So(len(na.NatRules), ShouldEqual, 4)
				Printf("\n        And it will create a snat for the new network ")
				So(na.NatRules[1].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[1].OriginIP, ShouldEqual, "10.2.0.0/24")
				So(na.NatRules[1].OriginPort, ShouldEqual, "any")
				So(na.NatRules[1].Type, ShouldEqual, "snat")
				So(na.NatRules[1].TranslationIP, ShouldEqual, "172.16.186.44")
				So(na.NatRules[1].TranslationPort, ShouldEqual, "any")
				So(na.NatRules[1].Protocol, ShouldEqual, "any")
				So(na.RouterIP, ShouldEqual, "172.16.186.44")
				So(na.RouterName, ShouldEqual, "vse2")
				So(na.RouterType, ShouldEqual, "$(datacenters.items.0.type)")
			})

			nsub.Unsubscribe()
			asub.Unsubscribe()
		})

		Convey("When I apply a valid novse9.yml definition", func() {
			icsub, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			iusub, _ := n.ChanSubscribe("instance.update.vcloud-fake", inUpdateSub)

			f := getDefinitionPath("novse9.yml", service)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				i := instanceEvent{}
				msg, err := waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)
				iu := instanceEvent{}
				msg, err = waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &iu)

				Info("And I should receive a valid instance.create.vcloud-fake", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-db-1")
				So(i.Cpus, ShouldEqual, 1)
				So(len(i.Disks), ShouldEqual, 0)
				So(i.IP, ShouldEqual, "10.2.0.11")
				So(i.Memory, ShouldEqual, 1024)
				So(i.ReferenceCatalog, ShouldEqual, "r3")
				So(i.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service+"-db")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")

				Info("And I should receive a valid instance.update.vcloud-fake", " ", 8)
				So(iu.DatacenterName, ShouldEqual, "fake")
				So(iu.DatacenterPassword, ShouldEqual, default_pwd)
				So(iu.DatacenterType, ShouldEqual, "vcloud-fake")
				So(iu.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(iu.InstanceName, ShouldEqual, "fake-"+service+"-db-1")
				So(iu.Cpus, ShouldEqual, 1)
				So(len(iu.Disks), ShouldEqual, 0)
				So(iu.IP, ShouldEqual, "10.2.0.11")
				So(iu.Memory, ShouldEqual, 1024)
				So(iu.ReferenceCatalog, ShouldEqual, "r3")
				So(iu.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(iu.InstanceType, ShouldEqual, "vcloud-fake")
				So(iu.NetworkName, ShouldEqual, "fake-"+service+"-db")
				So(iu.RouterIP, ShouldEqual, "")
				So(iu.RouterName, ShouldEqual, "")
				So(iu.RouterType, ShouldEqual, "")
			})

			icsub.Unsubscribe()
			iusub.Unsubscribe()
		})

		Convey("When I apply a valid novse10.yml definition", func() {
			isub, _ := n.ChanSubscribe("instance.delete.vcloud-fake", inDeleteSub)

			f := getDefinitionPath("novse10.yml", service)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				event := instanceEvent{}
				msg, err := waitMsg(inDeleteSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)

				Info("And I should receive a valid instance.delete.vcloud-fake", " ", 8)
				So(event.DatacenterName, ShouldEqual, "fake")
				So(event.DatacenterPassword, ShouldEqual, default_pwd)
				So(event.DatacenterType, ShouldEqual, "vcloud-fake")
				So(event.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(event.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				So(event.Cpus, ShouldEqual, 2)
				So(len(event.Disks), ShouldEqual, 1)
				So(event.Disks[0].ID, ShouldEqual, 1)
				So(event.Disks[0].Size, ShouldEqual, 10240)
				So(event.IP, ShouldEqual, "10.1.0.12")
				So(event.Memory, ShouldEqual, 2048)
				So(event.ReferenceCatalog, ShouldEqual, "r3")
				So(event.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(event.InstanceType, ShouldEqual, "vcloud-fake")
				So(event.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(event.RouterIP, ShouldEqual, "")
				So(event.RouterName, ShouldEqual, "")
				So(event.RouterType, ShouldEqual, "")
			})

			isub.Unsubscribe()
		})

		Convey("When I apply a valid novse11.yml definition", func() {
			isub, _ := n.ChanSubscribe("instance.delete.vcloud-fake", inDeleteSub)

			f := getDefinitionPath("novse11.yml", service)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				event := instanceEvent{}
				msg, err := waitMsg(inDeleteSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)

				Info("And I should receive a valid instance.delete.vcloud-fake", " ", 8)
				So(event.DatacenterName, ShouldEqual, "fake")
				So(event.DatacenterPassword, ShouldEqual, default_pwd)
				So(event.DatacenterType, ShouldEqual, "vcloud-fake")
				So(event.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(event.InstanceName, ShouldEqual, "fake-"+service+"-db-1")
				So(event.Cpus, ShouldEqual, 1)
				So(len(event.Disks), ShouldEqual, 0)
				So(event.IP, ShouldEqual, "10.2.0.11")
				So(event.Memory, ShouldEqual, 1024)
				So(event.ReferenceCatalog, ShouldEqual, "r3")
				So(event.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(event.InstanceType, ShouldEqual, "vcloud-fake")
				So(event.NetworkName, ShouldEqual, "fake-"+service+"-db")
				So(event.RouterIP, ShouldEqual, "")
				So(event.RouterName, ShouldEqual, "")
				So(event.RouterType, ShouldEqual, "")
			})

			isub.Unsubscribe()
		})

		Convey("When I apply a valid novse12.yml definition", func() {
			nsub, _ := n.ChanSubscribe("network.create.vcloud-fake", nwCreateSub)
			isub, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			fsub, _ := n.ChanSubscribe("firewall.create.vcloud-fake", fwCreateSub)
			asub, _ := n.ChanSubscribe("nat.create.vcloud-fake", ntCreateSub)
			bsub, _ := n.ChanSubscribe("bootstrap.create.fake", boCreateSub)
			esub, _ := n.ChanSubscribe("execution.create.fake", exCreateSub)

			f := getDefinitionPath("novse12.yml", service2)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				n1 := networkEvent{}
				msg, err := waitMsg(nwCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &n1)
				n2 := networkEvent{}
				msg, err = waitMsg(nwCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &n2)
				i := instanceEvent{}
				msg, err = waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)
				i2 := instanceEvent{}
				msg, err = waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i2)
				f := firewallEvent{}
				msg, err = waitMsg(fwCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &f)
				na := natEvent{}
				msg, err = waitMsg(ntCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &na)
				ex := executionEvent{}
				msg, err = waitMsg(boCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &ex)
				ex2 := executionEvent{}
				msg, err = waitMsg(exCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &ex2)

				Info("And I should receive a valid network.create.vcloud-fake", " ", 8)
				Info("And it should create the salt master network", " ", 8)
				So(n1.DatacenterName, ShouldEqual, "fake")
				So(n1.DatacenterPassword, ShouldEqual, default_pwd)
				So(n1.DatacenterType, ShouldEqual, "vcloud-fake")
				So(n1.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(n1.NetworkName, ShouldEqual, "fake-"+service2+"-salt")
				So(n1.NetworkGateway, ShouldEqual, "10.254.254.1")
				So(n1.NetworkNetmask, ShouldEqual, "255.255.255.0")
				So(n1.NetworkStartAddress, ShouldEqual, "10.254.254.5")
				So(n1.NetworkEndAddress, ShouldEqual, "10.254.254.250")

				Info("And it should create the user defined network", " ", 8)
				So(n2.DatacenterName, ShouldEqual, "fake")
				So(n2.DatacenterPassword, ShouldEqual, default_pwd)
				So(n2.DatacenterType, ShouldEqual, "vcloud-fake")
				So(n2.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(n2.NetworkName, ShouldEqual, "fake-"+service2+"-web")
				So(n2.NetworkGateway, ShouldEqual, "10.1.0.1")
				So(n2.NetworkNetmask, ShouldEqual, "255.255.255.0")
				So(n2.NetworkStartAddress, ShouldEqual, "10.1.0.5")
				So(n2.NetworkEndAddress, ShouldEqual, "10.1.0.250")

				Info("And I should receive a valid instance.create.vcloud-fake", " ", 8)
				Info("And it should create the salt master instance", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service2+"-salt-master")
				So(i.Cpus, ShouldEqual, 1)
				So(len(i.Disks), ShouldEqual, 0)
				So(i.IP, ShouldEqual, "10.254.254.100")
				So(i.Memory, ShouldEqual, 2048)
				So(i.ReferenceCatalog, ShouldEqual, "r3")
				So(i.ReferenceImage, ShouldEqual, "r3-salt-master")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service2+"-salt")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")

				Info("And it should create the user defined instance ", " ", 8)
				So(i2.DatacenterName, ShouldEqual, "fake")
				So(i2.DatacenterPassword, ShouldEqual, default_pwd)
				So(i2.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i2.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i2.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i2.InstanceName, ShouldEqual, "fake-"+service2+"-web-1")
				So(i2.Cpus, ShouldEqual, 1)
				So(len(i2.Disks), ShouldEqual, 0)
				So(i2.IP, ShouldEqual, "10.1.0.11")
				So(i2.Memory, ShouldEqual, 1024)
				So(i2.ReferenceCatalog, ShouldEqual, "r3")
				So(i2.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(i2.InstanceType, ShouldEqual, "vcloud-fake")
				So(i2.NetworkName, ShouldEqual, "fake-"+service2+"-web")
				So(i2.RouterIP, ShouldEqual, "")
				So(i2.RouterName, ShouldEqual, "")
				So(i2.RouterType, ShouldEqual, "")

				Info("Then I should receive a valid firewall.create.vcloud-fake", " ", 8)
				So(f.DatacenterName, ShouldEqual, "fake")
				So(f.DatacenterPassword, ShouldEqual, default_pwd)
				So(f.DatacenterType, ShouldEqual, "vcloud-fake")
				So(f.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(f.Type, ShouldEqual, "vcloud-fake")
				So(f.RouterIP, ShouldEqual, "172.16.186.44")
				So(f.RouterName, ShouldEqual, "vse2")
				So(f.RouterType, ShouldEqual, "$(datacenters.items.0.type)")
				So(len(f.Rules), ShouldEqual, 8)

				Info("And it will allow 10.254.254.0/24:any to any:22 ", " ", 8)
				So(f.Rules[0].SourcePort, ShouldEqual, "any")
				So(f.Rules[0].SourceIP, ShouldEqual, "10.254.254.0/24")
				So(f.Rules[0].DestinationIP, ShouldEqual, "any")
				So(f.Rules[0].DestinationPort, ShouldEqual, "22")
				So(f.Rules[0].Protocol, ShouldEqual, "tcp")

				Info("And it will allow 10.254.254.0/24:any to any:5985 ", " ", 8)
				So(f.Rules[1].SourcePort, ShouldEqual, "any")
				So(f.Rules[1].SourceIP, ShouldEqual, "10.254.254.0/24")
				So(f.Rules[1].DestinationIP, ShouldEqual, "any")
				So(f.Rules[1].DestinationPort, ShouldEqual, "5985")
				So(f.Rules[1].Protocol, ShouldEqual, "tcp")

				Info("And it will allow internal:any to external:any ", " ", 8)
				So(f.Rules[2].SourcePort, ShouldEqual, "any")
				So(f.Rules[2].SourceIP, ShouldEqual, "internal")
				So(f.Rules[2].DestinationIP, ShouldEqual, "external")
				So(f.Rules[2].DestinationPort, ShouldEqual, "any")
				So(f.Rules[2].Protocol, ShouldEqual, "any")

				Info("And it will allow 172.17.241.95/24:any to 172.16.186.44:8000 ", " ", 8)
				So(f.Rules[3].SourcePort, ShouldEqual, "any")
				So(f.Rules[3].SourceIP, ShouldEqual, "172.17.241.95")
				So(f.Rules[3].DestinationIP, ShouldEqual, "172.16.186.44")
				So(f.Rules[3].DestinationPort, ShouldEqual, "8000")
				So(f.Rules[3].Protocol, ShouldEqual, "tcp")

				Info("And it will allow 10.1.0.0/24:any to 10.254.254.100:any ", " ", 8)
				So(f.Rules[4].SourcePort, ShouldEqual, "any")
				So(f.Rules[4].SourceIP, ShouldEqual, "10.1.0.0/24")
				So(f.Rules[4].DestinationIP, ShouldEqual, "10.254.254.100")
				So(f.Rules[4].DestinationPort, ShouldEqual, "4505")
				So(f.Rules[4].Protocol, ShouldEqual, "tcp")

				Info("And it will allow 10.1.0.0/24:any to 10.254.254.100:4506 ", " ", 8)
				So(f.Rules[5].SourcePort, ShouldEqual, "any")
				So(f.Rules[5].SourceIP, ShouldEqual, "10.1.0.0/24")
				So(f.Rules[5].DestinationIP, ShouldEqual, "10.254.254.100")
				So(f.Rules[5].DestinationPort, ShouldEqual, "4506")
				So(f.Rules[5].Protocol, ShouldEqual, "tcp")

				Info("And it will allow internal:any to internal:any ", " ", 8)
				So(f.Rules[6].SourcePort, ShouldEqual, "any")
				So(f.Rules[6].SourceIP, ShouldEqual, "internal")
				So(f.Rules[6].DestinationIP, ShouldEqual, "internal")
				So(f.Rules[6].DestinationPort, ShouldEqual, "any")
				So(f.Rules[6].Protocol, ShouldEqual, "any")

				Info("And it will allow internal:any to external:any ", " ", 8)
				So(f.Rules[7].SourcePort, ShouldEqual, "any")
				So(f.Rules[7].SourceIP, ShouldEqual, "internal")
				So(f.Rules[7].DestinationIP, ShouldEqual, "external")
				So(f.Rules[7].DestinationPort, ShouldEqual, "any")
				So(f.Rules[7].Protocol, ShouldEqual, "any")

				Info("And I should receive a valid nat.create.vcloud-fake", " ", 8)
				So(na.DatacenterName, ShouldEqual, "fake")
				So(na.DatacenterPassword, ShouldEqual, default_pwd)
				So(na.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(na.DatacenterType, ShouldEqual, "vcloud-fake")
				So(na.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(na.NatName, ShouldEqual, "fake-"+service2+"-vse2")
				So(len(na.NatRules), ShouldEqual, 4)

				Info("And it will forward 172.16.186.44:8000 to 10.254.254.100:8000 ", " ", 8)
				So(na.NatRules[0].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[0].OriginIP, ShouldEqual, "172.16.186.44")
				So(na.NatRules[0].OriginPort, ShouldEqual, "8000")
				So(na.NatRules[0].Type, ShouldEqual, "dnat")
				So(na.NatRules[0].TranslationIP, ShouldEqual, "10.254.254.100")
				So(na.NatRules[0].TranslationPort, ShouldEqual, "8000")
				So(na.NatRules[0].Protocol, ShouldEqual, "tcp")

				Info("And it will forward 172.16.186.44:22 to 10.254.254.100:22 ", " ", 8)
				So(na.NatRules[1].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[1].OriginIP, ShouldEqual, "172.16.186.44")
				So(na.NatRules[1].OriginPort, ShouldEqual, "22")
				So(na.NatRules[1].Type, ShouldEqual, "dnat")
				So(na.NatRules[1].TranslationIP, ShouldEqual, "10.254.254.100")
				So(na.NatRules[1].TranslationPort, ShouldEqual, "22")
				So(na.NatRules[1].Protocol, ShouldEqual, "tcp")

				Info("And it will pat 10.254.254.0/24:any to 172.16.186.44 ", " ", 8)
				So(na.NatRules[2].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[2].OriginIP, ShouldEqual, "10.254.254.0/24")
				So(na.NatRules[2].OriginPort, ShouldEqual, "any")
				So(na.NatRules[2].Type, ShouldEqual, "snat")
				So(na.NatRules[2].TranslationIP, ShouldEqual, "172.16.186.44")
				So(na.NatRules[2].TranslationPort, ShouldEqual, "any")
				So(na.NatRules[2].Protocol, ShouldEqual, "any")

				Info("And it will pat 10.1.0.0/24:any to 172.16.186.44 ", " ", 8)
				So(na.NatRules[3].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[3].OriginIP, ShouldEqual, "10.1.0.0/24")
				So(na.NatRules[3].OriginPort, ShouldEqual, "any")
				So(na.NatRules[3].Type, ShouldEqual, "snat")
				So(na.NatRules[3].TranslationIP, ShouldEqual, "172.16.186.44")
				So(na.NatRules[3].TranslationPort, ShouldEqual, "any")
				So(na.NatRules[3].Protocol, ShouldEqual, "any")
				So(na.RouterIP, ShouldEqual, "172.16.186.44")
				So(na.RouterName, ShouldEqual, "vse2")
				So(na.RouterType, ShouldEqual, "$(datacenters.items.0.type)")

				Info("And I should receive a valid execution.create.fake", " ", 8)
				Info("And it will bootstrap the web node ", " ", 8)
				So(ex.Name, ShouldEqual, "Bootstrap fake-"+service2+"-web-1")
				So(ex.ExecutionType, ShouldEqual, "fake")
				So(ex.ExecutionPayload, ShouldContainSubstring, "-host 10.1.0.11")
				So(ex.ExecutionTarget, ShouldEqual, "list:salt-master.localdomain")
				So(ex.ServiceOptions.User, ShouldEqual, salt.User)
				So(ex.ServiceOptions.Password, ShouldEqual, salt.Password)

				Info("And it will run the execution on the web node", " ", 8)
				So(ex2.Name, ShouldEqual, "Execution web 1")
				So(ex2.ExecutionType, ShouldEqual, "fake")
				So(ex2.ExecutionPayload, ShouldEqual, "date")
				So(ex2.ExecutionTarget, ShouldEqual, "list:fake-"+service2+"-web-1")
				So(ex2.ServiceOptions.User, ShouldEqual, salt.User)
				So(ex2.ServiceOptions.Password, ShouldEqual, salt.Password)
			})

			nsub.Unsubscribe()
			isub.Unsubscribe()
			fsub.Unsubscribe()
			asub.Unsubscribe()
			esub.Unsubscribe()
			bsub.Unsubscribe()
		})

		Convey("When I apply a valid novse13.yml definition", func() {
			isub, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			esub, _ := n.ChanSubscribe("execution.create.fake", exCreateSub)
			bsub, _ := n.ChanSubscribe("bootstrap.create.fake", boCreateSub)

			f := getDefinitionPath("novse13.yml", service2)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				i := instanceEvent{}
				msg, err := waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)
				ex := executionEvent{}
				msg, err = waitMsg(boCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &ex)
				ex2 := executionEvent{}
				msg, err = waitMsg(exCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &ex2)

				Info("And I should receive a valid instance.create.vcloud-fake", " ", 8)
				Info("And it should create the second user defined instance ", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service2+"-web-2")
				So(i.Cpus, ShouldEqual, 1)
				So(len(i.Disks), ShouldEqual, 0)
				So(i.IP, ShouldEqual, "10.1.0.12")
				So(i.Memory, ShouldEqual, 1024)
				So(i.ReferenceCatalog, ShouldEqual, "r3")
				So(i.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service2+"-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")

				Info("Then I should receive a valid execution.create.fake", " ", 8)
				Info("And it will bootstrap the web node ", " ", 8)
				So(ex.Name, ShouldEqual, "Bootstrap fake-"+service2+"-web-2")
				So(ex.ExecutionType, ShouldEqual, "fake")
				So(ex.ExecutionPayload, ShouldContainSubstring, "-host 10.1.0.12")
				So(ex.ExecutionTarget, ShouldEqual, "list:salt-master.localdomain")
				So(ex.ServiceOptions.User, ShouldEqual, salt.User)
				So(ex.ServiceOptions.Password, ShouldEqual, salt.Password)

				Info("And it will run the execution on the web node ", " ", 8)
				So(ex2.Name, ShouldEqual, "Execution web 1")
				So(ex2.ExecutionType, ShouldEqual, "fake")
				So(ex2.ExecutionPayload, ShouldEqual, "date")
				So(ex2.ExecutionTarget, ShouldEqual, "list:fake-"+service2+"-web-2")
				So(ex2.ServiceOptions.User, ShouldEqual, salt.User)
				So(ex2.ServiceOptions.Password, ShouldEqual, salt.Password)
			})

			isub.Unsubscribe()
			esub.Unsubscribe()
			bsub.Unsubscribe()
		})

		Convey("When I apply a valid novse14.yml definition", func() {
			esub, _ := n.ChanSubscribe("execution.create.fake", exCreateSub)

			f := getDefinitionPath("novse14.yml", service2)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				event := executionEvent{}
				msg, err := waitMsg(exCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)

				Info("And I should receive a valid execution.create.fake", " ", 8)
				Info("And it will run the updated execution on both web nodes ", " ", 8)
				So(event.Name, ShouldEqual, "Execution web 1")
				So(event.ExecutionType, ShouldEqual, "fake")
				So(event.ExecutionPayload, ShouldEqual, "date; uptime")
				So(event.ExecutionTarget, ShouldEqual, "list:fake-"+service2+"-web-1,fake-"+service2+"-web-2")
				So(event.ServiceOptions.User, ShouldEqual, salt.User)
				So(event.ServiceOptions.Password, ShouldEqual, salt.Password)
			})

			esub.Unsubscribe()
		})

		Convey("When I apply a valid novse15.yml definition", func() {
			isub, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			esub, _ := n.ChanSubscribe("execution.create.fake", exCreateSub)
			bsub, _ := n.ChanSubscribe("bootstrap.create.fake", boCreateSub)

			f := getDefinitionPath("novse15.yml", service2)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				Info("And I should receive a valid instance.create.vcloud-fake", " ", 8)
				i := instanceEvent{}
				msg, err := waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)

				Info("And it should create the third user defined instance ", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service2+"-db-1")
				So(i.Cpus, ShouldEqual, 1)
				So(len(i.Disks), ShouldEqual, 0)
				So(i.IP, ShouldEqual, "10.1.0.21")
				So(i.Memory, ShouldEqual, 1024)
				So(i.ReferenceCatalog, ShouldEqual, "r3")
				So(i.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service2+"-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")

				Info("And I should receive a valid execution.create.fake", " ", 8)
				ex := executionEvent{}
				msg, err = waitMsg(boCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &ex)
				ex2 := executionEvent{}
				msg, err = waitMsg(exCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &ex2)

				Info("And it will bootstrap the db node ", " ", 8)
				So(ex.Name, ShouldEqual, "Bootstrap fake-"+service2+"-db-1")
				So(ex.ExecutionType, ShouldEqual, "fake")
				So(ex.ExecutionPayload, ShouldContainSubstring, "-host 10.1.0.21")
				So(ex.ExecutionTarget, ShouldEqual, "list:salt-master.localdomain")
				So(ex.ServiceOptions.User, ShouldEqual, salt.User)
				So(ex.ServiceOptions.Password, ShouldEqual, salt.Password)
				Printf("\n        And it will run the execution on the db node ")
				So(ex2.Name, ShouldEqual, "Execution db 1")
				So(ex2.ExecutionType, ShouldEqual, "fake")
				So(ex2.ExecutionPayload, ShouldEqual, "date")
				So(ex2.ExecutionTarget, ShouldEqual, "list:fake-"+service2+"-db-1")
				So(ex2.ServiceOptions.User, ShouldEqual, salt.User)
				So(ex2.ServiceOptions.Password, ShouldEqual, salt.Password)
			})

			isub.Unsubscribe()
			esub.Unsubscribe()
			bsub.Unsubscribe()
		})

		Convey("When I apply a valid novse16.yml definition", func() {
			idsub, _ := n.ChanSubscribe("instance.delete.vcloud-fake", inDeleteSub)
			esub, _ := n.ChanSubscribe("execution.create.fake", exCreateSub)

			f := getDefinitionPath("novse16.yml", service2)

			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				Info("And I should receive a valid instance.delete.vcloud-fake", " ", 8)
				i := instanceEvent{}
				msg, err := waitMsg(inDeleteSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)

				Info("And it should create the third user defined instance ", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service2+"-web-2")
				So(i.Cpus, ShouldEqual, 1)
				So(len(i.Disks), ShouldEqual, 0)
				So(i.IP, ShouldEqual, "10.1.0.12")
				So(i.Memory, ShouldEqual, 1024)
				So(i.ReferenceCatalog, ShouldEqual, "r3")
				So(i.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service2+"-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")

				Info("And I should receive a valid execution.create.fake", " ", 8)
				ex := executionEvent{}
				msg, err = waitMsg(exCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &ex)

				Info("And it will remove web-2's key from the salt master ", " ", 8)
				So(ex.Name, ShouldEqual, "Cleanup Bootstrap fake-"+service2+"-web-2")
				So(ex.ExecutionType, ShouldEqual, "fake")
				So(ex.ExecutionPayload, ShouldEqual, "salt-key -y -d fake-"+service2+"-web-2")
				So(ex.ExecutionTarget, ShouldEqual, "list:salt-master.localdomain")
				So(ex.ServiceOptions.User, ShouldEqual, salt.User)
				So(ex.ServiceOptions.Password, ShouldEqual, salt.Password)
			})

			idsub.Unsubscribe()
			esub.Unsubscribe()
		})
	})
}
