package utils

import "github.com/LjungErik/zetra-lucene/lucene/internal"

const (
	CodecMagic  = 0x3fd76c17
	FooterMagic = ^CodecMagic
)

func WriteHeader(out *internal.OutputStream, codec string, version int) error {
	if err := writeBEInt(out, CodecMagic); err != nil {
		return err
	}

	if _, err := out.Write([]byte(codec)); err != nil {
		return err
	}

	if err := writeBEInt(out, version); err != nil {
		return err
	}

	return nil
}

func WriteIndexHeader(out *internal.OutputStream, codec string, version int, id []byte, suffix string) error {
	if err := WriteHeader(out, codec, version); err != nil {
		return err
	}

	if _, err := out.Write(id); err != nil {
		return err
	}

	if _, err := out.Write([]byte{byte(len(suffix))}); err != nil {
		return err
	}

	if _, err := out.Write([]byte(suffix)); err != nil {
		return err
	}

	return nil
}

func writeBEInt(out *internal.OutputStream, i int) error {
	_, err := out.Write([]byte{
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i & 0xFF),
	})

	return err
}
