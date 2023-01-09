# caddy-security-creds-aws-secrets-manager

[Caddy Security](https://github.com/greenpau/caddy-security) Credentials Plugin
for AWS Secrets Manager Integration.

<!-- begin-markdown-toc -->
## Table of Contents

* [Getting Started](#getting-started)
  * [AWS Secrets Manager](#aws-secrets-manager)
  * [Caddyfile Usage](#caddyfile-usage)

<!-- end-markdown-toc -->

## Getting Started

### AWS Secrets Manager

Please follow this [doc](https://github.com/greenpau/go-authcrunch-creds-aws-secrets-manager#getting-started)
to set up AWS IAM Policy, Rolem and Secrets.

### Caddyfile Usage

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

Now, here is the configuration using `credentials` retrieved from AWS Secrets Manager:

```
{
        security {
                credentials access_token_secret {
                        driver aws_secrets_manager
                        path authcrunch/caddy/access_token
                }

                credentials users_jsmith {
                        driver aws_secrets_manager
                        path authcrunch/caddy/users/jsmith
                }

                local identity store localdb {
                        realm local
                        path /etc/caddy/users.json
                        user jsmith {
                                name "credentials:users_jsmith:name"
                                email "credentials:users_jsmith:email"
                                password "credentials:users_jsmith:password" overwrite
                                api_key "credentials:users_jsmith:api_key" overwrite
                                roles authp/admin authp/user
                        }
                }

                authentication portal myportal {
                        crypto default token lifetime 3600
                        crypto key sign-verify "credentials:access_token_secret:value"
                        enable identity store localdb
                }
        }
}
```
