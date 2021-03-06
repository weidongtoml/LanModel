package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	module = flag.String("module", "segmenter", "segmenter: chinese setence segmenter; ngram: ngram model.")
	action = flag.String("action", "segment", "evaluate: evaluate the given segmenter;segment: do segmentation.")
	corpus = flag.String("corpus", "../data/training/hkcu/hkcu_corpus.txt",
		"path to the segmented corpus file.")
	corpusCharSet = flag.String("corpus_charset", "Big5", "Encoding type of the corpus.")
	unigramModel  = flag.String("unigram", "../data/model/unigram.dat",
		"path to output the unigram model, use empty string to disable unigram.")
	bigramModel = flag.String("bigram", "../data/model/bigram.dat",
		"path to output the bigram model, use empty string to disable bigram.")
	dictPath = flag.String("dict", "../data/dict/cedict_ts.u8.txt",
		"path to locate the Chinese dictionary.")
)

func printHelp() {
	fmt.Printf("Usage: %s command\n", os.Args[0])
	fmt.Printf("\twhere command can be one of the following:\n")
	fmt.Printf("\tngram	for ngram generation and evaluation\n")
	fmt.Printf("\tsegment for Chinese sentence segmentation and evaluation\n")
}

func main() {
	flag.Parse()
	switch *module {
	case "ngram":
		doLanguageModel()
		break
	case "segmenter":
		doSegmenter()
		break
	default:
		printHelp()
	}
}
