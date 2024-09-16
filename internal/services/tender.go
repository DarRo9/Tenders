package service

import (
	"context"

	"github.com/DarRo9/Tenders/models"
)

func (s *Service) RefreshTenderStatus(ctx context.Context, tenderID, username string, status models.TenderStatus) (*models.TenderResponse, error) {
	if err := s.repo.ControlTendersCreationByName(ctx, tenderID, username); err != nil {
		return nil, err
	}

	return s.repo.RefreshTenderStatus(ctx, tenderID, status)
}

func (s *Service) ChangeTender(ctx context.Context, tenderID string, username string, tender *models.TenderEdit) (*models.TenderResponse, error) {
	if err := s.repo.ControlTendersCreationByName(ctx, tenderID, username); err != nil {
		return nil, err
	}

	return s.repo.UpdateTender(ctx, tenderID, tender)
}

func (s *Service) RollbackTender(ctx context.Context, tenderID string, version int32, username string) (*models.TenderResponse, error) {
	if err := s.repo.ControlTendersCreationByName(ctx, tenderID, username); err != nil {
		return nil, err
	}

	return s.repo.RollbackTender(ctx, tenderID, version)
}

func (s *Service) GetAllTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error) {
	return s.repo.GetAllTenders(ctx, serviceType, limit, offset)
}

func (s *Service) BuildTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error) {
	err := s.repo.ControlOrganizationPermission(ctx, &tender.OrganizationID, tender.CreatorUsername)
	if err != nil {
		return nil, err
	}

	return s.repo.BuildTender(ctx, tender)
}

func (s *Service) GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error) {
	_, err := s.repo.ControlBidCreationByID(ctx, username)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserTenders(ctx, username, limit, offset)
}

func (s *Service) GetStatusOfTender(ctx context.Context, tenderID, username string) (*models.TenderStatus, error) {
	status, organizationID, err := s.repo.GetStatusOfTender(ctx, tenderID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.ControlOrganizationPermission(ctx, organizationID, username); err != nil {
		return nil, err
	}

	return status, nil
}