package middleware

import "context"

type contextKey string

const (
	userIDKey   contextKey = "user_id"
	usernameKey contextKey = "username"
)

func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}

func UsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey).(string)
	return username, ok
}
