package models

import "time"

type BidFeedback string

type BidStatus string


type BidReviewResponse struct {
	ID             string      `json:"id"`
	BidID          string      `json:"-"`
	Description    BidFeedback `json:"description"`
	CreatedAt      time.Time   `json:"createdAt"`
}


var (
	BidStatusCreated   BidStatus = "Created"
	BidStatusPublished BidStatus = "Published"
	BidStatusCanceled  BidStatus = "Canceled"
	BidStatusApproved  BidStatus = "Approved"
	BidStatusRejected  BidStatus = "Rejected"
)

type BidDecision string

var (
	BidDecisionApproved BidDecision = "Approved"
	BidDecisionRejected BidDecision = "Rejected"
)

type BidAuthorType string

const (
	BidAuthorTypeOrganization BidAuthorType = "Organization"
	BidAuthorTypeUser         BidAuthorType = "User"
)

type BidCreate struct {
	Name        string        `json:"name" binding:"required,max=100"`
	Description string        `json:"description" binding:"required,max=500"`
	TenderID    string        `json:"tenderId" binding:"required,max=100,uuid"`
	AuthorType  BidAuthorType `json:"authorType" binding:"required,oneof=Organization User"`
	AuthorId    string        `json:"authorId" binding:"required,max=100,uuid"`
}

type BidEdit struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
}

type BidResponse struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      BidStatus     `json:"status"`
	TenderID    string        `json:"tenderId"`
	AuthorType  BidAuthorType `json:"authorType"`
	AuthorID    string        `json:"authorId"`
	Version     int           `json:"version"`
	CreatedAt   time.Time     `json:"createdAt"`
}



func (b *BidEdit) IsEmpty() bool {
	return b.Name == nil && b.Description == nil
}
