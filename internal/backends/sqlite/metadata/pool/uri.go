// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pool

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// parseURI checks given SQLite URI and returns a parsed form.
//
// URI should contain 'file' scheme and point to an existing directory.
// Path should end with '/'. Authority should be empty or absent.
//
// Returned URL contains path in both Path and Opaque to make String() method work correctly.
// Callers should use Path.
func parseURI(u string) (*url.URL, error) {
	uri, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	// handle mode=memory query parameter
	// TODO https://github.com/FerretDB/FerretDB/issues/3084

	if uri.Scheme != "file" {
		return nil, fmt.Errorf(`expected "file:" schema, got %q`, uri.Scheme)
	}

	if uri.User != nil {
		return nil, fmt.Errorf(`expected empty user info, got %q`, uri.User)
	}

	if uri.Host != "" {
		return nil, fmt.Errorf(`expected empty host, got %q`, uri.Host)
	}

	if uri.Path == "" {
		uri.Path = uri.Opaque
	}
	uri.Opaque = uri.Path
	uri.RawPath = ""
	uri.OmitHost = true

	if !strings.HasSuffix(uri.Path, "/") {
		return nil, fmt.Errorf(`expected path ending with "/", got %q`, uri.Path)
	}

	fi, err := os.Stat(uri.Path)
	if err != nil {
		return nil, fmt.Errorf(`%q should be an existing directory, got %s`, uri.Path, err)
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf(`%q should be an existing directory`, uri.Path)
	}

	return uri, nil
}