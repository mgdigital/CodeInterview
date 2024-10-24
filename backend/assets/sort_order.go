package assets

type SortOrder string

const (
	SortOrderID    SortOrder = "id"
	SortOrderHost  SortOrder = "host"
	SortOrderOwner SortOrder = "owner"
)

func ParseSortOrder(s string) (SortOrder, bool) {
	switch s {
	case "id":
		return SortOrderID, true
	case "host":
		return SortOrderHost, true
	case "owner":
		return SortOrderOwner, true
	}
	return SortOrderHost, false
}
