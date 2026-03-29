package helper

// library pendukung lainnya

import (
	"github.com/chromedp/chromedp"
)

func takeScreenshot(url string, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		// Tunggu sampai elemen dashboard (grafik) muncul
		chromedp.WaitVisible(sel, chromedp.ByQuery),
		// Ambil screenshot dari elemen spesifik atau seluruh halaman
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
