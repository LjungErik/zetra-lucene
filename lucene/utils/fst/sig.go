package fst

// func commonPrefixLen(s1, s2 string) int {
// 	n := len(s1)
// 	if len(s2) < n {
// 		n = len(s2)
// 	}

// 	for i := range n {
// 		if s1[i] != s2[i] {
// 			return i
// 		}
// 	}

// 	return n
// }

// func (b *Builder) prevSignature() string {
// 	n := len(b.registry)
// 	if n == 0 {
// 		return ""
// 	}

// 	return b.registry[n-1].key
// }

// func (b *Builder) nextSignature(key string) string {
// 	// find shared prefix between current key and last key
// 	prefixLen := commonPrefixLen(key, b.lastKey)
// 	if prefixLen == 0 {
// 		prefixLen = commonPrefixLen(b.prevSignature(), b.lastKey)
// 		return b.lastKey[:prefixLen+1]
// 	}

// 	sig := b.lastKey[:prefixLen]
// 	if sig == b.prevSignature() {
// 		sig = b.lastKey[:prefixLen+1]
// 	}

// 	return sig
// }

// func (b *Builder) freeze(key string) {
// 	if b.lastKey == "" {
// 		return
// 	}

// 	b.registry = append(b.registry, entry{
// 		offset: b.lastOffset,
// 		key:    b.nextSignature(key),
// 	})
// }
