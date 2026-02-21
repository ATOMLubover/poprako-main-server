package repo

import "poprako-main-server/internal/model/po"

type InvitationRepo interface {
	RetrieveInvitations(ex Exct) ([]po.BasicInvitation, error)
	GetInvitationByQQ(ex Exct, inviteeQQ string) (*po.BasicInvitation, error)

	CreateInvitations(ex Exct, newInvitation *po.NewInvitation) error
	MarkInvitationAsUsed(ex Exct, invitationID string) error
}

type invitationRepo struct {
	ex Exct
}

func NewInvitationRepo(ex Exct) InvitationRepo {
	return &invitationRepo{
		ex: ex,
	}
}

func (ir *invitationRepo) Exec() Exct { return ir.ex }

func (ir *invitationRepo) withTrx(tx Exct) Exct {
	if tx != nil {
		return tx
	}

	return ir.ex
}

func (ir *invitationRepo) RetrieveInvitations(
	ex Exct,
) (
	[]po.BasicInvitation,
	error,
) {
	ex = ir.withTrx(ex)

	var invitations []po.BasicInvitation

	if err := ex.
		Where("pending = ?", true).
		Order("created_at DESC").
		Find(&invitations).
		Error; err != nil {
		return invitations, err
	}

	return invitations, nil
}

func (ir *invitationRepo) GetInvitationByQQ(
	ex Exct,
	inviteeQQ string,
) (
	*po.BasicInvitation,
	error,
) {
	ex = ir.withTrx(ex)

	var invitation po.BasicInvitation

	result := ex.
		Where("invitee_qq = ? AND pending = ?", inviteeQQ, true).
		First(&invitation)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, REC_NOT_FOUND
	}

	return &invitation, nil
}

func (ir *invitationRepo) CreateInvitations(
	ex Exct,
	newInvitation *po.NewInvitation,
) error {
	ex = ir.withTrx(ex)

	return ex.Create(newInvitation).Error
}

func (ir *invitationRepo) MarkInvitationAsUsed(
	ex Exct,
	invitationID string,
) error {
	ex = ir.withTrx(ex)

	return ex.
		Model(&po.BasicInvitation{}).
		Where("id = ?", invitationID).
		Update("pending", false).
		Error
}
