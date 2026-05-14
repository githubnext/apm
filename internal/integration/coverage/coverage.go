// Package coverage provides primitive dispatch coverage validation.
package coverage

import "fmt"

// DispatchEntry holds integrator method names for a primitive.
type DispatchEntry struct {
Targets []string
Methods []string
}

// CheckPrimitiveCoverage validates that every primitive has a handler and vice versa.
func CheckPrimitiveCoverage(knownPrimitives []string, dispatchTable map[string]DispatchEntry, specialCases map[string]bool) error {
handled := map[string]bool{}
for k := range dispatchTable {
handled[k] = true
}
for k := range specialCases {
handled[k] = true
}

for _, p := range knownPrimitives {
if !handled[p] {
return fmt.Errorf("primitive %q is registered but has no integrator in dispatch table", p)
}
}

primSet := map[string]bool{}
for _, p := range knownPrimitives {
primSet[p] = true
}
for k := range dispatchTable {
if !primSet[k] && !specialCases[k] {
return fmt.Errorf("dispatch table entry %q has no corresponding primitive in KNOWN_TARGETS", k)
}
}
return nil
}
