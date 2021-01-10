package auth

import (
	"github.com/zapkub/cftl/internal/repository"
)

type Authenticator struct {
	db *repository.DB
}

func New(db *repository.DB) *Authenticator {
	return &Authenticator{
		db: db,
	}
}
