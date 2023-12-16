package sync

import "sync"

type Flag struct {
	sync.Mutex
	value bool
}

// Set sets the boolean flag to true atomically.
func (f *Flag) Set() {
	f.Lock()
	defer f.Unlock()
	f.value = true
}

// Unset sets the boolean flag to false atomically.
func (f *Flag) Unset() {
	f.Lock()
	defer f.Unlock()
	f.value = false
}

// IsSet returns the current value of the boolean flag.
func (f *Flag) IsSet() bool {
	f.Lock()
	defer f.Unlock()
	return f.value
}
