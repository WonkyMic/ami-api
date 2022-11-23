package main

import (
	"context"
	"net/http"
	"log"
	"os"
	"strconv"

	"wonky/ami-api/domain"
	"wonky/ami-api/data/author"
	"wonky/ami-api/data/message"
	"wonky/ami-api/data/reaction"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)
// TODO :: Return 404s instead of 'null' 200s
// 		:: Null checks
func setupRouter(dbpool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()
	
	r.GET("/ami/health", func(c *gin.Context) {
		c.String(http.StatusOK, "we good")
	})
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
		msg.DELETE("/:id", func(c *gin.Context) {
			id := c.Params.ByName("id")
			message.DeleteMessageAndReactions(dbpool, id)
			c.JSON(http.StatusOK, gin.H{"message": "Message Deleted"})
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
		msg.GET("/:id/reactions", func(c *gin.Context) {
			id := c.Params.ByName("id")
			reactions := reaction.GetMessageReactions(dbpool, id)
			c.JSON(http.StatusOK, reactions)
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
		author_route.GET("/search", func(c *gin.Context) {
			platform_id, err := strconv.ParseUint(c.Query("platformAliasId"), 10, 64)
			if err != nil {
				log.Fatal("error parsing query param: ", err)
			}
			author := author.GetByPlatformAliasId(dbpool, &platform_id)
			if (domain.AuthorRes{}) == author {
				c.JSON(http.StatusNotFound, "message: No Author Found")
			} else {
				c.JSON(http.StatusOK, author)
			}
		})
		author_route.GET("/:id", func(c *gin.Context) {
			id := c.Params.ByName("id")
			author := author.Get(dbpool, id)
			c.JSON(http.StatusOK, author)
		})
		author_route.GET("/:id/messages", func(c *gin.Context) {
			id := c.Params.ByName("id")
			messages := message.GetByAuthor(dbpool, id)
			c.JSON(http.StatusOK, messages)
		})
		author_route.GET("/:id/reactions", func(c *gin.Context) {
			id := c.Params.ByName("id")
			reactions := reaction.GetAuthorReactions(dbpool, id)
			c.JSON(http.StatusOK, reactions)
		})
		author_route.DELETE("/:id", func(c *gin.Context) {
			id := c.Params.ByName("id")
			// Delete Reactions for Author
			reaction.DeleteByAuthor(dbpool, id)
			// Delete messages and related message reactions by Author
			message.DeleteByAuthor(dbpool, id)
			// Delete Author
			author.Delete(dbpool, id)
			c.JSON(http.StatusOK, gin.H{"message": "Author Deleted"})
		})
	}

	reaction_route := r.Group("/ami/reaction")
	{
		reaction_route.POST("/", func(c *gin.Context) {
			var add_reaction_req domain.Reaction
			c.BindJSON(&add_reaction_req)

			// Check if message exists
			message_check := message.Get(dbpool, add_reaction_req.MessageId)
			if (domain.MessageRes{}) == message_check {
				c.JSON(http.StatusNotFound, gin.H{"message": "MessageId for Reaction not found"})
				return
			}
			// Check if author exists
			author_check := author.Get(dbpool, add_reaction_req.AuthorId)
			if (domain.AuthorRes{}) == author_check {
				c.JSON(http.StatusNotFound, gin.H{"message": "AuthorId for Reaction not found"})
				return
			}
			// Check if reaction relationship already exists
			reaction_check := reaction.Get(dbpool, add_reaction_req)
			if (domain.Reaction{}) != reaction_check {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Reaction already exists"})
				return
			}

			reaction_res := reaction.Add(dbpool, add_reaction_req)
			c.JSON(http.StatusCreated, reaction_res)
		})
		reaction_route.DELETE("/", func(c *gin.Context) {
			var reaction_req domain.Reaction
			c.BindJSON(&reaction_req)
			reaction.Delete(dbpool, reaction_req)
			c.JSON(http.StatusOK, gin.H{"message": "Reaction Deleted"})
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