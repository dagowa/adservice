package advertmanager

// SortCriteria contains
type SortCriteria struct {
	priceAsc bool
	dateAsv  bool
}

// SortOrder is ...
type SortOrder string

// SortType
const (
	SortOrderASC  SortOrder = "ascending"
	SortOrderDESC SortOrder = "descending"
)

func (sc *SortCriteria) SetPriceSortOrder(st SortOrder) {
	if st == SortOrderASC {
		sc.priceAsc = true
	}
	if st == SortOrderDESC {
		sc.priceAsc = false
	}
}

func (sc *SortCriteria) SetDateSortOrder(st SortOrder) {
	if st == SortOrderASC {
		sc.priceAsc = true
	}
	if st == SortOrderDESC {
		sc.priceAsc = false
	}
}
