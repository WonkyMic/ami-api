package message

import (
	"context"
	"log"
	"strconv"
	"wonky/ami-api/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4"
)

func Add(dbpool *pgxpool.Pool, message domain.MessageReq) domain.MessageRes {
	q := "INSERT INTO message(id, author_id, content, platform) VALUES($1, $2, $3, $4) RETURNING id, author_id, content, platform"
	ctx := context.Background()
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})

	var message_res domain.MessageRes
	err = tx.QueryRow(ctx, q, uuid.New(), message.AuthorId, message.Content, message.Platform).Scan(&message_res.Id, &message_res.AuthorId, &message_res.Content, &message_res.Platform)
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}
	return message_res
}

func Get(dbpool *pgxpool.Pool, id string) domain.MessageRes {
	q := "SELECT id, author_id, content, platform FROM message WHERE id = $1"
	var message_res domain.MessageRes
	err := dbpool.QueryRow(context.Background(), q, id).Scan(&message_res.Id, &message_res.AuthorId, &message_res.Content, &message_res.Platform)
	if err != nil {
		log.Print("error: ", err)
	}
	return message_res
}

func GetByAuthor(dbpool *pgxpool.Pool, id string) []domain.MessageRes {
	q := "SELECT id, author_id, content, platform FROM message WHERE author_id = $1"
	rows, err := dbpool.Query(context.Background(), q, id)
	if err != nil {
		log.Print("error: ", err)
	}
	defer rows.Close()
	var messages []domain.MessageRes
	for rows.Next() {
		var message domain.MessageRes
		err := rows.Scan(&message.Id, &message.AuthorId, &message.Content, &message.Platform)
		if err != nil {
			log.Fatal("error: ", err)
		}
		messages = append(messages, message)
	}
	return messages
}

func GetIdList(dbpool *pgxpool.Pool) []domain.MessageId {
	q := "SELECT id FROM message"
	rows, err := dbpool.Query(context.Background(), q)
	if err != nil {
		log.Fatal("error: ", err)
	}
	defer rows.Close()
	var message_ids []domain.MessageId
	for rows.Next() {
		var message_id domain.MessageId
		err := rows.Scan(&message_id.Id)
		if err != nil {
			log.Fatal("error: ", err)
		}
		message_ids = append(message_ids, message_id)
	}
	return message_ids
}
/**
	When a Message is deleted all related reactions must be deleted
**/
func DeleteMessageAndReactions(dbpool *pgxpool.Pool, id string) {
	qDelMessage := "DELETE FROM message WHERE id = $1"
	qDelReactions := "DELETE FROM reactions WHERE message_id = $1"
	ctx := context.Background()
	// TODO - update Isolation level (pgx.TxOptions{})
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Fatal("error: ", err)
	}
	_, err = tx.Exec(ctx, qDelMessage, id)
	if err != nil {
		log.Fatal("error: ", err)
	}
	_, err = tx.Exec(ctx, qDelReactions, id)
	if err != nil {
		log.Fatal("error: ", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func DeleteByAuthor(dbpool *pgxpool.Pool, id string) {
	messages := GetByAuthor(dbpool, id)
	deletedCount := 0
	for _, message := range messages {
		log.Print("Deleting Message and Reactions for Message ID: " + message.Id.String())
		DeleteMessageAndReactions(dbpool, message.Id.String())
		deletedCount += 1
	}
	log.Print(strconv.Itoa(deletedCount) + " messages deleted for Author ID: " + id)
}