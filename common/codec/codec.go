package codec

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"reflect"

	"github.com/c3systems/go-substrate/common/u8compact"
)

func writeBinary(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, ErrNilKind
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func encodeStruct(v *reflect.Value) ([]byte, error) {
	if v == nil {
		return nil, ErrNilKind
	}

	if v.Kind() != reflect.Struct {
		return encode(v)
	}

	// TODO: use go routine?
	var (
		ret, tmpBytes []byte
		val           reflect.Value
		err           error
	)
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Struct:
			{
				val = v.Field(i)
				tmpBytes, err = encodeStruct(&val)
				if err != nil {
					return nil, err
				}

			}
		default:
			{
				val = v.Field(i)
				tmpBytes, err = encode(&val)
				if err != nil {
					return nil, err
				}
			}
		}

		ret = append(ret, tmpBytes...)
	}

	return ret, nil
}

func encode(v *reflect.Value) ([]byte, error) {
	if v == nil {
		return nil, ErrNilKind
	}

	var (
		ret []byte
		err error
	)
	switch v.Kind() {
	case reflect.String:
		{
			ret, err = writeBinary([]byte(v.Interface().(string)))
			if err != nil {
				return nil, err
			}
			leader := u8compact.CompactToUint8Slice(big.NewInt(int64(len(ret))), int(v.Type().Size()))
			ret = append(leader, ret...)
		}
	case reflect.Int:
		{
			return writeBinary(int32(v.Interface().(int)))
		}
	case reflect.Uint:
		{
			return writeBinary(uint32(v.Interface().(uint)))
		}
	case reflect.Ptr, reflect.UnsafePointer, reflect.Uintptr:
		{
			if v.IsNil() {
				// TODO: is this correct?
				break
			}

			// note: already checked for v == nil, so should not panic
			return Encode(reflect.Indirect(*v))
		}
	case reflect.Struct:
		{
			return encodeStruct(v)
		}
	case reflect.Array, reflect.Slice:
		{
			s := reflect.ValueOf(v)

			// TODO: is this correct? What about len?
			var tmpBytes []byte
			for i := 0; i < s.Len(); i++ {
				tmpBytes, err = Encode(s.Index(i))
				if err != nil {
					return nil, err
				}
				ret = append(ret, tmpBytes...)
			}
		}
	case reflect.Invalid, reflect.Chan, reflect.Func:
		{
			// note: also Complex64, Complex128, Interface, Map
			return nil, ErrInvalidKind
		}
	default:
		{
			return writeBinary(v.Interface())
		}
	}

	return ret, nil
}

// Encode ...
func Encode(input interface{}) ([]byte, error) {
	if input == nil {
		return nil, ErrNilInput
	}

	v := reflect.ValueOf(input)
	return encode(&v)
}

// Decode ...
// TODO: this is an incomplete implementation
func Decode(input []byte, target interface{}) error {
	if input == nil {
		return ErrNilInput
	}
	if target == nil {
		return ErrNilTarget
	}

	switch v := reflect.ValueOf(target); v.Kind() {
	case reflect.Ptr, reflect.UnsafePointer, reflect.Uintptr:
		{
			return binary.Read(bytes.NewReader(input), binary.LittleEndian, target)
		}

	default:
		{
			return ErrNonTargetPointer
		}
	}
}