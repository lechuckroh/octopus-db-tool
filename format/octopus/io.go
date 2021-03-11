package octopus

import (
	"fmt"
	"github.com/lechuckroh/octopus-db-tools/format/common"
)

func LoadSchema(filename string) (*Schema, error) {
	inputFormat := common.GetFileFormat(filename)
	if inputFormat != common.FormatOctopus1 && inputFormat != common.FormatOctopus2 {
		return nil, fmt.Errorf("'%s' is not octopus file", filename)
	}
	reader := &Schema{}
	if err := reader.FromFile(filename); err != nil {
		return nil, err
	}
	if schema, err := reader.ToSchema(); err != nil {
		return nil, err
	} else {
		schema.Normalize()
		return schema, nil
	}
}
