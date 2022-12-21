package uri

import (
	"fmt"
	"regexp"
	"strings"
)

type URI struct {
	Authority string
	Path      string
	Uri       string
}

var (
	authorities = map[string]bool{"ipfs": true, "ens": true, "fs": true, "uns": true}
	re          = regexp.MustCompile("^(wrap://)*(?P<authority>[a-z][a-z0-9-_]+)/(?P<path>.*)$")
)

func New(uri string) (*URI, error) {
	exp := regexp.MustCompile("^/*")
	processed := exp.Split(uri, 1)[0]
	if !strings.HasPrefix(uri, "wrap://") {
		processed = fmt.Sprintf("wrap://%s", uri)
	}
	res := re.FindStringSubmatch(processed)
	if res == nil || len(res) != 4 || res[3] == "" {
		return nil, fmt.Errorf("Invalid URI Received: %s ", uri)
	}
	if _, ok := authorities[res[2]]; !ok {
		return nil, fmt.Errorf("Invalid authority: %s ", res[2])
	}
	return &URI{Authority: res[2], Path: res[3], Uri: processed}, nil
}

func (u *URI) PartialEq(uri *URI) bool {
	return u.Uri == uri.Uri
}

func (u *URI) String() string {
	return u.Uri
}
