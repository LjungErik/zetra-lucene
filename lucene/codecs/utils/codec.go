package utils

import (
	"fmt"

	"github.com/LjungErik/zetra-lucene/lucene/internal"
)

const (
	CodecMagic  = 0x3fd76c17
	FooterMagic = ^CodecMagic
)

func WriteHeader(out internal.DataOutputStream, codec string, version int) error {
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

func WriteIndexHeader(out internal.DataOutputStream, codec string, version int, id []byte, suffix string) error {
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

func WriteFooter(out internal.DataOutputStream) error {
	if err := writeBEInt(out, FooterMagic); err != nil {
		return err
	}

	if err := writeBEInt(out, 0); err != nil {
		return err
	}

	if err := writeCRC(out); err != nil {
		return err
	}

	return nil
}

func writeBEInt(out internal.DataOutputStream, i int) error {
	_, err := out.Write([]byte{
		byte(i >> 24),
		byte(i >> 16),
		byte(i >> 8),
		byte(i & 0xFF),
	})

	if err != nil {
		return err
	}

	return nil
}

func writeBEInt64(out internal.DataOutputStream, i int64) error {
	if err := writeBEInt(out, int(i>>32)); err != nil {
		return err
	}

	if err := writeBEInt(out, int(i)); err != nil {
		return err
	}

	return nil
}

func writeCRC(out internal.DataOutputStream) error {
	value := out.GetCheckSum()
	if (value & uint64(0xFFFFFFFF00000000)) != 0 {
		return fmt.Errorf("Illegal CRC-32 checksum: %d", value)
	}

	if err := writeBEInt64(out, int64(value)); err != nil {
		return err
	}

	return nil
}
