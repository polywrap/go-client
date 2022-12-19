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

func New(uri string) (*URI, error) {
	re := regexp.MustCompile("^/*")
	processed := re.Split(uri, 1)[0]
	if strings.Index(uri, "wrap://") == -1 {
		processed = fmt.Sprintf("wrap://%s", uri)
	}
	re = regexp.MustCompile("^(wrap://)*(?P<authority>[a-z][a-z0-9-_]+)/(?P<path>.*)$")
	res := re.FindStringSubmatch(processed)
	if res == nil || len(res) != 4 || res[2] == "" || res[3] == "" {
		return nil, fmt.Errorf("Invalid URI Received: %s ", uri)
	}
	return &URI{Authority: res[2], Path: res[3], Uri: processed}, nil
}

func (u *URI) PartialEq(uri *URI) bool {
	return u.Uri == uri.Uri
}

func (u *URI) String() string {
	return u.Uri
}
