package domain

import (
	"github.com/google/uuid"
)

type MessageReq struct {
	Author string
	Content	string
	Platform string
}

type MessageRes struct {
	Id string
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