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
	Features    []string `json:"features"`
	IsEarlyAdopter bool  `json:"isEarlyAdopter"`
	Access      []string `json:"access"`
}

// OrganizationStats represents organization event statistics
type OrganizationStats struct {
	Start    time.Time      `json:"start"`
	End      time.Time      `json:"end"`
	Projects []ProjectStats `json:"projects"`
	Totals   struct {
		Sum       int64 `json:"sum"`
		TimesSeen int64 `json:"times_seen"`
	} `json:"totals"`
}