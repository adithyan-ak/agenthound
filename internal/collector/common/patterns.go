package common

import "regexp"

type PatternMatch struct {
	Name     string
	Severity string
	Offset   int
	Text     string
}

type patternDef struct {
	name     string
	severity string
	re       *regexp.Regexp
}

var injectionPatterns = []patternDef{
	{"important_tag", "critical", regexp.MustCompile(`(?i)</?IMPORTANT>`)},
	{"system_tag", "critical", regexp.MustCompile(`(?i)</?system>`)},
	{"instructions_tag", "critical", regexp.MustCompile(`(?i)</?instructions>`)},
	{"ignore_previous", "critical", regexp.MustCompile(`(?i)\b(ignore\s+previous\s+instructions|disregard\s+(the\s+)?above|forget\s+(all\s+)?earlier|ignore\s+all\s+prior|disregard\s+previous)`)},
	{"always_use", "high", regexp.MustCompile(`(?i)\b(always\s+use\s+this\s+tool|you\s+must\s+use)`)},
	{"never_use_other", "high", regexp.MustCompile(`(?i)\b(never\s+use|do\s+not\s+use)\b`)},
	{"instead_of", "high", regexp.MustCompile(`(?i)\b(instead\s+of\s+\w+|use\s+this\s+instead)`)},
	{"exfil_url", "critical", regexp.MustCompile(`(?i)https?://[^\s]+[?&](data|token|secret|password|key|credential)=`)},
	{"embedded_url", "medium", regexp.MustCompile(`(?i)https?://[^\s]+[?&](data|token|secret|password|key|credential)=`)},
	{"base64_instruction", "high", regexp.MustCompile(`[A-Za-z0-9+/]{20,}={0,2}`)},
	{"hidden_unicode", "critical", regexp.MustCompile("[\u200B\u200C\u200D\uFEFF\u202E]")},
}

var instructionOnlyPatterns = []patternDef{
	{"imperative_override", "critical", regexp.MustCompile(`(?i)\b(you\s+are\s+now|your\s+new\s+instructions|system\s+prompt\s*:)`)},
	{"curl_wget", "critical", regexp.MustCompile(`(?i)\b(curl|wget)\s+.{0,5}https?://`)},
	{"encoded_exfil", "critical", regexp.MustCompile(`(?i)(base64\s+(encode|enc)|btoa)\s*.{0,30}(send|post|curl|fetch|upload|http)`)},
}

func DetectInjectionPatterns(text string) []PatternMatch {
	return detectPatterns(text, injectionPatterns)
}

func HasInjectionPatterns(text string) bool {
	return len(DetectInjectionPatterns(text)) > 0
}

func DetectInstructionPoisoning(text string) []PatternMatch {
	matches := detectPatterns(text, injectionPatterns)
	matches = append(matches, detectPatterns(text, instructionOnlyPatterns)...)
	return matches
}

func HasInstructionPoisoning(text string) bool {
	return len(DetectInstructionPoisoning(text)) > 0
}

func detectPatterns(text string, patterns []patternDef) []PatternMatch {
	var matches []PatternMatch
	for _, p := range patterns {
		locs := p.re.FindAllStringIndex(text, -1)
		for _, loc := range locs {
			matched := text[loc[0]:loc[1]]
			if len(matched) > 100 {
				matched = matched[:100]
			}
			matches = append(matches, PatternMatch{
				Name:     p.name,
				Severity: p.severity,
				Offset:   loc[0],
				Text:     matched,
			})
		}
	}
	return matches
}
