package httphandler

import "github.com/DarRo9/Tenders/models"

type UsernameRequest struct {
	Username string `form:"username" binding:"required,max=50"`
}

type errorResponse struct {
	Reason string `json:"reason"`
}

type PaginationRequest struct {
	Limit  int32 `form:"limit,default=5" binding:"omitempty,min=1"`
	Offset int32 `form:"offset,default=0" binding:"omitempty,min=0"`
}

type cancelTenderUri struct {
	ID      string `uri:"tenderId" binding:"required,uuid"`
	Version int32  `uri:"version" binding:"required,min=1"`
}

type onesRequest struct {
	PaginationRequest
	UsernameRequest
}

type tenderIdURI struct {
	ID string `uri:"tenderId" binding:"required,uuid"`
}

type allTenderRequests struct {
	Limit       int32                      `form:"limit,default=5" binding:"omitempty,min=1"`
	Offset      int32                      `form:"offset,default=0" binding:"omitempty,min=0"`
	ServiceType []models.TenderServiceType `form:"service_type" binding:"omitempty,dive,oneof=Construction Delivery Manufacture"`
}

type updateTenderStatusRequests struct {
	Status   models.TenderStatus `form:"status" binding:"required,oneof=Created Published Closed"`
	Username string              `form:"username" binding:"required,max=50"`
}

type bidTenderIdURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type bidIdURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type feedbackRequest struct {
	BidFeedback models.BidFeedback `form:"bidFeedback" binding:"required,max=500"`
	Username    string             `form:"username" binding:"required,max=50"`
}

type refreshBidStatusRequest struct {
	Status   models.BidStatus `form:"status" binding:"required,oneof=Created Published Closed"`
	Username string           `form:"username" binding:"required,max=50"`
}

type reviewsRequest struct {
	AuthorUsername    string `form:"authorUsername" binding:"required,max=50"`
	RequesterUsername string `form:"requesterUsername" binding:"required,max=50"`
	PaginationRequest
}

type cancelBidUri struct {
	ID      string `uri:"id" binding:"required,uuid"`
	Version int32  `uri:"version" binding:"required,min=1"`
}

type decisionRequest struct {
	Decision models.BidDecision `form:"decision" binding:"required,oneof=Approved Rejected"`
	UsernameRequest
}


