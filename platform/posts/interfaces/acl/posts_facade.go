package acl


import "context"

// PostsFacade exposes posts operations to other bounded contexts.
type PostsFacade interface {
	PostExists(ctx context.Context, postID string) (bool, error)
}
