package lucene

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
)

var (
	stopwords = map[string]bool{
		"the":  true,
		"a":    true,
		"an":   true,
		"on":   true,
		"in":   true,
		"and":  true,
		"or":   true,
		"that": true,
	}
	tokenRegexStr = "\\w+"
)

type LuceneTermFreq struct {
	DocumentID string
	Count      int
}

type LuceneIndex struct {
	index          map[string][]LuceneTermFreq
	docs           map[string]string
	doc_length     map[string]int
	avg_doc_length float64
	token_regex    *regexp.Regexp
}

type LuceneQueryDocument struct {
	Score      float64
	DocumentID string
	Document   string
}

func NewIndex() *LuceneIndex {
	r, _ := regexp.Compile(tokenRegexStr)

	return &LuceneIndex{
		index:          make(map[string][]LuceneTermFreq),
		docs:           make(map[string]string),
		doc_length:     make(map[string]int),
		avg_doc_length: 0.0,
		token_regex:    r,
	}
}

func (l *LuceneIndex) extractTerms(data string) []string {
	tokens := l.token_regex.FindAllString(strings.ToLower(data), -1)
	ret := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopwords[token]; !ok {
			ret = append(ret, token)
		}
	}

	return ret
}

func (l *LuceneIndex) Add(doc_id, data string) error {
	if _, ok := l.doc_length[doc_id]; ok {
		return fmt.Errorf("document already indexed")
	}

	terms := l.extractTerms(data)
	doc_length := len(data)
	l.docs[doc_id] = data

	term_count := make(map[string]int)
	for _, term := range terms {
		term_count[term]++
	}

	for term, count := range term_count {
		l.index[term] = append(l.index[term], LuceneTermFreq{
			DocumentID: doc_id,
			Count:      count,
		})
	}

	l.doc_length[doc_id] = doc_length
	l.avg_doc_length += (float64(doc_length) - l.avg_doc_length) / float64(len(l.docs))

	return nil
}

func (l *LuceneIndex) Search(query string, k1, b float64, total int) []*LuceneQueryDocument {
	terms := l.extractTerms(query)
	scores := make(map[string]float64)
	n := float64(len(l.docs))

	for _, term := range terms {
		postings, ok := l.index[term]
		if !ok {
			continue
		}

		df := float64(len(postings))
		idf := math.Log((n-df+0.5)/(df+0.5) + 1.0)

		for _, posting := range postings {
			dl := float64(l.doc_length[posting.DocumentID])

			numerator := float64(posting.Count) * (k1 + 1.0)
			denominator := float64(posting.Count) + k1*(1.0-b+b*dl/float64(l.avg_doc_length))
			scores[posting.DocumentID] += idf * (numerator / denominator)
		}
	}

	foundDocs := make([]*LuceneQueryDocument, 0, len(scores))
	for id, score := range scores {
		foundDocs = append(foundDocs, &LuceneQueryDocument{
			Score:      score,
			DocumentID: id,
		})
	}

	sort.SliceStable(foundDocs, func(i, j int) bool {
		return foundDocs[i].Score > foundDocs[j].Score
	})

	if total > len(foundDocs) {
		total = len(foundDocs)
	}

	ret := foundDocs[:total]
	for i := range ret {
		ret[i].Document = l.docs[ret[i].DocumentID]
	}

	return ret
}
