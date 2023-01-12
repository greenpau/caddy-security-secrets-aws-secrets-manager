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

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func findReplaceAll(repl *caddy.Replacer, arr []string) (output []string) {
	for _, item := range arr {
		output = append(output, repl.ReplaceAll(item, "CADDY_REPLACEMENT_FAILED"))
	}
	return output
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (p *Plugin) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	repl := caddy.NewReplacer()

	if !d.Next() {
		return d.Err("unexpected end of configuration")
	}

	p.Config.ID = d.Val()

	for d.NextBlock(0) {
		k := d.Val()
		v := findReplaceAll(repl, d.RemainingArgs())
		switch k {
		case "path":
			if len(v) != 1 {
				return d.Errf("field %q of %q secret with value of %q has invalid syntax", k, p.Name, v)
			}
			p.Config.Path = v[0]
		case "region":
			if len(v) != 1 {
				return d.Errf("field %q of %q secret with value of %q has invalid syntax", k, p.Name, v)
			}
			p.Config.Region = v[0]
		default:
			return d.Errf("unsupported %q field of %q secret with value of %q", k, p.Name, v)
		}
	}

	cfg, _ := json.Marshal(p.Config)
	p.ConfigRaw = json.RawMessage(cfg)

	if err := p.ValidateConfig(); err != nil {
		return d.Errf("%v", err)
	}

	return nil
}
