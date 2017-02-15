package plugin

import (
	"fmt"
	"plugin"
)

// Load opens a plugin, but panics if it cannot be loaded.
func Load(filename string) (plugin.Symbol, error) {
	p, err := plugin.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("load: %s", err)
	}

	sym, err := p.Lookup("New")
	if err != nil {
		return nil, fmt.Errorf("load: %s", err)
	}

	return sym, nil
}

// LoadFilter creates a filter from a plugin. If the plugin does not
// exist or implement a New function that returns a Filter, then an
// error will be returned.
func LoadFilter(filename string, opts Options) (Filter, error) {
	v, err := Load(filename)
	if err != nil {
		return nil, fmt.Errorf("load filter: %s", err)
	}

	// Ensure that the loaded symbol matches our interface.
	fn, ok := v.(func(Options) (Filter, error))
	if !ok {
		return nil, fmt.Errorf("plugin cannot create filters: %v", v)
	}

	return fn(opts)
}
