package utils

import "github.com/LjungErik/zetra-lucene/lucene/internal"

const blockSize = 256

const blockSizeLog2 = 8

func mask32(bitsPerValue int) uint32 {
	return (1 << uint(bitsPerValue)) - 1
}

func mask16(bitsPerValue int) uint32 {
	m := (1 << uint(bitsPerValue)) - 1
	return expandMask16(uint32(m))
}

func mask8(bitsPerValue int) uint32 {
	m := (1 << uint(bitsPerValue)) - 1
	return expandMask8(uint32(m))
}

func expandMask16(m uint32) uint32 { return m | (m << 16) }
func expandMask8(m uint32) uint32  { return expandMask16(m | (m << 8)) }

var (
	masks8  [8]uint32
	masks16 [16]uint32
	masks32 [32]uint32
)

func init() {
	for i := 0; i < 8; i++ {
		masks8[i] = mask8(i)
	}
	for i := 0; i < 16; i++ {
		masks16[i] = mask16(i)
	}
	for i := 0; i < 32; i++ {
		masks32[i] = mask32(i)
	}
}

var (
	mask8_1 = masks8[1]
	mask8_2 = masks8[2]
	mask8_3 = masks8[3]
	mask8_4 = masks8[4]
	mask8_5 = masks8[5]
	mask8_6 = masks8[6]
	mask8_7 = masks8[7]

	mask16_1  = masks16[1]
	mask16_2  = masks16[2]
	mask16_3  = masks16[3]
	mask16_4  = masks16[4]
	mask16_5  = masks16[5]
	mask16_6  = masks16[6]
	mask16_7  = masks16[7]
	mask16_8  = masks16[8]
	mask16_9  = masks16[9]
	mask16_10 = masks16[10]
	mask16_11 = masks16[11]
	mask16_12 = masks16[12]
	mask16_13 = masks16[13]
	mask16_14 = masks16[14]
	mask16_15 = masks16[15]

	mask32_1  = masks32[1]
	mask32_2  = masks32[2]
	mask32_3  = masks32[3]
	mask32_4  = masks32[4]
	mask32_5  = masks32[5]
	mask32_6  = masks32[6]
	mask32_7  = masks32[7]
	mask32_8  = masks32[8]
	mask32_9  = masks32[9]
	mask32_10 = masks32[10]
	mask32_11 = masks32[11]
	mask32_12 = masks32[12]
	mask32_13 = masks32[13]
	mask32_14 = masks32[14]
	mask32_15 = masks32[15]
	mask32_16 = masks32[16]
)

func expand8(arr []uint32) {
	for i := 0; i < 64; i++ {
		v := arr[i]
		arr[i] = (v >> 24) & 0xFF
		arr[64+i] = (v >> 16) & 0xFF
		arr[128+i] = (v >> 8) & 0xFF
		arr[192+i] = v & 0xFF
	}
}

func collapse8(arr []uint32) {
	for i := 0; i < 64; i++ {
		arr[i] = (arr[i] << 24) | (arr[64+i] << 16) | (arr[128+i] << 8) | arr[192+i]
	}
}

func expand16(arr []uint32) {
	for i := 0; i < 128; i++ {
		v := arr[i]
		arr[i] = (v >> 16) & 0xFFFF
		arr[128+i] = v & 0xFFFF
	}
}

func collapse16(arr []uint32) {
	for i := 0; i < 128; i++ {
		arr[i] = (arr[i] << 16) | arr[128+i]
	}
}

type ForUtil struct {
	tmp [blockSize]uint32
}

func NumBytes(bitsPerValue int) int {
	return bitsPerValue << (blockSizeLog2 - 3)
}

func (f *ForUtil) Encode(ints []uint32, bitsPerValue int, out internal.DataOutputStream) error {
	var nextPrimitive int
	switch {
	case bitsPerValue <= 8:
		nextPrimitive = 8
		collapse8(ints)
	case bitsPerValue <= 16:
		nextPrimitive = 16
		collapse16(ints)
	default:
		nextPrimitive = 32
	}
	return encode(ints, bitsPerValue, nextPrimitive, out, f.tmp[:])
}

