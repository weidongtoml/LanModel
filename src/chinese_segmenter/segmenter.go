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
	LongestWordLen = 4 // The maximum length of the longest word possible,
	// which is mainly used to speed up segmentation.
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
//
// Implementation note:
// Current implemetation uses unigram and maximum entropy model together with
// dynamic programming to implement the segmentation.
// Here, from the set of all possible segmentation of sentence S, denoted by
// seg(S), we find the set G such that the joint probability is the maximum,
// i.e. the maximum likelihood obtain under the unigram model is given by
//
//   M(S) = max_(G \member_of seg(S)) {Product_(1 <= k <= |G|){P(g_k} }
//
//        = max_(1 <= i <= |S|) {
//               max_(G' \member_of seg(S[i+1..n]){
//                       P(S[1..i]) * Product_(1 <= k <= |G'|){P(g'_k} }  }
//
//        = max_(1 <= i <= |S|) {
//               P(S[1..i]) * max_(G' \member_of seg(S[i+1..n]){
//                                     Product_(1 <= k <= |G'|){P(g'_k} } }
//
//        = max_(i <= i <= |S|) { P(S[1..i]) * M(S[i+1..n]) }
//
// From the above, we can implement the segmentation in the following way:
// First, we segment S to two parts at j = argmax(S[1..n]),
//   then we obtain the first word segment S[1..j], and remaining sentence
//   S[j+1..n],
// repeat the above process with S[j+1..n] as S to obtain the second, third,
// and the rest to the word segments until sentence S is completely segmented.
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

	// p[i] is the  probability of the most likely sequence of words in
	// substr[i..n], i.e. M(S[i+1..n])
	p := make([]float64, n+1)

	// q[i] is the first word in that sentence
	q := make([]string, n)

	skip := make([]int, n)
	p[n] = 1.0
	// Calculate M(S[i..n]), and the best binary segmentation point j.
	for i := n - 1; i >= 0; i-- {
		skip[i] = 1
		// Calculate:
		//   p[i] = max_(i+1 <= j <= n){ P(S[i..j]) * M(S[j+1..n]) } )
		//   q[i] = S[i..j], where j results in the maximum value of p[i]
		//   skip[i] = len(q[i])
		for j := i + 1; j <= n; j++ {
			word := chars[i:j]
			word_len := len(word)
			if word_len > LongestWordLen {
				// Longer than the maximum possible, skip it
				continue
			}
			word_str := string(word)
			word_p := s.model.Unigram[word_str] // P(S[i..j])
			if word_p <= 0.0 {
				// Unknown word, assign a small probability
				word_p = 0.000000000000001
			}
			new_p := word_p * p[j] // P(S[i..j]) * M(S[j+1..n])
			//fmt.Printf("P(%s)*M(%s) = %f [%d]\n", word_str, string(chars[j:]), new_p, i)

			if new_p > p[i] {
				p[i] = new_p
				q[i] = word_str
				skip[i] = word_len
			}
		}
		//fmt.Printf("%v: %v %v %v\n", i, p[i], q[i], skip[i])
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
