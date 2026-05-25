package codecs

type CodecsWriter interface {
	GetTermsWriter() TermsWriter
	GetPostingsWriter() PostingsWriter
	GetStoredFieldWriter()
}
