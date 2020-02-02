package advert

type Advert struct {
	ID          int      `json:"id,omitempty"`
	Title       string   `json:"title"`
	Price       int      `json:"price"`
	Date        string   `json:"date,omitempty"`
	Description *string  `json:"description,omitepmty"`
	Gallery     *[]Photo `json:"gallery"`
}

type Photo struct {
	Index int    `json:"index"`
	Link  string `json:"photo"`
}
