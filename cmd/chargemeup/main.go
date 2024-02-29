package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	openchargemap "github.com/neomantra/go-openchargemap"
	"github.com/spf13/pflag"
	nominatim "github.com/yuriizinets/go-nominatim"
)

/////////////////////////////////////////////////////////////////////////////////////

var usageFormat string = `usage:  %s <options> [input]

"chargemeup" assists with queries to OpenChargeMap.

Around Newark is:
chargemeup -p "(40.63010790372053,-74.2775717248681),(40.7356464076158,-74.09370618215354)"

chargemup -a "Newark, NJ" -r 10

chargemeup --lat 40.7356464076158 --lon -74.09370618215354 --radius 5

`

const defaultServer = "https://api.openchargemap.io/v3"

type ChargeMeUpConfig struct {
	Server  string
	APIKey  string
	Verbose bool

	// Area to search, in order of priority:
	//   BoundingBox
	//   Latitude + Longitude + Radius
	BoundingBox string
	Latitude    float32
	Longitude   float32
	Radius      float32
}

// Returns nil if the Search Area in the Config is 'roughly valid', otherwise an error describing the issue with it.
// 'Roughly valid' means that it has the appropriate data, not that the data makes any sense or that OCM will accept it.
func (c ChargeMeUpConfig) IsAreaValid() error {
	// see order in ChargeMeUp struct above
	if c.BoundingBox != "" {
		return nil
	}
	if c.Latitude == 0 && c.Longitude == 0 {
		return fmt.Errorf("either --bbox, --address, or --lat/--lon is required")
	}
	if c.Radius == 0 {
		return fmt.Errorf("--radius is required with --address or --latitude / --longitude")
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////

func main() {
	// Set up configuration
	var config ChargeMeUpConfig
	var address string
	var showHelp bool

	pflag.StringVarP(&config.BoundingBox, "bbox", "b", "", "bounding box for query, \"(lat1,lon1),(lat2,lon2)\"")
	pflag.StringVarP(&address, "address", "a", "", "address to query (requires --radius)")
	pflag.Float32VarP(&config.Radius, "radius", "r", 0, "radial distance to query, in kilometers (requires --address)")
	pflag.Float32VarP(&config.Latitude, "lat", "", 0, "latitude to query (requires --lon and --radius)")
	pflag.Float32VarP(&config.Longitude, "lon", "", 0, "longitude to query (requires --lon and --radius)")

	pflag.StringVarP(&config.Server, "server", "s", defaultServer, "API Server for OpenChargeMap, env var OCM_SERVER")
	pflag.StringVarP(&config.APIKey, "key", "k", "", "API key for OpenChargeMap, env var OCM_KEY)")
	pflag.BoolVarP(&config.Verbose, "verbose", "v", false, "verbose output to stderr")
	pflag.BoolVarP(&showHelp, "help", "h", false, "show help")
	pflag.Parse()

	if showHelp {
		fmt.Fprintf(os.Stdout, usageFormat, os.Args[0])
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if address != "" {
		n := nominatim.Nominatim{}
		results, err := n.Search(nominatim.SearchParameters{ // Check SearchResult struct for details
			Query:          address,
			IncludeGeoJSON: true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error finding address '%s': %s\n", address, err.Error())
			os.Exit(1)
		}
		if len(results) == 0 {
			fmt.Fprintf(os.Stderr, "no location found for address '%s'\n", address)
			os.Exit(1)
		}
		config.Latitude = float32(results[0].Lat)
		config.Longitude = float32(results[0].Lng)
	}

	if err := config.IsAreaValid(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	if config.Server == "" {
		config.Server = os.Getenv("OCM_SERVER")
		if config.Server == "" {
			config.Server = defaultServer
		}
	}

	if config.APIKey == "" {
		config.APIKey = os.Getenv("OCM_KEY")
		if config.APIKey == "" {
			fmt.Fprintf(os.Stderr, "must set OCM_KEY environment variable or pass --key\n")
			os.Exit(1)
		}
	}

	ocmClient, err := makeClient(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create client: %s\n", err.Error())
		os.Exit(1)
	}

	if err := lookupChargePoints(config, ocmClient); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

///////////////////////////////////////////////////////////////////////////////

func makeClient(config ChargeMeUpConfig) (*openchargemap.ClientWithResponses, error) {
	// Create OCM client
	apiKeyProvider, err := securityprovider.NewSecurityProviderApiKey("query", "key", config.APIKey)
	if err != nil {
		return nil, fmt.Errorf("NewSecurityProviderApiKey error: %w", err)
	}

	ocmClient, err := openchargemap.NewClientWithResponses(config.Server, openchargemap.WithRequestEditorFn(apiKeyProvider.Intercept))
	if err != nil {
		return nil, fmt.Errorf("creating openchargemap client: %w", err)
	}

	return ocmClient, nil
}

///////////////////////////////////////////////////////////////////////////////

func lookupChargePoints(config ChargeMeUpConfig, ocmClient *openchargemap.ClientWithResponses) error {
	ctx := context.Background()
	params := &openchargemap.GetPoiParams{}
	if config.BoundingBox != "" {
		params.Boundingbox = &config.BoundingBox
	} else {
		kmUnit := "km"
		params.Distance = &config.Radius
		params.Distanceunit = &kmUnit
		params.Latitude = &config.Latitude
		params.Longitude = &config.Longitude
	}

	resp, err := ocmClient.GetPoiWithResponse(ctx, params)
	if err != nil {
		return fmt.Errorf("openchargemap.GetPoi: %w", err)
	}
	if resp.JSON200 == nil {
		if resp.HTTPResponse != nil {
			return fmt.Errorf("openchargemap.GetPoi: response is nil, status: %d", resp.HTTPResponse.StatusCode)
		}
		return fmt.Errorf("openchargemap.GetPoi: response is nil, no status")
	}

	// Request the POI
	if config.Verbose {
		fmt.Fprintf(os.Stderr, "found %d responses\n", len(*resp.JSON200))
	}
	jstr, _ := json.Marshal(*resp.JSON200)
	fmt.Fprint(os.Stdout, string(jstr))
	return nil
}
