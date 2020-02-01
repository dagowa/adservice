package advert_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dagowa/adservice/internal/models/advert"

	"github.com/dagowa/adservice/internal/store"
)

type adverts struct {
	Adverts []advert.Advert `json:"adverts"`
}

func TestAddBatch(t *testing.T) {
	pconn, err := store.NewPSQLConnection("")
	if err != nil {
		t.Fatalf("Cannot init psql conenction; err: %v", err)
	}
	ppool := pconn.Pool

	jsonFile, err := os.Open("../../../test/source_adverts.json")
	if err != nil {
		t.Fatalf("Cannot open the source file; err: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var advList adverts
	if err := json.Unmarshal(byteValue, &advList); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if err := advert.AddBatch(ppool, &(advList.Adverts)); err != nil {
		t.Fatalf("Cannot add batch; err: %v", err)
	}
}

func TestInsert(t *testing.T) {
	pconn, err := store.NewPSQLConnection("")
	if err != nil {
		t.Fatalf("Cannot init psql conenction; err: %v", err)
	}
	ppool := pconn.Pool
	keks := "Чтобы был пушистенький и мягОнький :3 примерно как на фото"
	a := advert.Advert{
		Title:       "Куплю кота",
		Price:       5000,
		Description: &keks,
		Gallery: &[]advert.Photo{
			{
				Index: 0,
				Link:  "",
			},
			{
				Index: 1,
				Link:  "",
			},
			{
				Index: 2,
				Link:  "",
			},
		},
	}

	id, err := a.Add(ppool)
	if err != nil || id == 0 {
		t.Fatalf("Cannot add new advert; err: %v", err)
	}
	if err := advert.Delete(ppool, id); err != nil {
		t.Fatalf("Cannot delete test row; err: %v", err)
	}
}
