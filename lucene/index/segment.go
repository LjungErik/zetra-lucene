package index

type Segments struct {
	NextSegment   int               `json:"next_segment"`
	SegmentCount  int               `json:"segment_count"`
	Segments      []SegmentMetadata `json:"segments"`
	LatestSegment int               `json:"latest:segment"`
}

type SegmentMetadata struct {
	SegmentName string `json:"segment_name"`
	SegmentID   int    `json:"segment_id"`
}

type SegementWriteState struct {
	Segments Segments
}
