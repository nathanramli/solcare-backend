package gorm

import (
	"context"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
)

type proposalRepo struct {
	db *gorm.DB
}

func NewProposalRepo(db *gorm.DB) repositories.ProposalRepo {
	return &proposalRepo{
		db: db,
	}
}

func (r *proposalRepo) SaveProposal(ctx context.Context, proposal *models.Proposal) error {
	return r.db.WithContext(ctx).Save(proposal).Error
}

func (r *proposalRepo) FindProposalByAddress(ctx context.Context, address string) (*models.Proposal, error) {
	proposal := new(models.Proposal)
	err := r.db.WithContext(ctx).Where("address = ?", address).Take(proposal).Error
	return proposal, err
}
