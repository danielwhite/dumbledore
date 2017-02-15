package main

// Event describes data passing through a pipeline.
type Event map[string]interface{}

// Filter modifies an event in the pipeline.
type Filter interface {
	Filter(v Event, out chan<- Event)
}

func startFilter(filter Filter, in <-chan Event) <-chan Event {
	out := make(chan Event)
	go func() {
		for event := range in {
			filter.Filter(event, out)
		}
		close(out)
	}()
	return out
}

type pruneFilter struct {
	BlacklistNames []string
}

func (f *pruneFilter) Filter(v Event, out chan<- Event) {
	for _, name := range f.BlacklistNames {
		delete(v, name)
	}
	out <- v
}

type cloneFilter struct {
	// A new clone will be created with the given type for each
	// type in this list.
	Clones []string
}

func (f *cloneFilter) Filter(event Event, out chan<- Event) {
	for _, name := range f.Clones {
		// KLUDGE: Shallow copy; real thing should do a deep copy.
		clone := Event{}
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
