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
	inputSchema, err := cmd.inputToSchema(input)
	if err != nil {
		return err
	}
	log.Printf("[READ] %s\n", input.Filename)

	// Write Output
	err = cmd.schemaToOutput(inputSchema, output)
	if err == nil {
		log.Printf("[WRITE] %s\n", output.Filename)
	}
	return err
}

func (cmd *ConvertCmd) inputToSchema(input *Input) (*Schema, error) {
	var reader FormatReader

	switch input.Format {
	case FormatOctopus:
		reader = &Schema{}
	case FormatStaruml2:
		reader = &StarUML2{}
	case FormatXlsx:
		reader = &Xlsx{}
	}

	if reader == nil {
		return nil, fmt.Errorf("unsupported input format: %s", input.Format)
	}
	if err := reader.FromFile(input.Filename); err != nil {
		return nil, err
	}
	if schema, err := reader.ToSchema(); err != nil {
		return nil, err
	} else {
		schema.Normalize()
		return schema, nil
	}
}

func (cmd *ConvertCmd) schemaToOutput(schema *Schema, output *Output) error {
	var writer FormatWriter

	switch output.Format {
	case FormatOctopus:
		return schema.ToFile(output.Filename)
	case FormatDbdiagramIo:
		writer = &DBDiagramIO{}
	case FormatPlantuml:
		writer = &PlantUML{}
	case FormatQuickdbd:
		writer = &QuickDBD{}
	case FormatXlsx:
		writer = &Xlsx{}
	}

	if writer == nil {
		return fmt.Errorf("unsupported output format: %s", output.Format)
	}
	return writer.ToFile(schema, output.Filename)
}
