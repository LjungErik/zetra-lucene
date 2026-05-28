package filenames

import "strings"

func SegmentFileName(segmentName, segmentSuffix, ext string) string {
	if len(ext) > 0 || len(segmentSuffix) > 0 {
		var sb strings.Builder
		sb.Grow(len(segmentName) + 2 + len(segmentSuffix) + len(ext))
		sb.WriteString(segmentName)
		if len(segmentSuffix) > 0 {
			sb.WriteByte('_')
			sb.WriteString(segmentSuffix)
		}
		if len(ext) > 0 {
			sb.WriteByte('.')
			sb.WriteString(ext)
		}
		return sb.String()
	}
	return segmentName
}
