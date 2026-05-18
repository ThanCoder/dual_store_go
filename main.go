package main

import (
	"fmt"

	dualstore "github.com/ThanCoder/dual_store_go/dual_store"
)

func main() {
	ds := dualstore.NewDualStore()
	if err := ds.Load("apyar.dual.db-no-gzip"); err != nil {
		fmt.Printf("Load Store Error: %v\n", err)
		return
	}
	ds.GetAllRecords()

	fmt.Printf("lastIndex: %v\nDel Count: %v\nDel Size: %v\n", ds.LastIndex(), ds.DeletedCount(), ds.DeletedSize())

	ds.Close()

	// ds, err := dualstore.NewDualStoreIndexed("apyar.dual.db-no-gzip")
	// if err != nil {
	// 	fmt.Printf("Load Store Error: %v\n", err)
	// 	return
	// }
	// ds.Load()

	// for _, meta := range ds.All_Records {
	// 	fmt.Printf("ID: %v\nAdapterTypID: %v\nParentID: %v\n", meta.ID, meta.AdapterTypID, meta.ParentID)
	// }

	// fmt.Printf("lastIndex: %v\nDel Count: %v\nDel Size: %v\n", ds.Last_Index, ds.Deleted_Count, ds.Deleted_Size)

	// ds.Close()

}
