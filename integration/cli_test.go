/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats"
)

var salt struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

var adminUsr = "ci_admin"
var adminPwd = "pwd"
var defaultUsr = "usr"
var defaultPwd = "pwd"
var defaultOrg = "org"
var ernestInstance = "https://ernest.local/"
var endSub = make(chan *nats.Msg, 1)

var setup = false
var n *nats.Conn

func wait(ch chan bool) error {
	return waitTime(ch, 500*time.Millisecond)
}

func waitTime(ch chan bool, timeout time.Duration) error {
	select {
	case <-ch:
		return nil
	case <-time.After(timeout):
	}
	return errors.New("timeout")
}

func waitMsg(ch chan *nats.Msg) (*nats.Msg, error) {
	select {
	case msg := <-ch:
		return msg, nil
	case <-time.After(time.Millisecond * 10000):
	}
	return nil, errors.New("timeout")
}

func waitToDone() {
	subEnd, _ := n.ChanSubscribe("service.create.done", endSub)
	_, _ = waitMsg(endSub)
	_ = subEnd.Unsubscribe()
}

func getDefinitionPath(def string, service string) string {
	finalPath := "/tmp/currentTest.yml"

	_, filename, _, _ := runtime.Caller(1)
	filePath := path.Join(path.Dir(filename), "definitions", def)

	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
	var finalLines []string

	for _, line := range lines {
		if strings.Contains(line, "name: my_service") {
			finalLines = append(finalLines, "name: "+service)
		} else if strings.Contains(line, "datacenter: r3-dc2") {
			finalLines = append(finalLines, "datacenter: fake")
		} else {
			finalLines = append(finalLines, line)
		}
	}
	output := strings.Join(finalLines, "\n")
	err = ioutil.WriteFile(finalPath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	return finalPath
}

func getDefinitionPathAWS(def string, service string) string {
	finalPath := "/tmp/currentTest.yml"

	_, filename, _, _ := runtime.Caller(1)
	filePath := path.Join(path.Dir(filename), "definitions", def)

	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
	var finalLines []string

	for _, line := range lines {
		if strings.Contains(line, "name: my_service") {
			finalLines = append(finalLines, "name: "+service)
		} else if strings.Contains(line, "datacenter: r3-dc2") {
			finalLines = append(finalLines, "datacenter: fakeaws")
		} else {
			finalLines = append(finalLines, line)
		}
	}
	output := strings.Join(finalLines, "\n")
	err = ioutil.WriteFile(finalPath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	return finalPath
}

func basicSetup(provider string) {
	if setup == false {
		var err error
		n, err = nats.Connect(os.Getenv("NATS_URI"))
		if err != nil {
			panic(err)
		}

		msg, err := n.Request("config.get.salt", []byte("{}"), time.Second)
		if err != nil {
			panic("Salt config not accessible")
		}
		_ = json.Unmarshal(msg.Data, &salt)

		if os.Getenv("CURRENT_INSTANCE") != "" {
			ernestInstance = os.Getenv("CURRENT_INSTANCE")
		}
		_, _ = ernest("target", ernestInstance)
		_, _ = ernest("login", "--user", adminUsr, "--password", adminPwd)

		// Create user
		_, _ = ernest("user", "create", defaultUsr, defaultPwd)
		_, _ = ernest("group", "create", "test")
		_, _ = ernest("group", "add-user", defaultUsr, "test")

		// Login as this user
		login()

		// Create a datacenter
		_, _ = ernest("datacenter", "create", "vcloud", "fake", "--vcloud-url", "https://myvdc.me.com", "--fake", "--user", defaultUsr, "--password", defaultPwd, "--org", defaultOrg, "--vse-url", "http://localhost", "--public-network", "NETWORK")
		_, _ = ernest("datacenter", "create", "aws", "fakeaws", "--region", "fake", "--secret_access_key", "fake_up_to_16_chars", "--access_key_id", "secret_up_to_16_chars", "--fake")

		setup = true
	} else {
		// Login as this user
		login()
	}
}

func login() {
	_, _ = ernest("login", "--user", defaultUsr, "--password", defaultPwd)
}

func deleteConfig() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove(usr.HomeDir + "/.ernest")
}

func ernest(cmdArgs ...string) (string, error) {
	if cmdArgs[1] == "apply" {
		if delay := os.Getenv("ERNEST_APPLY_DELAY"); delay != "" {
			if t, err := strconv.Atoi(delay); err == nil {
				println("\nWaiting " + delay + " seconds...")
				time.Sleep(time.Duration(t) * time.Second)
			}
		}
	}
	cmd := exec.Command("ernest-cli", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		println(err.Error())
	}
	return string(output), nil
}

func Info(str, pad string, l int) {
	for i := 0; i < l; i++ {
		str = pad + str
	}

	print("\n" + str + " ")
}

func CheckOutput(ol, cl []string) bool {
	for i, v := range ol {
		if i < len(cl) {
			if cl[i] != "" {
				if cl[i] != v {
					fmt.Printf("\nOutput line %d expected to be: \n%s\nbut found\n%s", i, cl[i], v)
					return false
				}
			}
		}
	}
	return true
}

func main() {}
