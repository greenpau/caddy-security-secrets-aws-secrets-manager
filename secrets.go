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
	"context"
)

// GetSecret return a secret in the form of a key-value map.
func (p *Plugin) GetSecret(_ context.Context) (map[string]interface{}, error) {
	secret := map[string]interface{}{
		"foo": "bar",
	}
	return secret, nil
}

// GetSecretByKey return a value of key in the secret key-value map.
func (p *Plugin) GetSecretByKey(_ context.Context, key string) (interface{}, error) {
	return "bar", nil
}
