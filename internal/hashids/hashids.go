package hashids

import "github.com/speps/go-hashids/v2"

var Client = New()

func New() *hashids.HashID {
	hashID, _ := hashids.New()
	return hashID
}
