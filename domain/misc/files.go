package misc

import "time"

type File struct {
	ID        uint64
	ParentID  uint64
	CreatedAt time.Time
	Name      string
	URL       string
}
