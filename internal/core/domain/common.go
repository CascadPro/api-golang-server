package domain

type SortType string

const (
	SortTypeNewest  = SortType("newest")
	SortTypeOldest  = SortType("oldest")
	SortTypePopular = SortType("popular")
	SortTypeNil     = SortType("")
)

var (
	SortTypes = []SortType{SortTypeNewest, SortTypeOldest, SortTypePopular}
)
