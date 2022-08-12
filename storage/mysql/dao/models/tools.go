package models

import "time"

func OrNil64(u uint64) *uint64 {
	if u == 0 {
		return nil
	}
	return &u
}

func OrNil32(u uint32) *uint32 {
	if u == 0 {
		return nil
	}
	return &u
}

func OrNil16(u uint16) *uint16 {
	if u == 0 {
		return nil
	}
	return &u
}

func OrNil8(u uint8) *uint8 {
	if u == 0 {
		return nil
	}
	return &u
}

func Val64(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

func Val32(v *uint32) uint32 {
	if v == nil {
		return 0
	}
	return *v
}

func Val16(v *uint16) uint16 {
	if v == nil {
		return 0
	}
	return *v
}

func Val8(v *uint8) uint8 {
	if v == nil {
		return 0
	}
	return *v
}

// func OrNil[V uint8 | uint64 | uint32 | uint16](u V) *V {
// 	if u == 0 {
// 		return nil
// 	}
// 	return &u
// }

// func Val[V uint8 | uint64 | uint32 | uint16](v *V) V {
// 	if v == nil {
// 		return 0
// 	}
// 	return *v
// }

func TimeOrNil(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func Time(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
