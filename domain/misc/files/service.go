package files

import "context"

type Repository interface {
	CreateFile(ctx context.Context, f FileCreateReq) (File, error)
	DeleteFile(ctx context.Context, id uint64) error
	GetFile(ctx context.Context, id uint64) (File, error)
	GetFilesByParent(ctx context.Context, parentID uint64) ([]File, error)
}

//go:generate mockery --name Storage --filename FilesStorageService.go --structname FilesStorageService --output ../../../mocks
type Storage interface {
	Delete(ctx context.Context, url string) error
	Upload(ctx context.Context, req FileAddReq) (FileCreateReq, error)
	PresignUrl(url string) (string, error)
}

//go:generate mockery --name Service --filename FilesService.go --structname FilesService --output ../../../mocks
type Service interface {
	Upload(ctx context.Context, req FileAddReq) (File, error)
	Delete(ctx context.Context, id uint64) error
	GetUrl(ctx context.Context, id uint64) (string, error)
	GetByParent(ctx context.Context, parentID uint64) ([]File, error)
}

type service struct {
	r Repository
	s Storage
}

func NewService(r Repository, s Storage) *service {
	return &service{r, s}
}

func (s *service) Upload(ctx context.Context, req FileAddReq) (File, error) {
	createRequest, err := s.s.Upload(ctx, req)
	if err != nil {
		return File{}, err
	}
	return s.r.CreateFile(ctx, createRequest)
}

func (s *service) Delete(ctx context.Context, id uint64) error {
	f, err := s.r.GetFile(ctx, id)
	if err != nil {
		return err
	}
	if err := s.r.DeleteFile(ctx, f.ID); err != nil {
		return err
	}
	return s.s.Delete(ctx, f.URL)
}

func (s *service) GetUrl(ctx context.Context, id uint64) (string, error) {
	f, err := s.r.GetFile(ctx, id)
	if err != nil {
		return "", err
	}
	return s.s.PresignUrl(f.URL)
}

func (s *service) GetByParent(ctx context.Context, parentID uint64) ([]File, error) {
	return s.r.GetFilesByParent(ctx, parentID)
}
