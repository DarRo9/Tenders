package service

import (
	"context"
	"errors"

	"github.com/DarRo9/Service_for_creating_and_processing_tenders/internal/repository"
	"github.com/DarRo9/Service_for_creating_and_processing_tenders/models"
)

func (s *Service) ConstructBid(ctx context.Context, bid *models.BidCreate) (*models.BidResponse, error) {

	if err := s.repo.IsTenderPudlished(ctx, bid.TenderID); err != nil {
		return nil, err
	}


	if err := s.repo.ControlUserResponsibility(ctx, bid.AuthorId); err != nil {
		return nil, err
	}

	return s.repo.ConstructBid(ctx, bid)
}


func (s *Service) GetBidsOfUser(ctx context.Context, username string, limit, offset int32) ([]*models.BidResponse, error) {

	userId, err := s.repo.ControlBidCreationByID(ctx, username)
	if err != nil {
		return nil, err
	}

	return s.repo.GetBidsOfUser(ctx, userId, limit, offset)
}

func (s *Service) CancelChangesOfBid(ctx context.Context, bidID, username string, version int32) (*models.BidResponse, error) {
	if err := s.repo.ControlBidCreationByName(ctx, bidID, username); err != nil {
		return nil, err
	}

	return s.repo.CancelChangesOfBid(ctx, bidID, version)
}

func (s *Service) GetCommentsOfBid(ctx context.Context, tenderID, authorUsername, requesterUsername string, limit, offset int32) ([]*models.BidReviewResponse, error) {
	if err := s.repo.ControlUserResponsibilityForTender(ctx, tenderID, requesterUsername); err != nil {
		return nil, err
	}

	return s.repo.GetCommentsOfBid(ctx, tenderID, authorUsername, limit, offset)
}

func (s *Service) GetBidsOfTender(ctx context.Context, tenderID, username string, limit, offset int32) ([]*models.BidResponse, error) {

	if err := s.repo.ControlUserResponsibilityForTender(ctx, tenderID, username); err != nil {
		return nil, err
	}
	return s.repo.GetBidsOfTender(ctx, tenderID, limit, offset)
}


func (s *Service) GetStatusOfBids(ctx context.Context, bidID string, username string) (*models.BidStatus, error) {
	bid, err := s.repo.GetBidsWithID(ctx, bidID)
	if err != nil {
		return nil, err
	}


	err = s.repo.ControlUserResponsibilityForAuthorBid(ctx, bidID, username)
	switch {
	case err == nil:
		return &bid.Status, nil
	case !errors.Is(err, repository.ErrRelationNotExist):
		return nil, err
	}

	if err := s.repo.ControlUserResponsibilityForTender(ctx, bid.TenderID, username); err != nil {
		return nil, err
	}

	if bid.Status == models.BidStatusCreated {
		return nil, repository.ErrRelationNotExist
	}

	return &bid.Status, nil
}

func (s *Service) RenewStatusOfBid(ctx context.Context, bidID, username string, status *models.BidStatus) (*models.BidResponse, error) {
	if err := s.repo.ControlBidCreationByName(ctx, bidID, username); err != nil {
		return nil, err
	}

	return s.repo.RenewStatusOfBid(ctx, bidID, username, status)
}

func (s *Service) ChangeBid(ctx context.Context, bidID, username string, bid *models.BidEdit) (*models.BidResponse, error) {
	if err := s.repo.ControlBidCreationByName(ctx, bidID, username); err != nil {
		return nil, err
	}

	return s.repo.ChangeBid(ctx, bidID, bid)
}

func (s *Service) ApplyBidDecision(ctx context.Context, bidID, username string, decision *models.BidDecision) (*models.BidResponse, error) {
	if err := s.repo.ControlUserResponsibilityForTenderByBidID(ctx, bidID, username); err != nil {
		return nil, err
	}

	bid, err := s.repo.ApplyBidDecision(ctx, bidID, username, decision)
	if err != nil {
		return nil, err
	}

	if *decision == models.BidDecisionRejected {
		return s.repo.RenewStatusOfBid(ctx, bidID, username, &models.BidStatusCanceled)
	}

	quorum, err := s.getQuorum(ctx, bidID)
	if err != nil {
		return nil, err
	}

	approvedCount, err := s.repo.CountApplyedDecisions(ctx, bidID)
	if err != nil {
		s.log.Info(err)
		return nil, err
	}

	if approvedCount >= quorum {
		bid, err := s.repo.RenewStatusOfBid(ctx, bidID, username, &models.BidStatusApproved)
		if err != nil {
			s.log.Info(err)
			return nil, err
		}
		_, err = s.repo.RefreshTenderStatus(ctx, bid.TenderID, models.TenderStatusClosed)
		s.log.Info(err)
		return bid, err

	}

	return bid, nil
}

func (s *Service) getQuorum(ctx context.Context, bidID string) (int, error) {
	count, err := s.repo.CountOrganizationsByBid(ctx, bidID)
	return min(3, count), err
}


func (s *Service) ApplyBidFeedback(ctx context.Context, bidID, username string, feedback *models.BidFeedback) (*models.BidResponse, error) {
	bid, err := s.repo.GetBidsWithID(ctx, bidID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.ControlUserResponsibilityForTender(ctx, bid.TenderID, username); err != nil {
		return nil, err
	}

	return bid, s.repo.ApplyBidFeedback(ctx, bidID, feedback)
}

