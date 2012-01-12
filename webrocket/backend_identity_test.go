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

package webrocket

import (
	zmq "../gozmq"
	"testing"
)

func TestParseBackendIdentityWithValidIdentity(t *testing.T) {
	raw := "dlr:/hello/there:1234567890abcdefghij1234567890abcdefghij:12345678-90ab-cdef-ghij-1234567890ab"
	identity, err := parseBackendIdentity([]byte(raw))
	if err != nil {
		t.Errorf("Expected to parse a valid identity")
		return
	}
	if identity.Type != zmq.DEALER {
		t.Errorf("Expected to parse identity type")
	}
	if identity.AccessToken != "1234567890abcdefghij1234567890abcdefghij" {
		t.Errorf("Expected to parse access token")
	}
	if identity.Vhost != "/hello/there" {
		t.Errorf("Expected to parse vhost")
	}
	if identity.Id != "12345678-90ab-cdef-ghij-1234567890ab" {
		t.Errorf("Expected to parse client id")
	}
}

func TestParseBackendIdentityWhenInvalidIdentity(t *testing.T) {
	_, err := parseBackendIdentity([]byte("invalid"))
	if err == nil {
		t.Errorf("Expected an error while parsing invalid identity")
	}
}
