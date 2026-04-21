package protocol

import (
	"fmt"
	"io"
	"reflect"
)

func serialize(w io.Writer, v any) (int64, error) {
	queue := make([]reflect.Value, 0, 10)
	queue = append(queue, reflect.ValueOf(v)) // append the current struct to the queue

	total := int64(0)

	for len(queue) > 0 { // as long as the queue is not empty
		top := queue[0]   // take the first element from the queue
		queue = queue[1:] // remove the first element from the queue

		nilPointer := false

		for top.Kind() == reflect.Ptr || top.Kind() == reflect.Interface {
			if top.IsNil() {
				nilPointer = true
				break
			}
			top = top.Elem()
		}

		if nilPointer {
			continue
		}

		ptr := reflect.New(top.Type())
		ptr.Elem().Set(top)

		if wt, ok := ptr.Interface().(io.WriterTo); ok {
			n, err := wt.WriteTo(w)
			total += n
			if err != nil {
				return total, err
			}
			continue
		}

		if top.Kind() == reflect.Struct { // if the element is a struct
			for i := 0; i < top.NumField(); i++ {
				queue = append(queue, top.Field(i))
			}
		} else if top.Kind() == reflect.Slice {
			for i := 0; i < top.Len(); i++ {
				queue = append(queue, top.Index(i))
			}
		} else {
			n, err := writeField(w, top) // otherwise, write the field
			total += n
			if err != nil {
				return total, err
			}
		}
	}

	return total, nil
}

func deserialize(r io.Reader, v *any) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

func writeField(w io.Writer, v reflect.Value) (int64, error) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return 0, nil // skip nil fields
		}

		v = v.Elem()
	}

	switch v.Kind() {
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
