package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_read(t *testing.T) {
	r := require.New(t)

	got, err := read(strings.NewReader(example))
	r.NoError(err)

	r.Equal(10, len(got))
}

func Test_isOpen(t *testing.T) {
	tt := []struct {
		in   byte
		want bool
	}{
		{'(', true},
		{'[', true},
		{'{', true},
		{'<', true},

		{')', false},
		{']', false},
		{'}', false},
		{'>', false},

		{'.', false},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(string([]byte{tc.in}), func(t *testing.T) {
			t.Parallel()
			got := isOpen(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_isClose(t *testing.T) {
	tt := []struct {
		in   byte
		want bool
	}{
		{'(', false},
		{'[', false},
		{'{', false},
		{'<', false},

		{')', true},
		{']', true},
		{'}', true},
		{'>', true},

		{'.', false},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(string([]byte{tc.in}), func(t *testing.T) {
			t.Parallel()
			got := isClose(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_parse(t *testing.T) {
	tt := []struct {
		name       string
		in         string
		wantErr    error
		wantChunks []chunk
	}{
		{
			name:    "corrupt example 1",
			in:      "(]",
			wantErr: errCorrupted{pos: 2, want: ')', got: ']'},
		},
		{
			name:    "corrupt example 2",
			in:      "{()()()>",
			wantErr: errCorrupted{pos: 8, want: '}', got: '>'},
		},
		{
			name:    "corrupt example 3",
			in:      "(((()))}",
			wantErr: errCorrupted{pos: 8, want: ')', got: '}'},
		},
		{
			name:    "corrupt example 4",
			in:      "<([]){()}[{}])",
			wantErr: errCorrupted{pos: 14, want: '>', got: ')'},
		},
		{
			name:       "valid example 1",
			in:         "([])",
			wantChunks: []chunk{"([])"},
		},
		{
			name:       "valid example 2",
			in:         "{()()()}",
			wantChunks: []chunk{"{()()()}"},
		},
		{
			name:       "valid example 3",
			in:         "<([{}])>",
			wantChunks: []chunk{"<([{}])>"},
		},
		{
			name:       "valid example 4",
			in:         "[<>({}){}[([])<>]]",
			wantChunks: []chunk{"[<>({}){}[([])<>]]"},
		},
		{
			name:       "valid example 5",
			in:         "(((((((((())))))))))",
			wantChunks: []chunk{"(((((((((())))))))))"},
		},
		{
			name:       "empty string is a valid set of brackets",
			in:         "",
			wantErr:    nil,
			wantChunks: nil,
		},
		{
			name:       "consecutive single pairs",
			in:         "(){}[]<>",
			wantChunks: []chunk{"()", "{}", "[]", "<>"},
		},
		{
			name:       "out of order pairs",
			in:         "(){[}]<>",
			wantErr:    errCorrupted{pos: 5, want: ']', got: '}'},
			wantChunks: []chunk{"()"},
		},
		{
			name:       "incomplete chunk",
			in:         "(<>){[]}(<>",
			wantErr:    errIncomplete{pos: 12},
			wantChunks: []chunk{"(<>)", "{[]}"},
		},
		{
			name:    "starting with a closing bracket is invalid",
			in:      ">",
			wantErr: errCorrupted{pos: 1, got: '>'},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := line{raw: []byte(tc.in)}
			err := l.parse()
			require.Equal(t, tc.wantErr, err)

			if len(tc.wantChunks) == 0 {
				assert.Empty(t, l.chunks)
			} else {
				assert.Equal(t, tc.wantChunks, l.chunks)
			}
		})
	}
}
