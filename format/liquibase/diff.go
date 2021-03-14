package liquibase

import (
	"github.com/google/go-cmp/cmp"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type DiffOption struct {
	TableFilter      octopus.TableFilterFn
	DiffFrom         *octopus.Schema
	DiffTo           *octopus.Schema
	Author           string
	UniqueNameSuffix string
	UseComments      bool
}

type Diff struct {
	option *DiffOption
}

func NewDiff(option *DiffOption) *Diff {
	return &Diff{option}
}

func (c *Diff) Generate(outputPath string) error {
	if outputBytes, err := c.generate(); err != nil {
		return err
	} else {
		// Write file
		if err := ioutil.WriteFile(outputPath, outputBytes, 0644); err != nil {
			return err
		}
		log.Printf("[WRITE] %s", outputPath)

		return nil
	}
}

func (c *Diff) generate() ([]byte, error) {
	result := newLqYaml()

	option := c.option
	uniqueNameSuffix := option.UniqueNameSuffix
	useComments := option.UseComments

	fromTableByName := option.DiffFrom.TablesByName()
	author := util.IfThenElseString(option.Author != "", option.Author, option.DiffTo.Author)

	var addedTables []*octopus.Table
	renamedTableMap := make(map[*octopus.Table]*octopus.Table)
	removedTableNames := util.NewStringSet()
	for name := range fromTableByName {
		removedTableNames.Add(name)
	}

	id := newLqId()
	for _, table := range option.DiffTo.Tables {
		// filter table
		if option.TableFilter != nil && !option.TableFilter(table) {
			continue
		}

		tableName := table.Name
		removedTableNames.Remove(tableName)

		oldTable, ok := fromTableByName[tableName]
		if !ok {
			// added
			addedTables = append(addedTables, table)
			continue
		}

		// diff table
		id.bumpMajor()
		if diffChangeSet, err := c.diffTable(id, author, table, oldTable, useComments, uniqueNameSuffix); err != nil {
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
			if removedTable := fromTableByName[removedTableName]; removedTable != nil {
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
		changeSet := newLqChangeSet(id.version(), author)
		changeSet.Append("dropTable", &LqDropTable{TableName: tableName})
		result.AddChangeSet(changeSet)
	}

	// renamed tables
	for newTable, oldTable := range renamedTableMap {
		id.bumpMajor()

		// drop old unique constraint
		if oldTable.UniqueKeyNameSet().Size() > 0 {
			changeSet := newLqChangeSet(id.version(), author)
			changeSet.Append("dropUniqueConstraint", newDropUniqueConstraint(oldTable, oldTable.Name+uniqueNameSuffix))
			result.AddChangeSet(changeSet)
		}

		// rename table
		{
			changeSet := newLqChangeSet(id.bumpMinor(), author)
			changeSet.Append("renameTable", newRenameTable(newTable.Name, oldTable.Name))
			result.AddChangeSet(changeSet)
		}

		// add new unique constraint
		newUqSet := newTable.UniqueKeyNameSet()
		if newUqSet.Size() > 0 {
			uniqueColumeNames := newUqSet.Join(", ")
			uniqueConstraintName := newTable.Name + uniqueNameSuffix

			changeSet := newLqChangeSet(id.bumpMinor(), author)
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
		if changeSets, err := newCreateTableChangeSet(id, author, table, uniqueNameSuffix, useComments); err != nil {
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
func (c *Diff) diffTable(
	id *LqId,
	author string,
	table *octopus.Table,
	oldTable *octopus.Table,
	useComments bool,
	uniqueNameSuffix string,
) ([]*LqChangeSet, error) {
	var changeSets []*LqChangeSet

	if useComments && table.Description != oldTable.Description {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		changeSet.Append("setTableRemarks", newSetTableRemarks(table))
		changeSets = append(changeSets, changeSet)
	}

	oldColumnByName := oldTable.ColumnNameMap()

	var addedColumns []*octopus.Column
	renamedColumnMap := make(map[*octopus.Column]*octopus.Column)
	removedColumnNameSet := util.NewStringSet()
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
		if changes, err := c.diffColumn(id, author, table, column, oldColumn, useComments); err != nil {
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

		// not null constraint is removed after renameColumn. (fixed in liquibase v4.0)
		if oldColumn.NotNull {
			changeSet = newLqChangeSet(id.bumpMinor(), author)
			changeSet.Append("addNotNullConstraint", newAddNotNullConstraint(table, newColumn))
			changeSets = append(changeSets, changeSet)
		}
	}

	// added columns
	var filteredAddedColumns []*octopus.Column
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
func (c *Diff) diffColumn(
	id *LqId,
	author string,
	table *octopus.Table,
	column *octopus.Column,
	oldColumn *octopus.Column,
	useComments bool,
) ([]*LqChangeSet, error) {
	var changeSets []*LqChangeSet

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

	if column.NotNull != oldColumn.NotNull {
		changeSet := newLqChangeSet(id.bumpMinor(), author)
		if column.NotNull {
			changeSet.Append("addNotNullConstraint", newAddNotNullConstraint(table, column))
		} else {
			changeSet.Append("dropNotNullConstraint", newDropNotNullConstraint(table, column))
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
