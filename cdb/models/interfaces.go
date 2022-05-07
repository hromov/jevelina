package models

type FilesService interface {
	Delete(ID uint64) error
	Upload(req *FileAddReq) (*File, error)
	GetUrl(ID uint64) (string, error)
	List(filter ListFilter) ([]*File, error)
}
