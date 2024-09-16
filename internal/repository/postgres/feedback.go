package postgres

import (
	"context"

	"github.com/DarRo9/Service_for_creating_and_processing_tenders/internal/repository"
	"github.com/DarRo9/Service_for_creating_and_processing_tenders/models"
)

func (p *Postgres) ApplyBidFeedback(ctx context.Context, bidID string, feedback *models.BidFeedback) error {
	pgCmd, err := p.DB.Exec(ctx, `
	INSERT INTO bid_feedback 
		(bid_id, description) 
    VALUES ($1, $2) `, bidID, feedback)
	if pgCmd.RowsAffected() == 0 {
		return repository.ErrBidNotFound
	}

	return err
}

func (p *Postgres) GetCommentsOfBid(ctx context.Context, tenderID, authorUsername string, limit, offset int32) ([]*models.BidReviewResponse, error) {
	rows, err := p.DB.Query(ctx, `
	SELECT bf.*
		FROM bid_feedback bf
		JOIN bid b ON bf.bid_id = b.id
		WHERE b.tender_id = $1
		AND EXISTS (
			SELECT 1
			FROM employee e
				WHERE e.id = b.author_id
				AND e.username = $2  
		ORDER BY created_at ASC
		LIMIT $3
		OFFSET $4);`, tenderID, authorUsername, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.BidReviewResponse
	for rows.Next() {
		var review models.BidReviewResponse
		if err := rows.Scan(&review.ID, &review.BidID, &review.Description, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}

	if len(reviews) == 0 {
		return nil, repository.ErrBidReviewsNotFound
	}

	return reviews, err
}
