package diff

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
)

type ChangeSetID struct {
	major int
	minor int
}

func newChangeSetID() *ChangeSetID {
	return &ChangeSetID{
		major: 0,
		minor: 0,
	}
}

func (l *ChangeSetID) bumpMajor() {
	l.major++
	l.minor = 0
}

func (l *ChangeSetID) revertMajor() {
	l.major--
	l.minor = 0
}

func (l *ChangeSetID) bumpMinor() string {
	l.minor++
	return l.version()
}

func (l *ChangeSetID) version() string {
	if l.minor == 0 {
		return fmt.Sprintf("%d", l.major)
	} else {
		return fmt.Sprintf("%d-%d", l.major, l.minor)
	}
}

type Result struct {
	From       *octopus.Schema
	To         *octopus.Schema
	ChangeSets []*ChangeSet
}

func (r *Result) Add(changeSet *ChangeSet) {
	r.ChangeSets = append(r.ChangeSets, changeSet)
}

type ChangeSet struct {
	ID      string
	Author  string
	Changes []interface{}
}

func (s *ChangeSet) Add(change interface{}) {
	s.Changes = append(s.Changes, change)
}

type CreateTable struct {
	Table *octopus.Table
}

type DropTable struct {
	Table *octopus.Table
}

type RenameTable struct {
	OldTable *octopus.Table
	NewTable *octopus.Table
}

type UpdatePrimaryKey struct {
	ConstraintName string
	OldTable       *octopus.Table
	NewTable       *octopus.Table
}

type UpdateUniqueConstraint struct {
	ConstraintName string
	Table          *octopus.Table
}

type SetTableComment struct {
	Table *octopus.Table
}

type AddColumn struct {
	Table        *octopus.Table
	Column       *octopus.Column
	BeforeColumn *octopus.Column
	AfterColumn  *octopus.Column
}

type DropColumn struct {
	Table      *octopus.Table
	ColumnName string
}

type SetColumnComment struct {
	Table  *octopus.Table
	Column *octopus.Column
}

type ChangeColumnType struct {
	Table     *octopus.Table
	OldColumn *octopus.Column
	NewColumn *octopus.Column
}

type RenameColumn struct {
	Table     *octopus.Table
	OldColumn *octopus.Column
	NewColumn *octopus.Column
}

type SetNotNullConstraint struct {
	Table  *octopus.Table
	Column *octopus.Column
}

type SetAutoIncrement struct {
	Table  *octopus.Table
	Column *octopus.Column
}

type SetDefaultValue struct {
	Table  *octopus.Table
	Column *octopus.Column
}
