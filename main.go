package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

// Supported Viper config file types, by opinionated order of popularity
var viperTypes = []string{
	"json",
	"yaml",
	"hcl",
	"toml",
	"properties",
}

// GetViperConfigType returns a guess of the file extension. Viper has no way of
// guessing the file type without trying to load the data from the config: we're
// loading in a stream of bytes rather than a file with an extension. If we
// can't find the config type it returns an empty string.
//
// Dumping viper.AllSettings() and attempting to marshal it to JSON does give us
// the ability to check if Viper found any properties, so we try every
// acceptable Viper type to see if we can grab any data successfully. If we do,
// we've (probably) found the right config type.
//
// This method isn't guaranteed to guess the correct file type if the config
// file is empty, but if the config file is empty we don't have any properties,
// which would yield an empty JSON object anyway.
func GetConfigType(b *[]byte) string {
	// This is the []bytes result of (1) loading a config file as the incorrect
	// type or (2) loading an empty config file as its correct type.
	emptyViperConfig := []byte("{}")
	testConfig := viper.New()
	for _, j := range viperTypes {
		testConfig.SetConfigType(j)
		testConfig.ReadConfig(bytes.NewReader(*b))
		json, _ := json.Marshal(testConfig.AllSettings())
		if !bytes.Equal(json, emptyViperConfig) {
			return j
		}
	}
	// Config file type not found
	return ""
}

func GetViperConfigFromBytes(b []byte) *viper.Viper {
	language := GetConfigType(&b)
	conf := viper.New()
	conf.SetConfigType(language)
	conf.ReadConfig(bytes.NewReader(b))
	return conf
}

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic("ioutil: could not read stdin")
	}

	// Copy the bytes read from stdin and store them since buffers are consumed
	// when read.
	buf := bytes.NewBuffer(stdin)
	data := buf.Bytes()

	config := GetViperConfigFromBytes(data)
	bytes, err := json.Marshal(config.AllSettings())
	if err != nil {
		panic("json: error marshaling JSON")
	}
	fmt.Printf("%s\n", string(bytes))
}
