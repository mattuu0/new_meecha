package main

import (
	"meecha/database"
)

func main() {
	/*
	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run("127.0.0.1:12222") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	*/

	//初期設定
	database.DBpath = "./test.db"
	database.Init()
}