func encode(ints []uint32, bitsPerValue, primitiveSize int, out internal.DataOutputStream, tmp []uint32) error {
	numInts := blockSize * primitiveSize / 32
	numIntsPerShift := bitsPerValue * 8
	idx := 0
	shift := uint(primitiveSize - bitsPerValue)

	for i := 0; i < numIntsPerShift; i++ {
		tmp[i] = ints[idx] << shift
		idx++
	}
	for s := int(shift) - bitsPerValue; s >= 0; s -= bitsPerValue {
		for i := 0; i < numIntsPerShift; i++ {
			tmp[i] |= ints[idx] << uint(s)
			idx++
		}
	}

	remainingBitsPerInt := int(shift) - (int(shift)/bitsPerValue)*bitsPerValue + bitsPerValue
	{
		s := primitiveSize - bitsPerValue
		for s-bitsPerValue >= 0 {
			s -= bitsPerValue
		}
		remainingBitsPerInt = s + bitsPerValue
	}

	var maskRemainingBitsPerInt uint32
	switch primitiveSize {
	case 8:
		maskRemainingBitsPerInt = masks8[remainingBitsPerInt]
	case 16:
		maskRemainingBitsPerInt = masks16[remainingBitsPerInt]
	default:
		maskRemainingBitsPerInt = masks32[remainingBitsPerInt]
	}

	tmpIdx := 0
	remainingBitsPerValue := bitsPerValue
	for idx < numInts {
		if remainingBitsPerValue >= remainingBitsPerInt {
			remainingBitsPerValue -= remainingBitsPerInt
			tmp[tmpIdx] |= (ints[idx] >> uint(remainingBitsPerValue)) & maskRemainingBitsPerInt
			tmpIdx++
			if remainingBitsPerValue == 0 {
				idx++
				remainingBitsPerValue = bitsPerValue
			}
		} else {
			var m1, m2 uint32
			switch primitiveSize {
			case 8:
				m1 = masks8[remainingBitsPerValue]
				m2 = masks8[remainingBitsPerInt-remainingBitsPerValue]
			case 16:
				m1 = masks16[remainingBitsPerValue]
				m2 = masks16[remainingBitsPerInt-remainingBitsPerValue]
			default:
				m1 = masks32[remainingBitsPerValue]
				m2 = masks32[remainingBitsPerInt-remainingBitsPerValue]
			}
			tmp[tmpIdx] |= (ints[idx] & m1) << uint(remainingBitsPerInt-remainingBitsPerValue)
			idx++
			remainingBitsPerValue = bitsPerValue - remainingBitsPerInt + remainingBitsPerValue
			tmp[tmpIdx] |= (ints[idx] >> uint(remainingBitsPerValue)) & m2
			tmpIdx++
		}
	}

	for i := 0; i < numIntsPerShift; i++ {
		if err := out.WriteUInt32(tmp[i]); err != nil {
			return err
		}
	}
	return nil
}

