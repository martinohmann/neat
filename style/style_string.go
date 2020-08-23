package style

import (
	"strings"
	"sync"
)

const (
	// Delimiters for style attribute blocks in strings.
	attrBlockStart = '{'
	attrBlockEnd   = '}'
)

var (
	// sequenceCache is a map of raw style attribute names such as
	// "yellow,bold" to the resulting escape sequencesafe for concurrent use.
	// This is used to reduce the amount of heavy lifting during style
	// replacements in strings.
	sequenceCache sync.Map
)

// StyleString replaces all supported style attributes of the from "{attr1,attr2}"
// in s with the corresponding ANSI escape sequences. If an attribute is not
// recognized it is not replaced. If styles are disabled, style attributes are
// replaced with empty strings. This supports all attribute names defined in
// AttributeMap.
func StyleString(s string) string {
	return styleString(s)
}

// styleString replaces style attribute blocks of the form
// `{attr1[,attr2]*}` with style escape sequences. This avoids the usage of
// regular expressions for performance reasons.
func styleString(s string) string {
	// // Fast path
	if len(s) <= 2 || !strings.ContainsRune(s, attrBlockStart) {
		return s
	}

	var sb strings.Builder
	var inBlock bool

	sb.Grow(len(s))

	for {
		if !inBlock {
			// Search for the next attribute block start.
			start := strings.IndexRune(s, attrBlockStart)
			if start == -1 {
				// If there are no more attribute blocks, write out the rest of
				// the string and break.
				sb.WriteString(s)
				break
			}

			// Consume s until the block start.
			sb.WriteString(s[:start])
			s = s[start+1:]
			inBlock = true
			continue
		}

		// We are inside an attribute block and try to find its end.
		end := strings.IndexRune(s, attrBlockEnd)
		if end == -1 {
			// If there is no matching }, this is no attribute block, write
			// out the rest of the string and break.
			sb.WriteString(s)
			break
		}

		nextStart := strings.IndexRune(s, attrBlockStart)
		if nextStart != -1 && nextStart < end {
			// If there is another { before the next }. Skip forward to it
			// before consuming the attribute block content.
			sb.WriteRune(attrBlockStart)
			sb.WriteString(s[:nextStart])
			s = s[nextStart+1:]
			continue
		}

		// Consume attribute block from s.
		rawBlock := s[:end]
		s = s[end+1:]
		inBlock = false

		// Fast path.
		if len(rawBlock) == 0 {
			sb.WriteRune(attrBlockStart)
			sb.WriteRune(attrBlockEnd)
			continue
		}

		// Write out the ANSI escape sequence if it could be built or just
		// append the original attribute block unaltered.
		sequence, ok := resolveEscapeSequence(rawBlock)
		if ok {
			sb.WriteString(sequence)
			continue
		}

		sb.WriteRune(attrBlockStart)
		sb.WriteString(rawBlock)
		sb.WriteRune(attrBlockEnd)
	}

	return sb.String()
}

func resolveEscapeSequence(raw string) (string, bool) {
	raw = strings.ToLower(raw)

	val, ok := sequenceCache.Load(raw)
	if ok {
		return val.(string), true
	}

	attrNames := strings.Split(raw, ",")

	attrs := make([]Attribute, 0, len(attrNames))

	for _, name := range attrNames {
		attr, ok := AttributeMap[name]
		if !ok {
			return raw, false
		}

		attrs = append(attrs, attr)
	}

	sequence := EscapeString(&Style{attrs})

	sequenceCache.Store(raw, sequence)

	return sequence, true
}
