package fields

type FieldsInfo struct {
	hasFreq         bool
	hasPostings     bool
	hasProx         bool
	hasPayload      bool
	hasOffsets      bool
	hasTermVectors  bool
	hasNorms        bool
	hasDocValues    bool
	hasPointValues  bool
	hasVectorValues bool
}

func NewFieldsInfo() *FieldsInfo {
	return &FieldsInfo{}
}

func (f *FieldsInfo) HasFreq() bool {
	return f.hasFreq
}

func (f *FieldsInfo) HasPostings() bool {
	return f.hasPostings
}

func (f *FieldsInfo) HasProx() bool {
	return f.hasProx
}

func (f *FieldsInfo) HasPayload() bool {
	return f.hasPayload
}

func (f *FieldsInfo) HasOffsets() bool {
	return f.hasOffsets
}

func (f *FieldsInfo) HasTermVectors() bool {
	return f.hasTermVectors
}

func (f *FieldsInfo) HasNorms() bool {
	return f.hasNorms
}

func (f *FieldsInfo) HasDocValues() bool {
	return f.hasDocValues
}

func (f *FieldsInfo) HasPointValues() bool {
	return f.hasPointValues
}

func (f *FieldsInfo) HasVectorValues() bool {
	return f.hasVectorValues
}
