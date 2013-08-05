// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package chinese_segmenter provides ways to train
// and apply Chinese segmentation program.
package chinese_segmenter

import (
	"common/util"
	"fmt"
	"log"
	"strings"
)

// CEDict content associated with a key.
type Term struct {
	Alternative string   // Alternative representation
	Pinyin      []string // pinyin of each Chinese character
	English     []string // Possible English translation
}

// CEDict implements methods for loading and retriving
// information from CC-CEDict, the community maintained
// free Chinese-English dictionary.
type CEDict struct {
	keyType  int
	keyTerms map[string]Term
}

// Flags to pass to CEDict.Load to determine the lookup key type.
const (
	SIMPLE_CHINESE    int = 0 // Use simplified Chinese words as lookup keys.
	TRADITION_CHINESE int = 1 // Use traditional Chinese words as lookup key.
)

// Function Lookup returns the Term associated with the given key.
func (dict *CEDict) Lookup(key string) *Term {
	term, found := dict.keyTerms[key]
	if !found {
		return nil
	}
	return &term
}

type cedictFields struct {
	traditional string
	simplified  string
	pinyin      []string
	english     []string
}

func fieldExtractor(line string) interface{} {
	keysAndRest := strings.SplitN(line, " ", 3)
	if len(keysAndRest) != 3 {
		log.Printf("Invalid line (expected at least 2 keys): %s", line)
		return nil
	}
	pinyinAndRest := strings.SplitN(keysAndRest[2], "]", 2)
	if len(pinyinAndRest) != 2 {
		log.Printf("Invalid line (expected pinyin field): %s (%v)\n", line, pinyinAndRest)
		return nil
	}
	var english []string
	for _, e := range strings.Split(pinyinAndRest[1], "/") {
		s := strings.Trim(e, " \n\t")
		if s != "" {
			english = append(english, s)
		}
	}
	return cedictFields{
		keysAndRest[0],
		keysAndRest[1],
		strings.Split(pinyinAndRest[0][1:], " "),
		english,
	}
}

// Method LoadCEDict loads the CC-CEDict from the given path.
func LoadCEDict(path string, keyType int) (*CEDict, error) {
	lineHandler := util.NewPrefixDispatcher([]util.PrefixHandler{
		{"# ", nil},          //skip comments
		{"#! ", nil},         //skip meta information
		{"", fieldExtractor}, //default: extract fields
	})
	dict := CEDict{
		keyType,
		make(map[string]Term),
	}
	lineProcessor := func(line string) (bool, error) {
		fieldsI := lineHandler.Process(line)
		if fieldsI != nil {
			fields, fieldsOk := fieldsI.(cedictFields)
			if !fieldsOk {
				panic("Logic Error in the program, expected line handler to return key & term.")
			}
			var key, alter string
			switch dict.keyType {
			case SIMPLE_CHINESE:
				key = fields.simplified
				alter = fields.traditional
				break
			case TRADITION_CHINESE:
				key = fields.traditional
				alter = fields.simplified
				break
			default:
				panic(fmt.Sprintf("Invalid key type value: %d", dict.keyType))
				break
			}
			term := Term{
				alter,
				fields.pinyin,
				fields.english,
			}
			if oldTerm, found := dict.keyTerms[key]; found {
				log.Printf("Found duplicate definition for key %s, old value is %v, new value is %v",
					key, oldTerm, term)
			}
			dict.keyTerms[key] = term
		}
		return true, nil
	}
	if err := util.ForEachLineInFile(path, lineProcessor); err != nil {
		return nil, err
	}
	return &dict, nil
}
