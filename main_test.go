package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestBnew(t *testing.T) {
	cases := []struct {
		name    string
		in1     string
		in2     string
		n       int
		want    string
		failure bool
	}{
		{
			name: "newlines",
			in1:  "foo\nbar\nbaz\n",
			in2:  "bar\nbaz\nqux\nquux\n",
			want: "qux\nquux\n",
		}, {
			name: "no newlines",
			in1:  "foo\nbar\nbaz",
			in2:  "bar\nbaz\nqux\nquux",
			want: "qux\nquux\n",
		}, {
			name:    "negative size hint",
			n:       -1,
			failure: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			w := new(bytes.Buffer)
			r1 := strings.NewReader(c.in1)
			r2 := strings.NewReader(c.in2)
			err := Bnew(w, r1, r2, c.n)
			if !c.failure && err != nil {
				t.Errorf("got %q; want nil", err)
			}
			if c.failure && err == nil {
				t.Error("got nil; want some error")
			}
			got := w.String()
			if got != c.want {
				t.Errorf("got %q; want %q", got, c.want)
			}
		})
	}
}
