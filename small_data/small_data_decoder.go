package smalldata

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type SmallDataDecoder struct {
	data        []byte
	offset      int
	decodedData map[int]any // Dart ရဲ့ Map<int, dynamic> နေရာမှာ သုံးတာပါ
}

func NewSmallDataDecoder(data []byte) (*SmallDataDecoder, error) {
	d := &SmallDataDecoder{
		data:        data,
		offset:      0,
		decodedData: make(map[int]any),
	}
	err := d.decodeAll()
	return d, err

}

func (d *SmallDataDecoder) decodeAll() error {
	dataLen := len(d.data)

	// Go ရဲ့ standard panic recovery (Dart ရဲ့ try-catch လိုမျိုး သုံးတာပါ)
	defer (func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	})()

	for d.offset < dataLen {
		// key အရင်ဖတ်
		key := int(d.data[d.offset])
		d.offset += 1
		// data type ဖတ်
		dataType := int(d.data[d.offset])
		d.offset += 1

		switch dataType {
		case TypeInt:
			val := int64(binary.LittleEndian.Uint64(d.data[d.offset : d.offset+8]))
			d.offset += 8
			d.decodedData[key] = val
		case TypeBool:
			val := d.data[d.offset] == 1
			d.offset += 1
			d.decodedData[key] = val
		case TypeDouble:
			// float64 ကို byte slice ကနေ တိုက်ရိုက် ဖတ်တာပါ
			var val float64
			buf := bytes.NewReader(d.data[d.offset : d.offset+8])
			d.offset += 8
			binary.Read(buf, binary.LittleEndian, &val)
			d.decodedData[key] = val
		case TypeString:
			strLen := int(binary.LittleEndian.Uint32(d.data[d.offset : d.offset+4]))
			d.offset += 4

			strBytes := d.data[d.offset : d.offset+strLen]
			str := string(strBytes)
			d.offset += strLen
			d.decodedData[key] = str
		}

	}

	return nil
}

func (d *SmallDataDecoder) PrintAllData() {
	// for key, value := range d.decodedData {
	println("------------")
	// fmt.Printf("Key: %d, Value: %v (Type: %T)\n", key, value, value)
	fmt.Printf("Data: %v\n", d.decodedData)
	println("------------")
	// }
	// fmt.Println(d.decodedData)
}

// default -> 0
func (d *SmallDataDecoder) GetIntDefault(key int) int64 {
	return d.GetInt(key, 0)
}
func (d *SmallDataDecoder) GetInt(key int, defVal int64) int64 {
	if val, ok := d.decodedData[key].(int64); ok {
		return val
	}
	return defVal
}

// default -> 0.0
func (d *SmallDataDecoder) GetDoubleDefault(key int) float64 {
	return d.GetDouble(key, 0.0)
}
func (d *SmallDataDecoder) GetDouble(key int, defVal float64) float64 {
	if val, ok := d.decodedData[key].(float64); ok {
		return val
	}
	return defVal
}

// default -> false
func (d *SmallDataDecoder) GetBoolDefault(key int) bool {
	return d.GetBool(key, false)
}
func (d *SmallDataDecoder) GetBool(key int, defVal bool) bool {
	if val, ok := d.decodedData[key].(bool); ok {
		return val
	}
	return defVal
}

// default -> empty
func (d *SmallDataDecoder) GetStringDefault(key int) string {
	return d.GetString(key, "")
}
func (d *SmallDataDecoder) GetString(key int, defVal string) string {
	if val, ok := d.decodedData[key].(string); ok {
		return val
	}
	return defVal
}
