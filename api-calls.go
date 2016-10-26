package rcAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Me returns the currently logged in recurser
func (t *Auth) Me() Recurser {
	me := t.Recurser("me")
	return me
}

// Recurser returns any recurser given their email address or id number
func (t *Auth) Recurser(id string) Recurser {
	url := t.BaseURL + t.RecurserPath + id + t.TokenParam
	res := makeRequest(url)
	recurser := Recurser{}
	err := json.Unmarshal(res, &recurser)
	if err != nil {
		fmt.Println(err)
	}
	return recurser
}

// Batchlist returns a list of all the batches
func (t *Auth) BatchList() []Batch {
	url := t.BaseURL + t.BatchPath + t.TokenParam
	res := makeRequest(url)
	var batchList []Batch
	err := json.Unmarshal(res, &batchList)
	if err != nil {
		fmt.Println(err)
	}
	return batchList
}

// Batch returns the current batch
func (t *Auth) Batch(id string) Batch {
	url := t.BaseURL + t.BatchPath + id + t.TokenParam
	res := makeRequest(url)
	batch := Batch{}
	err := json.Unmarshal(res, &batch)
	if err != nil {
		fmt.Println(err)
	}
	return batch
}

// BatchMembers returns a list of the members of a given batch
func (t *Auth) BatchMembers(id string) []Recurser {
	url := t.BaseURL + t.BatchPath + id + "/people" + t.TokenParam
	res := makeRequest(url)
	var batchMembers []Recurser
	err := json.Unmarshal(res, &batchMembers)
	if err != nil {
		fmt.Println(err)
	}
	return batchMembers
}

func makeRequest(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return body
}
