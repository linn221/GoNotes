package main

import (
	"linn221/shop/config"
	"linn221/shop/models"
	"linn221/shop/utils"
	"log"
)

func main() {
	db := config.ConnectDB()

	user := models.User{
		Username: "linn",
		Password: "secret44",
	}
	h, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Fatal(err.Error())
	}
	user.Password = string(h)
	db.Create(&user)
}
