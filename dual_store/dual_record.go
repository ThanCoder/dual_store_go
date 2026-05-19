package dualstore

import (
	"encoding/binary"
	"io"
	"os"
)

type DualRecord struct {
	Flag          RecordFlag
	ID            int64
	AdapterTypID  int32
	ParentID      int64
	BigDataSize   int64
	SmallData     []byte
	BigDataType   BigDataType
	BigDataStream io.Reader // Dart ရဲ့ Stream အစား io.Reader ကို သုံးပါတယ်
}

// offset -> header start offset
func (r *DualRecord) ToMeta(offset int64) *RecordMeta {
	smallDataSize := int32(len(r.SmallData))
	return &RecordMeta{
		Offset:             offset,
		Flag:               RecordFlagActive,
		ID:                 r.ID,
		ParentID:           -1,
		AdapterTypID:       -1,
		SmallData:          r.SmallData,
		SmallDataSize:      smallDataSize,
		BigDataSize:        r.BigDataSize,
		BigDataStartOffset: offset + int64(HeaderSize+smallDataSize),
		RecordSize:         int64(HeaderSize+smallDataSize) + r.BigDataSize,
	}
}

// Return -> header start offset
func (r *DualRecord) Write(file *os.File, onProgress func(progress float64)) (int64, error) {
	offset, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	header := make([]byte, HeaderSize)
	header[FlagPos] = byte(r.Flag)

	binary.LittleEndian.PutUint64(header[IdPos:], uint64(r.ID))
	binary.LittleEndian.PutUint32(header[AdapterTypIdPos:], uint32(r.AdapterTypID))
	binary.LittleEndian.PutUint64(header[ParentIdPos:], uint64(r.ParentID))
	binary.LittleEndian.PutUint32(header[SmallDataSizePos:], uint32(len(r.SmallData)))
	binary.LittleEndian.PutUint64(header[BigDataSizePos:], uint64(r.BigDataSize))

	// ၃။ Header ကော SmallData ပါ အတူတူ ပေါင်းရေးခြင်း
	// Go မှာ Buffer အသစ်ထပ်မဆောက်တော့ဘဲ Write ကို နှစ်ခါခေါ်ခြင်းက ပိုမြန်ဆန်ပါတယ်

	// write header
	if _, err := file.Write(header); err != nil {
		return 0, err
	}
	// write small data
	if _, err := file.Write(r.SmallData); err != nil {
		return 0, err
	}

	// ၄။ Big Data Stream ကို ဖတ်ပြီး ရေးခြင်း (OOM မဖြစ်အောင် Chunk အလိုက် ရေးသည်)
	if r.BigDataSize > 0 && r.BigDataStream != nil {
		var writtenBytes int64
		lastProgress := -1

		// 32KB buffer ဖြင့် chunk အလိုက် ဖတ်ရန်
		buf := make([]byte, 3*1024)
		for {
			n, readErr := r.BigDataStream.Read(buf)
			if n > 0 {
				if _, writeErr := file.Write(buf[:n]); writeErr != nil {
					return 0, writeErr
				}
				writtenBytes += int64(n)
				// progress တွက်မယ်
				currentProgress := int((float64(writtenBytes) / float64(r.BigDataSize)) * 100)
				if currentProgress != lastProgress {
					lastProgress = currentProgress
					onProgress(float64(currentProgress))
				}

			}

			// file အဆုံးသတ်သွားပြီ
			if readErr == io.EOF {
				break
			}
			// error
			if readErr != nil {
				return 0, readErr
			}
		}
	}
	return offset, nil
}
