package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tayl0r89/bgdict/api"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mysql_user := os.Getenv("MYSQL_USER")
	mysql_pass := os.Getenv("MYSQL_PASSWORD")
	mysql_db := os.Getenv("MYSQL_DATABASE")
	hostname := os.Getenv("HOSTNAME")
	port := os.Getenv("PORT")

	router := gin.Default()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(db)/%s", mysql_user, mysql_pass, mysql_db))

	if err != nil {
		log.Println(err.Error())
		return
	}

	wordRepository := api.NewWordRepository(db)

	router.GET("/find", api.FindWordHandler(wordRepository))
	router.GET("/id", api.GetWordByIdHandler(wordRepository))
	router.POST("/bulk", api.BulkSearchHandler(wordRepository))
	router.GET("/search", api.SearchHandler(wordRepository))
	router.POST("/bulkById", api.BulkByIdHandler(wordRepository))

	router.Run(fmt.Sprintf("%s:%s", hostname, port))
}
