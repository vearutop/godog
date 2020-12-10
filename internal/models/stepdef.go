package models

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/cucumber/messages-go/v10"

	"github.com/cucumber/godog/formatters"
)

var typeOfBytes = reflect.TypeOf([]byte(nil))

// matchable errors
var (
	ErrUnmatchedStepArgumentNumber = errors.New("func received more arguments than expected")
	ErrCannotConvert               = errors.New("cannot convert argument")
	ErrUnsupportedArgumentType     = errors.New("unsupported argument type")
)

// StepDefinition ...
type StepDefinition struct {
	formatters.StepDefinition

	Args         []interface{}
	HandlerValue reflect.Value

	// multistep related
	Nested    bool
	Undefined []string
}

// Run a step with the matched arguments using reflect
func (sd *StepDefinition) Run() interface{} {
	typ := sd.HandlerValue.Type()
	if len(sd.Args) < typ.NumIn() {
		return fmt.Errorf("%w: expected %d arguments, matched %d from step", ErrUnmatchedStepArgumentNumber, typ.NumIn(), len(sd.Args))
	}

	var values []reflect.Value
	for i := 0; i < typ.NumIn(); i++ {
		param := typ.In(i)
		switch param.Kind() {
		case reflect.Int:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return err
			}
			v, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				return fmt.Errorf(`%w %d: "%s" to int: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(int(v)))
		case reflect.Int64:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return err
			}
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return fmt.Errorf(`%w %d: "%s" to int64: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(int64(v)))
		case reflect.Int32:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return err
			}
			v, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return fmt.Errorf(`%w %d: "%s" to int32: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(int32(v)))
		case reflect.Int16:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return err
			}
			v, err := strconv.ParseInt(s, 10, 16)
			if err != nil {
				return fmt.Errorf(`%w %d: "%s" to int16: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(int16(v)))
		case reflect.Int8:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return err
			}
			v, err := strconv.ParseInt(s, 10, 8)
			if err != nil {
				return fmt.Errorf(`%w %d: "%s" to int8: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(int8(v)))
		case reflect.String:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return err
			}
			values = append(values, reflect.ValueOf(s))
		case reflect.Float64:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return err
			}
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return fmt.Errorf(`%w %d: "%s" to float64: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(v))
		case reflect.Float32:
			s, err := sd.shouldBeString(i)
			if err != nil {
				return err
			}
			v, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return fmt.Errorf(`%w %d: "%s" to float32: %s`, ErrCannotConvert, i, s, err)
			}
			values = append(values, reflect.ValueOf(float32(v)))
		case reflect.Ptr:
			arg := sd.Args[i]
			switch param.Elem().String() {
			case "messages.PickleStepArgument_PickleDocString":
				if v, ok := arg.(*messages.PickleStepArgument); ok {
					values = append(values, reflect.ValueOf(v.GetDocString()))
					break
				}

				if v, ok := arg.(*messages.PickleStepArgument_PickleDocString); ok {
					values = append(values, reflect.ValueOf(v))
					break
				}

				return fmt.Errorf(`%w %d: "%v" of type "%T" to *messages.PickleStepArgument_PickleDocString`, ErrCannotConvert, i, arg, arg)
			case "messages.PickleStepArgument_PickleTable":
				if v, ok := arg.(*messages.PickleStepArgument); ok {
					values = append(values, reflect.ValueOf(v.GetDataTable()))
					break
				}

				if v, ok := arg.(*messages.PickleStepArgument_PickleTable); ok {
					values = append(values, reflect.ValueOf(v))
					break
				}

				return fmt.Errorf(`%w %d: "%v" of type "%T" to *messages.PickleStepArgument_PickleTable`, ErrCannotConvert, i, arg, arg)
			default:
				return fmt.Errorf("%w: the argument %d type %T is not supported %s", ErrUnsupportedArgumentType, i, arg, param.Elem().String())
			}
		case reflect.Slice:
			switch param {
			case typeOfBytes:
				s, err := sd.shouldBeString(i)
				if err != nil {
					return err
				}
				values = append(values, reflect.ValueOf([]byte(s)))
			default:
				return fmt.Errorf("%w: the slice argument %d type %s is not supported", ErrUnsupportedArgumentType, i, param.Kind())
			}
		default:
			return fmt.Errorf("%w: the argument %d type %s is not supported", ErrUnsupportedArgumentType, i, param.Kind())
		}
	}

	res := sd.HandlerValue.Call(values)
	if len(res) == 0 {
		return nil
	}

	return res[0].Interface()
}

func (sd *StepDefinition) shouldBeString(idx int) (string, error) {
	arg := sd.Args[idx]
	s, ok := arg.(string)
	if !ok {
		return "", fmt.Errorf(`%w %d: "%v" of type "%T" to string`, ErrCannotConvert, idx, arg, arg)
	}

	return s, nil
}

// GetInternalStepDefinition ...
func (sd *StepDefinition) GetInternalStepDefinition() *formatters.StepDefinition {
	if sd == nil {
		return nil
	}

	return &sd.StepDefinition
}
