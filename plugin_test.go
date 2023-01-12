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
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/caddyserver/caddy/v2"
	"github.com/google/go-cmp/cmp"
)

func TestProvisionPlugin(t *testing.T) {
	testcases := []struct {
		name      string
		cfg       string
		want      Config
		shouldErr bool
		err       error
	}{
		{
			name: "test provisioning valid config",
			cfg:  `{"id":"foo","path":"foo/bar","region":"us-east-1"}`,
			want: Config{
				ID:     "foo",
				Path:   "foo/bar",
				Region: "us-east-1",
			},
		},
		{
			name:      "test provisioning malformed json config",
			cfg:       `{"id":"foo","path":"foo/bar","region":"us-east-1"`,
			shouldErr: true,
			err:       fmt.Errorf("unexpected end of JSON input"),
		},
		{
			name:      "test provisioning config with unmarshal",
			cfg:       `[]`,
			shouldErr: true,
			err:       fmt.Errorf("json: cannot unmarshal array into Go value of type secretsmanager.Config"),
		},
		{
			name:      "test provisioning config without id",
			cfg:       `{"path":"foo/bar","region":"us-east-1"}`,
			shouldErr: true,
			err:       fmt.Errorf("empty id"),
		},
		{
			name:      "test provisioning config with malformed region",
			cfg:       `{"id":"foo","path":"foo/bar","region":"foo-bar-baz"}`,
			shouldErr: true,
			err:       fmt.Errorf("malformed %q region", "foo-bar-baz"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Plugin{
				ConfigRaw: json.RawMessage(tc.cfg),
			}
			err := p.Provision(caddy.ActiveContext())
			if err != nil {
				if !tc.shouldErr {
					t.Fatalf("expected success, got: %v", err)
				}
				if diff := cmp.Diff(err.Error(), tc.err.Error()); diff != "" {
					t.Logf("unexpected error: %v", err)
					t.Fatalf("Provision() error mismatch (-want +got):\n%s", diff)
				}
				return
			}
			if tc.shouldErr {
				t.Fatalf("unexpected success, want: %v", tc.err)
			}

			if diff := cmp.Diff(tc.want, p.Config); diff != "" {
				t.Logf("JSON: %s", p.ConfigRaw)
				t.Fatalf("Provision() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidatePlugin(t *testing.T) {
	jsmith := map[string]interface{}{
		"api_key":  "bcrypt:10:$2a$10$TEQ7ZG9cAdWwhQK36orCGOlokqQA55ddE0WEsl00oLZh567okdcZ6",
		"email":    "jsmith@localhost.localdomain",
		"name":     "John Smith",
		"password": "bcrypt:10:$2a$10$iqq53VjdCwknBSBrnyLd9OH1Mfh6kqPezMMy6h6F41iLdVDkj13I6",
		"username": "jsmith",
	}

	testcases := []struct {
		name      string
		cfg       string
		secret    map[string]interface{}
		want      map[string]interface{}
		shouldErr bool
		err       error
	}{
		{
			name: "test validating valid config",
			cfg:  `{"id":"foo","path":"foo/bar","region":"us-east-1"}`,
			want: map[string]interface{}{
				"id":       "foo",
				"path":     "foo/bar",
				"region":   "us-east-1",
				"provider": "aws_secrets_manager",
			},
			secret: jsmith,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Plugin{
				ConfigRaw: json.RawMessage(tc.cfg),
			}
			err := p.Provision(caddy.ActiveContext())
			if err != nil {
				t.Fatalf("unexpected provisioning error: %v", err)
			}

			var mockClinet aws.HTTPClient = smithyhttp.ClientDoFunc(func(r *http.Request) (*http.Response, error) {
				response := packMapToJSON(t, map[string]interface{}{
					"SecretString": packMapToJSON(t, tc.secret),
				})
				return &http.Response{
					StatusCode: 200,
					Header:     http.Header{},
					Body:       ioutil.NopCloser(strings.NewReader(response)),
				}, nil
			})

			p.client.SetMockClient(mockClinet)

			err = p.Validate()
			if err != nil {
				if !tc.shouldErr {
					t.Fatalf("expected success, got: %v", err)
				}
				if diff := cmp.Diff(err.Error(), tc.err.Error()); diff != "" {
					t.Logf("unexpected error: %v", err)
					t.Fatalf("Validate() error mismatch (-want +got):\n%s", diff)
				}
				return
			}
			if tc.shouldErr {
				t.Fatalf("unexpected success, want: %v", tc.err)
			}

			got := p.GetConfig(caddy.ActiveContext())
			if diff := cmp.Diff(tc.want, got); diff != "" {
				// t.Logf("JSON: %s", p.ConfigRaw)
				t.Errorf("Validate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
