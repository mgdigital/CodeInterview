package assets

import (
	"context"
	"database/sql"
	"strings"
)

type lookup struct {
	db *sql.DB
}

const defaultLimit = 10

func (l lookup) Lookup(ctx context.Context, params Params) (Result, error) {
	assetSqlParts := []string{
		"SELECT id, host, comment, owner FROM assets",
	}
	var assetSqlArgs []any

	var conditions []string
	var conditionArgs []any
	if params.AssetID != "" {
		conditions = append(conditions, "id = ?")
		conditionArgs = append(conditionArgs, params.AssetID)
	}
	if params.Query != "" {
		conditions = append(conditions, "host LIKE concat('%', ?, '%')")
		conditionArgs = append(conditionArgs, params.Query)
	}

	if len(conditions) > 0 {
		assetSqlParts = append(assetSqlParts, "WHERE "+strings.Join(conditions, " AND "))
		assetSqlArgs = append(assetSqlArgs, conditionArgs...)
	}

	direction := "ASC"
	if params.SortDesc {
		direction = "DESC"
	}
	assetSqlParts = append(assetSqlParts, "ORDER BY "+string(params.SortOrder)+" "+direction+" LIMIT ? OFFSET ?")
	limit := params.Limit
	if limit == 0 {
		limit = defaultLimit
	}
	assetSqlArgs = append(assetSqlArgs, limit, params.Offset)
	assetRows, err := l.db.QueryContext(ctx, strings.Join(assetSqlParts, " "), assetSqlArgs...)
	if err != nil {
		return Result{}, err
	}
	defer func(assetRows *sql.Rows) {
		_ = assetRows.Close()
	}(assetRows)
	assetMap := make(map[int]*Asset, limit)
	assetIds := make([]int, 0, limit)
	assetIdsAny := make([]any, 0, limit)

	var id int
	var host, comment, owner string

	for assetRows.Next() {
		if err := assetRows.Scan(&id, &host, &comment, &owner); err != nil {
			return Result{}, err
		}
		assetMap[id] = newAsset(id, host, comment, owner)
		assetIds = append(assetIds, id)
		assetIdsAny = append(assetIdsAny, id)
	}

	totalCount := params.Offset + len(assetIds)
	if len(assetIds) >= limit {
		countSql := "SELECT COUNT(*) FROM assets"
		var countSqlArgs []any
		if len(conditions) > 0 {
			countSql += " WHERE " + strings.Join(conditions, " AND ")
			countSqlArgs = append(countSqlArgs, conditionArgs...)
		}
		countRow := l.db.QueryRowContext(ctx, countSql, countSqlArgs...)
		if err := countRow.Scan(&totalCount); err != nil {
			return Result{}, err
		}
	}

	if len(assetIds) == 0 {
		return Result{
			Assets:     []Asset{},
			TotalCount: totalCount,
		}, nil
	}

	portRows, err := l.db.QueryContext(ctx, "SELECT asset_id, port FROM ports WHERE asset_id IN (?"+strings.Repeat(",?", len(assetIds)-1)+") ORDER BY port", assetIdsAny...)
	if err != nil {
		return Result{}, err
	}
	defer func(portRows *sql.Rows) {
		_ = portRows.Close()
	}(portRows)

	var assetId, port int

	for portRows.Next() {
		if err := portRows.Scan(&assetId, &port); err != nil {
			return Result{}, err
		}
		assetMap[assetId].addPort(port)
	}

	ipRows, err := l.db.QueryContext(ctx, "SELECT asset_id, address FROM ips WHERE asset_id IN (?"+strings.Repeat(",?", len(assetIds)-1)+")", assetIdsAny...)
	if err != nil {
		return Result{}, err
	}
	defer func(ipRows *sql.Rows) {
		_ = ipRows.Close()
	}(ipRows)

	var address string
	for ipRows.Next() {
		if err := ipRows.Scan(&assetId, &address); err != nil {
			return Result{}, err
		}
		assetMap[assetId].addIP(address)
	}

	assets := make([]Asset, 0, len(assetIds))
	for _, id := range assetIds {
		assets = append(assets, *assetMap[id])
	}
	return Result{
		Assets:     assets,
		TotalCount: totalCount,
	}, nil
}
