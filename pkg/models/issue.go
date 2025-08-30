package models

import "time"

// Issue represents a Sentry issue
type Issue struct {
	ID               string        `json:"id"`
	ShareID          string        `json:"shareId,omitempty"`
	ShortID          string        `json:"shortId"`
	Title            string        `json:"title"`
	Level            string        `json:"level"`
	Status           string        `json:"status"`
	Substatus        string        `json:"substatus,omitempty"`
	StatusDetails    interface{}   `json:"statusDetails"`
	Priority         string        `json:"priority,omitempty"`
	PriorityLockedAt *time.Time    `json:"priorityLockedAt,omitempty"`
	IsPublic         bool          `json:"isPublic"`
	Platform         string        `json:"platform"`
	Project          IssueProject  `json:"project"`
	Type             string        `json:"type"`
	IssueType        string        `json:"issueType,omitempty"`
	IssueCategory    string        `json:"issueCategory,omitempty"`
	Count            string        `json:"count"`
	UserCount        int           `json:"userCount"`
	FirstSeen        time.Time     `json:"firstSeen"`
	LastSeen         time.Time     `json:"lastSeen"`
	Permalink        string        `json:"permalink"`
	Logger           string        `json:"logger,omitempty"`
	Culprit          string        `json:"culprit,omitempty"`
	Metadata         interface{}   `json:"metadata"`
	NumComments      int           `json:"numComments"`
	AssignedTo       *IssueUser    `json:"assignedTo,omitempty"`
	Owners           []IssueOwner  `json:"owners,omitempty"`
	IsBookmarked     bool          `json:"isBookmarked"`
	IsSubscribed     bool          `json:"isSubscribed"`
	IsUnhandled      bool          `json:"isUnhandled,omitempty"`
	SubscriptionDetails interface{} `json:"subscriptionDetails"`
	HasSeen          bool          `json:"hasSeen"`
	Annotations      []string      `json:"annotations,omitempty"`
	Activity         []IssueActivity `json:"activity,omitempty"`
}

// IssueProject represents project info in an issue
type IssueProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// IssueOwner represents issue ownership information
type IssueOwner struct {
	Type string `json:"type"` // "user", "team", etc.
	ID   string `json:"id"`
	Name string `json:"name"`
}

// IssueActivity represents issue activity/history
type IssueActivity struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	User      *IssueUser            `json:"user,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Datetime  time.Time             `json:"dateCreated"`
}

// IssueUser represents user assignment info
type IssueUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}