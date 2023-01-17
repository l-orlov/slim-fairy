package main

import (
	"context"
	"log"

	"github.com/l-orlov/slim-fairy/internal/db"
	"github.com/l-orlov/slim-fairy/internal/model"
)

func main() {
	ctx := context.Background()

	// connect to db
	database, err := db.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	//id := uuid.MustParse("bc294c5a-7e96-4a90-8670-f7b90f4b8faa")
	//userByID, err := database.GetUser(ctx, id)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("userByID: %+v\n", userByID)

	userToCreate := &model.User{
		Name: "Some user",
		Age:  23,
	}
	err = database.CreateUser(ctx, userToCreate)
	if err != nil {
		log.Fatal(err)
	}

	users, err := database.GetUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		log.Printf("user: %+v\n", user)
	}
}
