// Copyright 2022 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package secretsmanager

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/google/go-cmp/cmp"
)

func unpack(t *testing.T, i interface{}) (m map[string]interface{}) {
	switch v := i.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			t.Fatalf("failed to parse %q: %v", v, err)
		}
	default:
		b, err := json.Marshal(i)
		if err != nil {
			t.Fatalf("failed to marshal %T: %v", i, err)
		}
		if err := json.Unmarshal(b, &m); err != nil {
			t.Fatalf("failed to parse %q: %v", b, err)
		}
	}
	return m
}

func TestUnmarshalCaddyfile(t *testing.T) {
	testcases := []struct {
		name      string
		d         *caddyfile.Dispenser
		want      map[string]interface{}
		shouldErr bool
		err       error
	}{
		{
			name: "test valid config",
			d:    caddyfile.NewTestDispenser(testCfg1),
			want: map[string]interface{}{
				"id":     "access_token",
				"path":   "authcrunch/caddy/access_token",
				"region": "us-east-1",
			},
		},
		{
			name:      "test config without path",
			d:         caddyfile.NewTestDispenser(testCfg2),
			shouldErr: true,
			err:       fmt.Errorf("Testfile:%d - Error during parsing: secret %q has empty path", 4, "access_token"),
		},
		{
			name:      "test config without region",
			d:         caddyfile.NewTestDispenser(testCfg3),
			shouldErr: true,
			err:       fmt.Errorf("Testfile:%d - Error during parsing: secret %q has empty region", 4, "access_token"),
		},
		{
			name:      "test config without id",
			d:         caddyfile.NewTestDispenser(testCfg4),
			shouldErr: true,
			err:       fmt.Errorf("Testfile:%d - Error during parsing: secret %q has empty path", 2, "{"),
		},
		{
			name:      "test invalid region value",
			d:         caddyfile.NewTestDispenser(testCfg5),
			shouldErr: true,
			err: fmt.Errorf(
				"Testfile:%d - Error during parsing: field %q of %q secret with value of %q has invalid syntax",
				3, "region", "", []string{"us-east-1", "foo"},
			),
		},
		{
			name:      "test invalid path value",
			d:         caddyfile.NewTestDispenser(testCfg6),
			shouldErr: true,
			err: fmt.Errorf(
				"Testfile:%d - Error during parsing: field %q of %q secret with value of %q has invalid syntax",
				4, "path", "", []string{"authcrunch/caddy/access_token", "foo"},
			),
		},
		{
			name:      "test invalid key",
			d:         caddyfile.NewTestDispenser(testCfg7),
			shouldErr: true,
			err: fmt.Errorf(
				"Testfile:%d - Error during parsing: unsupported %q field of %q secret with value of %q",
				3, "foo", "", []string{"bar"},
			),
		},
		{
			name:      "test unexpected end of config",
			d:         caddyfile.NewTestDispenser(testCfg8),
			shouldErr: true,
			err:       fmt.Errorf(":0 - Error during parsing: unexpected end of configuration"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Plugin{}
			err := p.UnmarshalCaddyfile(tc.d)
			if err != nil {
				if !tc.shouldErr {
					t.Fatalf("expected success, got: %v", err)
				}
				if diff := cmp.Diff(err.Error(), tc.err.Error()); diff != "" {
					t.Logf("unexpected error: %v", err)
					t.Fatalf("UnmarshalCaddyfile() error mismatch (-want +got):\n%s", diff)
				}
				return
			}
			if tc.shouldErr {
				t.Fatalf("unexpected success, want: %v", tc.err)
			}

			got := unpack(t, string(p.ConfigRaw))

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Logf("JSON: %s", p.ConfigRaw)
				t.Errorf("UnmarshalCaddyfile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

var testCfg1 = `
access_token {
	region us-east-1
	path authcrunch/caddy/access_token
}
`

var testCfg2 = `
access_token {
	region us-east-1
}
`

var testCfg3 = `
access_token {
	path authcrunch/caddy/access_token
}
`

var testCfg4 = `
{
	region us-east-1
	path authcrunch/caddy/access_token
}
`

var testCfg5 = `
access_token {
	region us-east-1 foo
	path authcrunch/caddy/access_token
}
`

var testCfg6 = `
access_token {
	region us-east-1
	path authcrunch/caddy/access_token foo
}
`

var testCfg7 = `
access_token {
	foo bar
}
`

var testCfg8 = `

`
