package common

import "strings"

type AnnotationMapper struct {
	anno    string
	annoMap map[string][]string
	useMap  bool
}

func NewAnnotationMapper(annotation string) *AnnotationMapper {
	annoMap := make(map[string][]string)

	// populate annoMap
	if strings.Contains(annotation, ":") {
		for _, annoToken := range strings.Split(annotation, ",") {
			kv := strings.SplitN(annoToken, ":", 2)
			group := kv[0]
			annotations := strings.Split(kv[1], ";")
			annoMap[group] = annotations
		}
	}

	return &AnnotationMapper{
		anno:    annotation,
		annoMap: annoMap,
		useMap:  len(annoMap) > 0,
	}
}

func (m *AnnotationMapper) GetAnnotations(group string) []string {
	if m.useMap {
		if annotations, ok := m.annoMap[group]; ok {
			return annotations
		}
		return []string{}
	} else {
		return []string{m.anno}
	}
}
