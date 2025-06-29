package redis

import (
	"strings"
)

// Ex: BuildKey("comboname", "1") => combo:comboname:1
func BuildKey(parts ...string) string {
	return "combo:" + strings.Join(parts, ":")
}
