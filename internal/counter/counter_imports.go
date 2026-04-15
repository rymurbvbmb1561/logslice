package counter

import (
	"fmt"
	"sort"
)

// This file exists solely to satisfy imports used across counter.go.
// Go requires all imports to be declared in the same file or a sibling file
// within the same package. Keeping them here avoids cluttering the main file.

var (
	_ = fmt.Sprint
	_ = sort.Slice
)
