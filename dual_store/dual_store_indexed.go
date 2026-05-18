package dualstore

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// DualStoreIndexed Struct ကို Go ပုံစံဖြင့် သတ်မှတ်ခြင်း
type DualStoreIndexed struct {
	DB_Path       string   // File Metadata သို့မဟုတ် File Path ဆိုင်ရာ စစ်ဆေးရန်
	readRaf       *os.File // ဖတ်ရန်အတွက် သီးသန့် File Handler
	writeRaf      *os.File // ရေးရန်အတွက် သီးသန့် File Handler
	All_Records   map[int64]*RecordMeta
	Last_Index    int64
	Deleted_Count int
	Deleted_Size  int64
}

func NewDualStoreIndexed(dbPath string) (*DualStoreIndexed, error) {
	//read file
	rFile, err := os.OpenFile(dbPath, os.O_RDONLY, 0666)

	if err != nil {
		return nil, fmt.Errorf("Failed to Open Read: %v\n", err)
	}
	wFile, err := os.OpenFile(dbPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("Failed to Open Write: %v\n", err)
	}

	return &DualStoreIndexed{
		DB_Path:       dbPath,
		readRaf:       rFile,
		writeRaf:      wFile,
		Last_Index:    0,
		Deleted_Size:  0,
		Deleted_Count: 0,
		All_Records:   make(map[int64]*RecordMeta),
	}, nil
}

func (ds *DualStoreIndexed) Close() {
	if ds.readRaf != nil {
		ds.readRaf.Close()
	}
	if ds.writeRaf != nil {
		ds.writeRaf.Close()
	}
}

func (ds *DualStoreIndexed) Load() error {
	info, err := ds.readRaf.Stat()
	if err != nil {
		return err
	}
	size := info.Size()

	if _, err := ds.readRaf.Seek(0, io.SeekStart); err != nil {
		return err
	}
	for {
		currentPos, err := ds.readRaf.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		// EOF
		if currentPos >= size {
			break
		}
		// read meta
		meta, err := ds.ReadRecordMeta()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if meta.Flag == RecordFlagActive {
			ds.All_Records[meta.ID] = meta
		} else {
			// delete
			ds.Deleted_Count++
			ds.Deleted_Size += meta.RecordSize
		}
		// last index ရယူမယ်
		if meta.ID > ds.Last_Index {
			ds.Last_Index = meta.ID
		}
	}
	return nil

}

func (ds *DualStoreIndexed) ReadRecordMeta() (*RecordMeta, error) {
	currentOffset, err := ds.readRaf.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	headerBuf := make([]byte, HeaderSize)
	if _, err := io.ReadFull(ds.readRaf, headerBuf); err != nil {
		return nil, err
	}

	meta := &RecordMeta{
		Offset:        currentOffset,
		Flag:          RecordFlag(headerBuf[FlagPos]),
		ID:            int64(binary.LittleEndian.Uint64(headerBuf[IdPos : IdPos+8])),
		AdapterTypID:  int32(binary.LittleEndian.Uint32(headerBuf[AdapterTypIdPos : AdapterTypIdPos+4])),
		ParentID:      int64(binary.LittleEndian.Uint64(headerBuf[ParentIdPos : ParentIdPos+8])),
		SmallDataSize: int32(binary.LittleEndian.Uint32(headerBuf[SmallDataSizePos : SmallDataSizePos+4])),
		BigDataSize:   int64(binary.LittleEndian.Uint64(headerBuf[BigDataSizePos : BigDataSizePos+8])),
	}
	// fetch small data
	if meta.SmallDataSize > 0 {
		meta.SmallData = make([]byte, meta.SmallDataSize)
		if _, err := io.ReadFull(ds.readRaf, meta.SmallData); err != nil {
			return nil, err
		}
	}
	// get big data offset
	bigDataStartOffset, err := ds.readRaf.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	meta.BigDataStartOffset = bigDataStartOffset

	// skip big data
	if meta.BigDataSize > 0 {
		if _, err := ds.readRaf.Seek(meta.BigDataSize, io.SeekCurrent); err != nil {
			return nil, err

		}
	}
	//cal meta size
	meta.RecordSize = HeaderSize + int64(meta.SmallDataSize) + int64(meta.BigDataSize)

	return meta, nil

}
