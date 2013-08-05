// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ngram_model

import (
	"code.google.com/p/mahonia"
	"common/util"
	"math"
	"strings"
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
	ResetIterator()
}

// Function Perplexity calculates the perplexity per word of the given language
// model under the given corpus, which is defined as
//   2 ^ - L, L = sum[s_i]{log(P(s_i))} / M,
// where M is the number of words in the corpus, P(s_i) is the probability of
// the sentence as predicted by the model.
func Perplexity(lm LanguageModeler, c CorpusSupplier) float64 {
	c.ResetIterator()
	sumLogP := float64(0)
	m := float64(0)
	for c.MoreSentence() {
		s := c.SegmentedSentence()
		p := lm.Probability(s)
		if p > 0 {
			sumLogP += math.Log(p)
		}
		m += float64(len(s))
	}
	return math.Exp2(-sumLogP / m)
}

type SegCNCorpus struct {
	sentences [][]string
	iter      int
	decoder   mahonia.Decoder
}

func NewSegCNCorpus(source_charset string) *SegCNCorpus {
	if source_charset != "" {
		return &SegCNCorpus{nil, 0, mahonia.NewDecoder(source_charset)}
	}
	return &SegCNCorpus{nil, 0, nil}

}

func (s *SegCNCorpus) MoreSentence() bool {
	return (*s).iter < len((*s).sentences)
}

func (s *SegCNCorpus) SegmentedSentence() []string {
	var sentence []string
	if (*s).iter < len((*s).sentences) {
		sentence = (*s).sentences[(*s).iter]
		(*s).iter++
	}
	return sentence
}

func (s *SegCNCorpus) ResetIterator() {
	(*s).iter = 0
}

func (s *SegCNCorpus) Load(path string) error {
	return util.ForEachLineInFile(path, func(line string) (bool, error) {
		line = strings.Trim(line, " \t\r\b\f")
		if (*s).decoder != nil {
			//TODO(weidoliang): Add conversion error check
			line = (*s).decoder.ConvertString(line)
		}
		(*s).sentences = append((*s).sentences, strings.Split(line, " "))
		return true, nil
	})
}
