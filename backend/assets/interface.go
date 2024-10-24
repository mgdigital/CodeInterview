package assets

import "context"

type Params struct {
	AssetID   string
	Limit     int
	Offset    int
	Query     string
	SortOrder SortOrder
	SortDesc  bool
}

type Result struct {
	Assets     []Asset `json:"assets"`
	TotalCount int     `json:"totalCount"`
}

type Lookup interface {
	Lookup(ctx context.Context, params Params) (Result, error)
}
