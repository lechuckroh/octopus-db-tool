package common

import "strings"

type PrefixMapper struct {
	prefix    string
	prefixMap map[string]string
	useMap    bool
}

func NewPrefixMapper(prefix string) *PrefixMapper {
	prefixMap := make(map[string]string)

	// populate prefixMap
	if strings.Contains(prefix, ":") {
		for _, prefixToken := range strings.Split(prefix, ",") {
			kv := strings.SplitN(prefixToken, ":", 2)
			group := kv[0]
			prefixValue := kv[1]
			prefixMap[group] = prefixValue
		}
	}

	return &PrefixMapper{
		prefix:    prefix,
		prefixMap: prefixMap,
		useMap:    len(prefixMap) > 0,
	}
}

func (p *PrefixMapper) GetPrefix(group string) string {
	if p.useMap {
		if prefix, ok := p.prefixMap[group]; ok {
			return prefix
		}
		return ""
	} else {
		return p.prefix
	}
}
