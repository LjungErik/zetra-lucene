package document

import (
	"container/list"
	"iter"
)

type Document interface {
	Add(field DocumentField)
	Iter() iter.Seq[DocumentField]
}

type IndexableDocument struct {
	fields list.List
}

var _ Document = (*IndexableDocument)(nil)

func NewDocument() *IndexableDocument {
	return &IndexableDocument{
		fields: *list.New(),
	}
}

func (d *IndexableDocument) Add(field DocumentField) {
	d.fields.PushBack(field)
}

func (d *IndexableDocument) Iter() iter.Seq[DocumentField] {
	return func(yield func(DocumentField) bool) {
		for e := d.fields.Front(); e != nil; e = e.Next() {
			if !yield(e.Value.(DocumentField)) {
				return
			}
		}
	}
}
