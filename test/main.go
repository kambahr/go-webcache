package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kambahr/go-webcache"
)

var mWebcache *webcache.Cache
var mInstallPath string

// This is a simple website to show the usage of the webcache package.

func main() {

	// This is the gloal cache duration, for all files.
	// You can also set the cache duration for each file.
	var cacheDuration time.Duration = 10 * time.Minute
	var portNo int = 8005

	// This is the default root directory, where your main is
	// executed from; but it could be any path on your system -- that
	// you've located your web-files.
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	mInstallPath = dir

	//------------------------------------------------------------
	// Initialize the webcache object.
	// Pass the full root physical path, and the cache duration.
	mWebcache = webcache.NewWebCache(cacheDuration)
	//------------------------------------------------------------

	http.HandleFunc("/", handleMyPage)

	svr := http.Server{
		Addr:           fmt.Sprintf(":%d", portNo),
		MaxHeaderBytes: 20480,
	}

	fmt.Printf("Listening to port %d\n", portNo)

	log.Fatal(svr.ListenAndServe())
}

func handleMyPage(w http.ResponseWriter, r *http.Request) {

	var b []byte

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Caching the page avoids excessive I/O reads.
	secToCache := 15
	d := time.Duration(secToCache) * time.Second

	// Relative path to store on the cache list.
	const uri string = "/mypage.html"
	rPath := strings.ToLower(r.URL.Path)

	pageValid := (strings.HasSuffix(rPath, "/") || strings.HasSuffix(rPath, uri))

	if !pageValid {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h3>Error 404 - Not Found</h3>")
		return
	}

	// See if the cache exists on the list.
	exists := mWebcache.Exists(uri)
	if !exists {
		//---------------------------------
		log.Println("Reading from disk...")
		//---------------------------------

		// Read the file bytes.
		tInshtmlPath := fmt.Sprintf("%s%s", mInstallPath, uri)

		fExistOnDisk := FileExists(tInshtmlPath)

		if !fExistOnDisk {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "<h3>Error 404 - Not Found</h3>")
			return
		}

		b, _ = ioutil.ReadFile(tInshtmlPath)

		// Replace strings.
		t := strings.Split(time.Now().String(), ".")[0]
		b = bytes.Replace(b, []byte("{{.TimeNow}}"), []byte(t), -1)
		b = bytes.Replace(b, []byte("{{.DurTime}}"), []byte(fmt.Sprintf("%d", secToCache)), -1)

		// Do not wait for this.
		go mWebcache.AddItem(uri, b, d)

	} else {
		//---------------------------------
		log.Println("Reading from cache...")
		//---------------------------------
		// Get the file bytes from the cache.
		b = mWebcache.GetItem(uri)
	}

	// Write the response.
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func FileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
