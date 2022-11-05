package message

import (
	"context"
	"log"
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

func Delete(dbpool *pgxpool.Pool, id string) {
	q := "DELETE FROM message WHERE id = $1"
	ctx := context.Background()
	// TODO - update Isolation level (pgx.TxOptions{})
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Fatal("error: ", err)
	}

	_, err = tx.Exec(ctx, q, id)
	if err != nil {
		log.Fatal("error: ", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}
}