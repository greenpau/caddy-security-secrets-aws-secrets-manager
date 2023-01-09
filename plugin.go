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

package credentials

import (
	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2"
)

func init() {
	caddy.RegisterModule(CredentialManager{})
}

// AwsSecretsManagerConfig represents provisioned configuration value of AwsSecretsManager.
type AwsSecretsManagerConfig struct {
	Path string `json:"path,omitempty" xml:"path,omitempty" yaml:"path,omitempty"`
}

// AwsSecretsManager represents a host module for the ddd
type CredentialManager struct {
	AwsSecretsManagerRaw json.RawMessage         `json:"credentials_aws_secrets_manager,omitempty" caddy:"namespace=security.credentials.aws_secrets_manager"`
	AwsSecretsManager    AwsSecretsManagerConfig `json:"-"`
}

// CaddyModule returns the Caddy module information.
func (CredentialManager) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "security.credentials.aws_secrets_manager",
		New: func() caddy.Module { return new(CredentialManager) },
	}
}

// Provision sets up CredentialManager and loads AwsSecretsManager.
func (m *CredentialManager) Provision(ctx caddy.Context) error {
	if m.AwsSecretsManagerRaw != nil {
		val, err := ctx.LoadModule(m, "AwsSecretsManagerRaw")
		if err != nil {
			return fmt.Errorf("loading security.credentials.aws_secrets_manager module: %v", err)
		}
		m.AwsSecretsManager = val.(AwsSecretsManagerConfig)
	}
	return nil
}

// Validate implements caddy.Validator.
func (m *CredentialManager) Validate() error {
	if m.AwsSecretsManager.Path == "" {
		return fmt.Errorf("empty path")
	}
	return nil
}

// Interface guards
var (
	_ caddy.Provisioner = (*CredentialManager)(nil)
	_ caddy.Validator   = (*CredentialManager)(nil)
)
