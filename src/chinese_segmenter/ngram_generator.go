// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chinese_segmenter

import (
	"bufio"
	"fmt"
	"strings"
)

// BiGramKey is used to represent the bigram "First Second"
type BiGramKey struct {
	First  string
	Second string
}

const (
	SentenceStartTag string = "SENT_START"
)

// NGramGenerator implements the necessary methods for processing segmented
// text files and generater the Mono Gram and BiGram frequency.
type NGramGenerator struct {
	uniGram      map[string]int
	uniGramCount int
	biGram       map[BiGramKey]int
	biGramCount  int
}

// Function NEwNGramGenerator Creates a new Initialzied NGramGenerator.
func NewNGramGenerator() *NGramGenerator {
	var generator NGramGenerator
	generator.uniGram = make(map[string]int)
	generator.biGram = make(map[BiGramKey]int)
	return &generator
}

// Function ProcessFile process the given file and incorporate the information
// into the NGramGenerator g for future N-Gram model generation.
func (g *NGramGenerator) ProcessFile(filename string) error {
	lineProcessor := func(line string) (bool, error) {
		tokens := strings.Split(line, " ")
		var prevToken string
		for i, t := range tokens {
			//Monogram frequency
			g.uniGram[t]++
			g.uniGramCount++
			//Bigram frequency
			if i == 0 {
				g.biGram[BiGramKey{SentenceStartTag, t}]++
			} else {
				g.biGram[BiGramKey{prevToken, t}]++
			}
			g.biGramCount++
			prevToken = t
		}
		return true, nil
	}
	return ForEachLineInFile(filename, lineProcessor)
}

// Method GenerateUnigramModel generates a unigram from the information collected
// so far and save it to the given file. The file format of the model consists
// of multiple lines of unigram and unigram frequency seperated by space.
func (g *NGramGenerator) GenerateUnigramModel(filename string) error {
	modelWriter := func(w *bufio.Writer) error {
		for k, c := range g.uniGram {
			//TODO(weidoliang): add smoothing to avoid zero probabilities
			p_k := float64(c) / float64(g.uniGramCount)
			w.WriteString(fmt.Sprintf("%s %f\n", k, p_k))
		}
		return nil
	}
	return WithNewOpenFileAsBufioWriter(filename, modelWriter)
}

// Method GenerateBigramModel generates a bigram model from the inforamtion
// collected so far and save it to the given file. The file format of the model
// consists of multiple lines, each line contains first term, second term,
// and bigram frequency representing P(first term | second term), i.e. the
// probability of second term immediately follows the first term as seen from
// the processed documents.
func (g *NGramGenerator) GenerateBigramModel(filename string) error {
	modelWriter := func(w *bufio.Writer) error {
		for k, c := range g.biGram {
			//TODO(weidoliang): add smoothing
			p_k := float64(c) / float64(g.biGramCount)
			w.WriteString(fmt.Sprintf("%s %s %f\n", k.First, k.Second, p_k))
		}
		return nil
	}
	return WithNewOpenFileAsBufioWriter(filename, modelWriter)
}
