package formatter

import (
	"encoding/json"
	"io"
	"reflect"
	"sentire/pkg/models"
)

// NDJSONFormatter outputs data as newline-delimited JSON (one object per line)
type NDJSONFormatter struct {
	writer io.Writer
	fields []string
}

// NewNDJSONFormatter creates a new NDJSON formatter
func NewNDJSONFormatter(writer io.Writer, fields []string) *NDJSONFormatter {
	return &NDJSONFormatter{writer: writer, fields: fields}
}

func (f *NDJSONFormatter) FormatEvent(event *models.Event) error {
	return f.writeLine(event)
}

func (f *NDJSONFormatter) FormatEvents(events []models.Event) error {
	for _, e := range events {
		if err := f.writeLine(e); err != nil {
			return err
		}
	}
	return nil
}

func (f *NDJSONFormatter) FormatIssue(issue *models.Issue) error {
	return f.writeLine(issue)
}

func (f *NDJSONFormatter) FormatIssues(issues []models.Issue) error {
	for _, i := range issues {
		if err := f.writeLine(i); err != nil {
			return err
		}
	}
	return nil
}

func (f *NDJSONFormatter) FormatProject(project *models.Project) error {
	return f.writeLine(project)
}

func (f *NDJSONFormatter) FormatProjects(projects []models.Project) error {
	for _, p := range projects {
		if err := f.writeLine(p); err != nil {
			return err
		}
	}
	return nil
}

func (f *NDJSONFormatter) FormatOrgStats(stats *models.OrganizationStats) error {
	return f.writeLine(stats)
}

func (f *NDJSONFormatter) FormatGeneric(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if err := f.writeLine(v.Index(i).Interface()); err != nil {
				return err
			}
		}
		return nil
	}
	return f.writeLine(data)
}

func (f *NDJSONFormatter) writeLine(data interface{}) error {
	data = filterFields(data, f.fields)
	return json.NewEncoder(f.writer).Encode(data)
}
