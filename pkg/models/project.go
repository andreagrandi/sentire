package models

import "time"

// Project represents a Sentry project
type Project struct {
	ID           string        `json:"id"`
	Slug         string        `json:"slug"`
	Name         string        `json:"name"`
	IsPublic     bool          `json:"isPublic"`
	IsBookmarked bool          `json:"isBookmarked"`
	Color        string        `json:"color"`
	DateCreated  time.Time     `json:"dateCreated"`
	FirstEvent   *time.Time    `json:"firstEvent,omitempty"`
	Platform     string        `json:"platform,omitempty"`
	Platforms    []string      `json:"platforms"`
	HasAccess    bool          `json:"hasAccess"`
	Features     []string      `json:"features"`
	Status       string        `json:"status"`
	Organization Organization  `json:"organization"`
	Team         *ProjectTeam  `json:"team,omitempty"`
	Teams        []ProjectTeam `json:"teams"`
}

// ProjectTeam represents team information
type ProjectTeam struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// ProjectStats represents project statistics
type ProjectStats struct {
	ProjectID string `json:"projectId"`
	Quantity  int64  `json:"quantity"`
	TimesSeen int64  `json:"times_seen"`
}