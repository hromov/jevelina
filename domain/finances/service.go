package finances

import (
	"context"
	"errors"
)

type Repository interface {
	CreateWallet(ctx context.Context, w Wallet) (Wallet, error)
	ChangeWalletName(ctx context.Context, id uint16, name string) error
	ChangeWalletState(ctx context.Context, id uint16, closed bool) error
	DeleteWallet(ctx context.Context, id uint16) error
	ListWallets(ctx context.Context) ([]Wallet, error)
	CreateTransfer(ctx context.Context, t Transfer) (Transfer, error)
	UpdateTransfer(ctx context.Context, t Transfer) error
	CompleteTransfer(ctx context.Context, id, userID uint64) error
	GetTransfer(ctx context.Context, id uint64) (Transfer, error)
	DeleteTransfer(ctx context.Context, id, userID uint64) error
	Transfers(ctx context.Context, filter Filter) (TransfersResponse, error)
	TransferCategories(ctx context.Context) ([]string, error)
	SumByCategory(ctx context.Context, filter Filter) (CategorisedCashflow, error)
}

type Service interface {
	CreateWallet(ctx context.Context, w Wallet) (Wallet, error)
	ChangeWalletName(ctx context.Context, id uint16, name string) error
	ChangeWalletState(ctx context.Context, id uint16, closed bool) error
	DeleteWallet(ctx context.Context, id uint16) error
	ListWallets(ctx context.Context) ([]Wallet, error)
	CreateTransfer(ctx context.Context, t Transfer) (Transfer, error)
	UpdateTransfer(ctx context.Context, t Transfer) error
	CompleteTransfer(ctx context.Context, id uint64, userID uint64) error
	DeleteTransfer(ctx context.Context, id uint64, userID uint64) error
	GetTransfer(ctx context.Context, id uint64) (Transfer, error)
	Transfers(ctx context.Context, filter Filter) (TransfersResponse, error)
	TransferCategories(ctx context.Context) ([]string, error)
	SumByCategory(ctx context.Context, filter Filter) (CategorisedCashflow, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) CreateWallet(ctx context.Context, w Wallet) (Wallet, error) {
	return s.r.CreateWallet(ctx, w)
}
func (s *service) ChangeWalletName(ctx context.Context, id uint16, name string) error {
	return s.r.ChangeWalletName(ctx, id, name)
}
func (s *service) ChangeWalletState(ctx context.Context, id uint16, closed bool) error {
	return s.r.ChangeWalletState(ctx, id, closed)
}
func (s *service) DeleteWallet(ctx context.Context, id uint16) error {
	return s.r.DeleteWallet(ctx, id)
}
func (s *service) ListWallets(ctx context.Context) ([]Wallet, error) {
	return s.r.ListWallets(ctx)
}

func (s *service) GetTransfer(ctx context.Context, id uint64) (Transfer, error) {
	return s.r.GetTransfer(ctx, id)
}

func (s *service) CreateTransfer(ctx context.Context, t Transfer) (Transfer, error) {
	t.Completed = false
	if err := valid(t); err != nil {
		return Transfer{}, err
	}
	return s.r.CreateTransfer(ctx, t)
}
func (s *service) UpdateTransfer(ctx context.Context, t Transfer) error {
	if err := valid(t); err != nil {
		return err
	}

	oldTransfer, err := s.GetTransfer(ctx, t.ID)
	if err != nil {
		return err
	}

	if oldTransfer.Completed || !oldTransfer.DeletedAt.IsZero() {
		return s.r.UpdateTransfer(ctx, Transfer{ID: t.ID, Description: t.Description, Category: t.Category})
	}
	return s.r.UpdateTransfer(ctx, t)
}

func (s *service) CompleteTransfer(ctx context.Context, id uint64, userID uint64) error {
	return s.r.CompleteTransfer(ctx, id, userID)
}
func (s *service) DeleteTransfer(ctx context.Context, id uint64, userID uint64) error {
	return s.r.DeleteTransfer(ctx, id, userID)
}

func (s *service) Transfers(ctx context.Context, filter Filter) (TransfersResponse, error) {
	return s.r.Transfers(ctx, filter)
}
func (s *service) TransferCategories(ctx context.Context) ([]string, error) {
	return s.r.TransferCategories(ctx)
}
func (s *service) SumByCategory(ctx context.Context, filter Filter) (CategorisedCashflow, error) {
	return s.r.SumByCategory(ctx, filter)
}

func valid(t Transfer) error {
	if t.From == t.To {
		return errors.New("Can't transfer to the same wallet")
	}
	return nil
}
