package models

import "time"

// Organization represents a Sentry organization
type Organization struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	DateCreated time.Time `json:"dateCreated"`
	Status      struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"status"`
	Avatar struct {
		AvatarType string `json:"avatarType"`
		AvatarUUID string `json:"avatarUuid"`
	} `json:"avatar"`
	Features       []string `json:"features"`
	IsEarlyAdopter bool     `json:"isEarlyAdopter"`
	Access         []string `json:"access"`
}

// OrganizationStats represents organization event statistics
type OrganizationStats struct {
	Start    time.Time            `json:"start"`
	End      time.Time            `json:"end"`
	Projects []ProjectStatsDetail `json:"projects"`
	Totals   struct {
		Sum       int64 `json:"sum"`
		TimesSeen int64 `json:"times_seen"`
	} `json:"totals"`
}

// ProjectStatsDetail represents detailed project statistics
type ProjectStatsDetail struct {
	ID    interface{}     `json:"id"` // Can be string or number
	Slug  string          `json:"slug,omitempty"`
	Stats []CategoryStats `json:"stats"`
}

// CategoryStats represents statistics by category
type CategoryStats struct {
	Category string        `json:"category"`
	Outcomes StatsOutcomes `json:"outcomes"`
	Totals   StatsTotals   `json:"totals"`
}

// StatsOutcomes represents event outcomes
type StatsOutcomes struct {
	Accepted           int64 `json:"accepted"`
	Filtered           int64 `json:"filtered"`
	RateLimited        int64 `json:"rate_limited"`
	Invalid            int64 `json:"invalid"`
	Abuse              int64 `json:"abuse"`
	ClientDiscard      int64 `json:"client_discard"`
	CardinalityLimited int64 `json:"cardinality_limited"`
}

// StatsTotals represents total statistics
type StatsTotals struct {
	Dropped   int64 `json:"dropped"`
	Sum       int64 `json:"sum(quantity)"`
	TimesSeen int64 `json:"times_seen"`
}
