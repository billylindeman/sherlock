package sherlock

import (
	"crypto/md5"
	"fmt"
	"testing"
)

type schemaTestDoc struct {
	Title string `sherlock:"term,weight=10"`
	Body  string `sherlock:"term"`

	meta string `sherlock:"-"`
}

func (d *schemaTestDoc) ID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(d.Title)))
}

func TestSchemaFromStruct(t *testing.T) {
	doc := schemaTestDoc{
		Title: "test",
		Body:  "this is a big test",
	}
	schema, err := NewSchemaFromStruct(doc)
	if err != nil {
		t.Fatal("Error creating schema from document: ", err)
	}

	fmt.Println("Schema: ", schema)

	schema.analyze(doc)

}
