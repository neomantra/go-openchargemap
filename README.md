# go-openchargemap

[go-openchargemap](https://www.github.com/neomantra/go-openchargemap) is a Golang interface to the [OpenChargeMap API](https://openchargemap.org/site).

The base [`openchargemap.openapi.yml`](./openchargemap.openapi.yml) comes from [OpenChargeMaps's Developer Documentation](https://openchargemap.org/site/develop/api#/).

## Installation

```sh
echo TODO hello world
```

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

Many thanks to [OpenChargeMap](https://openchargemap.org/) and the community behind it!

Released under the [MIT License](https://en.wikipedia.org/wiki/MIT_License), see [LICENSE.txt](./LICENSE.txt).
