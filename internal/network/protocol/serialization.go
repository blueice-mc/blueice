package protocol

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

// Deserialize reads from the given reader and writes the data into a pointer of the given value.
// The value to read into MUST be a pointer, otherwise the deserialize function will not be able to
// write any content into it.
// If any value read by this function implements the ReaderFrom interface, the ReadFrom for that type
// will be used. Otherwise, this function will try to read every field from the struct in the
// order the struct was implemented.
func Deserialize(r io.Reader, v any) (int64, error) {
	return deserialize(r, reflect.ValueOf(v))
}

// Serialize writes the given value into the given writer.
// If any value written by this function implements the WriterTo interface, the WriteTo for that type
// will be used. Otherwise, this function will try to write every field from the struct in the
// order the struct was implemented.
func Serialize(w io.Writer, v any) (int64, error) {
	return serialize(w, reflect.ValueOf(v))
}

// validPtr is a helper function to get the address of a non-pointer type
func validPtr(v reflect.Value) (reflect.Value, error) {

	// checking if the field is already a pointer
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Interface {
		if v.CanAddr() {
			// if it is not and can be addressed, use the pointer
			return v.Addr(), nil
		}

		// if it is not and cannot be addressed, return an error
		return v, errors.New("cannot deserialize into non-addressable field")
	}

	return v, nil
}

// deserialize is a helper function that uses recursion to deserialize the given value from the given reader
func deserialize(r io.Reader, v reflect.Value) (int64, error) {
	// pointer is required
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Interface {
		return 0, fmt.Errorf("expected pointer, got %v", v.Kind())
	}

	// unwrap until we reach pointer to concrete type
	for v.Elem().Kind() == reflect.Ptr || v.Elem().Kind() == reflect.Interface {
		// for nil pointer, create a new instance
		v = v.Elem()
	}

	// initialize if v is nil pointer
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	// check if the value implements the ReaderFrom interface
	if rf, ok := v.Interface().(io.ReaderFrom); ok {
		return rf.ReadFrom(r)
	}

	if v.Elem().Kind() == reflect.Struct {
		total := int64(0)
		// iterate through all fields of struct
		for i := 0; i < v.Elem().NumField(); i++ {
			f := v.Elem().Field(i)

			// a pointer is required
			f, err := validPtr(f)
			if err != nil {
				return total, err
			}

			// deserialize the current field of the struct
			n, err := deserialize(r, f)
			total += n
			if err != nil {
				return total, err
			}
		}
		return total, nil
	}

	if v.Elem().Kind() == reflect.Slice || v.Elem().Kind() == reflect.Array {
		total := int64(0)
		// iterate through all elements of slice
		for i := 0; i < v.Elem().Len(); i++ {
			f := v.Elem().Index(i)

			// a pointer to the field of the slice/array is required
			f, err := validPtr(f)
			if err != nil {
				return total, err
			}

			// deserialize the current element of the slice/array
			n, err := deserialize(r, f)
			total += n
			if err != nil {
				return total, err
			}
		}
		return total, nil
	}

	// for all other types, read the field from the reader
	return readField(r, v)
}

// readField is a helper function that reads a field from the given reader and writes it into the given value
// it requires a pointer to a concrete type
func readField(r io.Reader, v reflect.Value) (int64, error) {
	switch v.Elem().Kind() {
	// ignore functions
	case reflect.Func:
		return 0, nil
	case reflect.Bool:
		return ReadBool(r, v.Interface().(*bool))
	case reflect.Uint8:
		return ReadUint8(r, v.Interface().(*uint8))
	case reflect.Int8:
		return ReadInt8(r, v.Interface().(*int8))
	case reflect.Int16:
		return ReadInt16(r, v.Interface().(*int16))
	case reflect.Int32:
		return ReadInt32(r, v.Interface().(*int32))
	case reflect.Int64:
		return ReadInt64(r, v.Interface().(*int64))
	case reflect.Float32:
		return ReadFloat32(r, v.Interface().(*float32))
	case reflect.Float64:
		return ReadFloat64(r, v.Interface().(*float64))
	case reflect.String:
		return ReadString(r, v.Interface().(*string))
	default:
		return 0, fmt.Errorf("unknown field type: %v", v.Type())
	}
}

// serialize is a helper function that uses recursion to serialize the given value into the given writer
func serialize(w io.Writer, v reflect.Value) (int64, error) {
	// unwrap until we reach concrete type
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		// skip this field if it is a nil pointer
		if v.IsNil() {
			return 0, nil
		}
		// dereference the pointer
		v = v.Elem()
	}

	ptr := reflect.New(v.Type())
	ptr.Elem().Set(v)

	if wt, ok := ptr.Interface().(io.WriterTo); ok {
		return wt.WriteTo(w)
	}

	if v.Kind() == reflect.Struct {
		total := int64(0)
		// iterate through all fields of struct
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			// serialize the current field of the struct
			n, err := serialize(w, f)
			total += n
			if err != nil {
				return total, err
			}
		}
		return total, nil
	}

	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		total := int64(0)
		// iterate through all elements of slice/array
		for i := 0; i < v.Len(); i++ {
			f := v.Index(i)
			// serialize the current element of the slice/array
			n, err := serialize(w, f)
			total += n
			if err != nil {
				return total, err
			}
		}
		return total, nil
	}

	return writeField(w, v)
}

func writeField(w io.Writer, v reflect.Value) (int64, error) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return 0, nil // skip nil fields
		}

		v = v.Elem()
	}

	switch v.Kind() {
	// ignore functions
	case reflect.Func:
		return 0, nil
	case reflect.Bool:
		return WriteBool(w, v.Bool())
	case reflect.Uint8:
		return WriteUint8(w, uint8(v.Uint()))
	case reflect.Int8:
		return WriteInt8(w, int8(v.Int()))
	case reflect.Int16:
		return WriteInt16(w, int16(v.Int()))
	case reflect.Int32:
		return WriteInt32(w, int32(v.Int()))
	case reflect.Int64:
		return WriteInt64(w, v.Int())
	case reflect.Float32:
		return WriteFloat32(w, float32(v.Float()))
	case reflect.Float64:
		return WriteFloat64(w, v.Float())
	case reflect.String:
		return WriteString(w, v.String())
	default:
		return 0, fmt.Errorf("unknown field type: %v", v.Type())
	}
}
