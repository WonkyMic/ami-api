package reaction

import (
	"context"
	"log"
	"wonky/ami-api/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4"
)

func Add(dbpool *pgxpool.Pool, reaction_req domain.Reaction) domain.Reaction {
	q := "INSERT INTO reactions(message_id, author_id, reaction) VALUES($1, $2, $3) RETURNING message_id, author_id, reaction"
	ctx := context.Background()
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})

	var reaction_res domain.Reaction
	err = tx.QueryRow(ctx, q, reaction_req.MessageId, reaction_req.AuthorId, reaction_req.Reaction).Scan(&reaction_res.MessageId, &reaction_res.AuthorId, &reaction_res.Reaction)
	err = tx.Commit(ctx)
	if err != nil {
		log.Print("error: ", err)
	}
	return reaction_res
}

func Get(dbpool *pgxpool.Pool, reaction_req domain.Reaction) domain.Reaction {
	q := "SELECT message_id, author_id, reaction FROM reactions WHERE message_id = $1 AND author_id = $2 AND reaction = $3"
	var reaction_res domain.Reaction
	err := dbpool.QueryRow(context.Background(), q, &reaction_req.MessageId, &reaction_req.AuthorId, &reaction_req.Reaction).Scan(&reaction_res.MessageId, &reaction_res.AuthorId, &reaction_res.Reaction)
	if err != nil {
		log.Print("error: ", err)
	}
	return reaction_res
}

func GetAuthorReactions(dbpool *pgxpool.Pool, id string) []domain.Reaction {
	q := "SELECT message_id, author_id, reaction FROM reactions WHERE author_id = $1"
	rows, err := dbpool.Query(context.Background(), q, id)
	if err != nil {
		log.Print("error: ", err)
	}
	defer rows.Close()
	reactions := parseReactionRows(rows)
	return reactions
}

func GetMessageReactions(dbpool *pgxpool.Pool, id string) []domain.Reaction {
	q := "SELECT message_id, author_id, reaction FROM reactions WHERE message_id = $1"
	rows, err := dbpool.Query(context.Background(), q, id)
	if err != nil {
		log.Fatal("error: ", err)
	}
	defer rows.Close()
	reactions := parseReactionRows(rows)
	return reactions
}

func parseReactionRows(rows pgx.Rows) []domain.Reaction {
	var reactions []domain.Reaction
	for rows.Next() {
		var reaction domain.Reaction
		err := rows.Scan(&reaction.MessageId, &reaction.AuthorId, &reaction.Reaction)
		if err != nil {
			log.Fatal("error: ", err)
		}
		reactions = append(reactions, reaction)
	}
	return reactions
}

func Delete(dbpool *pgxpool.Pool, reaction_req domain.Reaction) {
	q := "DELETE FROM reactions WHERE message_id = $1 AND author_id = $2 AND reaction = $3"
	ctx := context.Background()
	// TODO - update Isolation level (pgx.TxOptions{})
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Fatal("error: ", err)
	}
	_, err = tx.Exec(ctx, q, reaction_req.MessageId, reaction_req.AuthorId, reaction_req.Reaction)
	if err != nil {
		log.Fatal("error: ", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}
}