package main

import (
	"fmt"

	smalldata "github.com/ThanCoder/dual_store_go/small_data"
)

///home/thancoder/Downloads/Apyar App/apyar.dual.db
///home/thancoder/Downloads/Apyar App/apyar.dual.db-no-gzip

// Header offsets
// const (
// 	HeaderSize       = 33
// 	FlagPos          = 0
// 	IdPos            = 1
// 	AdapterTypIdPos  = 9
// 	ParentIdPos      = 13
// 	SmallDataSizePos = 21
// 	BigDataSizePos   = 25
// )

// // RecordFlag enum အစား Go မှာ iota သုံးပါမယ်
// type RecordFlag byte

// type RecordMeta struct {
// 	Offset             int64
// 	BigDataStartOffset int64
// 	Flag               RecordFlag
// 	ID                 int64
// 	AdapterTypID       int32
// 	ParentID           int64
// 	SmallDataSize      int32
// 	SmallData          []byte
// 	BigDataSize        int64
// 	RecordSize         int64
// }

// type Person struct {
// 	Name string
// 	Age  int
// }

func main() {
	enc := smalldata.NewSmallDataEncoder()
	enc.WriteInt(1, 22)
	enc.WriteDouble(2, 5.6)
	enc.WriteBool(3, true)
	enc.WriteString(4, "Thancoder")

	fmt.Printf("encoder bytes len: %d\n", len(enc.FinishedBytes()))

	decoder, err := smalldata.NewSmallDataDecoder(enc.FinishedBytes())
	if err != nil {
		fmt.Printf("Dec Error: %v\n", err)
		return
	}
	decoder.PrintAllData()

	// fmt.Println("Int Key 1:", decoder.GetInt(1, 0))
	// fmt.Println("Double Key 2:", decoder.GetDouble(2, 0.0))
	// fmt.Println("Bool Key 3:", decoder.GetBool(3, false))
	// fmt.Println("String Key 4:", decoder.GetString(4, ""))

}
