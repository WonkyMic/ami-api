package data

import (
	"context"
	"log"
	"wonky/ami-api/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4"
)

func AddAuthor(dbpool *pgxpool.Pool, author domain.AddAuthorReq) domain.AuthorRes {
	q := "INSERT INTO author(id, alias, platform, platform_alias_id) VALUES($1, $2, $3, $4) RETURNING id, alias, platform, platform_alias_id"
	ctx := context.Background()
	tx, err := dbpool.BeginTx(ctx, pgx.TxOptions{})

	var author_res domain.AuthorRes
	err = tx.QueryRow(ctx, q, uuid.New(), author.Alias, author.Platform, author.PlatformAliasId).Scan(&author_res.Id, &author_res.Alias, &author_res.Platform, &author_res.PlatformAliasId)
	err = tx.Commit(ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}
	return author_res
}

func GetAuthors(dbpool *pgxpool.Pool) []domain.AuthorRes {
	q := "SELECT id, alias, platform, platform_alias_id FROM author"
	rows, err := dbpool.Query(context.Background(), q)
	if err != nil {
		log.Fatal("error: ", err)
	}
	defer rows.Close()
	var authors []domain.AuthorRes
	for rows.Next() {
		var author domain.AuthorRes
		err := rows.Scan(&author.Id, &author.Alias, &author.Platform, &author.PlatformAliasId)
		if err != nil {
			log.Fatal("error: ", err)
		}
		authors = append(authors, author)
	}
	return authors
}

func GetAuthor(dbpool *pgxpool.Pool, id string) domain.AuthorRes {
	q := "SELECT id, alias, platform, platform_alias_id FROM author WHERE id = $1"
	var author domain.AuthorRes
	err := dbpool.QueryRow(context.Background(), q, id).Scan(&author.Id, &author.Alias, &author.Platform, &author.PlatformAliasId)
	if err != nil {
		log.Fatal("error: ", err)
	}
	return author
}

func GetAuthorByPlatformAliasId(dbpool *pgxpool.Pool, platform_alias_id *uint64) domain.AuthorRes {
	q := "SELECT id, alias, platform, platform_alias_id FROM author WHERE platform_alias_id = $1"
	var author domain.AuthorRes
	_ = dbpool.QueryRow(context.Background(), q, platform_alias_id).Scan(&author.Id, &author.Alias, &author.Platform, &author.PlatformAliasId)
	
	// May return empty as part of author_check
	return author
}

func DeleteAuthor(dbpool *pgxpool.Pool, id string) {
	q := "DELETE FROM author WHERE id = $1"
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