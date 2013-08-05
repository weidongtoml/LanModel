// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ngram_model provides implementation of the N-Gram language model
// generation from text corpus, loading of the model for prediction and a set
// of evaluation utility for the NGram language model.
package ngram_model

import (
	"common/util"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// NGramModel provides methods for loading and accessing NGram models.
type NGramModel struct {
	Unigram map[string]float64
	Bigram  map[BiGramKey]float64
}

// LoadNGramModel loads the ngram models from the given paths.
func LoadNGramModel(unigram, bigram string) (*NGramModel, error) {
	var m NGramModel
	if unigram != "" {
		if err := m.loadUnigram(unigram); err != nil {
			return nil, err
		}
	}
	if bigram != "" {
		if err := m.loadBigram(bigram); err != nil {
			return nil, err
		}
	}
	return &m, nil
}

func (m *NGramModel) loadUnigram(unigram string) error {
	m.Unigram = make(map[string]float64)
	return util.ForEachLineInFile(unigram, func(line string) (bool, error) {
		isValidLine := true
		err_msg := ""
		defer func() {
			if !isValidLine {
				log.Printf("Invalid line in unigram model[%s]: %s, %s",
					unigram, line, err_msg)
			}
		}()
		if !strings.HasPrefix(line, "#") {
			line = strings.Trim(line, " \t\n\r\f")
			fields := strings.Split(line, " ")
			if len(fields) != 2 {
				err_msg = fmt.Sprintf("Expected number of fields to be 2 but got %d.", len(fields))
				isValidLine = false
				return true, nil
			}
			p, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				isValidLine = false
				err_msg = fmt.Sprintf("attempt to obtain probability from %s failed: %s", fields[1], err)
				return true, nil
			}
			old, found := m.Unigram[fields[0]]
			if found {
				log.Printf("Warning, duplicated p for key %s, old value is %v, new value is %v",
					fields[0], old, p)
			}
			m.Unigram[fields[0]] = p
		}
		return true, nil
	})
}

func (m *NGramModel) loadBigram(bigram string) error {
	m.Bigram = make(map[BiGramKey]float64)
	return util.ForEachLineInFile(bigram, func(line string) (bool, error) {
		isValidLine := true
		defer func() {
			if !isValidLine {
				log.Printf("Invalid line in bigram model[%s]: %s",
					bigram, line)
			}
		}()
		if !strings.HasPrefix(line, "#") {
			line = strings.Trim(line, " \t\f\r")
			fields := strings.Split(line, " ")
			if len(fields) != 3 {
				isValidLine = false
				return true, nil
			}
			p, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				isValidLine = false
				return true, nil
			}
			key := BiGramKey{fields[0], fields[1]}
			old, found := m.Bigram[key]
			if found {
				log.Printf("Warning, duplicated p for key %v, old value is %v, new value is %v",
					key, old, p)
			}
			m.Bigram[key] = p
		}
		return true, nil
	})
}

type SimpleUnigramPredictor struct {
	model *NGramModel
}

func NewSimpleUnigramPredictor(m *NGramModel) *SimpleUnigramPredictor {
	if m.Unigram == nil {
		panic(fmt.Sprint("NewSimpleUnigramPredictor(%v), Unigram is nil.", *m))
	}
	return &SimpleUnigramPredictor{m}
}

func (p *SimpleUnigramPredictor) Probability(s []string) float64 {
	prob := float64(1)
	for _, w := range s {
		cur_p, found := p.model.Unigram[w]
		if found {
			prob *= cur_p
		} else {
			log.Printf("Missing term: %s", w)
		}
	}
	return prob
}

type SimpleBigramPredictor struct {
	model *NGramModel
}

func NewSimpleBigramPredictor(m *NGramModel) *SimpleBigramPredictor {
	if m.Bigram == nil {
		panic(fmt.Sprintf("NewSimpleBigramPredictor(%v), Bigram is nil.", *m))
	}
	return &SimpleBigramPredictor{m}
}

func (p *SimpleBigramPredictor) Probability(s []string) float64 {
	prob := float64(1)
	for i, _ := range s {
		var key BiGramKey
		if i == 0 {
			key = BiGramKey{SentenceStartTag, s[i]}
		} else {
			key = BiGramKey{s[i-1], s[i]}
		}
		cur_p, found := p.model.Bigram[key]
		if !found {
			log.Printf("Missing %v", key)
		} else {
			prob *= cur_p
		}
	}
	return prob
}
