package models

import (
	"time"
)

// DB table
type Proposals struct {
	ChainName       string
	ProposalID      string
	Type            string
	Title           string
	Status          string
	SubmitTime      time.Time
	DepositEndTime  time.Time
	VotingStartTime time.Time
	VotingEndTime   time.Time
	// VoteOption      string
	// VoteTx          string
	UpdateDate  time.Time
	Description string
	MsgID       int
	IsAlert     bool
}
