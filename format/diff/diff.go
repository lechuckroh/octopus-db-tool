package diff

import (
	"github.com/google/go-cmp/cmp"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
	"github.com/lechuckroh/octopus-db-tools/util"
)

func newDiffChangeSet(id string, author string) *ChangeSet {
	return &ChangeSet{
		ID:      id,
		Author:  author,
		Changes: make([]Change, 0),
	}
}

type Option struct {
	TableFilter      octopus.TableFilterFn
	DiffFrom         *octopus.Schema
	DiffTo           *octopus.Schema
	Author           string
	UniqueNameSuffix string
	UseComments      bool
}

func getDiff(option *Option) (*Result, error) {
	result := Result{From: option.DiffFrom, To: option.DiffTo}

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

	id := newChangeSetID()
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
		if diffChangeSet, err := diffTable(id, author, table, oldTable, useComments, uniqueNameSuffix); err != nil {
			return nil, err
		} else if len(diffChangeSet) > 0 {
			for _, changeSet := range diffChangeSet {
				result.Add(changeSet)
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
		changeSet := newDiffChangeSet(id.version(), author)
		changeSet.Add(&DropTable{Table: fromTableByName[tableName]})
		result.Add(changeSet)
	}

	// renamed tables
	for newTable, oldTable := range renamedTableMap {
		id.bumpMajor()

		// drop old unique constraint
		if oldTable.UniqueKeyNameSet().Size() > 0 {
			changeSet := newDiffChangeSet(id.version(), author)
			changeSet.Add(&DropUniqueConstraint{
				ConstraintName: oldTable.Name + uniqueNameSuffix,
				Table:          oldTable,
			})
			result.Add(changeSet)
		}

		// rename table
		{
			changeSet := newDiffChangeSet(id.bumpMinor(), author)
			changeSet.Add(&RenameTable{OldTable: oldTable, NewTable: newTable})
			result.Add(changeSet)
		}

		// add new unique constraint
		newUqSet := newTable.UniqueKeyNameSet()
		if newUqSet.Size() > 0 {
			changeSet := newDiffChangeSet(id.bumpMinor(), author)
			changeSet.Add(&CreateUniqueConstraint{
				ConstraintName: newTable.Name + uniqueNameSuffix,
				Table:          newTable,
			})
			result.Add(changeSet)
		}
	}

	// added tables
	for _, table := range addedTables {
		// skip renamed table
		if _, ok := renamedTableMap[table]; ok {
			continue
		}

		// create table
		id.bumpMajor()
		changeSet := newDiffChangeSet(id.version(), author)
		changeSet.Add(&CreateTable{Table: table})
		result.Add(changeSet)
	}

	return &result, nil
}

// diffTable compares two tables.
func diffTable(
	id *ChangeSetID,
	author string,
	table *octopus.Table,
	oldTable *octopus.Table,
	useComments bool,
	uniqueNameSuffix string,
) ([]*ChangeSet, error) {
	var changeSets []*ChangeSet

	if useComments && table.Description != oldTable.Description {
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&SetTableComment{Table: table})
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
		if changes, err := diffColumn(id, author, table, column, oldColumn, useComments); err != nil {
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
	if uniqueChanged && oldUqSet.Size() > 0{
		uniqueConstraintName := table.Name + uniqueNameSuffix
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&DropUniqueConstraint{
			ConstraintName: uniqueConstraintName,
			Table:          oldTable,
		})
		changeSets = append(changeSets, changeSet)
	}

	// removed columns
	for _, columnName := range removedColumnNameSet.Slice() {
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&DropColumn{Table: table, ColumnName: columnName})
		changeSets = append(changeSets, changeSet)
	}

	// renamed columns
	for newColumn, oldColumn := range renamedColumnMap {
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&RenameColumn{Table: table, OldColumn: oldColumn, NewColumn: newColumn})
		changeSets = append(changeSets, changeSet)
	}

	// added columns
	var finalAddedColumns []*octopus.Column
	for _, col := range addedColumns {
		if _, ok := renamedColumnMap[col]; ok {
			continue
		}
		finalAddedColumns = append(finalAddedColumns, col)
	}

	if len(finalAddedColumns) > 0 {
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		for _, column := range finalAddedColumns {
			changeSet.Add(&AddColumn{
				Table:        table,
				Column:       column,
				BeforeColumn: nil,
				AfterColumn:  nil,
			})
		}
	}

	// primary key
	pkSet := table.PrimaryKeyNameSet()
	oldPkSet := oldTable.PrimaryKeyNameSet()
	if !pkSet.Equals(oldPkSet) {
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&UpdatePrimaryKey{OldTable: oldTable, NewTable: table})
		changeSets = append(changeSets, changeSet)
	}

	// unique constraint
	if uniqueChanged && uqSet.Size() > 0{
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&CreateUniqueConstraint{
			ConstraintName: table.Name + uniqueNameSuffix,
			Table:          table,
		})
		changeSets = append(changeSets, changeSet)
	}

	return changeSets, nil
}

// diffColumn compares two columns.
func diffColumn(
	id *ChangeSetID,
	author string,
	table *octopus.Table,
	column *octopus.Column,
	oldColumn *octopus.Column,
	useComments bool,
) ([]*ChangeSet, error) {
	var changeSets []*ChangeSet

	if column.Name != oldColumn.Name {
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&RenameColumn{Table: table, OldColumn: oldColumn, NewColumn: column})
		changeSets = append(changeSets, changeSet)
	}

	if column.Type != oldColumn.Type || column.Size != oldColumn.Size || column.Scale != oldColumn.Scale {
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&ChangeColumnType{Table: table, OldColumn: oldColumn, NewColumn: column})
		changeSets = append(changeSets, changeSet)
	}

	if useComments && column.Description != oldColumn.Description {
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&SetColumnComment{Table: table, Column: column})
		changeSets = append(changeSets, changeSet)
	}

	if column.NotNull != oldColumn.NotNull {
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&SetNotNullConstraint{Table: table, Column: column})
		changeSets = append(changeSets, changeSet)
	}

	if column.AutoIncremental != oldColumn.AutoIncremental {
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&SetAutoIncrement{Table: table, Column: column})
		changeSets = append(changeSets, changeSet)
	}

	if column.DefaultValue != oldColumn.DefaultValue {
		changeSet := newDiffChangeSet(id.bumpMinor(), author)
		changeSet.Add(&SetDefaultValue{Table: table, Column: column})
		changeSets = append(changeSets, changeSet)
	}

	return changeSets, nil
}
