package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/danielwhite/dumbledore/plugin"
	"github.com/hashicorp/hcl"
	"github.com/mitchellh/mapstructure"
)

type config struct {
	Filters   map[string]plugin.Options `mapstructure:"filter"`
	PluginDir string
}

func (cfg *config) loadFilters() []plugin.Filter {
	var filters []plugin.Filter
	for name, opts := range cfg.Filters {
		// Ensure the plugin can be loaded from our path.
		filename := filepath.Join(cfg.PluginDir, name+".so")

		// Attempt to load the filter; for now, this is just
		// simply skips unloadable filters.
		filter, err := plugin.LoadFilter(filename, opts)
		if err != nil {
			log.Printf("failed to load filter; skipping: %s", err)
			continue
		}
		filters = append(filters, filter)
	}
	return filters
}

// readConfig parses a HCL configuration file containing a set of
// filter configuration, and returns them in the same order.
func readConfig(filename string) *config {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Decode the HCL into maps for further unpacking.
	var parsed map[string][]map[string]interface{}
	if err := hcl.Unmarshal(b, &parsed); err != nil {
		panic(err)
	}

	// Merge all of the slices of maps into single sets of
	// configuration per filter plugin by weakly decoding the
	// parsed map.
	var cfg config
	if err := mapstructure.WeakDecode(parsed, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
