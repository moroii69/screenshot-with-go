package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
)

func takeScreenshot(url string, width, height int64, fullPage bool) ([]byte, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.WindowSize(int(width), int(height)),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 20*time.Second)
	defer cancelTimeout()

	var buf []byte
	tasks := []chromedp.Action{
		chromedp.Navigate(url),
		chromedp.Sleep(2 * time.Second),
	}
	if fullPage {
		tasks = append(tasks, chromedp.FullScreenshot(&buf, 90))
	} else {
		tasks = append(tasks, chromedp.CaptureScreenshot(&buf))
	}
	err := chromedp.Run(timeoutCtx, tasks...)
	return buf, err
}

func screenshotHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "missing ?url=", http.StatusBadRequest)
		return
	}

	width := int64(1920)
	height := int64(1080)
	if wParam := r.URL.Query().Get("width"); wParam != "" {
		if v, err := strconv.ParseInt(wParam, 10, 64); err == nil {
			width = v
		}
	}
	if hParam := r.URL.Query().Get("height"); hParam != "" {
		if v, err := strconv.ParseInt(hParam, 10, 64); err == nil {
			height = v
		}
	}
	fullPage := r.URL.Query().Get("fullpage") == "true"

	img, err := takeScreenshot(url, width, height, fullPage)
	if err != nil {
		http.Error(w, "screenshot failed :( : "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(img)
}

func main() {
	http.HandleFunc("/screenshot", screenshotHandler)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	fmt.Println("server running at http://localhost:54321")
	log.Fatal(http.ListenAndServe(":54321", nil))
}
