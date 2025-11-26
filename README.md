## screenshot tool
a simple tool to capture screenshots of websites.
built with Go.

### prerequisites
- Go 1.25.4 or later
- chrome or chromium browser
- run "go get github.com/chromedp/chromedp" to install dependencies

## Usage

Run "go run cmd/main.go" to start the server on port 54321.

## Docker

Build: "docker build -t screenshot-tool ."

Run: "docker run -p 54321:54321 screenshot-tool"