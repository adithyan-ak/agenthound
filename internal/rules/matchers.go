package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/adithyan-ak/agenthound/internal/collector/common"
)

type CompiledMatcher interface {
	Match(text string) []MatchResult
}

type MatchResult struct {
	Matched bool
	Offset  int
	Text    string
}

func compileMatcher(spec MatcherSpec) (CompiledMatcher, error) {
	switch spec.Type {
	case "regex":
		return compileRegex(spec)
	case "keyword":
		return compileKeyword(spec)
	case "compound":
		return compileCompound(spec)
	case "entropy":
		return compileEntropy(spec)
	case "prefix":
		return compilePrefix(spec)
	default:
		return nil, fmt.Errorf("unknown matcher type %q", spec.Type)
	}
}

type regexMatcher struct {
	re *regexp.Regexp
}

func compileRegex(spec MatcherSpec) (*regexMatcher, error) {
	pattern := spec.Pattern
	if spec.CaseInsensitive && !strings.HasPrefix(pattern, "(?i)") {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex: %w", err)
	}
	return &regexMatcher{re: re}, nil
}

func (m *regexMatcher) Match(text string) []MatchResult {
	locs := m.re.FindAllStringIndex(text, -1)
	if len(locs) == 0 {
		return nil
	}
	results := make([]MatchResult, len(locs))
	for i, loc := range locs {
		matched := text[loc[0]:loc[1]]
		if len(matched) > 100 {
			matched = matched[:100]
		}
		results[i] = MatchResult{Matched: true, Offset: loc[0], Text: matched}
	}
	return results
}

type keywordMatcher struct {
	keywords        []string
	caseInsensitive bool
	matchAll        bool
}

func compileKeyword(spec MatcherSpec) (*keywordMatcher, error) {
	kw := make([]string, len(spec.Keywords))
	for i, k := range spec.Keywords {
		if spec.CaseInsensitive {
			kw[i] = strings.ToLower(k)
		} else {
			kw[i] = k
		}
	}
	return &keywordMatcher{
		keywords:        kw,
		caseInsensitive: spec.CaseInsensitive,
		matchAll:        spec.MatchMode == "all",
	}, nil
}

func (m *keywordMatcher) Match(text string) []MatchResult {
	searchText := text
	if m.caseInsensitive {
		searchText = strings.ToLower(text)
	}
	var results []MatchResult
	for _, kw := range m.keywords {
		idx := strings.Index(searchText, kw)
		if idx >= 0 {
			matched := text[idx : idx+len(kw)]
			if len(matched) > 100 {
				matched = matched[:100]
			}
			results = append(results, MatchResult{Matched: true, Offset: idx, Text: matched})
			if !m.matchAll {
				return results
			}
		} else if m.matchAll {
			return nil
		}
	}
	return results
}

type compoundMatcher struct {
	children []CompiledMatcher
	andMode  bool
}

func compileCompound(spec MatcherSpec) (*compoundMatcher, error) {
	children := make([]CompiledMatcher, len(spec.Matchers))
	for i, sub := range spec.Matchers {
		cm, err := compileMatcher(sub)
		if err != nil {
			return nil, fmt.Errorf("compound child %d: %w", i, err)
		}
		children[i] = cm
	}
	return &compoundMatcher{
		children: children,
		andMode:  spec.Operator == "and",
	}, nil
}

func (m *compoundMatcher) Match(text string) []MatchResult {
	var allResults []MatchResult
	for _, child := range m.children {
		results := child.Match(text)
		if len(results) > 0 {
			allResults = append(allResults, results...)
			if !m.andMode {
				return allResults
			}
		} else if m.andMode {
			return nil
		}
	}
	return allResults
}

type entropyMatcher struct {
	charset   string
	threshold float64
	minLength int
}

func compileEntropy(spec MatcherSpec) (*entropyMatcher, error) {
	return &entropyMatcher{
		charset:   spec.Charset,
		threshold: spec.Threshold,
		minLength: spec.MinLength,
	}, nil
}

func (m *entropyMatcher) Match(text string) []MatchResult {
	if len(text) < m.minLength {
		return nil
	}
	switch m.charset {
	case "base64":
		if !common.IsBase64Charset(text) {
			return nil
		}
	case "hex":
		if !common.IsHexCharset(text) {
			return nil
		}
	default:
		return nil
	}
	entropy := common.ShannonEntropy(text)
	if entropy <= m.threshold {
		return nil
	}
	matched := text
	if len(matched) > 100 {
		matched = matched[:100]
	}
	return []MatchResult{{Matched: true, Offset: 0, Text: matched}}
}

type prefixMatcher struct {
	prefixes        []string
	caseInsensitive bool
}

func compilePrefix(spec MatcherSpec) (*prefixMatcher, error) {
	prefixes := make([]string, len(spec.Prefixes))
	for i, p := range spec.Prefixes {
		if spec.CaseInsensitive {
			prefixes[i] = strings.ToLower(p)
		} else {
			prefixes[i] = p
		}
	}
	return &prefixMatcher{
		prefixes:        prefixes,
		caseInsensitive: spec.CaseInsensitive,
	}, nil
}

func (m *prefixMatcher) Match(text string) []MatchResult {
	checkText := text
	if m.caseInsensitive {
		checkText = strings.ToLower(text)
	}
	for _, p := range m.prefixes {
		if strings.HasPrefix(checkText, p) {
			matched := text[:len(p)]
			return []MatchResult{{Matched: true, Offset: 0, Text: matched}}
		}
	}
	return nil
}
