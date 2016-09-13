package rc_api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (t *RCAuth) Me() Recurser {
	me := t.Recurser("me")
	return me
}

func (t *RCAuth) Recurser(id string) Recurser {
	url := t.BaseUrl + t.RecurserPath + id + t.TokenParam
	res := makeRequest(url)
	recurser := Recurser{}
	err := json.Unmarshal(res, &recurser)
	if err != nil {
		log.Fatal(err)
	}
	return recurser
}

func (t *RCAuth) BatchList() []Batch {
	url := t.BaseUrl + t.BatchPath + t.TokenParam
	res := makeRequest(url)
	var batchList []Batch
	err := json.Unmarshal(res, &batchList)
	if err != nil {
		log.Fatal(err)
	}
	return batchList
}

func (t *RCAuth) Batch(id string) Batch {
	url := t.BaseUrl + t.BatchPath + id + t.TokenParam
	res := makeRequest(url)
	batch := Batch{}
	err := json.Unmarshal(res, &batch)
	if err != nil {
		log.Fatal(err)
	}
	return batch
}

func (t *RCAuth) BatchMembers(id string) []Recurser {
	url := t.BaseUrl + t.BatchPath + id + "/people" + t.TokenParam
	res := makeRequest(url)
	var batchMembers []Recurser
	err := json.Unmarshal(res, &batchMembers)
	if err != nil {
		log.Fatal(err)
	}
	return batchMembers
}

func makeRequest(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
