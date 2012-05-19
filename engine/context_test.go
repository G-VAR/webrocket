// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package engine

import (
	"io"
	"os"
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext()
	if ctx.Log() == nil {
		t.Errorf("Expected context logger to be initialized")
	}
}

func TestContextSetLog(t *testing.T) {
	ctx := NewContext()
	ctx.SetLog(nil)
	if ctx.Log() != nil {
		t.Errorf("Expected to set other logger")
	}
}

func TestContextAddVhost(t *testing.T) {
	ctx := NewContext()
	v, err := ctx.AddVhost("/foo")
	if err != nil || v == nil {
		t.Errorf("Expected to add vhost")
	}
	_, err = ctx.AddVhost("/foo")
	if err == nil || err.Error() != "vhost already exists" {
		t.Errorf("Expected error while adding duplicated vhost")
	}
	_, ok := ctx.vhosts["/foo"]
	if !ok {
		t.Errorf("Expected to add vhost")
	}
}

func TestContextAddVhostWhenWebsocketEndpointPresent(t *testing.T) {
	ctx := NewContext()
	e := ctx.NewWebsocketEndpoint("localhost:3000")
	w := e.(*WebsocketEndpoint)
	v, _ := ctx.AddVhost("/foo")
	h := w.handlers.Match("/foo")
	if h == nil || h.vhost.Path() != v.Path() {
		t.Errorf("Expected to register vhost in websocket endpoint")
	}
}

func TestContextDeleteVhost(t *testing.T) {
	ctx := NewContext()
	ctx.AddVhost("/foo")
	err := ctx.DeleteVhost("/foo")
	if err != nil {
		t.Errorf("Expected to delete vhost without errors")
	}
	_, ok := ctx.vhosts["/foo"]
	if ok {
		t.Errorf("Expected to delete vhost")
	}
	err = ctx.DeleteVhost("/foo")
	if err == nil || err.Error() != "vhost doesn't exist" {
		t.Errorf("Expected an error while deleting non existent vhost")
	}
}

func TestContextDeleteVhostWhenWebsocketEndpointPresent(t *testing.T) {
	ctx := NewContext()
	e := ctx.NewWebsocketEndpoint("localhost:3000")
	w := e.(*WebsocketEndpoint)
	ctx.AddVhost("/foo")
	ctx.DeleteVhost("/foo")
	h := w.handlers.Match("/foo")
	if h != nil {
		t.Errorf("Expected to unregister vhost from websocket endpoint")
	}
}

func TestContextGetVhost(t *testing.T) {
	ctx := NewContext()
	ctx.AddVhost("/foo")
	v, err := ctx.Vhost("/foo")
	if err != nil || v == nil || v.Path() != "/foo" {
		t.Errorf("Expected to get vhost")
	}
	_, err = ctx.Vhost("/bar")
	if err == nil || err.Error() != "vhost doesn't exist" {
		t.Errorf("Expected an error getting non existent vhost")
	}
}

func TestContextVhostsList(t *testing.T) {
	ctx := NewContext()
	ctx.AddVhost("/foo")
	if len(ctx.Vhosts()) != 1 {
		t.Errorf("Expected vhosts list to contain one element")
	}
}

func TestContextCookiesGeneration(t *testing.T) {
	ctx := NewContext()
	ctx.SetStorageDir("/tmp")
	ctx.GenerateCookie(false)
	cookieFile := "/tmp/" + DefaultNodeName() + ".cookie"
	f, err := os.Open(cookieFile)
	if err != nil {
		t.Errorf("Expected to create a cookie file")
		return
	}
	n, err := io.ReadFull(f, make([]byte, 40)[:])
	if err != nil || n != CookieSize {
		t.Errorf("Expected to write cookie to the file")
	}
	f.Close()
	os.Remove(cookieFile)
}

func TestContextSetNodeName(t *testing.T) {
	ctx := NewContext()
	if err := ctx.SetNodeName("&**()"); err == nil {
		t.Errorf("Expected error while setting invalid node name")
	}
	if err := ctx.SetNodeName("foo"); err != nil {
		t.Errorf("Expected to set valid node name without errors")
	}
}

func TestContextNewWebsocketEndpoint(t *testing.T) {
	ctx := NewContext()
	ctx.NewWebsocketEndpoint(":9772")
	if ctx.websocket == nil || ctx.websocket.Addr() != ":9772" {
		t.Errorf("Expected to register new websocket endpoint")
	}
}

func TestContextNewBackendEndpoint(t *testing.T) {
	ctx := NewContext()
	ctx.NewBackendEndpoint(":9772")
	if ctx.backend == nil || ctx.backend.Addr() != ":9772" {
		t.Errorf("Expected to register new backend endpoint")
	}
}

func TestContextKill(t *testing.T) {
	ctx := NewContext()
	ctx.NewWebsocketEndpoint(":9772")
	if ctx.websocket == nil {
		t.Errorf("Expected to set websocket endpoint")
	}
	ctx.NewBackendEndpoint(":9773")
	if ctx.backend == nil {
		t.Errorf("Expected to set backend endpoint")
	}
	ctx.Kill()
	if ctx.backend.IsAlive() || ctx.websocket.IsAlive() {
		t.Errorf("Expected to close and kill all endpoints")
	}
}
