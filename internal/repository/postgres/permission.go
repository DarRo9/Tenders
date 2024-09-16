package postgres

import (
	"context"
	"errors"

	"github.com/DarRo9/Service_for_creating_and_processing_tenders/internal/repository"
	"github.com/DarRo9/Service_for_creating_and_processing_tenders/models"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) ControlOrganizationPermission(ctx context.Context, organizationID *models.OrganizationID, username string) error {
	var existsRelation bool

	err := p.DB.QueryRow(ctx, `
	SELECT
		EXISTS (
			SELECT 1
			FROM organization_responsible
			WHERE user_id = e.id
			AND organization_id = $2
		) AS exists_relation
	FROM employee e
	WHERE username = $1;`, username, organizationID).Scan(&existsRelation)

	switch {
	
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	
	case !existsRelation:
		return repository.ErrRelationNotExist
	}

	return err
}

func (p *Postgres) ControlBidCreationByName(ctx context.Context, bidID, creatorUsername string) error {
	var isCreator bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
		SELECT 1
		FROM bid b 
			WHERE author_id = e.id and id = $2
	) AS is_creator
	FROM employee e
		WHERE username = $1`, creatorUsername, bidID).Scan(&isCreator)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	case !isCreator:
		return repository.ErrRelationNotExist
	}

	return err
}

func (p *Postgres) ControlBidCreationByID(ctx context.Context, username string) (string, error) {
	var userId string
	err := p.DB.QueryRow(ctx, `
	select 
		id 
	from employee e 
	where e.username = $1;
	`, username).Scan(&userId)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", repository.ErrUserNotExist
	}

	return userId, err
}

func (p *Postgres) ControlUserResponsibility(ctx context.Context, userId string) error {
	var existsRelation bool

	err := p.DB.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1
            FROM organization_responsible orr
            WHERE orr.user_id = $1
        ) AS exists_relation
        FROM employee e
        WHERE e.id = $1`, userId).Scan(&existsRelation)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	case !existsRelation:
		return repository.ErrRelationNotExist
	}

	return err
}

func (p *Postgres) ControlTendersCreationByName(ctx context.Context, tenderId, creatorUsername string) error {
	var isCreator bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
		SELECT 1
		FROM tender
			WHERE creator_username = $1 AND id = $2
	) AS is_creator
	FROM employee
		WHERE username = $1;`, creatorUsername, tenderId).Scan(&isCreator)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	case !isCreator:
		return repository.ErrRelationNotExist
	}

	return err
}

func (p *Postgres) ControlTendersCreationByID(ctx context.Context, tenderId, creatorId string) error {
	var isCreator bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
        SELECT 1
        FROM tender t
        WHERE t.creator_username = e.username
        AND t.id = $1
    ) AS is_creator
    FROM employee e
    	WHERE e.id = $2;`, tenderId, creatorId).Scan(&isCreator)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ErrUserNotExist
	case !isCreator:
		return repository.ErrRelationNotExist
	}

	return err
}

func (p *Postgres) ControlUserResponsibilityForTender(ctx context.Context, tenderID, username string) error {
	var isRelated bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
		SELECT 1
		FROM tender t
		JOIN organization_responsible orr ON orr.user_id = e.id
		WHERE t.id = $2
		AND orr.organization_id = t.organization_id
	) AS is_related
	from employee e
		where username = $1; `, username, tenderID).Scan(&isRelated)

	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrUserNotExist
	}

	if !isRelated {
		return repository.ErrRelationNotExist
	}

	return err
}

func (p *Postgres) ControlUserResponsibilityForAuthorBid(ctx context.Context, bidID, username string) error {
	var isRelated bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
		SELECT 1
			FROM bid b
			JOIN organization_responsible orr ON orr.user_id = e.id
			JOIN organization_responsible o ON o.organization_id = orr.organization_id
				WHERE b.id = $2
				AND o.user_id = b.author_id
	) AS is_related
	FROM employee e
		WHERE e.username = $1;`, username, bidID).Scan(&isRelated)

	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrUserNotExist
	}

	if !isRelated {
		return repository.ErrRelationNotExist
	}

	return err
}


func (p *Postgres) ControlUserResponsibilityForTenderByBidID(ctx context.Context, bidID, username string) error {
	var isRelated bool

	err := p.DB.QueryRow(ctx, `
    SELECT EXISTS (
		SELECT 1
			FROM tender t
			JOIN organization_responsible orr ON orr.organization_id = t.organization_id
			JOIN employee e ON e.id = orr.user_id
			WHERE t.id = (
				SELECT tender_id
				FROM bid
				WHERE id = $2
		)
		AND e.username = $1
	) AS is_related;`, username, bidID).Scan(&isRelated)

	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrUserNotExist
	}

	if !isRelated {
		return repository.ErrRelationNotExist
	}

	return err
}

