package plugin

import "github.com/mitchellh/mapstructure"

// Event describes data passing through a pipeline.
type Event map[string]interface{}

// Options describes raw configuration that can be decoded into a more
// specialised structure.
type Options map[string]interface{}

// Decode options into the given Go native structure.
func (opts Options) Decode(v interface{}) error {
	return mapstructure.Decode(opts, v)
}

// Filter modifies an event in the pipeline.
type Filter interface {
	Filter(v Event, out chan<- Event)
}
