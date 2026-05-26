package codecs

type CodecsWriter interface {
	GetPostingsFormat() PostingsFormat
	GetStoredFieldsFormat() StoredFieldsFormat
}
