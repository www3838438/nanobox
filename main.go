// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/nanobox-io/nanobox/commands"
	"github.com/nanobox-io/nanobox/config"
)

// main
func main() {

	// global panic handler; this is done to avoid showing any panic output if
	// something happens to fail. The output is logged and "pretty" message is
	// shown
	defer func() {
		if r := recover(); r != nil {
			// put r into your log ( it contains the panic message)
			// Then log debug.Stack (from the runtime/debug package)

			stack := debug.Stack()

			fmt.Println("Nanobox encountered an unexpected error. Please see ~.nanobox/nanobox.log and submit the issue to us.")
			config.Log.Fatal(fmt.Sprintf("Cause of failure: %v", r))
			config.Log.Fatal(fmt.Sprintf("Error output:\n%v\n", string(stack)))
			config.Log.Close()
			os.Exit(1)
		}
	}()

	pass := true

	// ensure vagrant is installed
	if err := exec.Command("vagrant", "-v").Run(); err != nil {
		fmt.Println("Missing dependency 'Vagrant'. Please download and install it to continue (https://www.vagrantup.com/).")
		pass = false
	}

	// ensure virtualbox is installed
	if err := exec.Command("vboxmanage", "-v").Run(); err != nil {
		fmt.Println("Missing dependency 'Virtualbox'. Please download and install it to continue (https://www.virtualbox.org/wiki/Downloads).")
		pass = false
	}

	// if a dependency check fails, exit
	if !pass {
		return
	}

	// check for updates
	// checkUpdate()

	//
	commands.NanoboxCmd.Execute()
}
