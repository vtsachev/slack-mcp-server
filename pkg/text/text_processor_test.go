package text

import "testing"

func TestProcessText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Simple stopword removal",
			input: "This is a test string with some stopwords",
			want:  "test string stopwords", // Lowercase
		},
		{
			name:  "String with 'Another' (capitalized) word - now expects lowercase",
			input: "Another example", // "Another" is a stopword
			want:  "example",         // Lowercase
		},
		{
			name:  "String with 'another' (lowercase) word - now expects lowercase",
			input: "another example", // "another" is a stopword
			want:  "example",         // Lowercase
		},
		{
			name:  "Empty string",
			input: "",
			want:  "",
		},
		{
			name:  "String with only stopwords",
			input: "this is an a THE AnD", // Mixed case stopwords
			want:  "",
		},
		{
			name:  "String with leading/trailing spaces and stopwords - now expects lowercase",
			input: "  Yet another test with stopwords  ", // "Yet", "another", "with" are stopwords
			want:  "test stopwords", // Lowercase
		},
		{
			name:  "String with mixed case input and mixed case stopwords - now expects all lowercase",
			input: "An Example With THE stopwords", // "An", "With", "THE" are stopwords
			want:  "example stopwords",         // Lowercase
		},
		{
			name:  "No stopwords, mixed case input - now expects all lowercase",
			input: "HelloWorld MyFriend",
			want:  "helloworld myfriend", // Lowercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessText(tt.input); got != tt.want {
				t.Errorf("ProcessText() = '%s', want '%s'", got, tt.want)
			}
		})
	}
}
