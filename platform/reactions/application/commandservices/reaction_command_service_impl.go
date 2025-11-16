package commandservices


import (
	"context"
	"errors"
	"fmt"

	"Gommunity/platform/reactions/application/outboundservices/acl"
	"Gommunity/platform/reactions/domain/model/commands"
	"Gommunity/platform/reactions/domain/model/entities"
	"Gommunity/platform/reactions/domain/model/valueobjects"
	"Gommunity/platform/reactions/domain/repositories"
	"Gommunity/platform/reactions/domain/services"
)

type reactionCommandServiceImpl struct {
	reactionRepository repositories.ReactionRepository
	externalPostsService *acl.ExternalPostsService
	externalUsersService *acl.ExternalUsersService
}

// NewReactionCommandService constructs the reactions command service implementation.
func NewReactionCommandService(
	reactionRepository repositories.ReactionRepository,
	externalPostsService *acl.ExternalPostsService,
	externalUsersService *acl.ExternalUsersService,
) services.ReactionCommandService {
	return &reactionCommandServiceImpl{
		reactionRepository:   reactionRepository,
		externalPostsService: externalPostsService,
		externalUsersService: externalUsersService,
	}
}

// HandleAdd adds or updates a user's reaction to a post.
func (s *reactionCommandServiceImpl) HandleAdd(ctx context.Context, cmd commands.AddReactionCommand) (*valueobjects.ReactionID, error) {
	// Validate post exists
	postExists, err := s.externalPostsService.ValidatePostExists(ctx, cmd.PostID())
	if err != nil {
		return nil, fmt.Errorf("failed to validate post: %w", err)
	}
	if !postExists {
		return nil, errors.New("post not found")
	}

	// Validate user exists
	userExists, err := s.externalUsersService.ValidateUserExists(ctx, cmd.UserID())
	if err != nil {
		return nil, fmt.Errorf("failed to validate user: %w", err)
	}
	if !userExists {
		return nil, errors.New("user not found")
	}

	// Check if user already reacted to this post
	existingReaction, err := s.reactionRepository.FindByPostAndUser(ctx, cmd.PostID(), cmd.UserID())
	if err != nil {
		return nil, fmt.Errorf("failed to check existing reaction: %w", err)
	}

	// If reaction exists, update the type (user changed their reaction)
	if existingReaction != nil {
		if err := existingReaction.ChangeReactionType(cmd.ReactionType()); err != nil {
			return nil, err
		}
		if err := s.reactionRepository.Update(ctx, existingReaction); err != nil {
			return nil, fmt.Errorf("failed to update reaction: %w", err)
		}
		reactionID := existingReaction.ReactionID()
		return &reactionID, nil
	}

	// Create new reaction
	reaction, err := entities.NewReaction(
		cmd.PostID(),
		cmd.UserID(),
		cmd.ReactionType(),
	)
	if err != nil {
		return nil, err
	}

	if err := s.reactionRepository.Save(ctx, reaction); err != nil {
		return nil, fmt.Errorf("failed to persist reaction: %w", err)
	}

	reactionID := reaction.ReactionID()
	return &reactionID, nil
}

// HandleRemove removes a user's reaction from a post.
func (s *reactionCommandServiceImpl) HandleRemove(ctx context.Context, cmd commands.RemoveReactionCommand) error {
	// Check if reaction exists
	existingReaction, err := s.reactionRepository.FindByPostAndUser(ctx, cmd.PostID(), cmd.UserID())
	if err != nil {
		return fmt.Errorf("failed to find reaction: %w", err)
	}
	if existingReaction == nil {
		return errors.New("reaction not found")
	}

	// Delete the reaction
	if err := s.reactionRepository.DeleteByPostAndUser(ctx, cmd.PostID(), cmd.UserID()); err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}

	return nil
}
