# Web Cache Manager for Golang websites

## Webcache is a simple, lightweight Cache Manager for Golang websites.

Cache html, image, css, javascript, and other files.
Set the time duration globally and/or for each file.

### Run the test sample
Build the project and navigate to http:&#47;&#47;localhost:8005/mypage.htm.

#### Run the test sample outside of $GOPATH

- Start a shell window
- cd &lt;to any directory (other than $GOPATH)&gt;
- git clone https:&#47;&#47;github.com&#47;kambahr/go-webcache.git && cd go-webcache/test
- go mod init go-webcache/test
- go mod tidy
- go mod vendor
- go build -o cacheDemo && ./cacheDemo

Navigate to http:&#47;&#47;localhost:8005/mypage.html.
