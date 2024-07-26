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

func (m *MemoryMap) GetAction(name string) *Action {
	action, ok := m.Actions[name]
	if !ok {
		return nil
	}
	return action
}

func (m *MemoryMap) GetVariable(name string) *Variable {
	variable, ok := m.Variables[name]
	if !ok {
		return nil
	}
	return variable
}

func (m *MemoryMap) MakeVariable(name string, value interface{}) *Variable {
	variable := NewVariable(value, DetermineVariableType(value))
	m.Variables[name] = variable

	return variable
}

func (m *MemoryMap) SetVariable(name string, value interface{}, variableType VariableType) *Variable {
	variable := m.GetVariable(name)
	if variable == nil {
		return m.MakeVariable(name, value)
	}

	variable.Value = value
	variable.Type = variableType

	return variable
}