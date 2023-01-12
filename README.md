# caddy-security-secrets-aws-secrets-manager

<a href="https://github.com/greenpau/addy-security-secrets-aws-secrets-manager/actions/" target="_blank"><img src="https://github.com/greenpau/addy-security-secrets-aws-secrets-manager/workflows/build/badge.svg?branch=main"></a>
<a href="https://pkg.go.dev/github.com/greenpau/addy-security-secrets-aws-secrets-manager" target="_blank"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>

[Caddy Security](https://github.com/greenpau/caddy-security) Secrets Plugin
for AWS Secrets Manager Integration.

<!-- begin-markdown-toc -->
## Table of Contents

* [Getting Started](#getting-started)
  * [AWS Secrets Manager](#aws-secrets-manager)
  * [Building Caddy](#building-caddy)
  * [Caddyfile Usage](#caddyfile-usage)
    * [Without Plugin](#without-plugin)
    * [Plugin Configuration](#plugin-configuration)

<!-- end-markdown-toc -->

## Getting Started

### AWS Secrets Manager

Please follow this [doc](https://github.com/greenpau/go-authcrunch-secrets-aws-secrets-manager#getting-started)
to set up AWS IAM Policy, Rolem and Secrets.

### Building Caddy

For `secrets aws_secrets_manager` directives to work, build `caddy` with the
`latest` version of this plugin.

```bash
xcaddy build ... \
  --with go-authcrunch-secrets-aws-secrets-manager@latest
```

### Caddyfile Usage

#### Without Plugin

The following is a snippet of `Caddyfile` without the use of this plugin.

```
{
        security {
                local identity store localdb {
                        realm local
                        path /etc/caddy/users.json
                        user jsmith {
                                name John Smith
                                email jsmith@localhost.localdomain
                                password "bcrypt:10:$2a$10$iqq53VjdCwknBSBrnyLd9OH1Mfh6kqPezMMy6h6F41iLdVDkj13I6" overwrite
                                roles authp/admin authp/user
                        }
                }

                authentication portal myportal {
                        crypto default token lifetime 3600
                        crypto key sign-verify b006d65b-c923-46a1-8da1-7d52558508fe
                        enable identity store localdb
                }
        }
}
```

#### Plugin Configuration

Now, here is the configuration using `secrets` retrieved from AWS Secrets Manager:

```
{
	security {
		secrets aws_secrets_manager access_token {
			region us-east-1
			path authcrunch/caddy/access_token
		}

		secrets aws_secrets_manager users/jsmith {
			region us-east-1
			path authcrunch/caddy/users/jsmith
		}

		local identity store localdb {
			realm local
			path users.json
			user jsmith {
				name "secrets:users/jsmith:name"
				email "secrets:users/jsmith:email"
				password "secrets:users/jsmith:password" overwrite
				api_key "secrets:users/jsmith:api_key" overwrite
				roles authp/admin authp/user
			}
		}

		authentication portal myportal {
			crypto default token lifetime 3600
			crypto key sign-verify "secrets:access_token:value"
			enable identity store localdb
		}
	}
}
```
