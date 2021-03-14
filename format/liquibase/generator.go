package liquibase

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type LqYaml struct {
	DatabaseChangeLog []interface{} `yaml:"databaseChangeLog"`
}

func newLqYaml() *LqYaml {
	result := LqYaml{make([]interface{}, 0)}
	result.SetProperty("objectQuotingStrategy", "QUOTE_ALL_OBJECTS")
	return &result
}

func (y *LqYaml) add(key string, value interface{}) {
	y.DatabaseChangeLog = append(y.DatabaseChangeLog, map[string]interface{}{key: value})
}
func (y *LqYaml) SetProperty(key string, value string) {
	y.DatabaseChangeLog = append(y.DatabaseChangeLog, map[string]interface{}{key: value})
}
func (y *LqYaml) AddChangeSet(changeSet *LqChangeSet) {
	y.DatabaseChangeLog = append(y.DatabaseChangeLog, map[string]interface{}{"changeSet": changeSet})
}

type LqChangeSet struct {
	Id            string                   `yaml:"id"`
	Author        string                   `yaml:"author,omitempty"`
	PreConditions map[string]interface{}   `yaml:"preConditions,omitempty"`
	Changes       []map[string]interface{} `yaml:"changes,omitempty"`
}

func newLqChangeSet(id string, author string) *LqChangeSet {
	return &LqChangeSet{
		Id:            id,
		Author:        author,
		PreConditions: make(map[string]interface{}),
		Changes:       make([]map[string]interface{}, 0),
	}
}

func (s *LqChangeSet) Append(key string, change interface{}) {
	s.Changes = append(s.Changes, map[string]interface{}{key: change})
}

func (s *LqChangeSet) CreateTable(table *LqCreateTable) {
	s.Append("createTable", table)
}

func (s *LqChangeSet) DropTable(table *LqDropTable) {
	s.Append("dropTable", table)
}

func (s *LqChangeSet) AddPrimaryKey(pk *LqAddPrimaryKey) {
	s.Append("addPrimaryKey", pk)
}

func (s *LqChangeSet) AddUniqueConstraint(c *LqAddUniqueConstraint) {
	s.Append("addUniqueConstraint", c)
}

func newCreateTableChangeSet(
	id *LqId,
	author string,
	table *octopus.Table,
	uniqueNameSuffix string,
	useComments bool,
) ([]*LqChangeSet, error) {
	var result []*LqChangeSet

	uniqueNameSet := table.UniqueKeyNameSet()
	primaryKeySet := table.PrimaryKeyNameSet()
	pkCount := primaryKeySet.Size()
	uniqueCount := uniqueNameSet.Size()

	createTable := &LqCreateTable{
		TableName: table.Name,
		Remarks:   table.Description,
		Columns:   make([]map[string]*LqColumn, 0),
	}

	for _, column := range table.Columns {
		if lc, err := newLqColumn(column, pkCount > 1, uniqueCount > 1); err != nil {
			return nil, err
		} else {
			if !useComments {
				lc.Remarks = ""
			}
			createTable.AddColumn(lc)
		}
	}

	createTableChangeSet := newLqChangeSet(id.bumpMinor(), author)
	createTableChangeSet.CreateTable(createTable)
	result = append(result, createTableChangeSet)

	// Primary Key
	if pkCount >= 2 {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		changeSet.AddPrimaryKey(&LqAddPrimaryKey{
			TableName:   table.Name,
			ColumnNames: primaryKeySet.Join(", "),
		})
		result = append(result, changeSet)
	}
	// Unique Constraint
	if uniqueCount >= 1 {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		changeSet.AddUniqueConstraint(newAddUniqueConstraint(table, uniqueNameSet.Join(", "), table.Name+uniqueNameSuffix))

		if uniqueCount == 1 {
			changeSet.PreConditions = map[string]interface{}{
				"onError": "CONTINUE",
				"onFail":  "CONTINUE",
				"dbms": map[string]string{
					"type": "derby, h2, mssql, mariadb, mysql, postgresql, sqlite",
				},
			}
		}
		result = append(result, changeSet)
	}

	return result, nil
}

