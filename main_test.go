package main

import "testing"

func TestParseIssue(t *testing.T) {
	type args struct {
		description string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple-test",
			args: args{
				description: "foobar\nm\n\n```predicate\necho whatever\n```",
			},
			want: "echo whatever\n",
		},
		{
			name: "incorrect-codeblock",
			args: args{
				description: "foobar\nm\n\n```bash\necho whatever\n```",
			},
			want: "",
		},
		{
			name: "multiple-codeblocks",
			args: args{
				description: "foobar\nm\n\n```predicate\necho whatever\n```\n```predicate\necho blah\n```",
			},
			want: "echo whatever\n",
		},
		{
			name: "carriage-return",
			args: args{
				description: "foobar\r\nm\r\n\r\n```predicate\r\nexit 1\r\n```",
			},
			want: "exit 1\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseIssue(tt.args.description); got != tt.want {
				t.Errorf("ParseIssue() = %v, want %v", got, tt.want)
			}
		})
	}
}
