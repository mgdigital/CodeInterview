package assets

import (
	"database/sql"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"time"
)

func NewLookup(fileName string) (Lookup, error) {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}
	return lookupCache{
		baseLookup: lookup{
			db: db,
		},
		lru: expirable.NewLRU[Params, Result](1000, nil, time.Minute*10),
	}, nil
}
