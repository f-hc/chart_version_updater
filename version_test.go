// SPDX-License-Identifier: GPL-3.0-only
//
// Copyright (C) 2026 f-hc <207619282+f-hc@users.noreply.github.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import "testing"

func TestVersionLess(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"equal versions", "1.0.0", "1.0.0", false},
		{"major less", "1.0.0", "2.0.0", true},
		{"major greater", "2.0.0", "1.0.0", false},
		{"minor less", "1.1.0", "1.2.0", true},
		{"minor greater", "1.2.0", "1.1.0", false},
		{"patch less", "1.0.1", "1.0.2", true},
		{"patch greater", "1.0.2", "1.0.1", false},
		{"different lengths a shorter", "1.0", "1.0.1", true},
		{"different lengths b shorter", "1.0.1", "1.0", false},
		{"two digit versions", "1.10.0", "1.9.0", false},
		{"large versions", "10.20.30", "10.20.29", false},
		{"v prefix stripped externally", "1.19.1", "1.19.2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := versionLess(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("versionLess(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
