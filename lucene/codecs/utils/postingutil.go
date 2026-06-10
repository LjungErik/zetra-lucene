package utils

// import "github.com/LjungErik/zetra-lucene/lucene/internal"

// // ---------------------------------------------------------------------------
// // PostingDecodingUtil — mirrors Lucene's PostingDecodingUtil
// // ---------------------------------------------------------------------------

// // PostingDecodingUtil wraps a DataInput and provides the splitInts helper
// // used by all decode paths.
// type PostingDecodingUtil struct {
// 	In internal.DataInputStream
// }

// // SplitInts reads `count` uint32 values into b[0:count], then unpacks them:
// //
// //  1. If rightBitsPerInt % bitsPerValue != 0, the bottom remainder bits of
// //     each packed lane are extracted into c[cOff:] and the values are shifted
// //     right by that remainder width.
// //  2. Any full bitsPerValue-wide layers that fit in the remaining right bits
// //     are extracted into b[count:], b[2*count:], … (bottom layer at the
// //     highest offset, top layer at the lowest).
// //  3. The topmost bitsPerValue bits are left in b[0:count], masked with bMask.
// func (pdu *PostingDecodingUtil) SplitInts(
// 	count int,
// 	b []uint32,
// 	bShift int,
// 	dec int,
// 	bMask uint32,
// 	c []uint32,
// 	cIndex int,
// 	cMask uint32,
// ) error {
// 	if err := pdu.In.ReadUInts(c[:count]); err != nil {
// 		return err
// 	}

// 	remainderWidth := rightBitsPerInt % bitsPerValue
// 	fullChunks := rightBitsPerInt / bitsPerValue

// 	// 1. Extract remainder bits from the bottom into c.
// 	if remainderWidth > 0 {
// 		for i := 0; i < count; i++ {
// 			c[cOff+i] = b[i] & cMask
// 			b[i] >>= uint(remainderWidth)
// 		}
// 	}

// 	// 2. Extract full bitsPerValue-wide layers (bottom-up) into b at
// 	//    ascending offsets above count.
// 	for chunk := fullChunks - 1; chunk >= 0; chunk-- {
// 		off := count * (chunk + 1)
// 		for i := 0; i < count; i++ {
// 			b[off+i] = b[i] & bMask
// 			b[i] >>= uint(bitsPerValue)
// 		}
// 	}

// 	// 3. Top bits remain in b[0:count].
// 	for i := 0; i < count; i++ {
// 		b[i] &= bMask
// 	}
// 	return nil
// }
