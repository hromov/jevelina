package models

import "time"

func OrNil[V uint8 | uint64 | uint32 | uint16](u V) *V {
	if u == 0 {
		return nil
	}
	return &u
}

func TimeOrNil(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}
