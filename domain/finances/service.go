package finances

import "context"

type Repository interface {
	CreateWallet(ctx context.Context, w Wallet) (Wallet, error)
	ChangeWalletName(ctx context.Context, id uint16, name string) error
	ChangeWalletState(ctx context.Context, id uint16, closed bool) error
	DeleteWallet(ctx context.Context, id uint16) error
	ListWallets(ctx context.Context) ([]Wallet, error)
}

type Service interface {
	CreateWallet(ctx context.Context, w Wallet) (Wallet, error)
	ChangeWalletName(ctx context.Context, id uint16, name string) error
	ChangeWalletState(ctx context.Context, id uint16, closed bool) error
	DeleteWallet(ctx context.Context, id uint16) error
	ListWallets(ctx context.Context) ([]Wallet, error)
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
