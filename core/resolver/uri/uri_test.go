package uri

import (
	"errors"
	"reflect"
	"testing"
)

func TestURIConstructor(t *testing.T) {
	var tests = []struct {
		uri      string
		expected *URI
		err      error
	}{
		{"wrap://ipfs/QmHASH", &URI{
			Authority: "ipfs",
			Path:      "QmHASH",
			Uri:       "wrap://ipfs/QmHASH",
		}, nil},
		{"wrap://ens/domain.eth", &URI{
			Authority: "ens",
			Path:      "domain.eth",
			Uri:       "wrap://ens/domain.eth",
		}, nil},
		{"ens/domain.eth", &URI{
			Authority: "ens",
			Path:      "domain.eth",
			Uri:       "wrap://ens/domain.eth",
		}, nil},
		{"amazon/domain.eth", nil, errors.New("Invalid authority: amazon ")},
		{"www.google.com", nil, errors.New("Invalid URI Received: www.google.com ")},
		{"{**/}", nil, errors.New("Invalid URI Received: {**/} ")},
	}

	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
			result, err := New(tt.uri)
			if tt.err != nil {
				if !reflect.DeepEqual(err, tt.err) {
					t.Fatalf("Error expected: '%v', got '%v'", tt.err, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Fatalf("Unexpected URI: %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}
