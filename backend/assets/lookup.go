package assets

import (
	"context"
	"database/sql"
	"strings"
)

type lookup struct {
	db *sql.DB
}

func (l lookup) Lookup(ctx context.Context, params Params) (Result, error) {
	lc := &requestLifecycle{
		lookup:  l,
		Context: ctx,
		Params:  params,
	}
	if err := lc.queryAssets(); err != nil {
		return Result{}, err
	}
	if err := lc.queryTotalCount(); err != nil {
		return Result{}, err
	}
	if err := lc.queryPorts(); err != nil {
		return Result{}, err
	}
	if err := lc.queryIPs(); err != nil {
		return Result{}, err
	}
	return lc.createResult(), nil
}

const defaultLimit = 10

type requestLifecycle struct {
	lookup
	context.Context
	Params
	limit         int
	totalCount    int
	conditions    []string
	conditionArgs []any
	assetMap      map[int]*Asset
	assetIds      []int
	assetIdsAny   []any
}

func (lc *requestLifecycle) queryAssets() error {
	sqlParts := []string{
		"SELECT id, host, comment, owner FROM assets",
	}
	var sqlArgs []any
	if lc.Params.AssetID != "" {
		lc.conditions = append(lc.conditions, "id = ?")
		lc.conditionArgs = append(lc.conditionArgs, lc.Params.AssetID)
	}
	if lc.Params.Query != "" {
		lc.conditions = append(lc.conditions, "host LIKE concat('%', ?, '%')")
		lc.conditionArgs = append(lc.conditionArgs, lc.Params.Query)
	}
	if len(lc.conditions) > 0 {
		sqlParts = append(sqlParts, "WHERE "+strings.Join(lc.conditions, " AND "))
		sqlArgs = append(sqlArgs, lc.conditionArgs...)
	}
	direction := "ASC"
	if lc.Params.SortDesc {
		direction = "DESC"
	}
	sqlParts = append(sqlParts, "ORDER BY "+string(lc.Params.SortOrder)+" "+direction+" LIMIT ? OFFSET ?")
	lc.limit = lc.Params.Limit
	if lc.limit == 0 {
		lc.limit = defaultLimit
	}
	sqlArgs = append(sqlArgs, lc.limit, lc.Params.Offset)
	assetRows, err := lc.db.QueryContext(lc.Context, strings.Join(sqlParts, " "), sqlArgs...)
	if err != nil {
		return err
	}
	defer func(assetRows *sql.Rows) {
		_ = assetRows.Close()
	}(assetRows)
	lc.assetMap = make(map[int]*Asset, lc.limit)
	lc.assetIds = make([]int, 0, lc.limit)
	lc.assetIdsAny = make([]any, 0, lc.limit)
	var id int
	var host, comment, owner string
	for assetRows.Next() {
		if err := assetRows.Scan(&id, &host, &comment, &owner); err != nil {
			return err
		}
		lc.assetMap[id] = newAsset(id, host, comment, owner)
		lc.assetIds = append(lc.assetIds, id)
		lc.assetIdsAny = append(lc.assetIdsAny, id)
	}
	return nil
}

func (lc *requestLifecycle) queryTotalCount() error {
	var totalCount int
	if len(lc.assetIds) < lc.limit {
		totalCount = lc.Params.Offset + len(lc.assetIds)
	} else {
		countSql := "SELECT COUNT(*) FROM assets"
		var countSqlArgs []any
		if len(lc.conditions) > 0 {
			countSql += " WHERE " + strings.Join(lc.conditions, " AND ")
			countSqlArgs = append(countSqlArgs, lc.conditionArgs...)
		}
		countRow := lc.db.QueryRowContext(lc.Context, countSql, countSqlArgs...)
		if err := countRow.Scan(&totalCount); err != nil {
			return err
		}
	}
	lc.totalCount = totalCount
	return nil
}

func (lc *requestLifecycle) queryPorts() error {
	if len(lc.assetIds) == 0 {
		return nil
	}
	portRows, err := lc.db.QueryContext(lc.Context, "SELECT asset_id, port FROM ports WHERE asset_id IN (?"+strings.Repeat(",?", len(lc.assetIds)-1)+") ORDER BY port", lc.assetIdsAny...)
	if err != nil {
		return err
	}
	defer func(portRows *sql.Rows) {
		_ = portRows.Close()
	}(portRows)
	var assetId, port int
	for portRows.Next() {
		if err := portRows.Scan(&assetId, &port); err != nil {
			return err
		}
		lc.assetMap[assetId].addPort(port)
	}
	return nil
}

func (lc *requestLifecycle) queryIPs() error {
	if len(lc.assetIds) == 0 {
		return nil
	}
	ipRows, err := lc.db.QueryContext(lc.Context, "SELECT asset_id, address FROM ips WHERE asset_id IN (?"+strings.Repeat(",?", len(lc.assetIds)-1)+")", lc.assetIdsAny...)
	if err != nil {
		return err
	}
	defer func(ipRows *sql.Rows) {
		_ = ipRows.Close()
	}(ipRows)
	var (
		assetId int
		address string
	)
	for ipRows.Next() {
		if err := ipRows.Scan(&assetId, &address); err != nil {
			return err
		}
		lc.assetMap[assetId].addIP(address)
	}
	return nil
}

func (lc *requestLifecycle) createResult() Result {
	assets := make([]Asset, 0, len(lc.assetIds))
	for _, id := range lc.assetIds {
		assets = append(assets, *lc.assetMap[id])
	}
	return Result{
		Assets:     assets,
		TotalCount: lc.totalCount,
	}
}
