package common

import "testing"

func TestStripJSONComments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no comments",
			input: `{"key": "value"}`,
			want:  `{"key": "value"}`,
		},
		{
			name:  "line comment",
			input: "{\"key\": \"value\" // comment\n}",
			want:  "{\"key\": \"value\" \n}",
		},
		{
			name:  "block comment",
			input: `{"key": /* comment */ "value"}`,
			want:  `{"key":  "value"}`,
		},
		{
			name:  "slashes in string preserved",
			input: `{"url": "https://example.com"}`,
			want:  `{"url": "https://example.com"}`,
		},
		{
			name:  "line comment at start",
			input: "{// comment\n\"key\":\"val\"}",
			want:  "{\n\"key\":\"val\"}",
		},
		{
			name:  "comment-like content in string preserved",
			input: `{"key":"// not a comment"}`,
			want:  `{"key":"// not a comment"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(StripJSONComments([]byte(tt.input)))
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
