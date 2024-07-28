// memorer.go
package taskwrappr

type MemoryMap struct {
	Parent    *MemoryMap
	Actions   map[string]*Action
	Variables map[string]*Variable
}

func NewMemoryMap(parent *MemoryMap) *MemoryMap {
	return &MemoryMap{
		Parent:    parent,
		Actions:   make(map[string]*Action),
		Variables: make(map[string]*Variable),
	}
}

func (m *MemoryMap) GetAction(name string) *Action {
	for scope := m; scope != nil; scope = scope.Parent {
		if action, ok := scope.Actions[name]; ok {
			return action
		}
	}
	return nil
}

func (m *MemoryMap) GetVariable(name string) *Variable {
	for scope := m; scope != nil; scope = scope.Parent {
		if variable, ok := scope.Variables[name]; ok {
			return variable
		}
	}
	return nil
}

func (m *MemoryMap) MakeVariable(name string, value interface{}) *Variable {
	variable := NewVariable(value, DetermineVariableType(value))
	m.Variables[name] = variable
	return variable
}

func (m *MemoryMap) SetVariable(name string, value interface{}, variableType VariableType) *Variable {
	for scope := m; scope != nil; scope = scope.Parent {
		if variable := scope.GetVariable(name); variable != nil {
			variable.Value = value
			variable.Type = variableType
			return variable
		}
	}

	return m.MakeVariable(name, value)
}

func (m *MemoryMap) DeleteAction(name string) {
    delete(m.Variables, name)
}

func (m *MemoryMap) DeleteVariable(name string) {
    delete(m.Variables, name)
}

func (m *MemoryMap) Clear() {
	m.Actions = make(map[string]*Action)
	m.Variables = make(map[string]*Variable)
}