package blocktree

import "github.com/LjungErik/zetra-lucene/lucene/internal"

type StatsWriter struct {
	out            internal.DataOutputStream
	hasFreqs       bool
	singletonCount int
}

func NewStatsWriter(out internal.DataOutputStream, hasFreqs bool) *StatsWriter {
	return &StatsWriter{
		out:            out,
		hasFreqs:       hasFreqs,
		singletonCount: 0,
	}
}

func (w *StatsWriter) Add(df uint32, ttf uint64) error {
	if df == 1 && (!w.hasFreqs && ttf == 1) {
		w.singletonCount++
	} else {
		if err := w.Finish(); err != nil {
			return err
		}

		if err := w.out.WriteVInt(int(df << 1)); err != nil {
			return err
		}

		if w.hasFreqs {
			if err := w.out.WriteVUInt64(ttf - uint64(df)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *StatsWriter) Finish() error {
	if w.singletonCount > 0 {
		if err := w.out.WriteVInt(((w.singletonCount - 1) << 1) | 1); err != nil {
			return err
		}

		w.singletonCount = 0
	}

	return nil
}