// ----------------------------------------------------------------------------
// Generator Changes
// ----------------------------------------------------------------------------

type LqCreateTable struct {
	TableName string                 `yaml:"tableName"`
	Remarks   string                 `yaml:"remarks,omitempty"`
	Columns   []map[string]*LqColumn `yaml:"columns"`
}

func (t *LqCreateTable) AddColumn(col *LqColumn) {
	t.Columns = append(t.Columns, map[string]*LqColumn{"column": col})
}

type LqDropTable struct {
	TableName string `yaml:"tableName"`
}

type LqRenameTable struct {
	NewTableName string `yaml:"newTableName"`
	OldTableName string `yaml:"oldTableName"`
}

func newRenameTable(newTableName string, oldTableName string) *LqRenameTable {
	return &LqRenameTable{
		NewTableName: newTableName,
		OldTableName: oldTableName,
	}
}

type LqSetTableRemarks struct {
	TableName string `yaml:"tableName"`
	Remarks   string `yaml:"remarks,omitempty"`
}

func newSetTableRemarks(table *octopus.Table) *LqSetTableRemarks {
	return &LqSetTableRemarks{
		TableName: table.Name,
		Remarks:   table.Description,
	}
}

type LqAddColumn struct {
	TableName string                 `yaml:"tableName"`
	Columns   []map[string]*LqColumn `yaml:"columns"`
}

func newAddColumn(table *octopus.Table, columns []*octopus.Column, useComments bool) (*LqAddColumn, error) {
	var lqColumns []map[string]*LqColumn
	for _, col := range columns {
		if lc, err := newLqColumn(col, true, true); err != nil {
			return nil, err
		} else {
			lastColName := ""
			for _, c := range table.Columns {
				if c == col {
					break
				}
				lastColName = c.Name
			}

			if !useComments {
				lc.Remarks = ""
			}

			lc.AfterColumn = lastColName
			lqColumns = append(lqColumns, map[string]*LqColumn{"column": lc})
		}
	}
	return &LqAddColumn{
		TableName: table.Name,
		Columns:   lqColumns,
	}, nil
}

type LqDropColumn struct {
	TableName  string `yaml:"tableName"`
	ColumnName string `yaml:"columnName"`
}

func newDropColumn(table *octopus.Table, columnName string) *LqDropColumn {
	return &LqDropColumn{
		TableName:  table.Name,
		ColumnName: columnName,
	}
}

type LqSetColumnRemarks struct {
	TableName  string `yaml:"tableName"`
	ColumnName string `yaml:"columnName"`
	Remarks    string `yaml:"remarks,omitempty"`
}

func newSetColumnRemarks(table *octopus.Table, column *octopus.Column) *LqSetColumnRemarks {
	return &LqSetColumnRemarks{
		TableName:  table.Name,
		ColumnName: column.Name,
		Remarks:    column.Description,
	}
}

type LqModifyDataType struct {
	TableName   string `yaml:"tableName"`
	ColumnName  string `yaml:"columnName"`
	NewDataType string `yaml:"newDataType"`
}

func newModifyDataType(table *octopus.Table, column *octopus.Column) *LqModifyDataType {
	return &LqModifyDataType{
		TableName:   table.Name,
		ColumnName:  column.Name,
		NewDataType: getLiquibaseType(column),
	}
}

type LqRenameColumn struct {
	TableName      string `yaml:"tableName"`
	NewColumnName  string `yaml:"newColumnName"`
	OldColumnName  string `yaml:"oldColumnName"`
	ColumnDataType string `yaml:"columnDataType"`
}

func newRenameColumn(table *octopus.Table, newColumn *octopus.Column, oldColumn *octopus.Column) *LqRenameColumn {
	return &LqRenameColumn{
		TableName:      table.Name,
		NewColumnName:  newColumn.Name,
		OldColumnName:  oldColumn.Name,
		ColumnDataType: getLiquibaseType(newColumn),
	}
}

type LqAddNotNullConstraint struct {
	TableName        string `yaml:"tableName"`
	ColumnName       string `yaml:"columnName"`
	ColumnDataType   string `yaml:"columnDataType"`
	DefaultNullValue string `yaml:"defaultNullValue,omitempty"`
}

