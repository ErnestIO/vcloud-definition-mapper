/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package integration

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"

	aes "github.com/ernestio/crypto/aes"
	"github.com/nats-io/nats"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVSE(t *testing.T) {
	var service = "vse"
	var encryptedPwd string
	var encryptedUsr string

	crypto := aes.New()
	key := os.Getenv("ERNEST_CRYPTO_KEY")

	service = service + strconv.Itoa(rand.Intn(9999999))

	inCreateSub := make(chan *nats.Msg, 1)
	fiCreateSub := make(chan *nats.Msg, 1)
	roCreateSub := make(chan *nats.Msg, 1)
	neCreateSub := make(chan *nats.Msg, 1)
	naCreateSub := make(chan *nats.Msg, 1)
	fiUpdateSub := make(chan *nats.Msg, 1)
	naUpdateSub := make(chan *nats.Msg, 1)
	inUpdateSub := make(chan *nats.Msg, 1)
	inDeleteSub := make(chan *nats.Msg, 1)
	inDeleteSub2 := make(chan *nats.Msg, 1)
	roDeleteSub := make(chan *nats.Msg, 1)
	inMultipleUpdateSub := make(chan *nats.Msg, 2)
	basicSetup("vcloud")

	Convey("Given I have a configuraed ernest instance", t, func() {
		Convey("When I apply a valid vse1.yml definition", func() {
			subIn, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			subRo, _ := n.ChanSubscribe("router.create.vcloud-fake", roCreateSub)
			subNe, _ := n.ChanSubscribe("network.create.vcloud-fake", neCreateSub)
			subFi, _ := n.ChanSubscribe("firewall.create.vcloud-fake", fiCreateSub)
			subNa, _ := n.ChanSubscribe("nat.create.vcloud-fake", naCreateSub)

			f := getDefinitionPath("vse1.yml", service)

			Convey("Then I should successfully create a valid service", func() {

				Info("And user output should be correct", " ", 6)
				_, err := ernest("service", "apply", f)
				if err != nil {
					log.Println(err.Error())
				}
				r := routerEvent{}
				msg, err := waitMsg(roCreateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &r)
				n := networkEvent{}
				msg, err = waitMsg(neCreateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &n)

				i := instanceEvent{}
				msg, err = waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &i)

				f := firewallEvent{}
				fiMsg, err := waitMsg(fiCreateSub)
				So(err, ShouldBeNil)

				na := natEvent{}
				msg, err = waitMsg(naCreateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &na)

				Info("And it creates router vse4", " ", 8)
				So(r.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(r.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(r.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(r.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(r.DatacenterType, ShouldEqual, "vcloud-fake")
				So(r.RouterName, ShouldEqual, "vse4")
				So(r.RouterType, ShouldEqual, "vcloud-fake")
				So(r.VCloudURL, ShouldNotEqual, "")
				So(r.VseURL, ShouldNotEqual, "")
				So(r.Status, ShouldEqual, "processing")

				Info("And it creates network *-web", " ", 8)
				So(n.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(n.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(n.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(n.DatacenterType, ShouldEqual, "vcloud-fake")
				So(n.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(n.NetworkGateway, ShouldEqual, "10.1.0.1")
				So(n.NetworkNetmask, ShouldEqual, "255.255.255.0")
				So(n.NetworkStartAddress, ShouldEqual, "10.1.0.5")
				So(n.NetworkEndAddress, ShouldEqual, "10.1.0.250")

				Info("Then it creates instance *-web-1", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(i.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(i.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(i.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
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

				Info("Then it configures ACLs on router vse4", " ", 8)
				_ = json.Unmarshal(fiMsg.Data, &f)
				So(f.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(f.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(f.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(f.DatacenterType, ShouldEqual, "vcloud-fake")
				So(f.Type, ShouldEqual, "vcloud-fake")

				So(len(f.Rules), ShouldEqual, 4)
				So(f.RouterIP, ShouldEqual, "1.1.1.1")
				So(f.RouterName, ShouldEqual, "vse4")
				So(f.RouterType, ShouldEqual, "vcloud-fake")

				Info("And it will allow internal:any to internal:any", " ", 8)
				So(f.Rules[0].SourcePort, ShouldEqual, "any")
				So(f.Rules[0].SourceIP, ShouldEqual, "internal")
				So(f.Rules[0].DestinationIP, ShouldEqual, "internal")
				So(f.Rules[0].DestinationPort, ShouldEqual, "any")
				So(f.Rules[0].Protocol, ShouldEqual, "any")

				Info("And it will allow 172.18.143.3:any to internal:22 ", " ", 8)
				So(f.Rules[1].SourcePort, ShouldEqual, "any")
				So(f.Rules[1].SourceIP, ShouldEqual, "172.18.143.3")
				So(f.Rules[1].DestinationIP, ShouldEqual, "internal")
				So(f.Rules[1].DestinationPort, ShouldEqual, "22")
				So(f.Rules[1].Protocol, ShouldEqual, "tcp")

				Info("And it will allow 172.17.240.0/24:any to internal:22 ", " ", 8)
				So(f.Rules[2].SourcePort, ShouldEqual, "any")
				So(f.Rules[2].SourceIP, ShouldEqual, "172.17.240.0/24")
				So(f.Rules[2].DestinationIP, ShouldEqual, "internal")
				So(f.Rules[2].DestinationPort, ShouldEqual, "22")
				So(f.Rules[2].Protocol, ShouldEqual, "tcp")

				Info("And it will allow 172.19.186.30/24:any to internal:22 ", " ", 8)
				So(f.Rules[3].SourcePort, ShouldEqual, "any")
				So(f.Rules[3].SourceIP, ShouldEqual, "172.19.186.30")
				So(f.Rules[3].DestinationIP, ShouldEqual, "internal")
				So(f.Rules[3].DestinationPort, ShouldEqual, "22")
				So(f.Rules[3].Protocol, ShouldEqual, "tcp")

				Info("And it configures NATs on router vse4", " ", 6)
				So(na.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(na.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(na.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(na.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(na.DatacenterType, ShouldEqual, "vcloud-fake")
				So(na.NatName, ShouldEqual, "fake-"+service+"-vse4")
				So(len(na.NatRules), ShouldEqual, 2)
				So(na.RouterIP, ShouldEqual, "1.1.1.1")
				So(na.RouterName, ShouldEqual, "vse4")
				So(na.RouterType, ShouldEqual, "vcloud-fake")
				_, _ = Printf("\n        And it configures from 10.1.0.0/24:any to NAT-IP:any")
				So(na.NatRules[0].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[0].OriginIP, ShouldEqual, "10.1.0.0/24")
				So(na.NatRules[0].OriginPort, ShouldEqual, "any")
				So(na.NatRules[0].Type, ShouldEqual, "snat")
				So(na.NatRules[0].TranslationIP, ShouldEqual, "1.1.1.1")
				So(na.NatRules[0].TranslationPort, ShouldEqual, "any")
				So(na.NatRules[0].Protocol, ShouldEqual, "any")
				_, _ = Printf("\n        And it configures from NAT-IP:22 to 10.1.0.11:22")
				So(na.NatRules[1].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[1].OriginIP, ShouldEqual, "1.1.1.1")
				So(na.NatRules[1].OriginPort, ShouldEqual, "22")
				So(na.NatRules[1].Type, ShouldEqual, "dnat")
				So(na.NatRules[1].TranslationIP, ShouldEqual, "10.1.0.11")
				So(na.NatRules[1].TranslationPort, ShouldEqual, "22")
				So(na.NatRules[1].Protocol, ShouldEqual, "tcp")
			})

			_ = subIn.Unsubscribe()
			_ = subRo.Unsubscribe()
			_ = subNe.Unsubscribe()
			_ = subFi.Unsubscribe()
			_ = subNa.Unsubscribe()
		})

		Convey("When I apply a valid vse2.yml definition", func() {
			subFi, _ := n.ChanSubscribe("firewall.update.vcloud-fake", fiUpdateSub)
			f := getDefinitionPath("vse2.yml", service)

			Convey("Then I should successfully create a valid service", func() {

				Info("And I should get a valid output for a processed service", " ", 8)
				_, err := ernest("service", "apply", f)
				if err != nil {
					log.Println(err.Error())
				}

				event := firewallEvent{}
				msg, err := waitMsg(fiUpdateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &event)

				Info("And it modifies ACLs on router vse4", " ", 8)
				So(event.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(event.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(event.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(event.DatacenterType, ShouldEqual, "vcloud-fake")
				So(event.Type, ShouldEqual, "vcloud-fake")
				So(len(event.Rules), ShouldEqual, 5)
				So(event.RouterName, ShouldEqual, "vse4")
				So(event.RouterIP, ShouldEqual, "1.1.1.1")
				So(event.RouterType, ShouldEqual, "vcloud-fake")

				Info("And it will allow internal:any to external:any", " ", 8)
				So(event.Rules[4].SourcePort, ShouldEqual, "any")
				So(event.Rules[4].SourceIP, ShouldEqual, "172.19.186.30")
				So(event.Rules[4].DestinationIP, ShouldEqual, "internal")
				So(event.Rules[4].DestinationPort, ShouldEqual, "22")
				So(event.Rules[4].Protocol, ShouldEqual, "tcp")
			})

			_ = subFi.Unsubscribe()
		})

		Convey("When I apply a valid vse3.yml definition", func() {
			subNat, _ := n.ChanSubscribe("nat.update.vcloud-fake", naUpdateSub)
			f := getDefinitionPath("vse3.yml", service)

			Convey("Then I should modify vse service", func() {

				Info("And I should get a valid output for a processed service", " ", 8)
				_, err := ernest("service", "apply", f)
				if err != nil {
					log.Println(err.Error())
				}

				event := natEvent{}
				msg, err := waitMsg(naUpdateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &event)

				Info("And it modifies ACLs on router vse4", " ", 8)
				So(event.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(event.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(event.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(event.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(event.DatacenterType, ShouldEqual, "vcloud-fake")
				So(event.NatName, ShouldEqual, "fake-"+service+"-vse4")
				So(len(event.NatRules), ShouldEqual, 3)
				So(event.RouterIP, ShouldEqual, "1.1.1.1")
				So(event.RouterName, ShouldEqual, "vse4")
				So(event.RouterType, ShouldEqual, "vcloud-fake")

				Info("And it configures from 10.1.0.0/24:any to NAT-IP:any", " ", 8)
				So(event.NatRules[0].Network, ShouldEqual, "NETWORK")
				So(event.NatRules[0].OriginIP, ShouldEqual, "10.1.0.0/24")
				So(event.NatRules[0].OriginPort, ShouldEqual, "any")
				So(event.NatRules[0].Type, ShouldEqual, "snat")
				So(event.NatRules[0].TranslationIP, ShouldEqual, "1.1.1.1")
				So(event.NatRules[0].TranslationPort, ShouldEqual, "any")
				So(event.NatRules[0].Protocol, ShouldEqual, "any")

				Info("And it configures from NAT-IP:22 to 10.1.0.11:22", " ", 8)
				So(event.NatRules[1].Network, ShouldEqual, "NETWORK")
				So(event.NatRules[1].OriginIP, ShouldEqual, "1.1.1.1")
				So(event.NatRules[1].OriginPort, ShouldEqual, "22")
				So(event.NatRules[1].Type, ShouldEqual, "dnat")
				So(event.NatRules[1].TranslationIP, ShouldEqual, "10.1.0.11")
				So(event.NatRules[1].TranslationPort, ShouldEqual, "22")
				So(event.NatRules[1].Protocol, ShouldEqual, "tcp")

				Info("And it configures from NAT-IP:23 to 10.1.0.12:23", " ", 8)
				So(event.NatRules[2].Network, ShouldEqual, "NETWORK")
				So(event.NatRules[2].OriginIP, ShouldEqual, "1.1.1.1")
				So(event.NatRules[2].OriginPort, ShouldEqual, "23")
				So(event.NatRules[2].Type, ShouldEqual, "dnat")
				So(event.NatRules[2].TranslationIP, ShouldEqual, "10.1.0.12")
				So(event.NatRules[2].TranslationPort, ShouldEqual, "23")
				So(event.NatRules[2].Protocol, ShouldEqual, "tcp")
			})

			_ = subNat.Unsubscribe()
		})

		Convey("When I apply a valid vse4.yml definition", func() {
			subInCreate, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			subInUpdate, _ := n.ChanSubscribe("instance.update.vcloud-fake", inUpdateSub)

			f := getDefinitionPath("vse4.yml", service)
			Convey("Then I should get a valid output for a processed service", func() {

				Info("And I should get a valid output for a processed service", " ", 8)
				_, err := ernest("service", "apply", f)
				if err != nil {
					log.Println(err.Error())
				}

				Info("Then it will create web-2 instance", " ", 8)
				ic := instanceEvent{}
				msg, err := waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &ic)
				i := instanceEvent{}
				msg, err = waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &i)

				So(ic.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(ic.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(ic.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(ic.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(ic.DatacenterType, ShouldEqual, "vcloud-fake")
				So(ic.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				So(ic.Cpus, ShouldEqual, 1)
				So(len(ic.Disks), ShouldEqual, 0)
				So(ic.IP, ShouldEqual, "10.1.0.12")
				So(ic.Memory, ShouldEqual, 1024)
				So(ic.ReferenceCatalog, ShouldEqual, "r3")
				So(ic.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(ic.InstanceType, ShouldEqual, "vcloud-fake")
				So(ic.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(ic.RouterIP, ShouldEqual, "")
				So(ic.RouterName, ShouldEqual, "")
				So(ic.RouterType, ShouldEqual, "")

				Info("And it will update web-2 instance", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(i.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(i.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(i.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
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
			})

			_ = subInCreate.Unsubscribe()
			_ = subInUpdate.Unsubscribe()
		})

		Convey("When I apply a valid vse5.yml definition", func() {
			subInUpdate, _ := n.ChanSubscribe("instance.update.vcloud-fake", inMultipleUpdateSub)

			f := getDefinitionPath("vse5.yml", service)
			_, err := ernest("service", "apply", f)
			Convey("Then service should be successfully processed", func() {
				Info("And I should get a valid output for a processed service", " ", 8)
				if err != nil {
					log.Println(err.Error())
				}

				ui := instanceEvent{}
				msg, err := waitMsg(inMultipleUpdateSub)
				_ = subInUpdate.Unsubscribe()
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &ui)
				ui2 := instanceEvent{}
				msg, err = waitMsg(inMultipleUpdateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &ui2)

				Info("And it will update web-1 instance", " ", 8)
				So(ui.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(ui.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(ui.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(ui.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(ui.DatacenterType, ShouldEqual, "vcloud-fake")
				So(ui.InstanceName, ShouldEqual, "fake-"+service+"-web-1")
				_, _ = Printf("\n        And it will update CPU from 1 to 2")
				So(ui.Cpus, ShouldEqual, 2)
				So(len(ui.Disks), ShouldEqual, 0)
				So(ui.IP, ShouldEqual, "10.1.0.11")
				So(ui.Memory, ShouldEqual, 1024)
				So(ui.ReferenceCatalog, ShouldEqual, "r3")
				So(ui.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(ui.InstanceType, ShouldEqual, "vcloud-fake")
				So(ui.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(ui.RouterIP, ShouldEqual, "")
				So(ui.RouterName, ShouldEqual, "")
				So(ui.RouterType, ShouldEqual, "")

				Info("Then it will update web-2 instance", " ", 8)
				So(ui2.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(ui2.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(ui2.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(ui2.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(ui2.DatacenterType, ShouldEqual, "vcloud-fake")
				So(ui2.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				_, _ = Printf("\n        And it will update CPU from 1 to 2")
				So(ui2.Cpus, ShouldEqual, 2)
				So(len(ui2.Disks), ShouldEqual, 0)
				So(ui2.IP, ShouldEqual, "10.1.0.12")
				So(ui2.Memory, ShouldEqual, 1024)
				So(ui2.ReferenceCatalog, ShouldEqual, "r3")
				So(ui2.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(ui2.InstanceType, ShouldEqual, "vcloud-fake")
				So(ui2.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(ui2.RouterIP, ShouldEqual, "")
				So(ui2.RouterName, ShouldEqual, "")
				So(ui2.RouterType, ShouldEqual, "")
			})

			_ = subInUpdate.Unsubscribe()
		})

		Convey("When I apply a valid vse6.yml definition", func() {
			subInUpdate, _ := n.ChanSubscribe("instance.update.vcloud-fake", inMultipleUpdateSub)

			f := getDefinitionPath("vse6.yml", service)
			_, err := ernest("service", "apply", f)
			Convey("Then it should successfully process the service", func() {
				Info("And I should get a valid output for a processed service", " ", 8)
				if err != nil {
					log.Println(err.Error())
				}

				ui1 := instanceEvent{}
				msg, err := waitMsg(inMultipleUpdateSub)
				_ = subInUpdate.Unsubscribe()
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &ui1)
				ui2 := instanceEvent{}
				msg, err = waitMsg(inMultipleUpdateSub)
				So(err, ShouldBeNil)

				Info("And it will update web-1 instance", " ", 8)
				_ = json.Unmarshal(msg.Data, &ui2)
				So(ui1.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(ui1.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(ui1.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(ui1.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(ui1.DatacenterType, ShouldEqual, "vcloud-fake")
				So(ui1.InstanceName, ShouldEqual, "fake-"+service+"-web-1")
				So(ui1.Cpus, ShouldEqual, 2)

				Info("And adds a 10GB of disk", " ", 8)
				So(len(ui1.Disks), ShouldEqual, 1)
				So(ui1.Disks[0].ID, ShouldEqual, 1)
				So(ui1.Disks[0].Size, ShouldEqual, 10240)
				So(ui1.IP, ShouldEqual, "10.1.0.11")
				So(ui1.Memory, ShouldEqual, 1024)
				So(ui1.ReferenceCatalog, ShouldEqual, "r3")
				So(ui1.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(ui1.InstanceType, ShouldEqual, "vcloud-fake")
				So(ui1.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(ui1.RouterIP, ShouldEqual, "")
				So(ui1.RouterName, ShouldEqual, "")
				So(ui1.RouterType, ShouldEqual, "")

				Info("Then it will update web-2 instance", " ", 8)
				So(ui2.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(ui2.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(ui2.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(ui2.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(ui2.DatacenterType, ShouldEqual, "vcloud-fake")
				So(ui2.InstanceName, ShouldEqual, "fake-"+service+"-web-2")

				Info("And it will update CPU from 1 to 2", " ", 8)
				So(ui2.Cpus, ShouldEqual, 2)

				Info("And adds a 10GB of disk", " ", 8)
				So(len(ui2.Disks), ShouldEqual, 1)
				So(ui2.Disks[0].ID, ShouldEqual, 1)
				So(ui2.Disks[0].Size, ShouldEqual, 10240)
				So(ui2.IP, ShouldEqual, "10.1.0.12")
				So(ui2.Memory, ShouldEqual, 1024)
				So(ui2.ReferenceCatalog, ShouldEqual, "r3")
				So(ui2.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(ui2.InstanceType, ShouldEqual, "vcloud-fake")
				So(ui2.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(ui2.RouterIP, ShouldEqual, "")
				So(ui2.RouterName, ShouldEqual, "")
				So(ui2.RouterType, ShouldEqual, "")
			})

			_ = subInUpdate.Unsubscribe()
		})

		Convey("When I apply a valid vse7.yml definition", func() {
			subInUpdate, _ := n.ChanSubscribe("instance.update.vcloud-fake", inMultipleUpdateSub)

			f := getDefinitionPath("vse7.yml", service)
			_, err := ernest("service", "apply", f)
			Convey("Then it will successfully process the service", func() {
				Info("Then I should get a valid output for a processed service", " ", 8)
				if err != nil {
					log.Println(err.Error())
				}

				ui1 := instanceEvent{}
				msg, err := waitMsg(inMultipleUpdateSub)
				_ = subInUpdate.Unsubscribe()
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &ui1)
				ui2 := instanceEvent{}
				msg, err = waitMsg(inMultipleUpdateSub)
				So(err, ShouldBeNil)

				Info("Then it will update web-1 instance", " ", 8)
				_ = json.Unmarshal(msg.Data, &ui2)
				So(ui1.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(ui1.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(ui1.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(ui1.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(ui1.DatacenterType, ShouldEqual, "vcloud-fake")
				So(ui1.InstanceName, ShouldEqual, "fake-"+service+"-web-1")
				So(ui1.Cpus, ShouldEqual, 2)
				So(len(ui1.Disks), ShouldEqual, 1)
				So(ui1.Disks[0].ID, ShouldEqual, 1)
				So(ui1.Disks[0].Size, ShouldEqual, 10240)
				So(ui1.IP, ShouldEqual, "10.1.0.11")
				_, _ = Printf("\n        And upgrades to 2GB RAM")
				So(ui1.Memory, ShouldEqual, 2048)
				So(ui1.ReferenceCatalog, ShouldEqual, "r3")
				So(ui1.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(ui1.InstanceType, ShouldEqual, "vcloud-fake")
				So(ui1.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(ui1.RouterIP, ShouldEqual, "")
				So(ui1.RouterName, ShouldEqual, "")
				So(ui1.RouterType, ShouldEqual, "")

				Info("And it will update web-2 instance", " ", 8)
				So(ui2.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(ui2.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(ui2.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(ui2.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(ui2.DatacenterType, ShouldEqual, "vcloud-fake")
				So(ui2.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				So(ui2.Cpus, ShouldEqual, 2)
				So(len(ui2.Disks), ShouldEqual, 1)
				So(ui2.Disks[0].ID, ShouldEqual, 1)
				So(ui2.Disks[0].Size, ShouldEqual, 10240)
				So(ui2.IP, ShouldEqual, "10.1.0.12")
				_, _ = Printf("\n        And upgrades to 2GB RAM")
				So(ui2.Memory, ShouldEqual, 2048)
				So(ui2.ReferenceCatalog, ShouldEqual, "r3")
				So(ui2.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(ui2.InstanceType, ShouldEqual, "vcloud-fake")
				So(ui2.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(ui2.RouterIP, ShouldEqual, "")
				So(ui2.RouterName, ShouldEqual, "")
				So(ui2.RouterType, ShouldEqual, "")
			})

			_ = subInUpdate.Unsubscribe()
		})

		Convey("When I apply a valid vse8.yml definition", func() {
			subNeCreate, _ := n.ChanSubscribe("network.create.vcloud-fake", neCreateSub)
			subNaUpdate, _ := n.ChanSubscribe("nat.update.vcloud-fake", naUpdateSub)

			f := getDefinitionPath("vse8.yml", service)
			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				n := networkEvent{}
				msg, err := waitMsg(neCreateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &n)

				na := natEvent{}
				msg, err = waitMsg(naUpdateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &na)

				Info("And it will create new network", " ", 8)
				So(n.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(n.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(n.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(n.DatacenterType, ShouldEqual, "vcloud-fake")
				So(n.NetworkName, ShouldEqual, "fake-"+service+"-db")
				So(n.NetworkGateway, ShouldEqual, "10.2.0.1")
				So(n.NetworkNetmask, ShouldEqual, "255.255.255.0")
				So(n.NetworkStartAddress, ShouldEqual, "10.2.0.5")
				So(n.NetworkEndAddress, ShouldEqual, "10.2.0.250")

				Info("And it modifies ACLs on router vse4 to reconfigure new network", " ", 8)
				So(na.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(na.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(na.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(na.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(na.DatacenterType, ShouldEqual, "vcloud-fake")
				So(na.NatName, ShouldEqual, "fake-"+service+"-vse4")
				So(len(na.NatRules), ShouldEqual, 4)
				So(na.RouterIP, ShouldEqual, "1.1.1.1")
				So(na.RouterName, ShouldEqual, "vse4")
				So(na.RouterType, ShouldEqual, "vcloud-fake")
				So(na.NatRules[0].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[0].OriginIP, ShouldEqual, "10.1.0.0/24")
				So(na.NatRules[0].OriginPort, ShouldEqual, "any")
				So(na.NatRules[0].Type, ShouldEqual, "snat")
				So(na.NatRules[0].TranslationIP, ShouldEqual, "1.1.1.1")
				So(na.NatRules[0].TranslationPort, ShouldEqual, "any")
				So(na.NatRules[0].Protocol, ShouldEqual, "any")

				Info("And it configures from NAT-IP:23 to 10.1.0.12:23", " ", 8)
				So(na.NatRules[1].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[1].OriginIP, ShouldEqual, "10.2.0.0/24")
				So(na.NatRules[1].OriginPort, ShouldEqual, "any")
				So(na.NatRules[1].Type, ShouldEqual, "snat")
				So(na.NatRules[1].TranslationIP, ShouldEqual, "1.1.1.1")
				So(na.NatRules[1].TranslationPort, ShouldEqual, "any")
				So(na.NatRules[1].Protocol, ShouldEqual, "any")
				So(na.NatRules[2].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[2].OriginIP, ShouldEqual, "1.1.1.1")
				So(na.NatRules[2].OriginPort, ShouldEqual, "22")
				So(na.NatRules[2].Type, ShouldEqual, "dnat")
				So(na.NatRules[2].TranslationIP, ShouldEqual, "10.1.0.11")
				So(na.NatRules[2].TranslationPort, ShouldEqual, "22")
				So(na.NatRules[2].Protocol, ShouldEqual, "tcp")
				So(na.NatRules[3].Network, ShouldEqual, "NETWORK")
				So(na.NatRules[3].OriginIP, ShouldEqual, "1.1.1.1")
				So(na.NatRules[3].OriginPort, ShouldEqual, "23")
				So(na.NatRules[3].Type, ShouldEqual, "dnat")
				So(na.NatRules[3].TranslationIP, ShouldEqual, "10.1.0.12")
				So(na.NatRules[3].TranslationPort, ShouldEqual, "23")
				So(na.NatRules[3].Protocol, ShouldEqual, "tcp")
			})

			_ = subNeCreate.Unsubscribe()
			_ = subNaUpdate.Unsubscribe()
		})

		Convey("When I apply a valid vse9.yml definition", func() {
			subInCreate, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			subInUpdate, _ := n.ChanSubscribe("instance.update.vcloud-fake", inUpdateSub)

			f := getDefinitionPath("vse9.yml", service)
			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				i := instanceEvent{}
				msg, err := waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &i)
				iu := instanceEvent{}
				msg, err = waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &iu)

				Info("And it will create db-1 instance", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(i.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(i.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(i.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
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

				Info("And it will update db-1 instance", " ", 8)
				So(iu.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(iu.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(iu.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(iu.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(iu.DatacenterType, ShouldEqual, "vcloud-fake")
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

			_ = subInCreate.Unsubscribe()
			_ = subInUpdate.Unsubscribe()
		})

		Convey("When I apply a valid vse10.yml definition", func() {
			subInDelete, _ := n.ChanSubscribe("instance.delete.vcloud-fake", inDeleteSub)

			f := getDefinitionPath("vse10.yml", service)
			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				event := instanceEvent{}
				msg, err := waitMsg(inDeleteSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &event)

				Info("Then it will delete web-2 instance", " ", 8)
				So(event.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(event.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(event.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(event.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(event.DatacenterType, ShouldEqual, "vcloud-fake")
				So(event.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				So(event.Cpus, ShouldEqual, 2)
				So(len(event.Disks), ShouldEqual, 1)
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

			_ = subInDelete.Unsubscribe()
		})

		Convey("When I apply a valid vse11.yml definition", func() {
			subInDelete, _ := n.ChanSubscribe("instance.delete.vcloud-fake", inDeleteSub)

			f := getDefinitionPath("vse11.yml", service)
			_, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				event := instanceEvent{}
				msg, err := waitMsg(inDeleteSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &event)

				Info("Then it will delete db-1 instance", " ", 8)
				So(event.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(event.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(event.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(event.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(event.DatacenterType, ShouldEqual, "vcloud-fake")
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

			_ = subInDelete.Unsubscribe()
		})

		Convey("When I destroy the current service", func() {
			subInDelete, _ := n.ChanSubscribe("instance.delete.vcloud-fake", inDeleteSub2)
			subRoDelete, _ := n.ChanSubscribe("router.delete.vcloud-fake", roDeleteSub)
			_, err := ernest("service", "destroy", "--yes", service)

			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				i := instanceEvent{}
				msg, err := waitMsg(inDeleteSub2)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &i)
				r := routerEvent{}
				msg, err = waitMsg(roDeleteSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &r)

				Info("Then it will delete web-1 instance", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(i.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(i.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(i.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-web-1")
				So(i.Cpus, ShouldEqual, 2)
				So(len(i.Disks), ShouldEqual, 1)
				So(i.IP, ShouldEqual, "10.1.0.11")
				So(i.Memory, ShouldEqual, 2048)
				So(i.ReferenceCatalog, ShouldEqual, "r3")
				So(i.ReferenceImage, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")

				Info("Then it deletes router vse4", " ", 8)
				So(r.DatacenterName, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(r.DatacenterPassword, key)
				encryptedUsr, _ = crypto.Decrypt(r.DatacenterUsername, key)
				So(defaultPwd, ShouldEqual, encryptedPwd)
				So(defaultUsr+"@"+defaultOrg, ShouldEqual, encryptedUsr)
				So(r.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(r.DatacenterType, ShouldEqual, "vcloud-fake")
				So(r.RouterName, ShouldEqual, "vse4")
				So(r.RouterType, ShouldEqual, "vcloud-fake")
				So(r.VCloudURL, ShouldNotEqual, "")
				So(r.VseURL, ShouldNotEqual, "")
				So(r.Status, ShouldEqual, "processing")
			})

			_ = subInDelete.Unsubscribe()
			_ = subRoDelete.Unsubscribe()
		})
	})
}
