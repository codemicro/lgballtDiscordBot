package requests

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
)

// Record returns an http.RoundTripper that writes out its
// requests and their responses to text files in basepath.
// Requests are named according to a hash of their contents.
// Responses are named according to the request that made them.
func Record(rt http.RoundTripper, basepath string) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("problem while recording transport: %w", err)
			}
		}()
		_ = os.MkdirAll(basepath, 0755)
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		reqname, resname := buildName(b)
		name := filepath.Join(basepath, reqname)
		if err = os.WriteFile(name, b, 0644); err != nil {
			return nil, err
		}
		if res, err = rt.RoundTrip(req); err != nil {
			return
		}
		b, err = httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}
		name = filepath.Join(basepath, resname)
		if err = os.WriteFile(name, b, 0644); err != nil {
			return nil, err
		}
		return
	})
}

// Replay returns an http.RoundTripper that reads its
// responses from text files in basepath.
// Responses are looked up according to a hash of the request.
func Replay(basepath string) http.RoundTripper {
	return ReplayFS(os.DirFS(basepath))
}

// ReplayFS returns an http.RoundTripper that reads its
// responses from text files in the fs.FS.
// Responses are looked up according to a hash of the request.
// Response file names may optionally be prefixed with comments for better human organization.
func ReplayFS(fsys fs.FS) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("problem while replaying transport: %w", err)
			}
		}()
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		_, name := buildName(b)
		glob := "*" + name
		matches, err := fs.Glob(fsys, glob)
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("no replay file matches %q", glob)
		}
		if len(matches) > 1 {
			return nil, fmt.Errorf("ambiguous response: multiple replay files match %q", glob)
		}
		b, err = fs.ReadFile(fsys, matches[0])
		if err != nil {
			return nil, err
		}
		r := bufio.NewReader(bytes.NewReader(b))
		return http.ReadResponse(r, req)
	})
}

func buildName(b []byte) (reqname, resname string) {
	h := md5.New()
	h.Write(b)
	s := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return s[:8] + ".req.txt", s[:8] + ".res.txt"
}
