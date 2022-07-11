package grpctest

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Collections []*collectionItem

type collectionItem struct {
	Name   string   `json:"name"`
	Url    string   `json:"url"`
	Proto  string   `json:"proto"`
	Call   string   `json:"call"`
	Expect []string `json:"expect"`
}

func NewCollections(f string, patch func(*collectionItem)) (Collections, error) {
	collectionFile, err := os.Open("./testdata/grpc_collection.json")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(collectionFile)
	if err != nil {
		return nil, err
	}

	var res Collections
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	if patch != nil {
		for _, item := range res {
			patch(item)
		}
	}

	return res, nil
}
