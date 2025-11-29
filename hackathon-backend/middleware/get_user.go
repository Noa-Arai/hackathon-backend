package middleware

import "context"

// Context から userID を取得する補助関数
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return "", false
	}
	uid, ok := val.(string)
	return uid, ok
}
