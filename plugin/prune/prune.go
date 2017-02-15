package main

import "github.com/danielwhite/dumbledore/plugin"

func main() {}

// New creates a new filter configured with the provided options.
func New(opts plugin.Options) (plugin.Filter, error) {
	var f filter
	if err := opts.Decode(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

type filter struct {
	BlacklistNames []string `mapstructure:"blacklist_names"`
}

func (f *filter) Filter(v plugin.Event, out chan<- plugin.Event) {
	for _, name := range f.BlacklistNames {
		delete(v, name)
	}
	out <- v
}
