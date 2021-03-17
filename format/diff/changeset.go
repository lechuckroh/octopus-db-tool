package diff

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/octopus"
)

// TODO: set immutable struct
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

// ===========================================================================

type ChangeSet struct {
	ID      string
	Author  string
	Changes []Change
}

func (s *ChangeSet) Add(change Change) {
	s.Changes = append(s.Changes, change)
}

// ===========================================================================

type Change interface {
	DepTable() *octopus.Table
}

type CreateTable struct {
	Table *octopus.Table
}

func (c *CreateTable) DepTable() *octopus.Table {
	return c.Table
}

type DropTable struct {
	Table *octopus.Table
}

func (c *DropTable) DepTable() *octopus.Table {
	return c.Table
}

type RenameTable struct {
	OldTable *octopus.Table
	NewTable *octopus.Table
}

func (c *RenameTable) DepTable() *octopus.Table {
	return c.OldTable
}

type UpdatePrimaryKey struct {
	ConstraintName string
	OldTable       *octopus.Table
	NewTable       *octopus.Table
}

func (c *UpdatePrimaryKey) DepTable() *octopus.Table {
	return c.OldTable
}

type DropUniqueConstraint struct {
	ConstraintName string
	Table          *octopus.Table
}

func (c *DropUniqueConstraint) DepTable() *octopus.Table {
	return c.Table
}

type CreateUniqueConstraint struct {
	ConstraintName string
	Table          *octopus.Table
}

func (c *CreateUniqueConstraint) DepTable() *octopus.Table {
	return c.Table
}

type SetTableComment struct {
	Table *octopus.Table
}

func (c *SetTableComment) DepTable() *octopus.Table {
	return c.Table
}

type AddColumn struct {
	Table        *octopus.Table
	Column       *octopus.Column
	BeforeColumn *octopus.Column
	AfterColumn  *octopus.Column
}

func (c *AddColumn) DepTable() *octopus.Table {
	return c.Table
}

type DropColumn struct {
	Table      *octopus.Table
	ColumnName string
}

func (c *DropColumn) DepTable() *octopus.Table {
	return c.Table
}

type SetColumnComment struct {
	Table  *octopus.Table
	Column *octopus.Column
}

func (c *SetColumnComment) DepTable() *octopus.Table {
	return c.Table
}

type ChangeColumnType struct {
	Table     *octopus.Table
	OldColumn *octopus.Column
	NewColumn *octopus.Column
}

func (c *ChangeColumnType) DepTable() *octopus.Table {
	return c.Table
}

type RenameColumn struct {
	Table     *octopus.Table
	OldColumn *octopus.Column
	NewColumn *octopus.Column
}

func (c *RenameColumn) DepTable() *octopus.Table {
	return c.Table
}

type SetNotNullConstraint struct {
	Table  *octopus.Table
	Column *octopus.Column
}

func (c *SetNotNullConstraint) DepTable() *octopus.Table {
	return c.Table
}

type SetAutoIncrement struct {
	Table  *octopus.Table
	Column *octopus.Column
}

func (c *SetAutoIncrement) DepTable() *octopus.Table {
	return c.Table
}

type SetDefaultValue struct {
	Table  *octopus.Table
	Column *octopus.Column
}

func (c *SetDefaultValue) DepTable() *octopus.Table {
	return c.Table
}
