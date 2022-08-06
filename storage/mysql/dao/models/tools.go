package models

func OrNil[V uint8 | uint64](u V) *V {
	if u == 0 {
		return nil
	}
	return &u
}
