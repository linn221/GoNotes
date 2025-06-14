package main

import (
	"context"
	"fmt"
	"linn221/shop/config"
	"linn221/shop/models"
	"log"
)

func main() {
	db := config.ConnectDB()

	user := models.User{
		Username: "linn",
		Password: "secret44",
	}
	userService := models.NewUserService(db)
	_, err := userService.Register(context.Background(), &user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user created successfully")
}
