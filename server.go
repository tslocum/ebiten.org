// Copyright 2019 Hajime Hoshi
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

// +build ignore

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var (
	httpAddr = flag.String("http", ":8000", "HTTP address")
)

func init() {
	flag.Parse()
}

var rootPath = ""

func init() {
	_, path, _, _ := runtime.Caller(0)
	rootPath = filepath.Join(filepath.Dir(path), "docs")
}

type handler struct{}

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(rootPath, r.URL.Path[1:])
	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, filepath.Join(rootPath, "404.html"))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if f.IsDir() {
		path = filepath.Join(path, "index.html")
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				w.WriteHeader(http.StatusNotFound)
				http.ServeFile(w, r, filepath.Join(rootPath, "404.html"))
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.ServeFile(w, r, path)
}

func main() {
	http.Handle("/", handler{})
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
