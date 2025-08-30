package models

import "time"

// Project represents a Sentry project
type Project struct {
	ID           string        `json:"id"`
	Slug         string        `json:"slug"`
	Name         string        `json:"name"`
	IsPublic     bool          `json:"isPublic"`
	IsBookmarked bool          `json:"isBookmarked"`
	IsMember     bool          `json:"isMember,omitempty"`
	Color        string        `json:"color"`
	DateCreated  time.Time     `json:"dateCreated"`
	FirstEvent   *time.Time    `json:"firstEvent,omitempty"`
	Platform     string        `json:"platform,omitempty"`
	Platforms    []string      `json:"platforms"`
	HasAccess    bool          `json:"hasAccess"`
	Access       []string      `json:"access,omitempty"`
	Features     []string      `json:"features"`
	Status       string        `json:"status"`
	Organization Organization  `json:"organization"`
	Team         *ProjectTeam  `json:"team,omitempty"`
	Teams        []ProjectTeam `json:"teams"`

	// Event processing
	FirstTransactionEvent bool `json:"firstTransactionEvent,omitempty"`
	ProcessingIssues      int  `json:"processingIssues,omitempty"`

	// Capability flags
	HasMinifiedStackTrace bool `json:"hasMinifiedStackTrace,omitempty"`
	HasFeedbacks          bool `json:"hasFeedbacks,omitempty"`
	HasMonitors           bool `json:"hasMonitors,omitempty"`
	HasProfiles           bool `json:"hasProfiles,omitempty"`
	HasReplays            bool `json:"hasReplays,omitempty"`
	HasSessions           bool `json:"hasSessions,omitempty"`

	// Insights capability flags
	HasInsightsHttp          bool `json:"hasInsightsHttp,omitempty"`
	HasInsightsDb            bool `json:"hasInsightsDb,omitempty"`
	HasInsightsAssets        bool `json:"hasInsightsAssets,omitempty"`
	HasInsightsCaches        bool `json:"hasInsightsCaches,omitempty"`
	HasInsightsQueues        bool `json:"hasInsightsQueues,omitempty"`
	HasInsightsVitals        bool `json:"hasInsightsVitals,omitempty"`
	HasInsightsLlmMonitoring bool `json:"hasInsightsLlmMonitoring,omitempty"`

	// Additional metadata
	LatestRelease   *ProjectRelease        `json:"latestRelease,omitempty"`
	Options         map[string]interface{} `json:"options,omitempty"`
	ResolveAge      *int                   `json:"resolveAge,omitempty"`
	DataScrubber    bool                   `json:"dataScrubber,omitempty"`
	SensitiveFields []string               `json:"sensitiveFields,omitempty"`
	GroupingConfig  string                 `json:"groupingConfig,omitempty"`
}

// ProjectTeam represents team information
type ProjectTeam struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// ProjectRelease represents latest release information
type ProjectRelease struct {
	Version      string     `json:"version"`
	ShortVersion string     `json:"shortVersion,omitempty"`
	DateCreated  *time.Time `json:"dateCreated,omitempty"`
}

// ProjectStats represents project statistics
type ProjectStats struct {
	ProjectID string `json:"projectId"`
	Quantity  int64  `json:"quantity"`
	TimesSeen int64  `json:"times_seen"`
}
