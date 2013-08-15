// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package chinese_segmenter provides facility for Chinese language sentence
// segmentation.
//
// TODO(weidoliang)
//	1. Based on the implementation of unigram model, extend it to support NGram
//	2. Make use of the punctuation present in the sentence to divide it into
//		mini-sentences before doing segmentation to speed up the process.
package chinese_segmenter

import (
	"common/util"
	"errors"
	"fmt"
	"ngram_model"
)

const (
	LongestWordLen = 4
)

type Segmenter struct {
	dict  *CEDict
	model *ngram_model.NGramModel
}

// Function NewSegmenter creates a new segmenter from the given dictionary and
// language model.
func NewSegmenter(dict *CEDict, model *ngram_model.NGramModel) *Segmenter {
	return &Segmenter{dict, model}
}

// Method Segment split the given sentence into a list of terms, assuming the
// sentence given is in utf8 encoding.
func (s *Segmenter) Segment(sentence string) ([]string, error) {
	if sentence == "" {
		return nil, errors.New("Empty sentence cannot be segmented.")
	}
	// Convert to array of Chinese Characters.
	chars, err := util.Ut8StringToRuneArray(sentence)
	if err != nil {
		return nil, err
	}
	// Use dynamic programming to find the best segmentation.
	n := len(chars)
	p := make([]float64, n+1) //best probability for substr i..n
	q := make([]string, n)
	skip := make([]int, n)
	p[n] = 1.0
	for i := n - 1; i >= 0; i-- {
		skip[i] = 1
		for j := i + 1; j <= n; j++ {
			word := chars[i:j]
			word_c_count := len(word)
			if word_c_count > LongestWordLen {
				continue
			}
			word_str := string(word)
			word_p := s.model.Unigram[word_str]
			prev_p := p[j]
			if word_p <= 0.0 {
				continue
			}

			new_p := word_p * prev_p
			if new_p > p[i] {
				p[i] = new_p
				q[i] = word_str
				skip[i] = word_c_count
			}
		}
	}

	// Retrieve the segmentation result.
	segments := retrieveTerms(q, skip, 0, -1)
	return segments, nil
}

func retrieveTerms(terms []string, skip []int, start_index, limit int) []string {
	var segments []string
	cnt := 0
	for i := start_index; i < len(terms) && (limit <= 0 || cnt < limit); {
		segments = append(segments, terms[i])
		i += skip[i]
		cnt++
	}
	return segments
}

func (s *Segmenter) evalSegmentation(cur_term string, terms_follow []string, is_start bool) float64 {
	p := float64(1)
	if len(terms_follow) > 0 {
		if is_start {
			key := ngram_model.BiGramKey{ngram_model.SentenceStartTag, cur_term}
			p = s.model.Bigram[key]
			if p == 0.0 {
				p = s.model.Unigram[cur_term]
			}
		} else {
			key := ngram_model.BiGramKey{cur_term, terms_follow[0]}
			p = s.model.Bigram[key]
			if p == 0.0 {
				p = s.model.Unigram[cur_term]
			}
		}
	}
	fmt.Printf("%s, %v, %v, %f\n", cur_term, terms_follow, is_start, p)
	return p
	//return s.model.Unigram[cur_term]
}