func newAddNotNullConstraint(table *octopus.Table, column *octopus.Column) *LqAddNotNullConstraint {
	return &LqAddNotNullConstraint{
		TableName:        table.Name,
		ColumnName:       column.Name,
		ColumnDataType:   getLiquibaseType(column),
		DefaultNullValue: column.DefaultValue,
	}
}

type LqDropNotNullConstraint struct {
	TableName      string `yaml:"tableName"`
	ColumnName     string `yaml:"columnName"`
	ColumnDataType string `yaml:"columnDataType"`
}

func newDropNotNullConstraint(table *octopus.Table, column *octopus.Column) *LqDropNotNullConstraint {
	return &LqDropNotNullConstraint{
		TableName:      table.Name,
		ColumnName:     column.Name,
		ColumnDataType: getLiquibaseType(column),
	}
}

type LqAddAutoIncrement struct {
	TableName      string `yaml:"tableName"`
	ColumnName     string `yaml:"columnName"`
	ColumnDataType string `yaml:"columnDataType"`
}

func newAddAutoIncrement(table *octopus.Table, column *octopus.Column) *LqAddAutoIncrement {
	return &LqAddAutoIncrement{
		TableName:      table.Name,
		ColumnName:     column.Name,
		ColumnDataType: getLiquibaseType(column),
	}
}

type LqAddDefaultValue struct {
	TableName           string      `yaml:"tableName"`
	ColumnName          string      `yaml:"columnName"`
	ColumnDataType      string      `yaml:"columnDataType,omitempty"`
	DefaultValue        string      `yaml:"defaultValue,omitempty"`
	DefaultValueBoolean *bool       `yaml:"defaultValueBoolean,omitempty"`
	DefaultValueNumeric interface{} `yaml:"defaultValueNumeric,omitempty"`
	DefaultValueDate    string      `yaml:"defaultValueDate,omitempty"`
}

func newAddDefaultValue(table *octopus.Table, column *octopus.Column) (*LqAddDefaultValue, error) {
	if dv, err := newLqDefaultValue(column); err != nil {
		return nil, err
	} else {
		return &LqAddDefaultValue{
			TableName:           table.Name,
			ColumnName:          column.Name,
			ColumnDataType:      getLiquibaseType(column),
			DefaultValue:        dv.DefaultValue,
			DefaultValueBoolean: dv.DefaultValueBoolean,
			DefaultValueNumeric: dv.DefaultValueNumeric,
			DefaultValueDate:    dv.DefaultValueDate,
		}, nil
	}
}

type LqDropDefaultValue struct {
	TableName      string `yaml:"tableName"`
	ColumnName     string `yaml:"columnName"`
	ColumnDataType string `yaml:"columnDataType"`
}

func newDropDefaultValue(table *octopus.Table, column *octopus.Column) *LqDropDefaultValue {
	return &LqDropDefaultValue{
		TableName:      table.Name,
		ColumnName:     column.Name,
		ColumnDataType: getLiquibaseType(column),
	}
}

type LqAddPrimaryKey struct {
	TableName   string `yaml:"tableName"`
	ColumnNames string `yaml:"columnNames"`
}

func newAddPrimaryKey(table *octopus.Table, columnNames string) *LqAddPrimaryKey {
	return &LqAddPrimaryKey{
		TableName:   table.Name,
		ColumnNames: columnNames,
	}
}

type LqDropPrimaryKey struct {
	TableName string `yaml:"tableName"`
}

func newDropPrimaryKey(table *octopus.Table) *LqDropPrimaryKey {
	return &LqDropPrimaryKey{
		TableName: table.Name,
	}
}

type LqAddUniqueConstraint struct {
	TableName      string `yaml:"tableName"`
	ColumnNames    string `yaml:"columnNames"`
	ConstraintName string `yaml:"constraintName"`
}

