# ConsulShip

A tool to help relieve the growing pains of developing copious microservices
with Consul service discovery by:

* Providing defaults for a platform
* Allowing service-level overrides
* Connecting to dependencies in different environments

## Flags

* `-consulEnv`: Configuration of Consul environments
* `-baseDeps`: Default dependencies for a platform
* `-deps`: Service-level overrides of dependencies

## Files and directories

* `$PWD/configs/dependencies.yaml`: default configuration for dependencies (committed)
* `$PWD/configs/overrides/dependencies.yaml`: user-specific overrides (gitignored)
* `$PWD/.consulship/consul-env.yaml`: list of consul environments downloaded by the setup script (gitignored)
