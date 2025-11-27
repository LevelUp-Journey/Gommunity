package events

import (
	"time"

	"Gommunity/platform/community/domain/model/valueobjects"
)

type CommunityCreatedEvent struct {
	communityID valueobjects.CommunityID
	ownerID     valueobjects.OwnerID
	name        string
	occurredOn  time.Time
}

func NewCommunityCreatedEvent(
	communityID valueobjects.CommunityID,
	ownerID valueobjects.OwnerID,
	name string,
) CommunityCreatedEvent {
	return CommunityCreatedEvent{
		communityID: communityID,
		ownerID:     ownerID,
		name:        name,
		occurredOn:  time.Now(),
	}
}

func (e CommunityCreatedEvent) CommunityID() valueobjects.CommunityID {
	return e.communityID
}

func (e CommunityCreatedEvent) OwnerID() valueobjects.OwnerID {
	return e.ownerID
}

func (e CommunityCreatedEvent) Name() string {
	return e.name
}

func (e CommunityCreatedEvent) OccurredOn() time.Time {
	return e.occurredOn
}
