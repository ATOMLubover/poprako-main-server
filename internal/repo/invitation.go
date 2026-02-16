package repo

import "poprako-main-server/internal/model/po"

type InvitationRepo interface {
	RetrieveInvitations(ex Executor) ([]po.BasicInvitation, error)
	GetInvitationByQQ(ex Executor, inviteeQQ string) (*po.BasicInvitation, error)

	CreateInvitations(ex Executor, newInvitation *po.NewInvitation) error
	MarkInvitationAsUsed(ex Executor, invitationID string) error
}

type invitationRepo struct {
	ex Executor
}

func NewInvitationRepo(ex Executor) InvitationRepo {
	return &invitationRepo{
		ex: ex,
	}
}

func (ir *invitationRepo) Exec() Executor { return ir.ex }

func (ir *invitationRepo) withTrx(tx Executor) Executor {
	if tx != nil {
		return tx
	}

	return ir.ex
}

func (ir *invitationRepo) RetrieveInvitations(
	ex Executor,
) (
	[]po.BasicInvitation,
	error,
) {
	ex = ir.withTrx(ex)

	var invitations []po.BasicInvitation

	if err := ex.
		Where("pending = ?", true).
		Find(&invitations).
		Error; err != nil {
		return invitations, err
	}

	return invitations, nil
}

func (ir *invitationRepo) GetInvitationByQQ(
	ex Executor,
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
	ex Executor,
	newInvitation *po.NewInvitation,
) error {
	ex = ir.withTrx(ex)

	return ex.Create(newInvitation).Error
}

func (ir *invitationRepo) MarkInvitationAsUsed(
	ex Executor,
	invitationID string,
) error {
	ex = ir.withTrx(ex)

	return ex.
		Model(&po.BasicInvitation{}).
		Where("id = ?", invitationID).
		Update("pending", false).
		Error
}
