package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/silvergama/studying-golang/elton-minetto-course/core/beer"
	"github.com/stretchr/testify/assert"
)

type BeerServiceMock struct{}

func (t BeerServiceMock) GetAll() ([]*beer.Beer, error) {
	return []*beer.Beer{
		{
			ID:    10,
			Name:  "Heineken",
			Type:  beer.TypeLager,
			Style: beer.StylePale,
		},
		{
			ID:    20,
			Name:  "Skol",
			Type:  beer.TypeLager,
			Style: beer.StylePale,
		},
	}, nil
}

func (t BeerServiceMock) Get(ID int64) (*beer.Beer, error) {
	return &beer.Beer{
		ID:    10,
		Name:  "Heineken",
		Type:  beer.TypeLager,
		Style: beer.StylePale,
	}, nil
}

func (t BeerServiceMock) Store(b *beer.Beer) error {
	return nil
}

func (t BeerServiceMock) Update(b *beer.Beer) error {
	return nil
}

func (t BeerServiceMock) Remove(b int64) error {
	return nil
}

func Test_getAllBeer(t *testing.T) {
	service := &BeerServiceMock{}
	handler := getAllBeer(service)
	req, err := http.NewRequest("GET", "/v1/beer", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var result []*beer.Beer
	err = json.NewDecoder(rr.Body).Decode(&result)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(10), result[0].ID)
	assert.Equal(t, int64(20), result[1].ID)

}

func Test_getBeer(t *testing.T) {
	service := &BeerServiceMock{}
	handler := getBeer(service)
	r := mux.NewRouter()
	r.Handle("/v1/beer/{id}", handler)

	req, err := http.NewRequest("GET", "/v1/beer/10", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var result *beer.Beer
	err = json.NewDecoder(rr.Body).Decode(&result)
	assert.Nil(t, err)
	assert.Equal(t, int64(10), result.ID)

}

func Test_storeBeer(t *testing.T) {
	handler := storeBeer(&BeerServiceMock{})
	r := mux.NewRouter()
	r.Handle("/v1/beer", handler)

	payload := []byte(`{"ID": 1, "Name": "Heineken", "Type": 1, "Style": 1}`)
	req, err := http.NewRequest("POST", "/v1/beer", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
}
