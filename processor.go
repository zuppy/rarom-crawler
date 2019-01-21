package main

import (
	"bytes"
)

// i know i could use regexps, but it's faster this way for my case (even if it's hackish)
func findTextBetweenMarkers(text []byte, markerStart []byte, markerEnd []byte) (match []byte, newCursorPosition int) {
	startIndex := bytes.Index(text, markerStart)
	if startIndex < 0 {
		return nil, 0
	}
	endIndex   := bytes.Index(text[startIndex+len(markerStart):], markerEnd)
	if endIndex < 0 {
		return nil, 0
	}

	sliceFrom := startIndex+len(markerStart)
	sliceTo   := startIndex+len(markerStart)+endIndex
	sliceNext := startIndex+endIndex+len(markerEnd)+len(markerStart)

	return text[sliceFrom:sliceTo], sliceNext
}
