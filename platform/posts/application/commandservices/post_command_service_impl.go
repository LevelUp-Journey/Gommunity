package commandservices

import (
	"context"
	"errors"
	"fmt"

	"Gommunity/platform/posts/application/outboundservices/acl"
	"Gommunity/platform/posts/domain/model/commands"
	"Gommunity/platform/posts/domain/model/entities"
	"Gommunity/platform/posts/domain/model/valueobjects"
	"Gommunity/platform/posts/domain/repositories"
	"Gommunity/platform/posts/domain/services"
)

type postCommandServiceImpl struct {
	postRepository               repositories.PostRepository
	externalUsersService         *acl.ExternalUsersService
	externalCommunitiesService   *acl.ExternalCommunitiesService
	externalSubscriptionsService *acl.ExternalSubscriptionsService
}

// NewPostCommandService constructs the posts command service implementation.
func NewPostCommandService(
	postRepository repositories.PostRepository,
	externalUsersService *acl.ExternalUsersService,
	externalCommunitiesService *acl.ExternalCommunitiesService,
	externalSubscriptionsService *acl.ExternalSubscriptionsService,
) services.PostCommandService {
	return &postCommandServiceImpl{
		postRepository:               postRepository,
		externalUsersService:         externalUsersService,
		externalCommunitiesService:   externalCommunitiesService,
		externalSubscriptionsService: externalSubscriptionsService,
	}
}

// HandlePublish publishes a new post.
// Only community owners and admins can publish posts.
func (s *postCommandServiceImpl) HandlePublish(ctx context.Context, cmd commands.CreatePostCommand) (*valueobjects.PostID, error) {
	// Validate community exists
	exists, err := s.externalCommunitiesService.ValidateCommunityExists(ctx, cmd.CommunityID())
	if err != nil {
		return nil, fmt.Errorf("failed to validate community: %w", err)
	}
	if !exists {
		return nil, errors.New("community not found")
	}

	// Validate user exists
	userExists, err := s.externalUsersService.ValidateUserExists(ctx, cmd.AuthorID())
	if err != nil {
		return nil, fmt.Errorf("failed to validate author: %w", err)
	}
	if !userExists {
		return nil, errors.New("author not found")
	}

	// Check user's role in the community
	role, err := s.externalSubscriptionsService.GetUserRole(ctx, cmd.AuthorID(), cmd.CommunityID())
	if err != nil {
		return nil, fmt.Errorf("failed to verify membership: %w", err)
	}

	// If no subscription found, check if user is the community owner
	if role == nil {
		isOwner, err := s.externalCommunitiesService.ValidateUserIsOwner(ctx, cmd.CommunityID(), cmd.AuthorID())
		if err != nil {
			return nil, fmt.Errorf("failed to verify ownership: %w", err)
		}
		if !isOwner {
			return nil, errors.New("only community owners and admins can publish posts")
		}

		ownerRole, roleErr := valueobjects.NewCommunityRole(valueobjects.CommunityRoleOwner)
		if roleErr != nil {
			return nil, roleErr
		}
		role = &ownerRole
	}

	// Only owners and admins can publish. If the stored role is not admin/owner,
	// double-check ownership to tolerate cases where the owner got re-subscribed as member.
	if !role.IsAdminOrOwner() {
		isOwner, err := s.externalCommunitiesService.ValidateUserIsOwner(ctx, cmd.CommunityID(), cmd.AuthorID())
		if err != nil {
			return nil, fmt.Errorf("failed to verify ownership: %w", err)
		}

		if !isOwner {
			// Some deployments store the community owner using profileID instead of userID.
			// Fetch profileID and retry ownership check to avoid false negatives.
			if profileID, profileErr := s.externalUsersService.GetProfileIDByUserID(ctx, cmd.AuthorID()); profileErr == nil {
				if profileAuthorID, voErr := valueobjects.NewAuthorID(profileID); voErr == nil {
					isOwner, err = s.externalCommunitiesService.ValidateUserIsOwner(ctx, cmd.CommunityID(), profileAuthorID)
					if err != nil {
						return nil, fmt.Errorf("failed to verify ownership: %w", err)
					}
				}
			}
		}

		if !isOwner {
			return nil, errors.New("only community owners and admins can publish posts")
		}
	}

	// Create and save post
	post, err := entities.NewPost(
		cmd.CommunityID(),
		cmd.AuthorID(),
		cmd.Content(),
		cmd.Images(),
	)
	if err != nil {
		return nil, err
	}

	if err := s.postRepository.Save(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to persist post: %w", err)
	}

	postID := post.PostID()
	return &postID, nil
}

// HandleDelete removes an existing post if the requester has privileges.
func (s *postCommandServiceImpl) HandleDelete(ctx context.Context, cmd commands.DeletePostCommand) error {
	post, err := s.postRepository.FindByID(ctx, cmd.PostID())
	if err != nil {
		return fmt.Errorf("failed to retrieve post: %w", err)
	}
	if post == nil {
		return errors.New("post not found")
	}

	role, err := s.externalSubscriptionsService.GetUserRole(ctx, cmd.RequestedBy(), post.CommunityID())
	if err != nil {
		return fmt.Errorf("failed to verify requester role: %w", err)
	}

	if role == nil || !role.IsAdminOrOwner() {
		isOwner, ownerErr := s.externalCommunitiesService.ValidateUserIsOwner(ctx, post.CommunityID(), cmd.RequestedBy())
		if ownerErr != nil {
			return fmt.Errorf("failed to verify ownership: %w", ownerErr)
		}
		if !isOwner {
			return errors.New("only community admins or owners can delete posts")
		}
	}

	if err := s.postRepository.Delete(ctx, cmd.PostID()); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}
