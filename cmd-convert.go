package main

import (
	"errors"
	"fmt"
	"log"
)

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
	switch input.Format {
	case FORMAT_OCTOPUS:
		schema := &Schema{}
		if err := schema.FromFile(input.Filename); err != nil {
			return nil, nil
		}
		return schema, nil

	case FORMAT_STARUML2:
		staruml2 := &StarUML2{}
		if err := staruml2.FromFile(input.Filename); err != nil {
			return nil, err
		}
		return staruml2.ToSchema()
	default:
		return nil, fmt.Errorf("unhandled input format: %s", input.Format)
	}
}

func (cmd *ConvertCmd) schemaToOutput(schema *Schema, output *Output) error {
	switch output.Format {
	case FORMAT_OCTOPUS:
		return schema.ToFile(output.Filename)
	case FORMAT_QUICKDBD:
		quickdbd := &QuickDBD{}
		return quickdbd.ToFile(schema, output.Filename)
	default:
		return fmt.Errorf("unhandled output format: %s", output.Format)
	}
}
