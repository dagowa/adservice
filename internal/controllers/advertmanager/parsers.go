package advertmanager

import (
	"net/http"

	"github.com/dagowa/adservice/internal/models/advert"
	"github.com/dagowa/adservice/internal/models/page"
)

func (am *AdvertManager) ParseAdditionalFileds(r *http.Request) *AdditionalFileds {
	additionalFileds := r.Context().Value("additional_fileds").(*AdditionalFileds)
	return additionalFileds
}

func (am *AdvertManager) ParsePage(r *http.Request) *page.Page {
	page := r.Context().Value("page").(*page.Page)
	return page
}

func (am *AdvertManager) ParseAdvert(r *http.Request) *advert.Advert {
	advert := r.Context().Value("advert").(*advert.Advert)
	return advert
}
