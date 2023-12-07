package helper

import (
	"context"
	"github.com/google/uuid"
)

var clientIDKey = new(struct{})

func UserToContext(ctx context.Context, userId uuid.UUID) context.Context {
	return context.WithValue(ctx, "user_id", userId)
}

func UserFromContext(ctx context.Context) uuid.UUID {
	return ctx.Value("user_id").(uuid.UUID)
}
