// data.go
package taskwrappr

type VariableType int

const (
	VariableUndefined VariableType = iota
	VariableAction
	VariableBool
	VariableNumber
	VariableString
	VariableArray
	VariableObject
)

func (t VariableType) String() string {
	switch t {
	case VariableUndefined:
		return "undefined"
	case VariableAction:
		return "action"
	case VariableBool:
		return "bool"
	case VariableNumber:
		return "number"
	case VariableString:
		return "string"
	case VariableArray:
		return "array"
	case VariableObject:
		return "object"
	default:
		return "unknown"
	}
}

type SelectorType uint

const (
	IndexSelector SelectorType = iota
	KeySelector
	FieldSelector
)

// Represents an index for array, e.g. arr[0]
type IndexSelectorValue int

/*
Represents a key for objects, e.g. obj["key"]
(can be used interchangeably with FieldSelectorValue unless key has a space, e.g. obj["key with space"])
*/
type KeySelectorValue string

/*
Represents a field for objects, e.g. obj.field
(can be used interchangeably with KeySelectorValue unless field has a space, e.g. obj["field with space"])
*/
type FieldSelectorValue string

type Selector struct {
    Type  SelectorType
    Value interface{} // can be int (index), string (key/field), etc.
}

type Variable struct {
	BaseName  string
	Type      VariableType
	Value     interface{}
	Selectors []Selector
}