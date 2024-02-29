# go-openchargemap

[go-openchargemap](https://www.github.com/neomantra/go-openchargemap) is a Golang interface to the [OpenChargeMap API](https://openchargemap.org/site).

The base [`openchargemap.openapi.yml`](./openchargemap.openapi.yml) comes from [OpenChargeMaps's Developer Documentation](https://openchargemap.org/site/develop/api#/).

You will need an OpenChargeMap API key -- you can get one by [Signing Up](https://openchargemap.org/site/loginprovider/beginlogin) and then going to [My Apps](https://openchargemap.org/site/profile/applications) in your Profile.   The examples and tools look for `OCM_KEY` in your environment, so `export OCM_KEY=<your-key>`

## Examples

We have the following examples:

 * `chargemeup` shows how to use the [`go-openchargemap` library](./cmd/chargemeup/main.go)

 * [`ocum_fun.ipynb` Jupyter notebook](./examples/ocm_fun.ipynb) showing how to use the `chargemeup` CLI tool to query the OpenChargeMap API within Python notebooks.  

 * [`ocum_3d.ipynb` Jupyter notebook](./examples/ocm_3d.ipynb) is the same thing, but also generates an [OSMBuildings](https://osmbuildings.org) WebGL experience, along with wiht instructions to use it in VR.

 * [`examples/chargers_near_addr.py`](./examples/chargers_near_addr.py) is a simple Python script that exercises it too


## Go Library

You can include it in your Go programs with:

```sh
go get github.com/neomantra/go-openchargemap
```

See the `chargemeup` [sample program](./cmd/chargemeup/main.go) for an example with how to use the API.  The bindings are created with [oapi-codegen](https://github.com/deepmap/oapi-codegen).

## `chargemeup` CLI tool

`chargemeup` is a command-line tool for querying the OpenChargeMap API.  Although you can [build it yourself](#building), you can install it with [Homebrew](https://brew.sh):

```
brew tap neomantra/homebrew-tap
brew install neomantra/homebrew-tap/go-openchargemap
```

You can run `chargemeup` to get a list of charging stations near a bounding box `(lat1,lon1),(lat2,lon2)`:

```
$ chargemeup -b "(40.63010790372053,-74.2775717248681),(40.7356464076158,-74.09370618215354)" | jq length  
41

$ chargemeup -a "Carroll Gardens, NY" --radius 1 | jq '.[].AddressInfo.AddressLine1'
"383 Court St"
"365 Bond St"
"363 Bond St"
"214 3rd Street"
"353 4th Ave"
"302 2ND St"

$ chargemeup --lat 40.63010790372053 --lon -74.2775717248681 --radius 5 | jq length
20
```

The output is the JSON, although there we are using [`jq`](https://jqlang.github.io/jq/) to count the number of POIs returned.

```
$ chargemeup --help
usage:  ./chargemeup <options> [input]

"chargemeup" assists with queries to OpenChargeMap.

Around Newark is:
chargemeup -p "(40.63010790372053,-74.2775717248681),(40.7356464076158,-74.09370618215354)"

chargemup -a "Newark, NJ" -r 10

chargemeup --lat 40.7356464076158 --lon -74.09370618215354 --radius 5

  -a, --address string   address to query (requires --radius)
  -b, --bbox string      bounding box for query, "(lat1,lon1),(lat2,lon2)"
  -h, --help             show help
  -k, --key string       API key for OpenChargeMap, env var OCM_KEY)
      --lat float32      latitude to query (requires --lon and --radius)
      --lon float32      longitude to query (requires --lon and --radius)
  -r, --radius float32   radial distance to query in kilometers (requires --address)
  -s, --server string    API Server for OpenChargeMap, env var OCM_SERVER (default "https://api.openchargemap.io/v3")
  -v, --verbose          verbose output to stderr
```

# Building

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

Authored with :heart: and :fire: by Evan Wies.  Copyright (c) 2024 Neomantra BV.

Many thanks to [OpenChargeMap](https://openchargemap.org/) and the community behind it!  Thanks to [oapi-codegen](https://github.com/deepmap/oapi-codegen) for the heavy lifting.

Released under the [MIT License](https://en.wikipedia.org/wiki/MIT_License), see [LICENSE.txt](./LICENSE.txt).
