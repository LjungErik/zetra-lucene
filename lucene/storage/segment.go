package storage

type Segments struct {
	NextSegment  int64             `json:"next_segment"`
	SegmentCount int               `json:"segment_count"`
	Segments     []SegmentMetadata `json:"segments"`
}

type SegmentMetadata struct {
	SegmentName string `json:"segment_name"`
	SegmentID   int64  `json:"segment_id"`
}
