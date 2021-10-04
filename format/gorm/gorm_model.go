package gorm

import (
	"github.com/lechuckroh/octopus-db-tools/util"
	"regexp"
	"strings"
)

type EmbeddedModel struct {
	ColumnNames []string
}

func (m EmbeddedModel) Contains(columnName string) bool {
	for _, colName := range m.ColumnNames {
		if colName == columnName {
			return true
		}
	}
	return false
}

func parseEmbeddedModelDefinition(definition string) (string, []string) {
	r := regexp.MustCompile(`^([a-zA-Z0-9._]+):(.*)$`)
	groups, ok := util.MatchRegexGroups(r, definition)
	if !ok {
		return "", nil
	}

	name := groups[0]
	concatColumns := groups[1]
	var columnNames []string
	if concatColumns == "" {
		columnNames = []string{}
	} else {
		columnNames = strings.Split(concatColumns, ",")
	}

	return name, columnNames
}

type EmbeddedModels map[string]EmbeddedModel

func (m EmbeddedModels) ContainsColumnName(columnName string) bool {
	for _, model := range m {
		if model.Contains(columnName) {
			return true
		}
	}
	return false
}

func (m EmbeddedModels) Names() []string {
	var names []string
	for name, _ := range m {
		names = append(names, name)
	}
	return names
}

var embeddedModels EmbeddedModels

func registerEmbeddedModel(name string, columnNames []string) {
	if embeddedModels == nil {
		embeddedModels = make(EmbeddedModels)
	}
	embeddedModels[name] = EmbeddedModel{
		ColumnNames: columnNames,
	}
}

// registerEmbeddedModel registers default embedded struct definitions
func registerDefaultEmbeddedModels() {
	registerEmbeddedModel("gorm.Model", []string{"id", "created_at", "updated_at", "deleted_at"})
}

// extractEmbeddedModels extracts matching EmbeddedModels
func extractEmbeddedModels(fields []*GoField) EmbeddedModels {
	// create set of column names
	columnNameSet := util.NewStringSet()
	for _, field := range fields {
		columnNameSet.Add(field.Column.Name)
	}

	result := EmbeddedModels{}
	for name, model := range embeddedModels {
		if columnNameSet.ContainsAll(model.ColumnNames) {
			result[name] = model
			columnNameSet.RemoveAll(model.ColumnNames)
		}
	}

	return result
}
