package models

import "time"

// Issue represents a Sentry issue
type Issue struct {
	ID               string        `json:"id"`
	ShortID          string        `json:"shortId"`
	Title            string        `json:"title"`
	Level            string        `json:"level"`
	Status           string        `json:"status"`
	StatusDetails    interface{}   `json:"statusDetails"`
	IsPublic         bool          `json:"isPublic"`
	Platform         string        `json:"platform"`
	Project          IssueProject  `json:"project"`
	Type             string        `json:"type"`
	Count            string        `json:"count"`
	UserCount        int           `json:"userCount"`
	FirstSeen        time.Time     `json:"firstSeen"`
	LastSeen         time.Time     `json:"lastSeen"`
	Permalink        string        `json:"permalink"`
	Logger           string        `json:"logger,omitempty"`
	Metadata         IssueMetadata `json:"metadata"`
	NumComments      int           `json:"numComments"`
	AssignedTo       *IssueUser    `json:"assignedTo,omitempty"`
	IsBookmarked     bool          `json:"isBookmarked"`
	IsSubscribed     bool          `json:"isSubscribed"`
	SubscriptionDetails interface{} `json:"subscriptionDetails"`
	HasSeen          bool          `json:"hasSeen"`
}

// IssueProject represents project info in an issue
type IssueProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// IssueMetadata represents issue metadata
type IssueMetadata struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	Filename string `json:"filename,omitempty"`
	Function string `json:"function,omitempty"`
}

// IssueUser represents user assignment info
type IssueUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}