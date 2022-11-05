package main

import (
	"context"
	"net/http"
	"log"
	"os"
	"wonky/ami-api/domain"
	"wonky/ami-api/data/author"
	"wonky/ami-api/data/message"
	// "wonky/ami-api/data/reaction"

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

			// Check if author exists
			author_check := author.Get(dbpool, message_req.AuthorId)
			if (domain.AuthorRes{}) == author_check {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Author not found"})
			} else {
				message_res := message.Add(dbpool, message_req)
				c.JSON(http.StatusCreated, message_res)
			} 
		})
		msg.GET("/:id", func(c *gin.Context) {
			id := c.Params.ByName("id")
			message_res := message.Get(dbpool, id)
			c.JSON(http.StatusOK, message_res)
		})
		msg.GET("/", func(c *gin.Context) {
			message_ids := message.GetIdList(dbpool)
			c.JSON(http.StatusOK, message_ids)
		})
	}

	author_route := r.Group("/ami/author")
	{
		author_route.POST("/", func(c *gin.Context) {
			var add_author_req domain.AddAuthorReq
			c.BindJSON(&add_author_req)

			// Check if author already exists
			author_check := author.GetByPlatformAliasId(dbpool, &add_author_req.PlatformAliasId)
			if (domain.AuthorRes{}) != author_check {
				c.JSON(http.StatusFound, gin.H{"message": "Author already exists"})
			} else {
				author_res := author.Add(dbpool, add_author_req)
				c.JSON(http.StatusCreated, author_res)
			}
		})
		author_route.GET("/", func(c *gin.Context) {
			authors := author.GetAuthors(dbpool)
			c.JSON(http.StatusOK, authors)
		})
		author_route.GET("/:id", func(c *gin.Context) {
			id := c.Params.ByName("id")
			author := author.Get(dbpool, id)
			c.JSON(http.StatusOK, author)
		})
		author_route.DELETE("/:id", func(c *gin.Context) {
			id := c.Params.ByName("id")
			author.Delete(dbpool, id)
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