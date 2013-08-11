// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chinese_segmenter

import (
	"fmt"
	//	"math"
	"ngram_model"
	"unicode/utf8"
)

type Segmenter struct {
	dict  *CEDict
	model *ngram_model.NGramModel
}

func NewSegmenter(dict *CEDict, model *ngram_model.NGramModel) *Segmenter {
	return &Segmenter{dict, model}
}

// Method Segment split the given sentence into a list of terms.
func (s *Segmenter) Segment(sentence string) []string {
	// Generate all possible term-splits and return the split with the highest
	// score.
	if sentence == "" {
		return nil
	}
	//return s.recursiveSegment(sentence, "")
	return s.dpSegment(sentence)
}

func (s *Segmenter) recursiveSegment(sentence string, token string) []string {
	rv, sz := utf8.DecodeRuneInString(sentence)
	if rv == utf8.RuneError {
		panic(fmt.Sprintf("Cannot decode utf8 string: %s", sentence))
	}
	r := string(rv)
	rest := sentence[sz:]

	fmt.Printf("%v, %v\n", r, rest)

	var join_result, split_result []string
	if rest == "" {
		if token == "" {
			split_result = append(join_result, r)
			join_result = split_result
		} else {
			join_result = append(join_result, token+r)
			split_result = append(split_result, token, r)
		}

	} else {
		s_r_result := s.recursiveSegment(rest, r)
		if token != "" {
			split_result = append(split_result, token)
		}
		for _, r_s := range s_r_result {
			split_result = append(split_result, r_s)
		}

		j_r_result := s.recursiveSegment(rest, token+r)
		for _, j_s := range j_r_result {
			join_result = append(join_result, j_s)
		}
	}
	split_score := s.scoreSequence(split_result)
	join_score := s.scoreSequence(join_result)
	if split_score > join_score {
		return split_result
	} else {
		return join_result
	}

}

func (s *Segmenter) scoreSequence(segments []string) float64 {
	prob := float64(1)
	for _, seg := range segments {
		prob *= s.model.Unigram[seg]
	}
	fmt.Printf("Score(%v) = %f\n", segments, prob)
	return prob
}

func (s *Segmenter) dpSegment(sentence string) []string {
	var chars []rune
	for s_iter := sentence; s_iter != ""; {
		r, sz := utf8.DecodeRuneInString(s_iter)
		if r == utf8.RuneError {
			//TODO(weidoliang): Change this to debug log
			panic(fmt.Sprintf("Cannot decode utf8 string: %s", sentence))
		}
		chars = append(chars, r)
		s_iter = s_iter[sz:]
	}
	n := len(chars)
	p := make([]float64, n+1) //best probability for substr i..n
	q := make([]string, n)
	skip := make([]int, n)
	p[n] = 1.0
	for i := n - 1; i >= 0; i-- {
		for j := i + 1; j <= n; j++ {
			word := chars[i:j]
			word_p := s.model.Unigram[string(word)]
			prev_p := p[j]
			if word_p <= 0.0 {
				continue
			}
			new_p := word_p * prev_p
			if new_p > p[i] {
				p[i] = new_p
				q[i] = string(word)
				skip[i] = len(word)

			}
		}
	}
	var segments []string
	for i := 0; i < len(q); {
		segments = append(segments, q[i])
		i += skip[i]
	}
	return segments
}
