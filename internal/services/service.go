package service

import (
	"context"

	"github.com/DarRo9/Service_for_creating_and_processing_tenders/internal/repository"
	"github.com/DarRo9/Service_for_creating_and_processing_tenders/models"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repo repository.Repository
	log  *logrus.Logger
}

type BidService interface {
	ConstructBid(ctx context.Context, bid *models.BidCreate) (*models.BidResponse, error)
	GetBidsOfUser(ctx context.Context, username string, limit, offset int32) ([]*models.BidResponse, error)
	GetBidsOfTender(ctx context.Context, tenderID, username string, limit, offset int32) ([]*models.BidResponse, error)
	GetStatusOfBids(ctx context.Context, bidID string, username string) (*models.BidStatus, error)
	RenewStatusOfBid(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error)
	ChangeBid(ctx context.Context, bidID, username string, bid *models.BidEdit) (*models.BidResponse, error)
	ApplyBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error)
	ApplyBidFeedback(ctx context.Context, bidID, username string, feedback *models.BidFeedback) (*models.BidResponse, error)
	CancelChangesOfBid(ctx context.Context, bidID, username string, version int32) (*models.BidResponse, error)
	GetCommentsOfBid(ctx context.Context, tenderID, authorUsername, requesterUsername string, limit, offset int32) ([]*models.BidReviewResponse, error)
}

type TenderService interface {
	GetAllTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error)
	BuildTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error)
	GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error)
	GetStatusOfTender(ctx context.Context, tenderID, username string) (*models.TenderStatus, error)
	RefreshTenderStatus(ctx context.Context, tenderID, username string, status models.TenderStatus) (*models.TenderResponse, error)
	ChangeTender(ctx context.Context, tenderID string, username string, tender *models.TenderEdit) (*models.TenderResponse, error)
	RollbackTender(ctx context.Context, tenderID string, version int32, username string) (*models.TenderResponse, error)
}

func New(repo repository.Repository, log *logrus.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}
