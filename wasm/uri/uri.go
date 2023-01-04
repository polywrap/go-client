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
	authorities    = map[string]bool{"ipfs": true, "ens": true, "fs": true, "uns": true}
	processedExp   = regexp.MustCompile("^wrap://(?P<authority>[a-z][a-z0-9-_]+)/(?P<path>.*)$")
	unprocessedExp = regexp.MustCompile("^(wrap://)?(?P<authority>[a-z][a-z0-9-_]+)/(?P<path>.*)$")
)

func New(uri string) (*URI, error) {
	processed := strings.TrimLeft(uri, "/")
	if !strings.HasPrefix(uri, "wrap://") {
		processed = fmt.Sprintf("wrap://%s", uri)
	}
	res := processedExp.FindStringSubmatch(processed)
	if res == nil || len(res) != 3 || res[2] == "" {
		return nil, fmt.Errorf("Invalid URI Received: %s ", uri)
	}
	if _, ok := authorities[res[1]]; !ok {
		return nil, fmt.Errorf("Invalid authority: %s ", res[1])
	}
	return &URI{Authority: res[1], Path: res[2], Uri: processed}, nil
}

func IsValid(uri string) (bool, error) {
	res := unprocessedExp.FindStringSubmatch(uri)
	if res == nil || len(res) != 4 || res[3] == "" {
		return false, fmt.Errorf("Invalid URI Received: %s ", uri)
	}
	if _, ok := authorities[res[2]]; !ok {
		return false, fmt.Errorf("Invalid authority: %s ", res[2])
	}
	return true, nil
}

func (u *URI) PartialEq(uri *URI) bool {
	return u.Uri == uri.Uri
}

func (u *URI) String() string {
	return u.Uri
}
