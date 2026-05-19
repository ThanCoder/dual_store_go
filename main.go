package main

import (
	"fmt"

	dualstore "github.com/ThanCoder/dual_store_go/dual_store"
	smalldata "github.com/ThanCoder/dual_store_go/small_data"
)

type User struct {
	ID   int64
	Name string
	Age  int64
}

func main() {
	adapter := dualstore.NewAdapterDefault[*User](func(user *User) int64 { return user.ID }, func(value *User, generatedId int64, encoder *smalldata.SmallDataEncoder) []byte {
		encoder.WriteInt(1, generatedId)
		encoder.WriteString(2, value.Name)
		encoder.WriteInt(3, value.Age)

		return encoder.FinishedBytes()
	}, func(decoder smalldata.SmallDataDecoder) *User {
		return &User{
			ID:   decoder.GetIntDefault(1),
			Name: decoder.GetStringDefault(2),
			Age:  decoder.GetIntDefault(3),
		}
	})

	ds := dualstore.NewDualStore()
	if err := dualstore.RegisterAdapterNotExists(ds, adapter); err != nil {
		fmt.Printf("Register Adp Error: %v\n", err)
	}

	adp, err := dualstore.GetAdapter[*User](ds)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("adp: %v\n", adp.AdapterTypeId)

	if err := ds.Load("test.dual.db"); err != nil {
		fmt.Printf("Load Store Error: %v\n", err)
		return
	}

	ds.GetAllRecords()

	fmt.Printf("lastIndex: %v\nDel Count: %v\nDel Size: %v\n", ds.LastIndex(), ds.DeletedCount(), ds.DeletedSize())

	ds.Close()

}
