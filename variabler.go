// structurer.go
package taskwrappr

import (
	"fmt"
	"reflect"
	"strconv"
)

type Variable struct {
	Value interface{}
    Type  VariableType
}

type VariableType int

const (
	StringType VariableType = iota
	IntegerType
	FloatType
	BooleanType
	ArrayType
	NilType
	InvalidType
)

func (v VariableType) String() string {
    switch v {
    case StringType:
        return "string"
    case IntegerType:
        return "integer"
    case FloatType:
        return "float"
    case BooleanType:
        return "boolean"
	case ArrayType:
		return "array"
	case NilType:
		return "nil"
    default:
        return "invalid"
    }
}

func NewVariable(value interface{}, variableType VariableType) *Variable {
	return &Variable{
		Value: value,
        Type:  variableType,
	}
}

func DetermineVariableType(v interface{}) VariableType {
	if v == nil {
		return NilType
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		return StringType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return IntegerType
	case reflect.Float32, reflect.Float64:
		return FloatType
	case reflect.Bool:
		return BooleanType
	case reflect.Slice:
		return ArrayType
	default:
		return InvalidType
	}
}

func (v *Variable) CastTo(targetType VariableType) (interface{}, error) {
	switch targetType {
	case StringType:
		return v.toString()
	case IntegerType:
		return v.toInt()
	case FloatType:
		return v.toFloat()
	case BooleanType:
		return v.toBool()
	default:
		return nil, fmt.Errorf("invalid variable type: %v", targetType)
	}
}

func (v *Variable) toString() (string, error) {
	switch v.Type {
	case StringType:
		return v.Value.(string), nil
	case IntegerType:
		return strconv.Itoa(v.Value.(int)), nil
	case FloatType:
		return fmt.Sprintf("%g", v.Value.(float64)), nil
	case BooleanType:
		return strconv.FormatBool(v.Value.(bool)), nil
	default:
		return "", fmt.Errorf("cannot convert %v to string", v.Type)
	}
}

func (v *Variable) toInt() (int, error) {
	switch v.Type {
	case StringType:
		if i, err := strconv.Atoi(v.Value.(string)); err == nil {
			return i, nil
		} else if f, err := strconv.ParseFloat(v.Value.(string), 64); err == nil {
			return int(f), nil
		}
	case IntegerType:
		return v.Value.(int), nil
	case FloatType:
		return int(v.Value.(float64)), nil
	case BooleanType:
		if v.Value.(bool) {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("cannot convert %v to integer", v.Type)
}

func (v *Variable) toFloat() (float64, error) {
    switch v.Type {
    case StringType:
        if f, err := strconv.ParseFloat(v.Value.(string), 64); err == nil {
            return f, nil
        }
    case IntegerType:
        return float64(v.Value.(int)), nil
    case FloatType:
        return v.Value.(float64), nil
    case BooleanType:
        if v.Value.(bool) {
            return 1.0, nil
        }
        return 0.0, nil
    }
    return 0.0, fmt.Errorf("cannot convert %v to float", v.Type)
}

func (v *Variable) toBool() (bool, error) {
	switch v.Type {
	case StringType:
		if b, err := strconv.ParseBool(v.Value.(string)); err == nil {
			return b, nil
		}
	case IntegerType:
	case FloatType:
		var value float64
		if v.Type == IntegerType {
			value = float64(v.Value.(int))
		} else {
			value = v.Value.(float64)
		}
		return value != 0, nil
	case BooleanType:
		if b, ok := v.Value.(bool); ok {
			return b, nil
		}
	}
	return false, fmt.Errorf("cannot convert %v to boolean", v.Type)
}