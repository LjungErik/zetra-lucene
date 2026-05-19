package document

type DocumentField interface {
	Name() string
	Type() DocumentFieldType
	Stored() bool
	Indexable() bool
	ValueAsString() string
	ValueLength() int
}

type DocumentFieldType int

const (
	Text DocumentFieldType = iota
	Numeric
)
