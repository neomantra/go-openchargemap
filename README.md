# go-openchargemap

[go-openchargemap](https://www.github.com/neomantra/go-openchargemap) is a Golang interface to the [OpenChargeMap API](https://openchargemap.org/site).

The base [`openchargemap.openapi.yml`](./openchargemap.openapi.yml) comes from [OpenChargeMaps's Developer Documentation](https://openchargemap.org/site/develop/api#/).

## Installation

You can include it in your Go programs with:

```sh
go get github.com/neomantra/go-openchargemap
```

See the `chargemeup` sample program for an example with how to use the API.  The bindings are created with [oapi-codegen](https://github.com/deepmap/oapi-codegen).

You after building (see below), you can run `chargemeup` to get a list of charging stations near a bounding box:

```
./chargemeup -b "(40.63010790372053,-74.2775717248681),(40.7356464076158,-74.09370618215354)" | jq length  
41
```
The output is the JSON, although there we are using [`jq`](https://jqlang.github.io/jq/) to count the number of POIs returned.

## Building

Building is performed with [task](https://taskfile.dev/) and our [Taskfile.yml](./Taskfile.yml):

```
$ task --list-all
task: Available tasks for this project:
* build:              
* build-deps:         
* clean:              
* default:            
* oapi-codegen:       
* test:               
* tidy:

# needed when the upstream spec changes
$ task oapi-codegen

$ task build
task: [tidy] go mod tidy
task: [build] go build -o chargemeup cmd/chargemeup/main.go
```

## Credits and License

Author with :heart: and :fire: by Evan Wies.  Copyright (c) 2024 Neomantra BV.

Many thanks to [OpenChargeMap](https://openchargemap.org/) and the community behind it!  Thanks to [oapi-codegen](https://github.com/deepmap/oapi-codegen) for the heavy lifting.

Released under the [MIT License](https://en.wikipedia.org/wiki/MIT_License), see [LICENSE.txt](./LICENSE.txt).
