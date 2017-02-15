package main

import "github.com/danielwhite/dumbledore/plugin"

type cloneFilter struct {
	// A new clone will be created with the given type for each
	// type in this list.
	Clones []string
}

func (f *cloneFilter) Filter(event plugin.Event, out chan<- plugin.Event) {
	for _, name := range f.Clones {
		// KLUDGE: Shallow copy; real thing should do a deep copy.
		clone := plugin.Event{}
		for k, v := range event {
			clone[k] = v
		}

		// Set the type to the name of the clone.
		clone["type"] = name

		// Forward the clone.
		out <- clone
	}

	// Forward the original event after so we don't need to worry
	// about concurrent access.
	out <- event
}
