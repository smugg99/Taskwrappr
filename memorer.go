// memorer.go
package taskwrappr

type MemoryMap struct {
	Actions   map[string]*Action
	Variables map[string]*Variable
}

func NewMemoryMap() *MemoryMap {
	return &MemoryMap{
		Actions:   make(map[string]*Action),
		Variables: make(map[string]*Variable),
	}
}
