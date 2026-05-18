package dualstore

import "fmt"

type DualStore struct {
	indexed *DualStoreIndexed
}

func NewDualStore() *DualStore {
	return &DualStore{}
}

func (ds *DualStore) Load(dbPath string) error {
	ds.Close()

	dsi, err := NewDualStoreIndexed(dbPath)
	if err != nil {
		ds.Close()
		return err
	}
	if err := dsi.Load(); err != nil {
		return err
	}
	ds.indexed = dsi
	return nil
}

func (ds *DualStore) GetAllRecords() {
	if ds.indexed == nil {
		return
	}
	for _, meta := range ds.indexed.All_Records {
		fmt.Printf("ID: %v\nAdapterTypID: %v\nParentID: %v\n", meta.ID, meta.AdapterTypID, meta.ParentID)
	}
}
func (ds *DualStore) LastIndex() int64 {
	if ds.indexed == nil {
		return 0
	}
	return ds.indexed.Last_Index
}

func (ds *DualStore) DeletedSize() int64 {
	if ds.indexed == nil {
		return 0
	}
	return ds.indexed.Deleted_Size
}
func (ds *DualStore) DeletedCount() int64 {
	if ds.indexed == nil {
		return 0
	}
	return int64(ds.indexed.Deleted_Count)
}

func (ds *DualStore) Close() {
	if ds.indexed != nil {
		ds.indexed.Close()
		ds.indexed = nil // Garbage Collector အတွက် ကူညီပေးခြင်း
	}
}
