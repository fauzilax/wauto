package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"wauto/config"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {
	// --- KONFIGURASI ---
	config, err := config.LoadConfig()
	if err != nil {

		fmt.Println("❌ File .env tidak ditemukan oleh Viper!")
		os.Exit(1)
	}
	// i := 1
	// totalURL := 2
	// for i <= totalURL {
	// 	fmt.Println("COUNT : ", i)
	// 	i++
	// }
	urls := []string{
		config.GRAFANAURL1,
		config.GRAFANAURL2,
	}

	for idx, grafanaURL := range urls {
		myToken := config.MYTOKEN
		outputFile := "grafana_final" + strconv.Itoa(idx+1) + ".png"

		// 1. Setup Allocator (Headless: false agar Anda bisa memantau prosesnya)
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("ignore-certificate-errors", true),
		)
		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		// Beri timeout total 2 menit
		ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
		defer cancel()

		var buf []byte

		fmt.Println("🚀 Membuka browser...")

		err = chromedp.Run(ctx,
			// 1. Masukkan Token (Pastikan sudah benar glsa_...)
			chromedp.ActionFunc(func(ctx context.Context) error {
				return network.SetExtraHTTPHeaders(network.Headers{
					"Authorization": "Bearer " + myToken,
				}).Do(ctx)
			}),

			// 2. Set Viewport & Navigate
			chromedp.EmulateViewport(1920, 1080),
			chromedp.Navigate(grafanaURL),

			// 3. TUNGGU 'BODY' SAJA (Ganti dari .dashboard-grid)
			// Semua website pasti punya tag <body>. Jika ini tetap timeout,
			// berarti IP tersebut memang tidak bisa diakses sama sekali.
			chromedp.ActionFunc(func(ctx context.Context) error {
				fmt.Println("⏳ Menunggu respons dari server (Body)...")
				return nil
			}),
			chromedp.WaitVisible(`body`, chromedp.ByQuery),

			// 4. BERI WAKTU TETAP (FIXED SLEEP)
			chromedp.ActionFunc(func(ctx context.Context) error {
				fmt.Println("🎨 Halaman terdeteksi. Menunggu 3 detik agar grafik selesai render...")
				return nil
			}),
			chromedp.Sleep(10*time.Second),

			// 5. AMBIL SCREENSHOT
			chromedp.ActionFunc(func(ctx context.Context) error {
				fmt.Println("⏳📸 Mengambil screenshot sekarang...")
				return nil
			}),
			chromedp.FullScreenshot(&buf, 100),
		)

		if err != nil {
			log.Fatalf("❌ Gagal: %v", err)
		}

		// Simpan file
		if err := os.WriteFile(outputFile, buf, 0644); err != nil {
			log.Fatalf("❌ Gagal simpan file: %v", err)
		}

		fmt.Printf("✨ Selesai! Cek file: %s\n", outputFile)

	}

}
