// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//go:embed assets
var assets embed.FS

// create middleware that serves embedded assets.
func embeddedAssetsHandler(h http.Handler) http.HandlerFunc {
	fakeModTime := time.Now().UTC()

	// etagCache is a map of path to SHA256 hash
	etagCache := map[string]string{}

	// servedFiles is a map of path to file info
	servedFiles := map[string]bool{}

	return func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("%s %s: middleware: assets\n", r.Method, r.URL.Path)
		if r.Method != http.MethodGet {
			h.ServeHTTP(w, r)
			return
		}

		// we have to add the prefix to the embedded filesystem
		path := filepath.Join("assets", r.URL.Path)
		//log.Printf("%s %s: %q\n", r.Method, r.URL.Path, path)
		if sb, err := fs.Stat(assets, path); err != nil {
			//log.Printf("%s %s: %s: %v\n", r.Method, r.URL.Path, path, err)
			if os.IsNotExist(err) {
				//log.Printf("%s %s: %s: does not exist\n", r.Method, r.URL.Path, path)
				h.ServeHTTP(w, r)
				return
			}
			log.Printf("%s %s: %s: %v\n", r.Method, r.URL.Path, path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else if sb.IsDir() {
			if r.URL.Path == "/" {
				log.Printf("%s %s: %s: is the root\n", r.Method, r.URL.Path, path)
				h.ServeHTTP(w, r)
				return
			}
			log.Printf("%s %s: %s: is a directory\n", r.Method, r.URL.Path, path)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		//log.Printf("%s %s: %v\n", r.Method, r.URL.Path, sb.ModTime())

		// let http handle the file
		fp, err := assets.Open(path)
		if err != nil {
			log.Printf("%s %s: %v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer fp.Close()

		// clients can use two headers to avoid downloading.
		// The first, If-Modified-Since, is a timestamp and is set by http.ServeContent
		// from the file's ModTime. Since we're using an embedded file system and Go has
		// made the decision to ignore ModTime, we have to use a fake ModTime.
		// The second, If-None-Match, is an ETag and is set by us. We compute the ETag
		// from a SHA256 of the file's content.
		//
		// If client sends an ETag header, compute the hash to see if it's still the same.
		if match := r.Header.Get("If-None-Match"); match != "" {
			etag, ok := etagCache[path]
			if !ok {
				// new file, so compute ETag based on file content and cache it
				hasher := sha256.New()
				if _, err := io.Copy(hasher, fp); err != nil {
					log.Printf("%s %s: %v\n", r.Method, r.URL.Path, err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				etag = fmt.Sprintf(`"%x"`, hasher.Sum(nil))
				etagCache[path] = etag
			}
			if match == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		f, ok := fp.(io.ReadSeeker)
		if !ok {
			log.Printf("%s %s: %s: %v\n", r.Method, r.URL.Path, path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !servedFiles[path] {
			log.Printf("%s %s: %s: served\n", r.Method, r.URL.Path, path)
			servedFiles[path] = true
		}

		http.ServeContent(w, r, r.URL.Path, fakeModTime, f)
	}
}

// isFileExists returns true if the path exists and is a regular file.
func isFileExists(path string) (bool, error) {
	sb, err := fs.Stat(assets, path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	} else if sb.IsDir() {
		return false, nil
	}
	return sb.Mode().IsRegular(), nil
}