// func (f *ForUtil) Decode(bitsPerValue int, pdu *PostingDecodingUtil, ints []uint32) error {
// 	switch bitsPerValue {
// 	case 1:
// 		if err := decode1(pdu, ints); err != nil {
// 			return err
// 		}
// 		expand8(ints)
// 	case 2:
// 		if err := decode2(pdu, ints); err != nil {
// 			return err
// 		}
// 		expand8(ints)
// 	case 3:
// 		if err := decode3(pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 		expand8(ints)
// 	case 4:
// 		if err := decode4(pdu, ints); err != nil {
// 			return err
// 		}
// 		expand8(ints)
// 	case 5:
// 		if err := decode5(pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 		expand8(ints)
// 	case 6:
// 		if err := decode6(pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 		expand8(ints)
// 	case 7:
// 		if err := decode7(pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 		expand8(ints)
// 	case 8:
// 		if err := decode8(pdu, ints); err != nil {
// 			return err
// 		}
// 		expand8(ints)
// 	case 9:
// 		if err := decode9(pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 		expand16(ints)
// 	case 10:
// 		if err := decode10(pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 		expand16(ints)
// 	case 11:
// 		if err := decode11(pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 		expand16(ints)
// 	case 12:
// 		if err := decode12(pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 		expand16(ints)
// 	case 13:
// 		if err := decode13(pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 		expand16(ints)
// 	case 14:
// 		if err := decode14(pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 		expand16(ints)
// 	case 15:
// 		if err := decode15(pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 		expand16(ints)
// 	case 16:
// 		if err := decode16(pdu, ints); err != nil {
// 			return err
// 		}
// 		expand16(ints)
// 	default:
// 		if err := decodeSlow(bitsPerValue, pdu, f.tmp[:], ints); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func decodeSlow(bitsPerValue int, pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	numInts := bitsPerValue << 3
// 	msk := masks32[bitsPerValue]

// 	if err := pdu.SplitInts(numInts, ints, 32-bitsPerValue, 32, msk, tmp, 0, 0xFFFFFFFF); err != nil {
// 		return err
// 	}

// 	remainingBitsPerInt := 32 - bitsPerValue
// 	mask32Rem := masks32[remainingBitsPerInt]

// 	tmpIdx := 0
// 	remainingBits := remainingBitsPerInt
// 	for intsIdx := numInts; intsIdx < blockSize; intsIdx++ {
// 		b := bitsPerValue - remainingBits
// 		l := (tmp[tmpIdx] & masks32[remainingBits]) << uint(b)
// 		tmpIdx++
// 		for b >= remainingBitsPerInt {
// 			b -= remainingBitsPerInt
// 			l |= (tmp[tmpIdx] & mask32Rem) << uint(b)
// 			tmpIdx++
// 		}
// 		if b > 0 {
// 			l |= (tmp[tmpIdx] >> uint(remainingBitsPerInt-b)) & masks32[b]
// 			remainingBits = remainingBitsPerInt - b
// 		} else {
// 			remainingBits = remainingBitsPerInt
// 		}
// 		ints[intsIdx] = l
// 	}
// 	return nil
// }

// func decode1(pdu *PostingDecodingUtil, ints []uint32) error {
// 	return pdu.SplitInts(8, ints, 7, 1, mask8_1, ints, 56, mask8_1)
// }

// func decode2(pdu *PostingDecodingUtil, ints []uint32) error {
// 	return pdu.SplitInts(16, ints, 6, 2, mask8_2, ints, 48, mask8_2)
// }

// func decode3(pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	if err := pdu.SplitInts(24, ints, 5, 3, mask8_3, tmp, 0, mask8_2); err != nil {
// 		return err
// 	}
// 	for iter, tmpIdx, intsIdx := 0, 0, 48; iter < 8; iter, tmpIdx, intsIdx = iter+1, tmpIdx+3, intsIdx+2 {
// 		l0 := tmp[tmpIdx+0] << 1
// 		l0 |= (tmp[tmpIdx+1] >> 1) & mask8_1
// 		ints[intsIdx+0] = l0
// 		l1 := (tmp[tmpIdx+1] & mask8_1) << 2
// 		l1 |= tmp[tmpIdx+2] << 0
// 		ints[intsIdx+1] = l1
// 	}
// 	return nil
// }

// func decode4(pdu *PostingDecodingUtil, ints []uint32) error {
// 	return pdu.SplitInts(32, ints, 4, 4, mask8_4, ints, 32, mask8_4)
// }

