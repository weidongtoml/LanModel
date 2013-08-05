// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	cn_seg "chinese_segmenter"
	"flag"
	"fmt"
	"log"
)

var (
	action = flag.String("action", "evaluate", "evaluate: evaluate the given models;create: create ngram models.")
	corpus = flag.String("corpus", "../data/training/hkcu/hkcu_corpus.txt",
		"path to the segmented corpus file.")
	unigramModel = flag.String("unigram", "../data/model/unigram.dat",
		"path to output the unigram model, use empty string to disable unigram.")
	bigramModel = flag.String("bigram", "../data/model/bigram.dat",
		"path to output the bigram model, use empty string to disable bigram.")
)

func createNGramModel() {
	generator := cn_seg.NewNGramGenerator("Big5")
	err := generator.ProcessFile(*corpus)
	if err != nil {
		log.Printf("Failed to process corpus[%s]: %s", *corpus, err)
	} else {
		if *unigramModel != "" {
			err = generator.GenerateUnigramModel(*unigramModel)
			if err != nil {
				log.Printf("Failed to generate unigram model[%s]: %s",
					*unigramModel, err)
			}
		}
		if *bigramModel != "" {
			err = generator.GenerateBigramModel(*bigramModel)
			if err != nil {
				log.Printf("Failed to generate bigram model[%s]: %s",
					*bigramModel, err)
			}
		}
	}
}

func evaluateNGramModel() {
	model, err := cn_seg.LoadNGramModel(*unigramModel, *bigramModel)
	if err != nil {
		log.Printf("Failed to load NGram model: %s", err)
	}
	corpusSupplier := cn_seg.NewSegCNCorpus("Big5")
	err = (&corpusSupplier).Load(*corpus)
	if err != nil {
		log.Printf("Failed to load Corpus [%s]: %s", *corpus, err)
	} else {
		predictor := cn_seg.NewSimpleUnigramPredictor(model)
		perplexity := cn_seg.Perplexity(predictor, corpusSupplier)
		fmt.Printf("Unigram Model Perplexity: %f\n", perplexity)
		bigram_predictor := cn_seg.NewSimpleBigramPredictor(model)
		bigram_perplexity := cn_seg.Perplexity(bigram_predictor, corpusSupplier)
		fmt.Printf("Bigram Model Perplexity: %f\n", bigram_perplexity)
	}
}

func main() {
	flag.Parse()
	switch *action {
	case "create":
		createNGramModel()
		break
	case "evaluate":
		evaluateNGramModel()
		break
	default:
		log.Printf("Invalid action option: %s, expected [create|evaluate]", *action)
		break
	}
}
