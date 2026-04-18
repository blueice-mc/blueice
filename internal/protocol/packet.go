package protocol

import (
	"fmt"
	"io"
	"reflect"
)

type Packet interface {
	ID() VarInt
}

type ServerboundPacket interface {
	Packet
	ReadFrom(r io.Reader) (int64, error)
}

type ClientboundPacket interface {
	Packet
	WriteTo(w io.Writer) (int64, error)
}

func writeField(w io.Writer, v reflect.Value) (int64, error) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return 0, nil // skip nil fields
		}

		v = v.Elem()
	}

	// write logic not done yet
}

func WritePacket(w io.Writer, packet *ClientboundPacket) (int64, error) {
	v := reflect.ValueOf(packet)

	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return 0, fmt.Errorf("expected struct, got %v", v.Kind())
	}

	queue := make([]reflect.Value, 0, 10)
	queue = append(queue, v) // append the current struct to the queue

	total := int64(0)

	for len(queue) > 0 { // as long as the queue is not empty
		top := queue[0]   // take the first element from the queue
		queue = queue[1:] // remove the first element from the queue

		if top.Kind() == reflect.Struct { // if the element is a struct
			writeMethod := top.MethodByName("WriteTo")
			if writeMethod.IsValid() { // and has a valid WriteTo method
				r := writeMethod.Call([]reflect.Value{reflect.ValueOf(w)}) // call the WriteTo method
				total += r[0].Int()
				if r[1].Interface() != nil {
					return total, r[1].Interface().(error)
				}
			} else { // otherwise, add all fields of the struct to the queue
				for i := 0; i < top.NumField(); i++ {
					queue = append(queue, top.Field(i))
				}
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
