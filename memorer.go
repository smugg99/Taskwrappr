// memorer.go
package taskwrappr

type MemoryMap struct {
	Actions   map[string]*Action
	Variables map[string]*Variable
}

func NewMemoryMap() *MemoryMap {
	actions, variables := GetInternals()

	return &MemoryMap{
		Actions:   actions,
		Variables: variables,
	}
}
