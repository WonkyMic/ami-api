package domain

import (
	"github.com/google/uuid"
)

type MessageReq struct {
	AuthorId string
	Content	string
	Platform string
}

type MessageRes struct {
	Id uuid.UUID
	AuthorId uuid.UUID
	Content	string
	Platform string
}

type MessageId struct {
	Id uuid.UUID
}

type AddAuthorReq struct {
	Alias string
	Platform string
	PlatformAliasId uint64
}

type AuthorRes struct {
	Id uuid.UUID
	Alias string
	Platform string
	PlatformAliasId uint64
}