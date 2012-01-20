// This package implements executable for starting and preconfiguring
// single webrocket server node.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package main

import (
	stepper "../gostepper"
	"../webrocket"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// Configuration variables.
var (
	BackendAddr   string
	WebsocketAddr string
	AdminAddr     string
	NodeName      string
	CertFile      string
	KeyFile       string
	StorageDir    string
)

var (
	ctx *webrocket.Context
	s   stepper.Stepper
)

func init() {
	flag.StringVar(&WebsocketAddr, "websocket-addr", ":8080", "websocket endpoint address")
	flag.StringVar(&BackendAddr, "backend-addr", ":8081", "backend endpoint address")
	flag.StringVar(&AdminAddr, "admin-addr", ":8082", "admin endpoint address")
	flag.StringVar(&NodeName, "node-name", "", "name of the node")
	flag.StringVar(&CertFile, "cert", "", "path to server certificate")
	flag.StringVar(&KeyFile, "key", "", "private key")
	flag.StringVar(&StorageDir, "storage-dir", "/var/lib/webrocket", "path to webrocket's internal data-store")
	flag.Parse()

	StorageDir, _ = filepath.Abs(StorageDir)
}

func SetupContext() {
	s.Start("Initializing context")
	ctx = webrocket.NewContext()
	ctx.SetStorageDir(StorageDir)
	if NodeName != "" {
		if err := ctx.SetNodeName(NodeName); err != nil {
			s.Fail(err.Error(), true)
		}
	}
	s.Ok()
	s.Start("Locking node")
	if err := ctx.Lock(); err != nil {
		s.Fail(err.Error(), true)
	}
	s.Ok()
	s.Start("Loading configuration")
	if err := ctx.Load(); err != nil {
		s.Fail(err.Error(), true)
	}
	s.Ok()
	s.Start("Generating cookie")
	if err := ctx.GenerateCookie(false); err != nil {
		s.Fail(err.Error(), true)
	}
	s.Ok()
}

func SetupEndpoint(kind string, e webrocket.Endpoint) {
	go func() {
		var err error
		s.Start("Starting %s", kind)
		if CertFile != "" && KeyFile != "" {
			err = e.ListenAndServeTLS(CertFile, KeyFile)
		} else {
			err = e.ListenAndServe()
		}
		if err != nil {
			s.Fail(err.Error(), true)
		}
	}()
	for !e.IsAlive() {
		<-time.After(500 * time.Nanosecond)
	}
	s.Ok()
}

func SignalTrap() {
	for sig := range signal.Incoming {
		if usig, ok := sig.(os.UnixSignal); ok {
			switch usig {
			case os.SIGQUIT, os.SIGINT:
				fmt.Printf("\n\033[33mInterrupted\033[0m\n")
				if ctx != nil {
					fmt.Printf("\n")
					s.Start("Cleaning up")
					if err := ctx.Kill(); err != nil {
						s.Fail(err.Error(), true)
					}
					s.Ok()
				}
				os.Exit(0)
			case os.SIGTSTP:
				syscall.Kill(syscall.Getpid(), syscall.SIGSTOP)
			case os.SIGHUP:
				// TODO: reload configuration
			}
		}
	}
}

func SetupDaemon() {
	fmt.Printf("\n\033[32mWebRocket has been launched!\033[0m\n")
}

func DisplayAsciiArt() {
	fmt.Printf("\n")
	fmt.Printf(
		`            /\                                                                     ` + "\n" +
		`      ,    /  \      o               .        ___---___                    .       ` + "\n" +
		`          /    \            .              .--\        --.     .     .         .   ` + "\n" +
		`         /______\                        ./.;_.\     __/~ \.                       ` + "\n" +
		`   .    |        |                      /;  / '-'  __\    . \                      ` + "\n" +
		`        |        |    .        .       / ,--'     / .   .;   \        |            ` + "\n" +
		`        |________|                    | .|       /       __   |      -O-       .   ` + "\n" +
		`        |________|                   |__/    __ |  . ;   \ | . |      |            ` + "\n" +
		`       /|   ||   |\                  |      /  \\_    . ;| \___|                   ` + "\n" +
		`      / |   ||   | \    .    o       |      \  .~\\___,--'     |           .       ` + "\n" +
		`     /  |   ||   |  \                 |     | . ; ~~~~\_    __|                    ` + "\n" +
		`    /___|:::||:::|___\   |             \    \   .  .  ; \  /_/   .                 ` + "\n" +
		`        |::::::::|      -O-        .    \   /         . |  ~/                  .   ` + "\n" +
		`         \::::::/         |    .          ~\ \   .      /  /~          o           ` + "\n" +
		`   o      ||__||       .                   ~--___ ; ___--~                         ` + "\n" +
		`            ||                        .          ---         .              .      ` + "\n" +
		`            ''                                                                     ` + "\n")
	fmt.Printf("WebRocket v%s\n", webrocket.Version())
	fmt.Printf("Copyright (C) 2011-2012 by Krzysztof Kowalik and folks at Cubox.\n")
	fmt.Printf("Released under the AGPL. See http://www.webrocket.io/ for details.\n\n")
}

func DisplaySystemSettings() {
	fmt.Printf("\n")
	fmt.Printf("Node               : %s\n", ctx.NodeName())
	fmt.Printf("Cookie             : %s\n", ctx.Cookie())
	fmt.Printf("Data store dir     : %s\n", ctx.StorageDir())
	fmt.Printf("Websocket endpoint : ws://%s\n", WebsocketAddr)
	fmt.Printf("Backend endpoint   : wr://%s\n", BackendAddr)
	fmt.Printf("Admin endpoint     : http://%s\n", AdminAddr)
}

func main() {
	DisplayAsciiArt()
	SetupContext()
	SetupEndpoint("backend endpoint", ctx.NewBackendEndpoint(BackendAddr))
	SetupEndpoint("websocket endpoint", ctx.NewWebsocketEndpoint(WebsocketAddr))
	SetupEndpoint("admin endpoint", ctx.NewAdminEndpoint(AdminAddr))
	DisplaySystemSettings()
	SetupDaemon()
	SignalTrap()
}
