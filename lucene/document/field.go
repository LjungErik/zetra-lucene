package document

type DocumentField interface {
	Name() string
	Type() DocumentFieldType
	Stored() bool
	Indexable() bool
	ValueAsString() string
}

type DocumentFieldType int

const (
	Text DocumentFieldType = iota
	Numeric
)
