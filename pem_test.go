package filters

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestToPem(t *testing.T) {

	type args struct {
		r   io.Reader
		blk Block
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestOne",
			args: args{r: strings.NewReader("This is only a test"), blk: Block{Type: "Test One", Headers: map[string]string{"COUNT": "100"}}},
			want: "-----BEGIN Test One-----\nCOUNT: 100\n\nVGhpcyBpcyBvbmx5IGEgdGVzdA==\n-----END Test One-----\n",
		},
		{
			name: "TestTwo",
			args: args{
				r: strings.NewReader("The quick brown fox jumped over the lazy dog.  The quick brown fox jumped over the lazy dog."),
				blk: Block{
					Type:    "Test Two",
					Headers: map[string]string{"COUNT": "200"},
				},
			},
			want: "-----BEGIN Test Two-----\n" +
				"COUNT: 200\n" +
				"\n" +
				"VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wZWQgb3ZlciB0aGUgbGF6eSBkb2cuICBU\n" +
				"aGUgcXVpY2sgYnJvd24gZm94IGp1bXBlZCBvdmVyIHRoZSBsYXp5IGRvZy4=\n" +
				"-----END Test Two-----\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := io.ReadAll(ToPem(tt.args.r, tt.args.blk)); strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("ToPem() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestFromPem(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 Block
	}{
		{
			name:  "Test One",
			args:  args{r: strings.NewReader("-----BEGIN Test One-----\nCOUNT: 100\n\nVGhpcyBpcyBvbmx5IGEgdGVzdA==\n-----END Test One-----\n")},
			want:  "This is only a test",
			want1: Block{Type: "Test One", Headers: map[string]string{"COUNT": "100"}},
		},
		{
			name: "Test Two",
			args: args{r: strings.NewReader("-----BEGIN Test Two-----\n" +
				"COUNT: 200\n" +
				"\n" +
				"VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wZWQgb3ZlciB0aGUgbGF6eSBkb2cuICBU\n" +
				"aGUgcXVpY2sgYnJvd24gZm94IGp1bXBlZCBvdmVyIHRoZSBsYXp5IGRvZy4=\n" +
				"-----END Test Two-----\n")},
			want:  "The quick brown fox jumped over the lazy dog.  The quick brown fox jumped over the lazy dog.",
			want1: Block{Type: "Test Two", Headers: map[string]string{"COUNT": "200"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdr, got1 := FromPem(tt.args.r)
			got, _ := io.ReadAll(rdr)
			if strings.Compare(string(got), tt.want) != 0 {
				t.Errorf("FromPem() got = %v, want %v", string(got), tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("FromPem() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
