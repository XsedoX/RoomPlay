package tests

import (
	"context"
	"time"

	"github.com/google/uuid"
	"xsedox.com/main/domain/user"
)

func AddUserIdToContext(ctx context.Context) (user.Id, context.Context) {
	userId := user.Id(uuid.New())
	ctx = context.WithValue(ctx, user.IdClaimContextKeyName, &userId)
	return userId, ctx
}

type FakeValueProviders struct {
	Sentence string `faker:"sentence"`
	Word     string `faker:"word"`
	Name     string `faker:"name"`
	Url      string `faker:"url"`
	UUID     string `faker:"uuid"`
}

func PtrString(s string) *string     { return &s }
func PtrTime(t time.Time) *time.Time { return &t }
