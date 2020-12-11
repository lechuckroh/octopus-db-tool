package octopus

import (
	"strings"
)

type TableFilterFn func(*Table) bool

// getTableFilterFn returns TableFilterFn using 'goups' option.
// 'groups' is collection of group name separater by comma(,).
func GetTableFilterFn(groups string) TableFilterFn {
	if groups == "" {
		return nil
	}

	groupSlice := strings.Split(groups, ",")
	return func(table *Table) bool {
		for _, group := range groupSlice {
			if table.Group == group {
				return true
			}
		}
		return false
	}
}
