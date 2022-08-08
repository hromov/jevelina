package files

import "time"

type File struct {
	ID        uint64
	ParentID  uint64
	CreatedAt time.Time
	Name      string
	URL       string
}

type FileCreateReq struct {
	ParentID uint64
	Name     string
	URL      string
}

type FileAddReq struct {
	Parent uint64
	Name   string
	Type   string
	Value  string
}
