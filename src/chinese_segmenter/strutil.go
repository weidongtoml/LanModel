// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chinese_segmenter

import "strings"

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
