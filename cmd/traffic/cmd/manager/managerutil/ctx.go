package managerutil

import (
	"context"

	"github.com/datawire/ambassador/pkg/kates"
	"github.com/datawire/dlib/dlog"
	"github.com/telepresenceio/telepresence/rpc/v2/manager"
)

func WithSessionInfo(ctx context.Context, si *manager.SessionInfo) context.Context {
	if id := si.GetSessionId(); id != "" {
		return WithSessionID(ctx, id)
	}
	return ctx
}

func WithSessionID(ctx context.Context, sessionID string) context.Context {
	ctx = context.WithValue(ctx, sessionContextKey{}, sessionID)
	ctx = dlog.WithField(ctx, "session_id", sessionID)
	return ctx
}

func GetSessionID(ctx context.Context) string {
	id := ctx.Value(sessionContextKey{})
	if id == nil {
		return ""
	}
	return id.(string)
}

type sessionContextKey struct{}

func WithKatesClient(ctx context.Context, client *kates.Client) context.Context {
	return context.WithValue(ctx, katesClientContextKey{}, client)
}

func GetKatesClient(ctx context.Context) *kates.Client {
	client, ok := ctx.Value(katesClientContextKey{}).(*kates.Client)
	if !ok {
		return nil
	}
	return client
}

type katesClientContextKey struct{}
