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
	case FormatGraphql:
		graphql := &Graphql{}
		return graphql.Generate(schema, output)
	case FormatJpaKotlin:
		jpa := NewJPAKotlin()
		return jpa.Generate(schema, output, false)
	case FormatJpaKotlinData:
		jpa := &JPAKotlin{}
		return jpa.Generate(schema, output, true)
	case FormatLiquibase:
		liquibase := &Liquibase{}
		return liquibase.Generate(schema, output)
	}

	return fmt.Errorf("unsupported output format: %s", output.Format)
}
