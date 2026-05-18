package dualstore

import smalldata "github.com/ThanCoder/dual_store_go/small_data"

type DualModel any

// BigDataType enum အစား Go သတ်မှတ်ချက်
type BigDataType int

const (
	BigDataTypeNone BigDataType = iota
	BigDataTypeBytes
	BigDataTypeFile
)

type DualAdapter[T DualModel] interface {
	// to small data
	ToSmallData(value T, genratedId int64, encoder *smalldata.SmallDataEncoder) []byte

	FromSmallData(decoder *smalldata.SmallDataDecoder) T

	// Adapter သတ်မှတ်ချက်များ
	AdapterTypeId() int32
	GetParentId(value T) int64
	GetId(value T) int64
	BigDataType() BigDataType
}

// func (ds *DualStoreIndexed) getAllSmallData() {

// }
