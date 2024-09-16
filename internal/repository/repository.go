package repository

import (
	"context"

	"github.com/DarRo9/Tenders/models"
)

type TenderRepository interface {
	GetAllTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error)
	BuildTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error)
	GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error)
	GetStatusOfTender(ctx context.Context, tenderID string) (*models.TenderStatus, *models.OrganizationID, error)
	RefreshTenderStatus(ctx context.Context, tenderID string, status models.TenderStatus) (*models.TenderResponse, error)
	UpdateTender(ctx context.Context, tenderID string, tenderEdit *models.TenderEdit) (*models.TenderResponse, error)
	RollbackTender(ctx context.Context, tenderID string, version int32) (*models.TenderResponse, error)

	ControlOrganizationPermission(ctx context.Context, organizationID *models.OrganizationID, username string) error
	ControlTendersCreationByName(ctx context.Context, tenderId, creatorUsername string) error
	ControlUserResponsibility(ctx context.Context, userId string) error
	ControlBidCreationByID(ctx context.Context, username string) (string, error)
	IsTenderPudlished(ctx context.Context, tenderID string) error
}

type BidRepository interface {
	ConstructBid(ctx context.Context, bid *models.BidCreate) (*models.BidResponse, error)
	GetBidsOfUser(ctx context.Context, userID string, limit, offset int32) ([]*models.BidResponse, error)
	GetBidsOfTender(ctx context.Context, tenderID string, limit, offset int32) ([]*models.BidResponse, error)
	GetBidsWithID(ctx context.Context, bidID string) (*models.BidResponse, error)
	RenewStatusOfBid(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error)
	ChangeBid(ctx context.Context, bidID string, bidEdit *models.BidEdit) (*models.BidResponse, error)
	ApplyBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error)
	ApplyBidFeedback(ctx context.Context, bidID string, feedback *models.BidFeedback) error
	CancelChangesOfBid(ctx context.Context, bidID string, version int32) (*models.BidResponse, error)
	GetCommentsOfBid(ctx context.Context, tenderID, authorUsername string, limit, offset int32) ([]*models.BidReviewResponse, error)
	
	ControlBidCreationByName(ctx context.Context, bidID, creatorUsername string) error
	ControlUserResponsibilityForTender(ctx context.Context, tenderID, username string) error
	ControlUserResponsibilityForAuthorBid(ctx context.Context, bidID, username string) error
	ControlUserResponsibilityForTenderByBidID(ctx context.Context, bidID, username string) error
	CountOrganizationsByBid(ctx context.Context, bidID string) (int, error)
	CountApplyedDecisions(ctx context.Context, bidID string) (int, error)
}


type Repository interface {
	TenderRepository
	BidRepository
}
