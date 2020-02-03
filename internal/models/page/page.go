package page

type Page struct {
	Numb     int
	Size     int
	PriceAsc bool
	DateAsc  bool
}

const (
	SortOrderASC  string = "ascending"
	SortOrderDESC string = "descending"
)
