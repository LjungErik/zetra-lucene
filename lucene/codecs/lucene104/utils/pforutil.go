package utils

import (
	"errors"
	"math/bits"

	"github.com/LjungErik/zetra-lucene/lucene/internal"
)

var (
	ErrInvalidBlocksize = errors.New("ints array size smallar the required block size")
)

const (
	maxExceptions = 7
)

type PForUtil struct {
	forUtil *ForUtil
}

func NewPForUtil(forUtil *ForUtil) *PForUtil {
	return &PForUtil{
		forUtil: forUtil,
	}
}

func allEqual(l []uint32) bool {
	for i := 1; i < blockSize; i++ {
		if l[i] != l[0] {
			return false
		}
	}

	return true
}

func (f *PForUtil) Encode(ints []uint32, out internal.DataOutputStream) error {
	var (
		histogram       [32]int
		maxBitsRequired int = 0
	)

	if len(ints) < blockSize {
		return ErrInvalidBlocksize
	}

	for i := 0; i < blockSize; i++ {
		v := ints[i]
		bits := bits.LeadingZeros32(v)
		histogram[bits]++
		maxBitsRequired = max(maxBitsRequired, bits)
	}

	var (
		minBits              int = max(0, maxBitsRequired-8)
		cumulativeExceptions int = 0
		patchedBitsRequired  int = maxBitsRequired
		numExceptions        int = 0
	)

	for b := maxBitsRequired; b >= minBits; b-- {
		if cumulativeExceptions > maxExceptions {
			break
		}
		patchedBitsRequired = b
		numExceptions = cumulativeExceptions
		cumulativeExceptions += histogram[b]
	}

	var (
		maxUnpatchedValue uint32 = (1 << patchedBitsRequired) - 1
		exceptions        []byte = make([]byte, numExceptions*2)
	)

	if numExceptions > 0 {
		exceptionCount := 0
		for i := 0; i < blockSize; i++ {
			if ints[i] > maxUnpatchedValue {
				exceptions[exceptionCount*2] = byte(i)
				exceptions[exceptionCount*2+1] = byte(ints[i] >> uint32(patchedBitsRequired))
				ints[i] &= maxUnpatchedValue
				exceptionCount++
			}
		}
	}

	if allEqual(ints) && maxBitsRequired <= 8 {
		for i := 0; i < numExceptions; i++ {
			exceptions[2*i+1] = byte(uint32(exceptions[2*i+1]) << uint32(patchedBitsRequired))
		}
		if err := out.WriteByte(byte(numExceptions << 5)); err != nil {
			return err
		}
		if err := out.WriteVUInt32(ints[0]); err != nil {
			return err
		}
	} else {
		var token int = (numExceptions << 5) | patchedBitsRequired
		if err := out.WriteByte(byte(token)); err != nil {
			return err
		}

		if err := f.forUtil.Encode(ints, patchedBitsRequired, out); err != nil {
			return err
		}
	}

	if _, err := out.Write(exceptions); err != nil {
		return err
	}

	return nil
}