// func decode5(pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	if err := pdu.SplitInts(40, ints, 3, 5, mask8_5, tmp, 0, mask8_3); err != nil {
// 		return err
// 	}
// 	for iter, tmpIdx, intsIdx := 0, 0, 40; iter < 8; iter, tmpIdx, intsIdx = iter+1, tmpIdx+5, intsIdx+3 {
// 		l0 := tmp[tmpIdx+0] << 2
// 		l0 |= (tmp[tmpIdx+1] >> 1) & mask8_2
// 		ints[intsIdx+0] = l0

// 		l1 := (tmp[tmpIdx+1] & mask8_1) << 4
// 		l1 |= tmp[tmpIdx+2] << 1
// 		l1 |= (tmp[tmpIdx+3] >> 2) & mask8_1
// 		ints[intsIdx+1] = l1

// 		l2 := (tmp[tmpIdx+3] & mask8_2) << 3
// 		l2 |= tmp[tmpIdx+4] << 0
// 		ints[intsIdx+2] = l2
// 	}
// 	return nil
// }

// func decode6(pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	if err := pdu.SplitInts(48, ints, 2, 6, mask8_6, tmp, 0, mask8_2); err != nil {
// 		return err
// 	}
// 	for iter, tmpIdx, intsIdx := 0, 0, 48; iter < 16; iter, tmpIdx, intsIdx = iter+1, tmpIdx+3, intsIdx+1 {
// 		l0 := tmp[tmpIdx+0] << 4
// 		l0 |= tmp[tmpIdx+1] << 2
// 		l0 |= tmp[tmpIdx+2] << 0
// 		ints[intsIdx+0] = l0
// 	}
// 	return nil
// }

// func decode7(pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	if err := pdu.SplitInts(56, ints, 1, 7, mask8_7, tmp, 0, mask8_1); err != nil {
// 		return err
// 	}
// 	for iter, tmpIdx, intsIdx := 0, 0, 56; iter < 8; iter, tmpIdx, intsIdx = iter+1, tmpIdx+7, intsIdx+1 {
// 		l0 := tmp[tmpIdx+0] << 6
// 		l0 |= tmp[tmpIdx+1] << 5
// 		l0 |= tmp[tmpIdx+2] << 4
// 		l0 |= tmp[tmpIdx+3] << 3
// 		l0 |= tmp[tmpIdx+4] << 2
// 		l0 |= tmp[tmpIdx+5] << 1
// 		l0 |= tmp[tmpIdx+6] << 0
// 		ints[intsIdx+0] = l0
// 	}
// 	return nil
// }

// func decode8(pdu *PostingDecodingUtil, ints []uint32) error {
// 	return pdu.In.ReadUint32s(ints[:64])
// }

// func decode9(pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	if err := pdu.SplitInts(72, ints, 7, 9, mask16_9, tmp, 0, mask16_7); err != nil {
// 		return err
// 	}
// 	for iter, tmpIdx, intsIdx := 0, 0, 72; iter < 8; iter, tmpIdx, intsIdx = iter+1, tmpIdx+9, intsIdx+7 {
// 		l0 := tmp[tmpIdx+0] << 2
// 		l0 |= (tmp[tmpIdx+1] >> 5) & mask16_2
// 		ints[intsIdx+0] = l0

// 		l1 := (tmp[tmpIdx+1] & mask16_5) << 4
// 		l1 |= (tmp[tmpIdx+2] >> 3) & mask16_4
// 		ints[intsIdx+1] = l1

// 		l2 := (tmp[tmpIdx+2] & mask16_3) << 6
// 		l2 |= (tmp[tmpIdx+3] >> 1) & mask16_6
// 		ints[intsIdx+2] = l2

// 		l3 := (tmp[tmpIdx+3] & mask16_1) << 8
// 		l3 |= tmp[tmpIdx+4] << 1
// 		l3 |= (tmp[tmpIdx+5] >> 6) & mask16_1
// 		ints[intsIdx+3] = l3

// 		l4 := (tmp[tmpIdx+5] & mask16_6) << 3
// 		l4 |= (tmp[tmpIdx+6] >> 4) & mask16_3
// 		ints[intsIdx+4] = l4

