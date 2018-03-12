[![CircleCI](https://circleci.com/gh/fresh8/consulship.svg?style=svg)](https://circleci.com/gh/fresh8/consulship)[![Go Report Card](https://goreportcard.com/badge/github.com/fresh8/consulship)](https://goreportcard.com/report/github.com/fresh8/consulship)

# Consulship

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

* `$HOME/.consulship`: file containing user consulship configuration as environment variables
* `$PWD/configs/dependencies.yaml`: default configuration for dependencies (committed)
* `$PWD/configs/overrides/dependencies.yaml`: user-specific overrides (gitignored)
* `$PWD/.consulship/consul-env.yaml`: list of consul environments downloaded by the setup script (gitignored)

### Local config
First, create `$HOME/.consulship` and put the following line into it:
```sh
GS_CONSUL_ENVS=gs://<your bucket>/consul-env.yaml
```
This script requires you to store your list of consul environments/servers in a Google Cloud Storage bucket. It is currently the only method of obtaining consul envs using the script.

### Project dependencies
You may now update the dependencies file to correspond to registered consul services across your environments. See [dependencies.yaml example](/configs/dependencies.yaml).
Now create an empty file at `$PWD/configs/overrides/dependencies.yaml` and gitignore it. You will be able to override what is specified in the comitted dependencies file in this overrides file for your custom environment setup.

### Running your service
Run `start-consulship` to start consul in Docker and run consulship to copy and register your dependencies.

## Using consulship directly
The following approach is no longer recommended as we now have a script to automate these steps.

### Running consul

We recommend you install consul through a docker container for simplicity's sake. Ensure that the ports are
being mapped to the default ports (`8500` => `localhost:8500` etc.).

You can test it is mapped correctly and running by hitting `http://localhost:8500` in your browser, this should lead you
to a UI screen which shows registered services.

### Attaching configured services to your consul instance

Once you are running consul and have set up the dependency files, just run `consulship` in the root alongside your `configs` and `.consulship` directory!

We assume that you are running consul on `localhost:8500`, if this is not the case you may override it by passing `LOCAL_CONSUL_ADDR={your IP:PORT}` as an environment
variable on the command, your command will then become `LOCAL_CONSUL_ADDR={your IP:PORT} consulship`
