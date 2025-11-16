package events

import "time"

// ProfileUpdatedEvent represents the event when a user profile is updated
type ProfileUpdatedEvent struct {
	UserID     string  `json:"userId"`
	ProfileID  string  `json:"profileId"`
	Username   string  `json:"username"`
	ProfileURL *string `json:"profileUrl"`
	OccurredOn []int   `json:"occurredOn"`
}

// GetOccurredOn converts the array format to time.Time
func (e ProfileUpdatedEvent) GetOccurredOn() time.Time {
	if len(e.OccurredOn) >= 7 {
		return time.Date(
			e.OccurredOn[0],             // year
			time.Month(e.OccurredOn[1]), // month
			e.OccurredOn[2],             // day
			e.OccurredOn[3],             // hour
			e.OccurredOn[4],             // minute
			e.OccurredOn[5],             // second
			e.OccurredOn[6],             // nanosecond
			time.UTC,
		)
	}
	return time.Now()
}
