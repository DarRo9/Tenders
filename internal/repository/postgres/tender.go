package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DarRo9/Service_for_creating_and_processing_tenders/internal/repository"
	"github.com/DarRo9/Service_for_creating_and_processing_tenders/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)




func (p *Postgres) GetUserTenders(ctx context.Context, username string, limit, offset int32) ([]*models.TenderResponse, error) {
	rows, err := p.DB.Query(ctx, `
	SELECT *
	FROM tender
		WHERE creator_username = $1
	ORDER BY name ASC
	LIMIT $2 OFFSET $3;`, username, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenders := []*models.TenderResponse{}
	for rows.Next() {
		tender := &models.TenderResponse{}
		if err := rows.Scan(
			&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
			&tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.CreatorUsername); err != nil {
			return nil, err
		}

		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func (p *Postgres) GetAllTenders(ctx context.Context, serviceType []models.TenderServiceType, limit, offset int32) ([]*models.TenderResponse, error) {
	var filter string
	if len(serviceType) != 0 {
		var types []string
		for _, stype := range serviceType {
			types = append(types, fmt.Sprintf("'%v'::service_type", stype))
		}

		filter = fmt.Sprintf(
			"AND service_type = ANY (ARRAY[%s])",
			strings.Join(types, ","),
		)
	}

	query := fmt.Sprintf(`
	SELECT 
		id, name, description, service_type, status, organization_id, version, created_at
	FROM tender
	WHERE status = 'Published'
	%s
	ORDER BY name ASC 
	LIMIT $1 OFFSET $2;`, filter)

	rows, err := p.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenders := []*models.TenderResponse{}
	for rows.Next() {
		tender := &models.TenderResponse{}
		if err := rows.Scan(
			&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
			&tender.OrganizationID, &tender.Version, &tender.CreatedAt); err != nil {
			return nil, err
		}

		tenders = append(tenders, tender)
	}

	return tenders, nil
}


func (p *Postgres) GetStatusOfTender(ctx context.Context, tenderID string) (*models.TenderStatus, *models.OrganizationID, error) {
	var status *models.TenderStatus
	var organizationID *models.OrganizationID

	err := p.DB.QueryRow(ctx, `
	select 
		status, organization_id
	from tender 
	where id = $1;`, tenderID).Scan(&status, &organizationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, repository.ErrTenderNotFound
	}

	return status, organizationID, err
}

func (p *Postgres) BuildTender(ctx context.Context, tender *models.TenderCreate) (*models.TenderResponse, error) {
	tenderResp := &models.TenderResponse{}

	err := p.DB.QueryRow(ctx, `
	insert into tender 
		(name, description, service_type, organization_id, creator_username) 
	values ($1, $2, $3, $4, $5) returning *;`,
		tender.Name, tender.Description, tender.ServiceType, tender.OrganizationID, tender.CreatorUsername).Scan(
		&tenderResp.ID, &tenderResp.Name, &tenderResp.Description, &tenderResp.ServiceType,
		&tenderResp.Status, &tenderResp.OrganizationID, &tenderResp.Version, &tenderResp.CreatedAt, &tenderResp.CreatorUsername)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == repository.FKViolation {
			return nil, repository.ErrOrganizationDepencyNotFound
		}
	}

	return tenderResp, err
}


func (p *Postgres) RefreshTenderStatus(ctx context.Context, tenderID string, status models.TenderStatus) (*models.TenderResponse, error) {
	tender := &models.TenderResponse{}

	err := p.DB.QueryRow(ctx, `
	UPDATE tender
	SET status = $2::tender_status
	WHERE id = $1
	returning *;`, tenderID, status).Scan(
		&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
		&tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.CreatorUsername)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrTenderNotFound
	}

	return tender, err
}

func (p *Postgres) IsTenderPudlished(ctx context.Context, tenderID string) error {
	var status string

    err := p.DB.QueryRow(ctx, `
	SELECT 	
		status 
	FROM tender 
		WHERE id = $1`, tenderID).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrTenderNotFound
	}

    if status != "Published" {
        return repository.ErrTenderClosed
    }
	
	return nil
}

func (p *Postgres) UpdateTender(ctx context.Context, tenderID string, tenderEdit *models.TenderEdit) (*models.TenderResponse, error) {
	tx, err := p.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	pgCmd, err := tx.Exec(ctx, `
	INSERT INTO tender_version 
		(tender_id, name, description, service_type, status, organization_id, version, created_at, creator_username) 
	SELECT
    	id, name, description, service_type, status, organization_id, version, created_at, creator_username
	FROM tender
	WHERE id = $1;`, tenderID)
	if pgCmd.RowsAffected() == 0 {
		return nil, repository.ErrTenderNotFound
	}

	var keys []string
	var values []interface{}

	if tenderEdit.Name != nil {
		keys = append(keys, "name=$1")
		values = append(values, tenderEdit.Name)
	}

	if tenderEdit.Description != nil {
		keys = append(keys, fmt.Sprintf("description=$%d", len(values)+1))
		values = append(values, tenderEdit.Description)
	}

	if tenderEdit.ServiceType != nil {
		keys = append(keys, fmt.Sprintf("service_type=$%d::service_type", len(values)+1))
		values = append(values, tenderEdit.ServiceType)
	}

	values = append(values, tenderID)
	query := fmt.Sprintf(`update tender set %s, version = version + 1 where id = $%v returning *;`, strings.Join(keys, ", "), len(values))

	tender := &models.TenderResponse{}
	err = tx.QueryRow(ctx, query, values...).Scan(
		&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
		&tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.CreatorUsername)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrTenderNotFound
	}

	return tender, err
}

func (p *Postgres) RollbackTender(ctx context.Context, tenderID string, version int32) (*models.TenderResponse, error) {
	tx, err := p.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	pgCmd, err := tx.Exec(ctx, `
	INSERT INTO tender_version 
		(tender_id, name, description, service_type, status, organization_id, version, created_at, creator_username) 
	SELECT
    	id, name, description, service_type, status, organization_id, version, created_at, creator_username
	FROM tender
	WHERE id = $1;`, tenderID)
	if pgCmd.RowsAffected() == 0 {
		return nil, repository.ErrTenderNotFound
	}

	tender := &models.TenderResponse{}
	err = tx.QueryRow(ctx, `
	with tv as (
		select
			name, description, service_type, status, organization_id, version, created_at, creator_username
		from tender_version
			where tender_id = $1 and version = $2
	)
	update tender t
	set
		name = tv.name,
		description = tv.description,
		service_type = tv.service_type,
		status = tv.status,
		organization_id = tv.organization_id,
		version = t.version + 1,
		created_at = tv.created_at,
		creator_username = tv.creator_username
	from tv
		where t.id = $1 
	returning t.*;`, tenderID, version).Scan(
		&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status,
		&tender.OrganizationID, &tender.Version, &tender.CreatedAt, &tender.CreatorUsername)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrTenderORVersionNotFound
	}

	return tender, err
}

