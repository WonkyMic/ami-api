package main

import (
	"context"
	"net/http"
	"log"
	"os"
	"wonky/ami-api/domain"
	"wonky/ami-api/data"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func setupRouter(dbpool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()
	
	r.GET("/ami/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	msg := r.Group("/ami/message") 
	{
		
		msg.POST("/", func(c *gin.Context){
			var message_req domain.MessageReq
			c.BindJSON(&message_req)
			c.JSON(http.StatusCreated, message_req)
		})
		msg.GET("/:id/reply", func(c *gin.Context) {
			id := c.Params.ByName("id")
			c.JSON(http.StatusOK, gin.H{"message": "Reply for id: " + id})
		})
		msg.POST("/:id/feedback", func(c *gin.Context) {
			id := c.Params.ByName("id")
			c.JSON(http.StatusCreated, gin.H{"message": "Stored for id: " + id})
		})
	}

	author := r.Group("/ami/author")
	{
		author.POST("/", func(c *gin.Context) {
			var add_author_req domain.AddAuthorReq
			c.BindJSON(&add_author_req)

			// Check if author already exists
			author_check := data.GetAuthorByPlatformAliasId(dbpool, &add_author_req.PlatformAliasId)
			if (domain.AuthorRes{}) != author_check {
				c.JSON(http.StatusFound, gin.H{"message": "Author already exists"})
			} else {
				author_res := data.AddAuthor(dbpool, add_author_req)
				c.JSON(http.StatusCreated, author_res)
			}
		})
		author.GET("/", func(c *gin.Context) {
			authors := data.GetAuthors(dbpool)
			c.JSON(http.StatusOK, authors)
		})
		author.GET("/:id", func(c *gin.Context) {
			id := c.Params.ByName("id")
			author := data.GetAuthor(dbpool, id)
			c.JSON(http.StatusOK, author)
		})
		author.DELETE("/:id", func(c *gin.Context) {
			id := c.Params.ByName("id")
			data.DeleteAuthor(dbpool, id)
			c.JSON(http.StatusOK, gin.H{"message": "Author Deleted"})
		})
	}

	return r
}

func main() {
	database_url := os.Getenv("WOODCORD_DB_URL")
	config, err := pgxpool.ParseConfig(database_url)
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}

	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer dbpool.Close()

	r := setupRouter(dbpool)

	r.Run(":8080")
}