// 		l5 := (tmp[tmpIdx+6] & mask16_4) << 5
// 		l5 |= (tmp[tmpIdx+7] >> 2) & mask16_5
// 		ints[intsIdx+5] = l5

// 		l6 := (tmp[tmpIdx+7] & mask16_2) << 7
// 		l6 |= tmp[tmpIdx+8] << 0
// 		ints[intsIdx+6] = l6
// 	}
// 	return nil
// }

// func decode10(pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	if err := pdu.SplitInts(80, ints, 6, 10, mask16_10, tmp, 0, mask16_6); err != nil {
// 		return err
// 	}
// 	for iter, tmpIdx, intsIdx := 0, 0, 80; iter < 16; iter, tmpIdx, intsIdx = iter+1, tmpIdx+5, intsIdx+3 {
// 		l0 := tmp[tmpIdx+0] << 4
// 		l0 |= (tmp[tmpIdx+1] >> 2) & mask16_4
// 		ints[intsIdx+0] = l0

// 		l1 := (tmp[tmpIdx+1] & mask16_2) << 8
// 		l1 |= tmp[tmpIdx+2] << 2
// 		l1 |= (tmp[tmpIdx+3] >> 4) & mask16_2
// 		ints[intsIdx+1] = l1

// 		l2 := (tmp[tmpIdx+3] & mask16_4) << 6
// 		l2 |= tmp[tmpIdx+4] << 0
// 		ints[intsIdx+2] = l2
// 	}
// 	return nil
// }

// func decode11(pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	if err := pdu.SplitInts(88, ints, 5, 11, mask16_11, tmp, 0, mask16_5); err != nil {
// 		return err
// 	}
// 	for iter, tmpIdx, intsIdx := 0, 0, 88; iter < 8; iter, tmpIdx, intsIdx = iter+1, tmpIdx+11, intsIdx+5 {
// 		l0 := tmp[tmpIdx+0] << 6
// 		l0 |= tmp[tmpIdx+1] << 1
// 		l0 |= (tmp[tmpIdx+2] >> 4) & mask16_1
// 		ints[intsIdx+0] = l0

// 		l1 := (tmp[tmpIdx+2] & mask16_4) << 7
// 		l1 |= tmp[tmpIdx+3] << 2
// 		l1 |= (tmp[tmpIdx+4] >> 3) & mask16_2
// 		ints[intsIdx+1] = l1

// 		l2 := (tmp[tmpIdx+4] & mask16_3) << 8
// 		l2 |= tmp[tmpIdx+5] << 3
// 		l2 |= (tmp[tmpIdx+6] >> 2) & mask16_3
// 		ints[intsIdx+2] = l2

// 		l3 := (tmp[tmpIdx+6] & mask16_2) << 9
// 		l3 |= tmp[tmpIdx+7] << 4
// 		l3 |= (tmp[tmpIdx+8] >> 1) & mask16_4
// 		ints[intsIdx+3] = l3

// 		l4 := (tmp[tmpIdx+8] & mask16_1) << 10
// 		l4 |= tmp[tmpIdx+9] << 5
// 		l4 |= tmp[tmpIdx+10] << 0
// 		ints[intsIdx+4] = l4
// 	}
// 	return nil
// }

// func decode12(pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	if err := pdu.SplitInts(96, ints, 4, 12, mask16_12, tmp, 0, mask16_4); err != nil {
// 		return err
// 	}
// 	for iter, tmpIdx, intsIdx := 0, 0, 96; iter < 32; iter, tmpIdx, intsIdx = iter+1, tmpIdx+3, intsIdx+1 {
// 		l0 := tmp[tmpIdx+0] << 8
// 		l0 |= tmp[tmpIdx+1] << 4
// 		l0 |= tmp[tmpIdx+2] << 0
// 		ints[intsIdx+0] = l0
// 	}
// 	return nil
// }

