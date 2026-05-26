package index

import "iter"

type Fields interface {
	Iter() iter.Seq[string]
	Terms(field string) Terms
}
