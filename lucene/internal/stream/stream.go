package stream

import "github.com/LjungErik/zetra-lucene/lucene/internal"

func writeVInt(s internal.DataOutputStream, i int) error {
	for (i & ^0x7F) != 0 {
		if err := s.WriteByte(byte((i & 0x7F) | 0x80)); err != nil {
			return err
		}
		i = int(uint32(i) >> 7) // Special case for handling writing of negative numbers
	}

	if err := s.WriteByte(byte(i)); err != nil {
		return err
	}

	return nil
}

func writeVUInt64(s internal.DataOutputStream, i uint64) error {
	for (i & ^uint64(0x7F)) != 0 {
		if err := s.WriteByte(byte((i & 0x7F) | 0x80)); err != nil {
			return err
		}
		i = i >> 7
	}

	if err := s.WriteByte(byte(i)); err != nil {
		return err
	}

	return nil
}
