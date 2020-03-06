package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
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
	table *Table,
	uniqueNameSuffix string,
	useComments bool,
) ([]*LqChangeSet, error) {
	result := make([]*LqChangeSet, 0)

	uniqueNameSet := table.UniqueKeyNameSet()
	primaryKeySet := table.PrimaryKeyNameSet()
	pkCount := primaryKeySet.Size()
	uniqueCount := uniqueNameSet.Size()

	createTable := &LqCreateTable{
		TableName: table.Name,
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
// Liquibase Changes
// ----------------------------------------------------------------------------

type LqCreateTable struct {
	TableName string                 `yaml:"tableName"`
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

type LqAddColumn struct {
	TableName string      `yaml:"tableName"`
	Columns   []*LqColumn `yaml:"columns"`
}

func newAddColumn(table *Table, columns []*Column, useComments bool) (*LqAddColumn, error) {
	lqColumns := make([]*LqColumn, 0)
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
			lqColumns = append(lqColumns, lc)
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

func newDropColumn(table *Table, columnName string) *LqDropColumn {
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

func newSetColumnRemarks(table *Table, column *Column) *LqSetColumnRemarks {
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

func newModifyDataType(table *Table, column *Column) *LqModifyDataType {
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

func newRenameColumn(table *Table, newColumn *Column, oldColumn *Column) *LqRenameColumn {
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

func newAddNotNullConstraint(table *Table, column *Column) *LqAddNotNullConstraint {
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

func newDropNotNullConstraint(table *Table, column *Column) *LqDropNotNullConstraint {
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

func newAddAutoIncrement(table *Table, column *Column) *LqAddAutoIncrement {
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

func newAddDefaultValue(table *Table, column *Column) (*LqAddDefaultValue, error) {
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

func newDropDefaultValue(table *Table, column *Column) *LqDropDefaultValue {
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

func newAddPrimaryKey(table *Table, columnNames string) *LqAddPrimaryKey {
	return &LqAddPrimaryKey{
		TableName:   table.Name,
		ColumnNames: columnNames,
	}
}

type LqDropPrimaryKey struct {
	TableName string `yaml:"tableName"`
}

func newDropPrimaryKey(table *Table) *LqDropPrimaryKey {
	return &LqDropPrimaryKey{
		TableName: table.Name,
	}
}

type LqAddUniqueConstraint struct {
	TableName      string `yaml:"tableName"`
	ColumnNames    string `yaml:"columnNames"`
	ConstraintName string `yaml:"constraintName"`
}

func newAddUniqueConstraint(table *Table, columnNames string, uniqueConstraintName string) *LqAddUniqueConstraint {
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

func newDropUniqueConstraint(table *Table, uniqueConstraintName string) *LqDropUniqueConstraint {
	return &LqDropUniqueConstraint{
		TableName:      table.Name,
		ConstraintName: uniqueConstraintName,
	}
}

// ----------------------------------------------------------------------------
// Liquibase struct definitions
// ----------------------------------------------------------------------------

type LqDefaultValue struct {
	DefaultValue        string      `yaml:"defaultValue,omitempty"`
	DefaultValueBoolean *bool       `yaml:"defaultValueBoolean,omitempty"`
	DefaultValueNumeric interface{} `yaml:"defaultValueNumeric,omitempty"`
	DefaultValueDate    string      `yaml:"defaultValueDate,omitempty"`
}

func newLqDefaultValue(column *Column) (*LqDefaultValue, error) {
	result := LqDefaultValue{}

	if IsStringType(column.Type) {
		result.DefaultValue = column.DefaultValue
	} else if IsBooleanType(column.Type) {
		result.DefaultValueBoolean = NewBool(column.DefaultValue == "true")
	} else if IsNumericType(column.Type) {
		if IsIntType(column.Type) {
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
	} else if IsDateType(column.Type) {
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

func newLqColumn(column *Column, createSeparatePK bool, createSeparateUq bool) (*LqColumn, error) {
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
	} else if !column.Nullable {
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

type Liquibase struct {
}

func (l *Liquibase) Generate(
	schema *Schema,
	output *Output,
	tableFilterFn TableFilterFn,
) error {
	// Create directory
	if err := os.MkdirAll(output.FilePath, 0777); err != nil {
		return err
	}
	log.Printf("[MKDIR] %s", output.FilePath)

	var outputBytes []byte

	diffFilename := output.Get(FlagDiff)
	if diffFilename != "" {
		// diff mode
		if input, err := NewInput(diffFilename, ""); err != nil {
			return err
		} else {
			if targetSchema, err := input.ToSchema(); err != nil {
				return err
			} else {
				if bytes, err := l.generateDiff(schema, targetSchema, output, tableFilterFn); err != nil {
					return err
				} else {
					outputBytes = bytes
				}
			}
		}
	} else {
		// generate all
		if bytes, err := l.generateAll(schema, output, tableFilterFn); err != nil {
			return err
		} else {
			outputBytes = bytes
		}
	}

	// Write file
	outputFile := path.Join(output.FilePath,
		fmt.Sprintf("%s-%s.yaml", schema.Name, schema.Version))

	if err := ioutil.WriteFile(outputFile, outputBytes, 0644); err != nil {
		return err
	}
	log.Printf("[WRITE] %s", outputFile)

	return nil
}

func (l *Liquibase) generateAll(
	schema *Schema,
	output *Output,
	tableFilterFn TableFilterFn,
) ([]byte, error) {
	result := newLqYaml()

	uniqueNameSuffix := output.Get(FlagUniqueNameSuffix)
	useComments := output.GetBool(FlagUseComments)

	id := newLqId()
	for _, table := range schema.Tables {
		// filter table
		if tableFilterFn != nil && !tableFilterFn(table) {
			continue
		}

		id.bumpMajor()

		// create table
		if changeSets, err := newCreateTableChangeSet(id, schema.Author, table, uniqueNameSuffix, useComments); err != nil {
			return nil, err
		} else {
			for _, changeSet := range changeSets {
				result.AddChangeSet(changeSet)
			}
		}
	}

	return yaml.Marshal(&result)
}

func (l *Liquibase) generateDiff(
	schema *Schema,
	oldSchema *Schema,
	output *Output,
	tableFilterFn TableFilterFn,
) ([]byte, error) {
	result := newLqYaml()

	uniqueNameSuffix := output.Get(FlagUniqueNameSuffix)
	useComments := output.GetBool(FlagUseComments)

	oldTableByName := oldSchema.TableByName()

	addedTables := make([]*Table, 0)
	renamedTableMap := make(map[*Table]*Table)
	removedTableNames := NewStringSet()
	for name, _ := range oldTableByName {
		removedTableNames.Add(name)
	}

	id := newLqId()
	for _, table := range schema.Tables {
		// filter table
		if tableFilterFn != nil && !tableFilterFn(table) {
			continue
		}

		tableName := table.Name
		removedTableNames.Remove(tableName)

		oldTable, ok := oldTableByName[tableName]
		if !ok {
			// added
			addedTables = append(addedTables, table)
			continue
		}

		// diff table
		id.bumpMajor()
		if diffChangeSet, err := l.diffTable(id, schema.Author, table, oldTable, useComments, uniqueNameSuffix); err != nil {
			return nil, err
		} else if len(diffChangeSet) > 0 {
			for _, changeSet := range diffChangeSet {
				result.AddChangeSet(changeSet)
			}
		} else {
			id.revertMajor()
		}
	}

	// find renamed tables
	for _, addedTable := range addedTables {
		for _, removedTableName := range removedTableNames.Slice() {
			if removedTable := oldTableByName[removedTableName]; removedTable != nil {
				if diff := cmp.Diff(addedTable.Columns, removedTable.Columns); diff == "" {
					renamedTableMap[addedTable] = removedTable
					removedTableNames.Remove(removedTableName)
					break
				}
			}
		}
	}

	// removed tables
	for _, tableName := range removedTableNames.Slice() {
		id.bumpMajor()

		// drop table
		changeSet := newLqChangeSet(id.version(), schema.Author)
		changeSet.Append("dropTable", &LqDropTable{TableName: tableName})
		result.AddChangeSet(changeSet)
	}

	// renamed tables
	for newTable, oldTable := range renamedTableMap {
		id.bumpMajor()

		// drop old unique constraint
		if oldTable.UniqueKeyNameSet().Size() > 0 {
			changeSet := newLqChangeSet(id.version(), schema.Author)
			changeSet.Append("dropUniqueConstraint", newDropUniqueConstraint(oldTable, oldTable.Name+uniqueNameSuffix))
			result.AddChangeSet(changeSet)
		}

		// rename table
		{
			changeSet := newLqChangeSet(id.bumpMinor(), schema.Author)
			changeSet.Append("renameTable", newRenameTable(newTable.Name, oldTable.Name))
			result.AddChangeSet(changeSet)
		}

		// add new unique constraint
		newUqSet := newTable.UniqueKeyNameSet()
		if newUqSet.Size() > 0 {
			uniqueColumeNames := newUqSet.Join(", ")
			uniqueConstraintName := newTable.Name + uniqueNameSuffix

			changeSet := newLqChangeSet(id.bumpMinor(), schema.Author)
			changeSet.Append("addUniqueConstraint", newAddUniqueConstraint(newTable, uniqueColumeNames, uniqueConstraintName))
			result.AddChangeSet(changeSet)
		}
	}

	// added tables
	for _, table := range addedTables {
		// skip renamed table
		if _, ok := renamedTableMap[table]; ok {
			continue
		}

		id.bumpMajor()

		// create table
		if changeSets, err := newCreateTableChangeSet(id, schema.Author, table, uniqueNameSuffix, useComments); err != nil {
			return nil, err
		} else {
			for _, changeSet := range changeSets {
				result.AddChangeSet(changeSet)
			}
		}
	}

	return yaml.Marshal(&result)
}

// diffTable compares two tables.
func (l *Liquibase) diffTable(
	id *LqId,
	author string,
	table *Table,
	oldTable *Table,
	useComments bool,
	uniqueNameSuffix string,
) ([]*LqChangeSet, error) {
	changeSets := make([]*LqChangeSet, 0)

	oldColumnByName := oldTable.ColumnByName()

	addedColumns := make([]*Column, 0)
	renamedColumnMap := make(map[*Column]*Column)
	removedColumnNameSet := NewStringSet()
	for _, column := range oldTable.Columns {
		removedColumnNameSet.Add(column.Name)
	}

	for _, column := range table.Columns {
		oldColumn, ok := oldColumnByName[column.Name]
		if !ok {
			// added column
			addedColumns = append(addedColumns, column)
			continue
		}
		removedColumnNameSet.Remove(column.Name)

		// diff column
		if changes, err := l.diffColumn(id, author, table, column, oldColumn, useComments); err != nil {
			return nil, err
		} else {
			for _, change := range changes {
				changeSets = append(changeSets, change)
			}
		}
	}

	// find renamed column
	for _, addedColumn := range addedColumns {
		for _, removedColumnName := range removedColumnNameSet.Slice() {
			if removedColumn := oldColumnByName[removedColumnName]; removedColumn != nil {
				if addedColumn.IsRenamed(removedColumn, !useComments) {
					renamedColumnMap[addedColumn] = removedColumn
					removedColumnNameSet.Remove(removedColumnName)
					break
				}
			}
		}
	}

	// unique key
	uqSet := table.UniqueKeyNameSet()
	oldUqSet := oldTable.UniqueKeyNameSet()
	uniqueChanged := !uqSet.Equals(oldUqSet)

	// drop unique constraint
	if uniqueChanged {
		uniqueConstraintName := table.Name + uniqueNameSuffix
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		changeSet.Append("dropUniqueConstraint", newDropUniqueConstraint(table, uniqueConstraintName))
		changeSets = append(changeSets, changeSet)
	}

	// removed columns
	for _, columnName := range removedColumnNameSet.Slice() {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		changeSet.Append("dropColumn", newDropColumn(table, columnName))
		changeSets = append(changeSets, changeSet)
	}

	// renamed columns
	for newColumn, oldColumn := range renamedColumnMap {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		changeSet.Append("renameColumn", newRenameColumn(table, newColumn, oldColumn))
		changeSets = append(changeSets, changeSet)
	}

	// added columns
	filteredAddedColumns := make([]*Column, 0)
	for _, col := range addedColumns {
		if _, ok := renamedColumnMap[col]; ok {
			continue
		}
		filteredAddedColumns = append(filteredAddedColumns, col)
	}

	if len(filteredAddedColumns) > 0 {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		if lqAddColumn, err := newAddColumn(table, filteredAddedColumns, useComments); err != nil {
			return nil, err
		} else {
			changeSet.Append("addColumn", lqAddColumn)
			changeSets = append(changeSets, changeSet)
		}
	}

	// primary key
	pkSet := table.PrimaryKeyNameSet()
	oldPkSet := oldTable.PrimaryKeyNameSet()
	if !pkSet.Equals(oldPkSet) {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		if oldPkSet.Size() > 0 {
			changeSet.Append("dropPrimaryKey", newDropPrimaryKey(table))
		}
		if pkSet.Size() > 0 {
			pkColumeNames := pkSet.Join(", ")
			changeSet.Append("addPrimaryKey", newAddPrimaryKey(table, pkColumeNames))
		}
		changeSets = append(changeSets, changeSet)
	}

	// add unique constraint
	if uniqueChanged {
		uniqueConstraintName := table.Name + uniqueNameSuffix
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		if uqSet.Size() > 0 {
			uniqueColumeNames := uqSet.Join(", ")
			changeSet.Append("addUniqueConstraint", newAddUniqueConstraint(table, uniqueColumeNames, uniqueConstraintName))
		}
		changeSets = append(changeSets, changeSet)
	}

	return changeSets, nil
}

// diffColumn compares two columns.
func (l *Liquibase) diffColumn(
	id *LqId,
	author string,
	table *Table,
	column *Column,
	oldColumn *Column,
	useComments bool,
) ([]*LqChangeSet, error) {
	changeSets := make([]*LqChangeSet, 0)

	if column.Name != oldColumn.Name {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		changeSet.Append("renameColumn", newRenameColumn(table, column, oldColumn))
		changeSets = append(changeSets, changeSet)
	}

	if column.Type != oldColumn.Type || column.Size != oldColumn.Size || column.Scale != oldColumn.Scale {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		changeSet.Append("modifyDataType", newModifyDataType(table, column))
		changeSets = append(changeSets, changeSet)
	}

	if useComments && column.Description != oldColumn.Description {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		changeSet.Append("setColumnRemarks", newSetColumnRemarks(table, column))
		changeSets = append(changeSets, changeSet)
	}

	if column.Nullable != oldColumn.Nullable {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		if column.Nullable {
			changeSet.Append("dropNotNullConstraint", newDropNotNullConstraint(table, column))
		} else {
			changeSet.Append("addNotNullConstraint", newAddNotNullConstraint(table, column))
		}
		changeSets = append(changeSets, changeSet)
	}

	if column.AutoIncremental != oldColumn.AutoIncremental {
		if !column.AutoIncremental {
			log.Printf("liquibase does not support drop autoIncrement. column: %v", column)
		} else {
			changeSet := newLqChangeSet(id.bumpMinor(), author)
			changeSet.Append("addAutoIncrement", newAddAutoIncrement(table, column))
			changeSets = append(changeSets, changeSet)
		}
	}

	if column.DefaultValue != oldColumn.DefaultValue {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		if column.DefaultValue == "" {
			changeSet.Append("dropDefaultValue", newDropDefaultValue(table, column))
		} else {
			if change, err := newAddDefaultValue(table, column); err != nil {
				return nil, err
			} else {
				changeSet.Append("addDefaultValue", change)
			}
		}
		changeSets = append(changeSets, changeSet)
	}

	return changeSets, nil
}

func getLiquibaseType(column *Column) string {
	typ := ""
	switch strings.ToLower(column.Type) {
	case ColTypeString:
		typ = "varchar"
	case ColTypeText:
		typ = "clob"
	case ColTypeBoolean:
		typ = "boolean"
	case ColTypeLong:
		typ = "bigint"
	case ColTypeInt:
		typ = "int"
	case ColTypeDecimal:
		typ = "decimal"
	case ColTypeFloat:
		typ = "float"
	case ColTypeDouble:
		typ = "double"
	case ColTypeDateTime:
		typ = "datetime"
	case ColTypeDate:
		typ = "date"
	case ColTypeTime:
		typ = "time"
	case ColTypeBlob:
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
