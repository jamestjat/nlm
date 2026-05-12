package api

import (
	"reflect"
	"strings"
	"testing"
)

func TestDetectMIMEType(t *testing.T) {
	tests := []struct {
		name         string
		content      []byte
		filename     string
		providedType string
		want         string
	}{
		{
			name:     "XML file with .xml extension",
			content:  []byte(`<?xml version="1.0"?><root><item>test</item></root>`),
			filename: "test.xml",
			want:     "text/xml",
		},
		{
			name:     "XML file without extension",
			content:  []byte(`<?xml version="1.0"?><root><item>test</item></root>`),
			filename: "test",
			want:     "text/xml",
		},
		{
			name:         "XML file with provided type",
			content:      []byte(`<?xml version="1.0"?><root><item>test</item></root>`),
			filename:     "test.xml",
			providedType: "application/xml",
			want:         "application/xml",
		},
		{
			name:     "EPUB file uses NotebookLM upload MIME type",
			content:  []byte("PK\x03\x04epub payload"),
			filename: "book.epub",
			want:     "application/epub+zip",
		},
		{
			name:     "PDF file keeps upload MIME type",
			content:  []byte("%PDF-1.7\n..."),
			filename: "paper.pdf",
			want:     "application/pdf",
		},
		{
			name:     "DOCX file keeps upload MIME type",
			content:  []byte("PK\x03\x04docx payload"),
			filename: "report.docx",
			want:     "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectMIMEType(tt.content, tt.filename, tt.providedType)
			// Strip charset for comparison
			gotType := strings.Split(got, ";")[0]
			if gotType != tt.want {
				t.Errorf("detectMIMEType() = %v, want %v", gotType, tt.want)
			}
		})
	}
}

func TestBuildRegisterFileSourceArgs(t *testing.T) {
	for _, filename := range []string{"book.epub", "paper.pdf", "report.docx"} {
		t.Run(filename, func(t *testing.T) {
			got := buildRegisterFileSourceArgs("project-123", filename)
			want := []interface{}{
				[]interface{}{
					[]interface{}{filename},
				},
				"project-123",
				[]interface{}{2},
				[]interface{}{
					1, nil, nil, nil, nil, nil, nil, nil, nil, nil,
					[]interface{}{1},
				},
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("buildRegisterFileSourceArgs() = %#v, want %#v", got, want)
			}
		})
	}
}
