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

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"go.uber.org/zap"
)

var (
	pluginPrefix = "security.secrets"
	pluginName   = pluginPrefix + ".aws_secrets_manager"

	// Interface guards
	_ caddy.Provisioner     = (*Plugin)(nil)
	_ caddy.Validator       = (*Plugin)(nil)
	_ caddyfile.Unmarshaler = (*Plugin)(nil)
	_ caddy.Module          = (*Plugin)(nil)
)

func init() {
	caddy.RegisterModule(Plugin{})
}

// Config represents provisioned configuration value of AWS Secrets Manager.
type Config struct {
	ID     string `json:"id,omitempty" xml:"id,omitempty" yaml:"id,omitempty"`
	Region string `json:"region,omitempty" xml:"region,omitempty" yaml:"region,omitempty"`
	Path   string `json:"path,omitempty" xml:"path,omitempty" yaml:"path,omitempty"`
}

// Plugin manages AWS Secret Manager integration.
type Plugin struct {
	Name      string          `json:"-"`
	ConfigRaw json.RawMessage `json:"config,omitempty" caddy:"namespace=security.secrets.aws_secrets_manager"`
	Config    Config          `json:"-"`
	logger    *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Plugin) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "security.secrets.aws_secrets_manager",
		New: func() caddy.Module { return new(Plugin) },
	}
}

// Provision sets up Handler and loads AwsSecretsManager.
func (p *Plugin) Provision(ctx caddy.Context) error {
	p.Name = pluginName
	p.logger = ctx.Logger(p)

	p.logger.Info(
		"provisioning plugin instance",
		zap.String("plugin_name", p.Name),
	)

	if err := json.Unmarshal(p.ConfigRaw, &p.Config); err != nil {
		p.logger.Info(
			"failed configuring plugin instance",
			zap.String("plugin_name", p.Name),
			zap.Error(err),
		)
		return err
	}

	p.logger.Info(
		"provisioned plugin instance",
		zap.String("plugin_name", p.Name),
	)
	return nil
}

// Validate implements caddy.Validator.
func (p *Plugin) Validate() error {
	p.logger.Info(
		"validating plugin instance",
		zap.String("plugin_name", p.Name),
	)
	if p.Config.ID == "" {
		return fmt.Errorf("empty id")
	}
	if p.Config.Path == "" {
		return fmt.Errorf("secret %q has empty path", p.Config.ID)
	}
	if p.Config.Region == "" {
		return fmt.Errorf("secret %q has empty region", p.Config.ID)
	}
	p.logger.Info(
		"validated plugin instance",
		zap.String("plugin_name", p.Name),
		zap.String("secret_id", p.Config.ID),
	)
	return nil
}
