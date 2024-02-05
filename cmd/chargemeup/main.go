package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	openchargemap "github.com/neomantra/go-openchargemap"
	"github.com/spf13/pflag"
)

/////////////////////////////////////////////////////////////////////////////////////

var usageFormat string = `usage:  %s <options> [input]

"chargemeup" assists with queries to OpenChargeMap.

Around Newark is:

chargemeup -p "(40.63010790372053,-74.2775717248681),(40.7356464076158,-74.09370618215354)"

`

const defaultServer = "https://api.openchargemap.io/v3"

type ChargeMeUpConfig struct {
	Server  string
	APIKey  string
	Verbose bool
}

/////////////////////////////////////////////////////////////////////////////////////

func main() {
	// Set up configuration
	var config ChargeMeUpConfig
	var boundingBox string
	var showHelp bool

	pflag.StringVarP(&boundingBox, "bbox", "b", "", "bounding box for query, \"(lat1,lon1),(lat2,lon2)\"")
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

	if err := lookupChargePoints(config, ocmClient, boundingBox); err != nil {
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

func lookupChargePoints(config ChargeMeUpConfig, ocmClient *openchargemap.ClientWithResponses, place string) error {

	ctx := context.Background()
	params := &openchargemap.GetPoiParams{
		Boundingbox: &place,
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
