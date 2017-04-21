# ConsulShip

[![Dependency Status](https://dependencyci.com/github/fresh8/consulship/badge)](https://dependencyci.com/github/fresh8/consulship)[![CircleCI](https://circleci.com/gh/fresh8/consulship.svg?style=svg)](https://circleci.com/gh/fresh8/consulship)

A tool to help relieve the growing pains of developing copious microservices
with Consul service discovery by:

* Providing defaults for a platform
* Allowing service-level overrides
* Connecting to dependencies in different environments

## Flags

* `-consulEnv`: Configuration of Consul environments
* `-baseDeps`: Default dependencies for a platform
* `-deps`: Service-level overrides of dependencies

## Installing

```go
go get github.com/fresh8/consulship
```

This will install `consulship` and it's dependencies within your `GOPATH`.

## Using

### Files and directories

Create the following files within your project, it's advisible to copy them from this repo.

* `$PWD/configs/dependencies.yaml`: default configuration for dependencies (committed)
* `$PWD/configs/overrides/dependencies.yaml`: user-specific overrides (gitignored)
* `$PWD/.consulship/consul-env.yaml`: list of consul environments downloaded by the setup script (gitignored)

Update your `consul-env.yaml` file to point at the relevant consul servers for that environment. We have included advice on whether or not
you should gitignore the files, however it is ultimately up to you.

You may now update the dependencies file to correspond to registered consul services across your environments.

### Running consul

We recommend you install consul through a docker container for simplicity's sake. Ensure that the ports are
being mapped to the default ports (`8500` => `localhost:8500` etc.).

You can test it is mapped correctly and running by hitting `http://localhost:8500` in your browser, this should lead you
to a UI screen which shows registered services.

### Attaching configured services to your consul instance

Once you are running consul and have set up the dependency files, just run `consulship` in the root alongside your `configs` and `.consulship` directory!

We assume that you are running consul on `localhost:8500`, if this is not the case you may override it by passing `LOCAL_CONSUL_ADDR={your IP:PORT}` as an environment
variable on the command, your command will then become `LOCAL_CONSUL_ADDR={your IP:PORT} consulship`
