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
	case FORMAT_OCTOPUS:
		reader = &Schema{}
		break
	case FORMAT_STARUML2:
		reader = &StarUML2{}
		break
	}

	if reader == nil {
		return nil, fmt.Errorf("unsupported input format: %s", input.Format)
	}
	if err := reader.FromFile(input.Filename); err != nil {
		return nil, err
	}
	return reader.ToSchema()
}

func (cmd *ConvertCmd) schemaToOutput(schema *Schema, output *Output) error {
	var writer FormatWriter

	switch output.Format {
	case FORMAT_OCTOPUS:
		return schema.ToFile(output.Filename)
	case FORMAT_DBDIAGRAM_IO:
		writer = &DBDiagramIO{}
		break
	case FORMAT_QUICKDBD:
		writer = &QuickDBD{}
		break
	}

	if writer == nil {
		return fmt.Errorf("unsupported output format: %s", output.Format)
	}
	return writer.ToFile(schema, output.Filename)
}
