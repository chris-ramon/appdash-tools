package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"sourcegraph.com/sourcegraph/appdash"
)

func main() {
	f := "trace.json"
	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("failed to read file '%s', error: %v", f, err)
		return
	}
	var t appdash.Trace
	if err := json.Unmarshal(b, &t); err != nil {
		log.Printf("failed to unmarshal trace, error: %v", err)
		return
	}

	sWithoutParent := searchSpans(t, false)
	sWithParent := searchSpans(t, true)
	log.Printf("\ntotal without parent: %v\n\n", len(sWithoutParent))
	log.Printf("\ntotal with parent: %v\n\n", len(sWithParent))
	log.Printf("\nspans without parent: %+v\n\n", sWithoutParent)
	log.Printf("\nspans with parent: %+v\n\n", sWithParent)
}

func searchSpans(root appdash.Trace, withParent bool) []appdash.Span {
	var (
		spans []appdash.Span
		walk  func(root appdash.Trace)
	)
	walk = func(t appdash.Trace) {
		for _, sub := range t.Sub {
			p := findTraceParent(&root, sub)
			if p == nil && !withParent {
				spans = append(spans, sub.Span)
			}
			if p != nil && withParent {
				spans = append(spans, sub.Span)
			}
			walk(*sub)
		}
	}
	walk(root)
	return spans
}

func findTraceParent(root, child *appdash.Trace) *appdash.Trace {
	var walkToParent func(root, child *appdash.Trace) *appdash.Trace
	walkToParent = func(root, child *appdash.Trace) *appdash.Trace {
		if root.ID.Span == child.ID.Parent {
			return root
		}
		for _, sub := range root.Sub {
			if sub.ID.Span == child.ID.Parent {
				return sub
			}
			if r := walkToParent(sub, child); r != nil {
				return r
			}
		}
		return nil
	}
	return walkToParent(root, child)
}
