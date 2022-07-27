package wisdom

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
)

type Set struct {
	Quotes []string `json:"quotes"`
}

func NewSet(path string) (*Set, error) {
	file, _ := ioutil.ReadFile(path)

	var data Set
	_ = json.Unmarshal([]byte(file), &data)

	var quotes *Set
	if err := json.Unmarshal(file, &quotes); err != nil {
		return nil, err
	}
	return &data, nil
}

func (b *Set) GetRandQuote() string {
	i := rand.Intn(len(b.Quotes))
	return b.Quotes[i]
}
