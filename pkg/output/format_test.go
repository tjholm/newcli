// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/nitrictech/newcli/pkg/containerengine"
	"github.com/nitrictech/newcli/pkg/target"
)

func Test_printStruct(t *testing.T) {
	tests := []struct {
		name   string
		object interface{}
		expect string
	}{
		{
			name:   "json tags",
			object: target.Target{Name: "test", Provider: "azure", Region: "somewhere"},
			expect: `+----------+-----------+
| NAME     | test      |
| PROVIDER | azure     |
| REGION   | somewhere |
+----------+-----------+
`,
		},
		{
			name:   "yaml tags",
			object: containerengine.Image{ID: "34234242", Repository: "my-app", Tag: "latest"},
			expect: `+------------+----------+
| ID         | 34234242 |
| REPOSITORY | my-app   |
| TAG        | latest   |
| CREATEDAT  |          |
+------------+----------+
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			printStruct(tt.object, buf)
			if !cmp.Equal(tt.expect, buf.String()) {
				t.Error(cmp.Diff(tt.expect, buf.String()))
			}
		})
	}
}

func Test_printList(t *testing.T) {
	tests := []struct {
		name   string
		object interface{}
		expect string
	}{
		{
			name: "json tags",
			object: []target.Target{
				{Name: "test", Provider: "azure", Region: "somewhere"},
				{Name: "local", Provider: "local"},
			},
			expect: `+-------+----------+-----------+
| NAME  | PROVIDER | REGION    |
+-------+----------+-----------+
| test  | azure    | somewhere |
| local | local    |           |
+-------+----------+-----------+
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			printList(tt.object, buf)
			if !cmp.Equal(tt.expect, buf.String()) {
				t.Error(cmp.Diff(tt.expect, buf.String()))
			}
		})
	}
}

func Test_printMap(t *testing.T) {
	tests := []struct {
		name    string
		object  interface{}
		wantOut string
	}{
		{
			name: "json tags",
			object: map[string]target.Target{
				"t1":    {Name: "test", Provider: "azure", Region: "somewhere"},
				"local": {Name: "local", Provider: "local"},
			},
			wantOut: `+-------+-------+----------+-----------+
| KEY   | NAME  | PROVIDER | REGION    |
+-------+-------+----------+-----------+
| t1    | test  | azure    | somewhere |
| local | local | local    |           |
+-------+-------+----------+-----------+
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			printMap(tt.object, out)
			if !cmp.Equal(tt.wantOut, out.String()) {
				t.Error(cmp.Diff(tt.wantOut, out.String()))
			}
		})
	}
}
