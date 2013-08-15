// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"code.google.com/p/mahonia"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// A PrefixHandler specifies the handler function associated with the given
// prefix.
type PrefixHandler struct {
	Prefix  string
	Handler func(string) interface{}
}

// PrefixDispatcher provides facility for dispatching the corresponding method
// based on the prefix of the string.
type PrefixDispatcher struct {
	prefixHandlers []PrefixHandler
}

// NewPrefixDispatcher creates a PrefixDispather from the given prefix handlers.
func NewPrefixDispatcher(hanlders []PrefixHandler) *PrefixDispatcher {
	return &PrefixDispatcher{prefixHandlers: hanlders}
}

// Process the given string by dispatching the corresponding method associated
// with the prefix found in the string.
func (p *PrefixDispatcher) Process(s string) interface{} {
	var ret interface{}
	for _, h := range p.prefixHandlers {
		if h.Prefix == "" {
			if h.Handler != nil {
				ret = h.Handler(s)
			}
			break
		}
		if strings.HasPrefix(s, h.Prefix) {
			if h.Handler != nil {
				ret = h.Handler(s)
			}
			break
		}
	}
	return ret
}

// Function Utf8StringToRuneArray converts the given utf8 string to an array of
// rune.
func Ut8StringToRuneArray(sentence string) ([]rune, error) {
	var chars []rune
	for s_iter := sentence; s_iter != ""; {
		r, sz := utf8.DecodeRuneInString(s_iter)
		if r == utf8.RuneError {
			return nil, errors.New(fmt.Sprintf("Cannot decode utf8 substring: %s", s_iter))
		}
		chars = append(chars, r)
		s_iter = s_iter[sz:]
	}
	return chars, nil
}

type Utf8Converter struct {
	mahonia.Decoder
}

// Function NewUtf8Converter creates a Utf8 converter that can convert the
// source_charset to Utf8.
func NewUtf8Converter(source_charset string) *Utf8Converter {
	return &Utf8Converter{mahonia.NewDecoder(source_charset)}
}
