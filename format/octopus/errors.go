package octopus

import "fmt"

type EmptyColumnNameError struct {
	Column *Column
}

func (e *EmptyColumnNameError) Error() string {
	return fmt.Sprintf("column name is empty. type: %s", e.Column.Type)
}

type InvalidAutoIncrementalError struct {
	Column *Column
}

func (e *InvalidAutoIncrementalError) Error() string {
	return fmt.Sprintf("column '%s' type '%s' cannot be autoIncremental", e.Column.Name, e.Column.Type)
}

type InvalidColumnTypeError struct {
	Column *Column
}

func (e *InvalidColumnTypeError) Error() string {
	return fmt.Sprintf("column '%s' type '%s' is invalid", e.Column.Name, e.Column.Type)
}

type InvalidColumnValuesError struct {
	Column *Column
	Msg    string
}

func (e *InvalidColumnValuesError) Error() string {
	return fmt.Sprintf("column '%s' has invalid values: %s", e.Column.Name, e.Msg)
}
