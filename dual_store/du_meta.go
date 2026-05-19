package dualstore

// RecordFlag enum အစား Go မှာ iota သုံးပါမယ်
type RecordFlag byte

const (
	RecordFlagActive RecordFlag = 1 // သင့် logic အလိုက် active flag byte value ထည့်ပါ
	RecordFlagDelete RecordFlag = 2 // သင့် logic အလိုက် delete flag byte value ထည့်ပါ
)

type RecordMeta struct {
	Offset             int64
	BigDataStartOffset int64
	Flag               RecordFlag
	ID                 int64
	AdapterTypID       int32
	ParentID           int64
	SmallDataSize      int32
	SmallData          []byte
	BigDataSize        int64
	RecordSize         int64
}
