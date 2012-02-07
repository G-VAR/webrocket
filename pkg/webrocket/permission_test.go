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

import "testing"

func TestNewPermission(t *testing.T) {
	p, err := NewPermission("joe", ".*")
	if err != nil {
		t.Errorf("Expected to create a new permission, error: %v", err)
	}
	if len(p.Token()) != 128 {
		t.Errorf("Expected to generate single access token for the permission")
	}
	_, err = NewPermission("joe", "%%&**")
	if err == nil {
		t.Errorf("Expected an error when creating new permission with invalid regexp")
	}
}

func TestPermissionIsMatching(t *testing.T) {
	p, _ := NewPermission("joe", ".*foo|bar.*")
	for _, ch := range []string{"lefoo", "barle"} {
		if !p.IsMatching(ch) {
			t.Errorf("Expected permission to match the '%s' channel", ch)
		}
	}
	for _, ch := range []string{"lefoobar", "lebar"} {
		if p.IsMatching(ch) {
			t.Errorf("Expected permission to not match the '%s' channel", ch)
		}
	}
}
