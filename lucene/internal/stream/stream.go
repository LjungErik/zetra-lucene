package stream

import "github.com/LjungErik/zetra-lucene/lucene/internal"

func writeVInt(s internal.DataOutputStream, i int) error {
	for (i & ^0x7F) != 0 {
		if err := s.WriteByte(byte((i & 0x7F) | 0x80)); err != nil {
			return err
		}
		i = int(uint32(i) >> 7)
	}

	if err := s.WriteByte(byte(i)); err != nil {
		return err
	}

	return nil
}
