package main

import (
	"context"
	"go-service-template/internal/database/nosql"
	"log"
)

func main() {
	uri := "mongodb://admin:adminpassword@localhost:27017"
	dbName := "repo"
	collection := "records"
	type product struct {
		ID   string `bson:"_id"`
		Name string `bson:"name"`
	}

	nosql.InitNoSQLDatabase(uri)
	client := nosql.GetClient()
	log.Printf("client : %v", client)
	r, _ := nosql.NewNoSQLRepository[product](client, dbName, collection)
	list := []*product{&product{ID: "1", Name: "sss"}, &product{ID: "2", Name: "aaa"}, &product{ID: "1", Name: "mohamed"}, &product{ID: "3", Name: "mohamed"}}

	// r.Save(context.TODO(), &product{ID: "123", Name: "sdsdsds"})
	// err := r.SaveAll(context.TODO(), list)
	// log.Printf("err : %v", err)

	// failed, err := r.SavePartialSuccess(context.Background(), list)
	// for _, val := range failed {
	// 	log.Printf("%v\n", val)
	// }
	// log.Printf("err : %v", err)

	// t, err := r.SaveAtomic(context.Background(), list)
	// log.Printf("err : %v", err)
	// log.Printf("t : %v", t)

}
