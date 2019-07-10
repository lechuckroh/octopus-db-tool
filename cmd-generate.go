package main

import "fmt"

type GenerateCmd struct {
}

func (cmd *GenerateCmd) Generate(input *Input, output *GenOutput) error {
	schema, err := (&ConvertCmd{}).inputToSchema(input)
	if err != nil {
		return err
	}

	switch output.Format {
	case FormatJpaKotlin:
		jpa := &JPAKotlin{}
		return jpa.Generate(schema, output)
	case FormatLiquibase:
		liquibase := &Liquibase{}
		return liquibase.Generate(schema, output)
	}

	return fmt.Errorf("unsupported output format: %s", output.Format)
}
