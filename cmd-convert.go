package main

import (
	"errors"
	"fmt"
	"log"
)

type FormatReader interface {
	FromFile(filename string) error
	ToSchema() (*Schema, error)
}

type FormatWriter interface {
	ToFile(schema *Schema, filename string) error
}

type ConvertCmd struct {
}

func (cmd *ConvertCmd) Convert(input *Input, output *Output) error {
	if input == nil {
		return errors.New("input is nil")
	}
	if output == nil {
		return errors.New("output is nil")
	}

	// Read Input
	inputSchema, err := input.ToSchema()
	if err != nil {
		return err
	}
	log.Printf("[READ] %s\n", input.Filename)

	// Write Output
	err = cmd.schemaToOutput(inputSchema, output)
	if err == nil {
		log.Printf("[WRITE] %s\n", output.FilePath)
	}
	return err
}

func (cmd *ConvertCmd) schemaToOutput(schema *Schema, output *Output) error {
	var writer FormatWriter

	switch output.Format {
	case FormatOctopus:
		return schema.ToFile(output.FilePath)
	case FormatDbdiagramIo:
		writer = &DBDiagramIO{}
	case FormatPlantuml:
		writer = &PlantUML{}
	case FormatQuickdbd:
		writer = &QuickDBD{}
	case FormatXlsx:
		writer = &Xlsx{
			UseNotNullColumn: output.GetBool(FlagNotNull),
		}
	case FormatSqlMysql:
		writer = &Mysql{}
	}

	if writer == nil {
		return fmt.Errorf("unsupported output format: %s", output.Format)
	}
	return writer.ToFile(schema, output.FilePath)
}
