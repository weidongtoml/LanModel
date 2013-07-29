// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chinese_segmenter

import (
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
	return ForEachLineInFile(unigram, func(line string) (bool, error) {
		isValidLine := true
		defer func() {
			if !isValidLine {
				log.Printf("Invalid line in unigram model[%s]: %s",
					unigram, line)
			}
		}()
		if !strings.HasPrefix(line, "#") {
			fields := strings.SplitN(line, " ", 1)
			if len(fields) != 2 {
				isValidLine = false
				return true, nil
			}
			p, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				isValidLine = false
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
	return ForEachLineInFile(bigram, func(line string) (bool, error) {
		isValidLine := true
		defer func() {
			if !isValidLine {
				log.Printf("Invalid line in bigram model[%s]: %s",
					bigram, line)
			}
		}()
		if !strings.HasPrefix(line, "#") {
			fields := strings.SplitN(line, " ", 1)
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
		prob *= p.model.Unigram[w]
	}
	return prob
}
