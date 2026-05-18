package smalldata

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// dart code
//   static const int typeInt = 1;
//   static const int typeDouble = 2;
//   static const int typeBool = 3;
//   static const int typeString = 4;

const (
	TypeInt    = 1
	TypeDouble = 2
	TypeBool   = 3
	TypeString = 4
)

type SmallDataEncoder struct {
	builder bytes.Buffer
}

// NewEncoder - Initializer function
func NewSmallDataEncoder() *SmallDataEncoder {
	return &SmallDataEncoder{}
}

// write key
func (e *SmallDataEncoder) WriteKey(key int) error {
	if key < 0 || key > 255 {
		return errors.New("key must be between 0 and 255")
	}
	e.builder.WriteByte(byte(key))
	return nil
}

// write int
func (e *SmallDataEncoder) WriteInt(key int, value int64) error {
	if err := e.WriteKey(key); err != nil {
		return err
	}
	e.builder.WriteByte(TypeInt)

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))
	e.builder.Write(buf)

	return nil
}

// write double
func (e *SmallDataEncoder) WriteDouble(key int, value float64) error {
	if err := e.WriteKey(key); err != nil {
		return err
	}
	// write type
	e.builder.WriteByte(TypeDouble)

	err := binary.Write(&e.builder, binary.LittleEndian, value)
	return err
}

// write bool
func (e *SmallDataEncoder) WriteBool(key int, value bool) error {
	if err := e.WriteKey(key); err != nil {
		return err
	}

	// write type
	e.builder.WriteByte(TypeBool)

	if value {
		e.builder.WriteByte(1)
	} else {
		e.builder.WriteByte(0)
	}

	return nil
}

// write string
func (e *SmallDataEncoder) WriteString(key int, value string) error {
	if err := e.WriteKey(key); err != nil {
		return err
	}

	// write type
	e.builder.WriteByte(TypeString)

	stringBytes := []byte(value) // UTF-8 string to byte slice
	stringLen := int32(len(stringBytes))

	// Write string length (4 bytes)
	lengthBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBuf, uint32(stringLen))
	e.builder.Write(lengthBuf)

	// write string value
	e.builder.Write(stringBytes)

	return nil

}

func (e *SmallDataEncoder) FinishedBytes() []byte {
	return e.builder.Bytes()
}