func newAddUniqueConstraint(table *octopus.Table, columnNames string, uniqueConstraintName string) *LqAddUniqueConstraint {
	return &LqAddUniqueConstraint{
		TableName:      table.Name,
		ColumnNames:    columnNames,
		ConstraintName: uniqueConstraintName,
	}
}

type LqDropUniqueConstraint struct {
	TableName      string `yaml:"tableName"`
	ConstraintName string `yaml:"constraintName"`
}

func newDropUniqueConstraint(table *octopus.Table, uniqueConstraintName string) *LqDropUniqueConstraint {
	return &LqDropUniqueConstraint{
		TableName:      table.Name,
		ConstraintName: uniqueConstraintName,
	}
}

// ----------------------------------------------------------------------------
// Generator struct definitions
// ----------------------------------------------------------------------------

type LqDefaultValue struct {
	DefaultValue        string      `yaml:"defaultValue,omitempty"`
	DefaultValueBoolean *bool       `yaml:"defaultValueBoolean,omitempty"`
	DefaultValueNumeric interface{} `yaml:"defaultValueNumeric,omitempty"`
	DefaultValueDate    string      `yaml:"defaultValueDate,omitempty"`
}

func newLqDefaultValue(column *octopus.Column) (*LqDefaultValue, error) {
	result := LqDefaultValue{}

	if util.IsStringType(column.Type) {
		result.DefaultValue = column.DefaultValue
	} else if util.IsBooleanType(column.Type) {
		result.DefaultValueBoolean = NewBool(column.DefaultValue == "true")
	} else if util.IsNumericType(column.Type) {
		if util.IsIntType(column.Type) {
			if num, err := strconv.Atoi(column.DefaultValue); err != nil {
				return nil, err
			} else {
				result.DefaultValueNumeric = &num
			}
		} else {
			if num, err := strconv.ParseFloat(column.DefaultValue, 64); err != nil {
				return nil, err
			} else {
				result.DefaultValueNumeric = &num
			}
		}
	} else if util.IsDateType(column.Type) {
		result.DefaultValueDate = column.DefaultValue
	} else {
		result.DefaultValue = column.DefaultValue
	}

	return &result, nil
}

type LqColumn struct {
	Name                string         `yaml:"name"`
	Type                string         `yaml:"type"`
	AutoIncrement       *bool          `yaml:"autoIncrement,omitempty"`
	Constraints         *LqConstraints `yaml:"constraints,omitempty"`
	Remarks             string         `yaml:"remarks,omitempty"`
	DefaultValue        string         `yaml:"defaultValue,omitempty"`
	DefaultValueBoolean *bool          `yaml:"defaultValueBoolean,omitempty"`
	DefaultValueNumeric interface{}    `yaml:"defaultValueNumeric,omitempty"`
	DefaultValueDate    string         `yaml:"defaultValueDate,omitempty"`
	AfterColumn         string         `yaml:"afterColumn,omitempty"`
}

func newLqColumn(column *octopus.Column, createSeparatePK bool, createSeparateUq bool) (*LqColumn, error) {
	lc := LqColumn{
		Name:    column.Name,
		Type:    getLiquibaseType(column),
		Remarks: column.Description,
	}

	// auto_incremental
	if column.AutoIncremental {
		lc.AutoIncrement = NewBool(true)
	}

	// constraints
	hasConstraints := false
	constraints := LqConstraints{}
	if column.PrimaryKey && !createSeparatePK {
		constraints.PrimaryKey = NewBool(true)
		hasConstraints = true
	} else if column.NotNull {
		constraints.Nullable = NewBool(false)
		hasConstraints = true
	}
	if column.UniqueKey && !createSeparateUq {
		constraints.Unique = NewBool(true)
		hasConstraints = true
	}
	if hasConstraints {
		lc.Constraints = &constraints
	}

	// default value
	if column.DefaultValue != "" {
		if dv, err := newLqDefaultValue(column); err != nil {
			return nil, err
		} else {
			lc.DefaultValue = dv.DefaultValue
			lc.DefaultValueBoolean = dv.DefaultValueBoolean
			lc.DefaultValueNumeric = dv.DefaultValueNumeric
			lc.DefaultValueDate = dv.DefaultValueDate
		}
	}
	return &lc, nil
}

