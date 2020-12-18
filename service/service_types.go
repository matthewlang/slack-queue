package service

import (
	"github.com/slack-go/slack"

	"time"
)

type EnqueueRequest struct {
	User *slack.User
}

type EnqueueResponse struct {
	User      *slack.User
	Ok        bool
	Pos       int
	Timestamp time.Time
}

type DequeueRequest struct {
	Place int
	Token int64
}

type DequeueResponse struct {
	User      *slack.User
	Timestamp time.Time
	Token     int64
}

type ListRequest struct {
}

type ListResponse struct {
	Users []*slack.User
	Times []time.Time
	Token int64
}

type RemoveRequest struct {
	Pos   int
	Token int64
}

type RemoveResponse struct {
	Err   error
	Token int64
}