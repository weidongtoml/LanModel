// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chinese_segmenter

import (
	"math"
)

// Interface LanguageModeler specifies a method for calculating the predicted
// probability of the given segmented sentence.
type LanguageModeler interface {
	Probability([]string) float64
}

// Interface CorpusSupplier specifies methods for accessing corpus inforamtion.
type CorpusSupplier interface {
	MoreSentence() bool
	SegmentedSentence() []string
}

// Function Perplexity calculates the perplexity per word of the given language
// model under the given corpus, which is defined as
//   2 ^ - L, L = sum[s_i]{log(P(s_i))} / M,
// where M is the number of words in the corpus, P(s_i) is the probability of
// the sentence as predicted by the model.
func Perplexity(lm LanguageModeler, c CorpusSupplier) float64 {
	sumLogP := float64(0)
	m := float64(0)
	for c.MoreSentence() {
		s := c.SegmentedSentence()
		sumLogP += math.Log(lm.Probability(s))
		m += float64(len(s))
	}
	return math.Exp2(-sumLogP / m)
}
