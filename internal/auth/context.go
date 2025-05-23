package auth

import "context"

type contextKey string

const ContextKeyClaims = contextKey("authClaims")

func SetClaims(ctx context.Context, claims *Auth52Claims) context.Context {
	return context.WithValue(ctx, ContextKeyClaims, claims)
}

func GetClaims(ctx context.Context) (*Auth52Claims, bool) {
	claims, ok := ctx.Value(ContextKeyClaims).(*Auth52Claims)
	return claims, ok
}
