package dualstore

import smalldata "github.com/ThanCoder/dual_store_go/small_data"

type DualModel any

// BigDataType - Don't use JSON for GB size data
type BigDataType int

const (
	// BigDataTypeNone - No Data
	BigDataTypeNone BigDataType = iota

	// BigDataTypeStringText - For small to medium text (up to few MBs).
	// Avoid using for GB size data to prevent OOM (Out of Memory) errors.
	BigDataTypeStringText

	// BigDataTypeFile - For Large Data (Videos, Databases, 1GB+ files).
	// Handled via Stream, so it's memory-safe for any size.
	BigDataTypeFile
)

type DualAdapter[T DualModel] struct {
	IDFunc        func(value T) int64
	ParentIdFunc  func(value T) int64
	ToSmallFunc   func(value T, generatedId int64, encoder *smalldata.SmallDataEncoder) []byte
	FromSmallFunc func(decoder smalldata.SmallDataDecoder) T
	// default func
	AdapterTypeId int32
	BigDataType   BigDataType
}

// AdapterTypeIdFunc: -1,
//
// ParentIdFunc:      -1,
//
// BigDataTypeFunc:   BigDataTypeNone,
func NewAdapterDefault[T DualModel](idFunc func(value T) int64,
	toSmallFunc func(value T, generatedId int64, encoder *smalldata.SmallDataEncoder) []byte,
	fromSmallFunc func(decoder smalldata.SmallDataDecoder) T) *DualAdapter[T] {
	return NewAdapter(idFunc, func(value T) int64 { return -1 }, -1, toSmallFunc, fromSmallFunc, -1, BigDataTypeNone)
}

func NewAdapter[T DualModel](idFunc func(value T) int64, parentIdFunc func(value T) int64, adpaterTypeId int32,
	toSmallFunc func(value T, generatedId int64, encoder *smalldata.SmallDataEncoder) []byte,
	fromSmallFunc func(decoder smalldata.SmallDataDecoder) T, adapterTypeId int32, bigDataType BigDataType) *DualAdapter[T] {
	return &DualAdapter[T]{
		IDFunc:        idFunc,
		ToSmallFunc:   toSmallFunc,
		FromSmallFunc: fromSmallFunc,
		ParentIdFunc:  parentIdFunc,
		AdapterTypeId: adpaterTypeId,
		BigDataType:   bigDataType,
	}
}
