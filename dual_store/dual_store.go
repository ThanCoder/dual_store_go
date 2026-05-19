package dualstore

import (
	"fmt"
	"reflect"
	"sync"

	smalldata "github.com/ThanCoder/dual_store_go/small_data"
)

type DualStore struct {
	indexed  *DualStoreIndexed
	adapters map[reflect.Type]any
	mu       sync.RWMutex
}

func NewDualStore() *DualStore {
	return &DualStore{
		adapters: make(map[reflect.Type]any),
	}
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
		dec, err := smalldata.NewSmallDataDecoder(meta.SmallData)
		if err != nil {
			fmt.Printf("Dec Error: %v\n", err)
			continue
		}
		dec.PrintAllData()
		// break

		// fmt.Printf("ID: %v\nAdapterTypID: %v\nParentID: %v\n", meta.ID, meta.AdapterTypID, meta.ParentID)
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

// register adatper
func RegisterAdapterNotExists[T DualModel](ds *DualStore, adapter *DualAdapter[T]) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var dummy T
	t := reflect.TypeOf(dummy)

	// fmt.Printf("T: %v\n",t)
	// type ကို အရင်စစ်
	if existing, exists := ds.adapters[t]; exists {
		// type ရှိနေလား စစ်
		if existingAdp, ok := existing.(*DualAdapter[T]); ok {
			// adapter id တူလား စစ်
			if existingAdp.AdapterTypeId == adapter.AdapterTypeId {
				return nil
			}
		}
	}
	// type မတူရင် | id တူနေလားစစ်
	for _, existing := range ds.adapters {
		v := reflect.ValueOf(existing).Elem()
		idField := v.FieldByName("AdapterTypeId")
		if idField.IsValid() && int32(idField.Int()) == adapter.AdapterTypeId {
			return fmt.Errorf("Duplicate Adapter ID: Unique id `%d` is already used by another type", adapter.AdapterTypeId)
		}
	}

	// မရှိဘူးဆိုရင် ထည့်မယ်
	ds.adapters[t] = adapter

	return nil
}

// get adapter
func GetAdapter[T DualModel](ds *DualStore) (*DualAdapter[T], error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var dummy T
	t := reflect.TypeOf(dummy)

	raw, exists := ds.adapters[t]
	if !exists {
		return nil, fmt.Errorf("No Adapter Registered for type `%s`", t.String())
	}
	adapter, ok := raw.(*DualAdapter[T])
	if !ok {
		return nil, fmt.Errorf("Adapter type assertion failed for `%s`", t.String())
	}
	return adapter, nil

}
