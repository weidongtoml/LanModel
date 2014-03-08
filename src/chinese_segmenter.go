// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	seg "chinese_segmenter"
	"common/util"
	"fmt"
	"ngram_model"
	"os"
	"strings"
)

func doSegmenter() {
	switch *action {
	case "evaluate":
		evaluateSegmenter()
		break
	case "segment":
		doSegment()
		break
	default:
		break
	}
}

func doSegment() {
	model, err := ngram_model.LoadNGramModel(*unigramModel, *bigramModel)
	if err != nil {
		fmt.Printf("Failed to load model[%s,%s]: %s", *unigramModel, *bigramModel, err)
		return
	}
	segmenter := seg.NewSegmenter(nil, model)
	reader := bufio.NewReader(os.Stdin)
	for {
		line, has_more, err := reader.ReadLine()
		if err == nil {
			clean_line := strings.Trim(string(line), " \t\n\r\f")
			if len(line) > 0 {
				result, ok := segmenter.Segment(clean_line)
				if ok == nil {
					fmt.Printf("%v\n", result)
				}
			}
		}
		if !has_more {
			break
		}
	}
}

func evaluateSegmenter() {
	//cedict, err := LoadCEDict(cedict_path, cedict_key_type)
	//if err != nil {
	//	t.Fatalf("Failed to load CEDict[%s]: %s", cedict_path, err)
	//}
	model, err := ngram_model.LoadNGramModel(*unigramModel, *bigramModel)
	if err != nil {
		fmt.Printf("Failed to load model[%s,%s]: %s", *unigramModel, *bigramModel, err)
		return
	}
	segmenter := seg.NewSegmenter(nil, model)
	converter := util.NewUtf8Converter(*corpusCharSet)
	err = util.ForEachLineInFile(*corpus, func(line string) (bool, error) {
		line = converter.ConvertString(strings.Trim(line, " \t\n\r\f"))
		sample := strings.Replace(line, " ", "", -1)
		exp_result := strings.Split(line, " ")
		result, _ := segmenter.Segment(sample)

		is_eqv := len(result) == len(exp_result)
		for i, r := range result {
			if r != exp_result[i] {
				is_eqv = false
				break
			}
		}
		if !is_eqv {
			fmt.Printf("Segment(%s) expect result to be:\n%v\nbut got:\n%v\n\n",
				sample, exp_result, result)
		}
		return true, nil
	})
	if err != nil {
		fmt.Printf("Error encountered when attempting to evaluate segmenter: %s", err)
	}
}