type LqConstraints struct {
	PrimaryKey *bool `yaml:"primaryKey,omitempty"`
	Nullable   *bool `yaml:"nullable,omitempty"`
	Unique     *bool `yaml:"unique,omitempty"`
}

type Option struct {
	TableFilter      octopus.TableFilterFn
	UniqueNameSuffix string
	UseComments      bool
}

type Generator struct {
	schema *octopus.Schema
	option *Option
}

func (c *Generator) Generate(outputPath string) error {
	var outputBytes []byte

	// generate all
	if bytes, err := c.generateAll(); err != nil {
		return err
	} else {
		outputBytes = bytes
	}

	// Write file
	if err := ioutil.WriteFile(outputPath, outputBytes, 0644); err != nil {
		return err
	}
	log.Printf("[WRITE] %s", outputPath)

	return nil
}

func (c *Generator) generateAll() ([]byte, error) {
	result := newLqYaml()

	option := c.option
	uniqueNameSuffix := option.UniqueNameSuffix
	useComments := option.UseComments

	id := newLqId()
	for _, table := range c.schema.Tables {
		// filter table
		if option.TableFilter != nil && !option.TableFilter(table) {
			continue
		}

		id.bumpMajor()

		// create table
		if changeSets, err := newCreateTableChangeSet(id, c.schema.Author, table, uniqueNameSuffix, useComments); err != nil {
			return nil, err
		} else {
			for _, changeSet := range changeSets {
				result.AddChangeSet(changeSet)
			}
		}
	}

	return yaml.Marshal(&result)
}

func getLiquibaseType(column *octopus.Column) string {
	typ := ""
	switch strings.ToLower(column.Type) {
	case octopus.ColTypeChar:
		typ = "char"
	case octopus.ColTypeVarchar:
		typ = "varchar"
	case octopus.ColTypeText8:
		fallthrough
	case octopus.ColTypeText16:
		fallthrough
	case octopus.ColTypeText24:
		fallthrough
	case octopus.ColTypeText32:
		typ = "clob"
	case octopus.ColTypeBoolean:
		typ = "boolean"
	case octopus.ColTypeInt8:
		typ = "tinyint"
	case octopus.ColTypeInt16:
		typ = "smallint"
	case octopus.ColTypeInt24:
		typ = "mediumint"
	case octopus.ColTypeInt32:
		typ = "int"
	case octopus.ColTypeInt64:
		typ = "bigint"
	case octopus.ColTypeDecimal:
		typ = "decimal"
	case octopus.ColTypeFloat:
		typ = "float"
	case octopus.ColTypeDouble:
		typ = "double"
	case octopus.ColTypeDateTime:
		typ = "datetime"
	case octopus.ColTypeDate:
		typ = "date"
	case octopus.ColTypeTime:
		typ = "time"
	case octopus.ColTypeBlob8:
		fallthrough
	case octopus.ColTypeBlob16:
		fallthrough
	case octopus.ColTypeBlob24:
		fallthrough
	case octopus.ColTypeBlob32:
		typ = "blob"
	default:
		typ = column.Type
	}
	if column.Size > 0 {
		if column.Scale > 0 {
			return fmt.Sprintf("%s(%d,%d)", typ, column.Size, column.Scale)
		} else {
			return fmt.Sprintf("%s(%d)", typ, column.Size)
		}
	} else {
		return typ
	}
}

type LqId struct {
	major int
	minor int
}

func newLqId() *LqId {
	return &LqId{
		major: 0,
		minor: 0,
	}
}

func (l *LqId) bumpMajor() {
	l.major++
	l.minor = 0
}

func (l *LqId) revertMajor() {
	l.major--
	l.minor = 0
}

func (l *LqId) bumpMinor() string {
	l.minor++
	return l.version()
}

func (l *LqId) version() string {
	if l.minor == 0 {
		return fmt.Sprintf("%d", l.major)
	} else {
		return fmt.Sprintf("%d-%d", l.major, l.minor)
	}
}

func NewBool(value bool) *bool {
	b := value
	return &b
}
