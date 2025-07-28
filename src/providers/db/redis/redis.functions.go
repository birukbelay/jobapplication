package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/projTemplate/goauth/src/providers"
)

//===================. blacklisted Related

func IsSessionBlackListed(keyServ providers.KeyValServ, ctx context.Context, sessionId string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", sessionId)
	return keyServ.Exists(ctx, key)
}
func BlacklistSession(keyServ providers.KeyValServ, ctx context.Context, sessionId string) error {
	// Store token in Redis with a prefix to identify blacklist
	key := fmt.Sprintf("blacklist:%s", sessionId)
	return keyServ.Set(ctx, key, "blacklisted", 1*time.Hour)
}