// func decode13(pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	if err := pdu.SplitInts(104, ints, 3, 13, mask16_13, tmp, 0, mask16_3); err != nil {
// 		return err
// 	}
// 	for iter, tmpIdx, intsIdx := 0, 0, 104; iter < 8; iter, tmpIdx, intsIdx = iter+1, tmpIdx+13, intsIdx+3 {
// 		l0 := tmp[tmpIdx+0] << 10
// 		l0 |= tmp[tmpIdx+1] << 7
// 		l0 |= tmp[tmpIdx+2] << 4
// 		l0 |= tmp[tmpIdx+3] << 1
// 		l0 |= (tmp[tmpIdx+4] >> 2) & mask16_1
// 		ints[intsIdx+0] = l0

// 		l1 := (tmp[tmpIdx+4] & mask16_2) << 11
// 		l1 |= tmp[tmpIdx+5] << 8
// 		l1 |= tmp[tmpIdx+6] << 5
// 		l1 |= tmp[tmpIdx+7] << 2
// 		l1 |= (tmp[tmpIdx+8] >> 1) & mask16_2
// 		ints[intsIdx+1] = l1

// 		l2 := (tmp[tmpIdx+8] & mask16_1) << 12
// 		l2 |= tmp[tmpIdx+9] << 9
// 		l2 |= tmp[tmpIdx+10] << 6
// 		l2 |= tmp[tmpIdx+11] << 3
// 		l2 |= tmp[tmpIdx+12] << 0
// 		ints[intsIdx+2] = l2
// 	}
// 	return nil
// }

// func decode14(pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	if err := pdu.SplitInts(112, ints, 2, 14, mask16_14, tmp, 0, mask16_2); err != nil {
// 		return err
// 	}
// 	for iter, tmpIdx, intsIdx := 0, 0, 112; iter < 16; iter, tmpIdx, intsIdx = iter+1, tmpIdx+7, intsIdx+1 {
// 		l0 := tmp[tmpIdx+0] << 12
// 		l0 |= tmp[tmpIdx+1] << 10
// 		l0 |= tmp[tmpIdx+2] << 8
// 		l0 |= tmp[tmpIdx+3] << 6
// 		l0 |= tmp[tmpIdx+4] << 4
// 		l0 |= tmp[tmpIdx+5] << 2
// 		l0 |= tmp[tmpIdx+6] << 0
// 		ints[intsIdx+0] = l0
// 	}
// 	return nil
// }

// func decode15(pdu *PostingDecodingUtil, tmp, ints []uint32) error {
// 	if err := pdu.SplitInts(120, ints, 1, 15, mask16_15, tmp, 0, mask16_1); err != nil {
// 		return err
// 	}
// 	for iter, tmpIdx, intsIdx := 0, 0, 120; iter < 8; iter, tmpIdx, intsIdx = iter+1, tmpIdx+15, intsIdx+1 {
// 		l0 := tmp[tmpIdx+0] << 14
// 		l0 |= tmp[tmpIdx+1] << 13
// 		l0 |= tmp[tmpIdx+2] << 12
// 		l0 |= tmp[tmpIdx+3] << 11
// 		l0 |= tmp[tmpIdx+4] << 10
// 		l0 |= tmp[tmpIdx+5] << 9
// 		l0 |= tmp[tmpIdx+6] << 8
// 		l0 |= tmp[tmpIdx+7] << 7
// 		l0 |= tmp[tmpIdx+8] << 6
// 		l0 |= tmp[tmpIdx+9] << 5
// 		l0 |= tmp[tmpIdx+10] << 4
// 		l0 |= tmp[tmpIdx+11] << 3
// 		l0 |= tmp[tmpIdx+12] << 2
// 		l0 |= tmp[tmpIdx+13] << 1
// 		l0 |= tmp[tmpIdx+14] << 0
// 		ints[intsIdx+0] = l0
// 	}
// 	return nil
// }

// func decode16(pdu *PostingDecodingUtil, ints []uint32) error {
// 	return pdu.In.ReadUint32s(ints[:128])
// }
