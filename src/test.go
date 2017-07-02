package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// InventoryData is the data type for inventory data produced by an integration
// data source and emitted to the agent's inventory data store
type InventoryData map[string]interface{}

// MetricData is the data type for events produced by an integration data source
// and emitted to the agent's metrics data store
type MetricData map[string]interface{}

// EventData is the data type for single shot events
type EventData map[string]interface{}

// IntegrationData defines the format of the output JSON integrations will return
type IntegrationData struct {
	Name               string                   `json:"name"`
	ProtocolVersion    string                   `json:"protocol_version"`
	IntegrationVersion string                   `json:"integration_version"`
	Metrics            []MetricData             `json:"metrics"`
	Inventory          map[string]InventoryData `json:"inventory"`
	Events             []EventData              `json:"events"`
}

// OutputJSON takes an object and prints it as a JSON string to the stdout.
// If the pretty attribute is set to true, the JSON will be indented for easy reading.
func OutputJSON(data interface{}, pretty bool) error {
	var output []byte
	var err error

	if pretty {
		output, err = json.MarshalIndent(data, "", "\t")
	} else {
		output, err = json.Marshal(data)
	}

	if err != nil {
		return fmt.Errorf("Error outputting JSON: %s", err)
	}

	if string(output) == "null" {
		fmt.Println("[]")
	} else {
		fmt.Println(string(output))
	}

	return nil
}

func main() {
	// Setup the integration's command line parameters
	verbose := flag.Bool("v", false, "Print more information to logs")
	pretty := flag.Bool("p", false, "Print pretty formatted JSON")
	pathtosearch := flag.String("s", "./", "Path to calculate, default './'")
	flag.Parse()

	// Setup logging, redirect logs to stderr and configure the log level.
	log.SetOutput(os.Stderr)
	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Initialize the output structure
	var data = IntegrationData{
		Name:               "test",
		ProtocolVersion:    "1",
		IntegrationVersion: "1.0.0",
		Inventory:          make(map[string]InventoryData),
		Metrics:            make([]MetricData, 0),
		Events:             make([]EventData, 0),
	}

	// Build the metrics dictionary with a valid event_type:
	//  * LoadBalancerSample
	//  * BlockDeviceSample
	//  * DatastoreSample
	//  * QueueSample
	//  * ComputeSample
	//  * IamAccountSummarySample
	//  * PrivateNetworkSample
	//  * ServerlessSample
	// Provider may be set no anything identifying the data provider
	var metric = map[string]interface{}{
		"event_type": "DatastoreSample",
		"provider":   "JoanmiTestExtension",
	}

	// Get ENVIRONMENT variable set by the agent
	env := os.Getenv("ENVIRONMENT")
	if env != "" {
		metric["environment"] = env
	}

	// Each metric specific to a provider should go prefixed with the provider namespace.
	keyList := []string{"provider.folderSize"} // , "provider.valueTwo", "provider.valueThree"}
	for _, key := range keyList {
		metric[key] = 0 //rand.Int()
		log.Debugf("Adding metric %s with value %d", key, metric[key])
	}

	data.Metrics = append(data.Metrics, metric)
	//error err1
	keyList = []string{"folderSize"} //, "valueTwo", "valueThree"}
	itemKeys := []string{"item1"}    //, "item2", "item3"}
	for _, item := range itemKeys {
		data.Inventory[item] = InventoryData{}
		for _, key := range keyList {
			data.Inventory[item][key] = dirsize(*pathtosearch) // checksize(pathtosearch) // rand.Int()
			//log.Debugf("Set inventory key %s=%d for %s", key, data.Inventory[item][key], item)
		}
	}

	fatalIfErr(OutputJSON(data, *pretty))
}

func fatalIfErr(err error) {
	if err != nil {
		log.WithError(err).Fatal("can't continue")
	}
}

func dirsize(path string) int64 { //  (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
			log.Debugf(".")
		}
		return err
	})
	log.Debugf("error %s", err)
	log.Debugf("Path to calculate %s, size: %d", path, size)
	return size //, err
}
