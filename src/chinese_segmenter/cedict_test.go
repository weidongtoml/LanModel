// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chinese_segmenter

import (
	"bufio"
	"os"
	"testing"
)

// Function AreStringSlicesEqual check to see if the given
// two string slices are equal or not.
func AreStringSlicesEqual(a, b []string) bool {
	isEqv := len(a) == len(b)
	if isEqv {
		for i, _ := range a {
			if a[i] != b[i] {
				isEqv = false
				break
			}
		}
	}
	return isEqv
}

func TestCEDict(t *testing.T) {
	lines := []string{
		"AA制 AA制 [A A zhi4] /to split the bill/to go Dutch/",
		"A咖 A咖 [A ka1] /class \"A\"/top grade/",
		"A片 A片 [A pian4] /adult movie/pornography/",
		"B型超聲 B型超声 [B xing2 chao1 sheng1] /type-B ultrasound/",
		"B超 B超 [B chao1] /type-B ultrasound/abbr. for B型超聲|B型超声[B xing2 chao1 sheng1]/",
		"C盤 C盘 [C pan2] /C drive or default startup drive (computing)/",
		"DNA鑒定 DNA鉴定 [D N A jian4 ding4] /DNA test/DNA testing/",
		"E仔 E仔 [e zai3] /MDMA (C11H15NO2)/",
		"G點 G点 [G dian3] /Gräfenberg Spot/G-Spot/",
		"K仔 K仔 [K zai3] /ketamine (slang)/",
	}
	dictResult := []struct {
		traditional string
		simplified  string
		pinyin      []string
		english     []string
	}{
		{"AA制", "AA制", []string{"A", "A", "zhi4"}, []string{"to split the bill", "to go Dutch"}},
		{"A咖", "A咖", []string{"A", "ka1"}, []string{"class \"A\"", "top grade"}},
		{"E仔", "E仔", []string{"e", "zai3"}, []string{"MDMA (C11H15NO2)"}},
	}

	dictPath := "cedict_test.txt"
	lineWriter := func(writer *bufio.Writer) error {
		for _, line := range lines {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				t.Errorf("Failed to write line[%s] to file [%s]", line, dictPath)
				return err
			}
		}
		return nil
	}
	err := WithNewOpenFileAsBufioWriter(dictPath, lineWriter)
	if err != nil {
		t.Errorf("Failed to create test file %s", dictPath)
		return
	}
	defer func() {
		os.Remove(dictPath)
	}()

	dict, err := LoadCEDict(dictPath, TRADITION_CHINESE)
	if err != nil {
		t.Errorf("Failed at LoadCEDict(%s)", dictPath)
	}

	for _, t_case := range dictResult {
		term := dict.Lookup(t_case.traditional)
		if term == nil {
			t.Errorf("Failed to Lookup key %s", t_case.traditional)
		}
		if term.Alternative != t_case.simplified {
			t.Errorf("Lookup returns incorrent alternative representation, expected [%v] but got [%v]",
				t_case.simplified, term.Alternative)
		}
		if !AreStringSlicesEqual(term.Pinyin, t_case.pinyin) {
			t.Errorf("Lookup returns incorrect pinyin representation. Expected [%v] but got [%v]",
				t_case.pinyin, term.Pinyin)
		}
		if !AreStringSlicesEqual(term.English, t_case.english) {
			t.Errorf("Lookup returns invalid english representation expected [%v] but got [%v]",
				t_case.english, term.English)
		}
	}
}
