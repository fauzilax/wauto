package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"wauto/config"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func mainChromeDP() {
	// --- KONFIGURASI ---

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("❌ File .env tidak ditemukan!")
		return
	}

	urls := []string{cfg.GRAFANAURL1, cfg.GRAFANAURL2}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("disable-gpu", true),
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	for idx, gURL := range urls {
		outputFile := fmt.Sprintf("grafana_final_%d.png", idx+1)
		fmt.Printf("\n🚀 [%d/%d] Menuju: %s\n", idx+1, len(urls), gURL)

		ctx, tabCancel := chromedp.NewContext(allocCtx)
		// Naikkan timeout ke 180 detik jika internet sedang tidak stabil
		ctx, timeCancel := context.WithTimeout(ctx, 180*time.Second)

		var buf []byte
		fmt.Println("🌐 Navigasi dan Auth...")

		err = chromedp.Run(ctx,
			network.Enable(),
			chromedp.ActionFunc(func(ctx context.Context) error {
				return network.SetExtraHTTPHeaders(network.Headers{
					"Authorization": "Bearer " + cfg.MYTOKEN,
				}).Do(ctx)
			}),
			chromedp.EmulateViewport(1920, 1080),
			chromedp.Navigate(gURL),

			// 1. TUNGGU elemen grafik spesifik muncul (Jangan cuma body)
			// Grafana biasanya pakai class '.dashboard-grid' atau '.panel-container'
			chromedp.ActionFunc(func(ctx context.Context) error {
				fmt.Println("⏳ Menunggu panel Grafana muncul...")
				return nil
			}),
			chromedp.WaitVisible(`.dashboard-grid`, chromedp.ByQuery),

			// 2. Jeda Render (Gunakan waktu yang cukup lama, misal 20 detik)
			// Ini untuk memastikan data benar-benar ditarik dari DB.
			chromedp.ActionFunc(func(ctx context.Context) error {
				fmt.Println("🎨 Panel terdeteksi. Sinkronisasi data (20 detik)...")
				return nil
			}),
			chromedp.Sleep(20*time.Second),

			// --- 🛠️ TRIK AMPUH ANTI-GELAP (WAJIB ADA) 🛠️ ---

			// A. Paksa browser fokus ke area konten
			chromedp.Click(`.dashboard-grid`, chromedp.ByQuery),

			// B. Jalankan JS untuk trigger "Repaint" (menggambar ulang layar)
			// Ini memaksa elemen transparan menjadi solid.
			chromedp.Evaluate(`
                window.scrollBy(0, 1); 
                window.scrollBy(0, -1);
                document.body.style.backgroundColor = '#161719'; // Paksa warna background Grafana Dark
            `, nil),

			// C. Beri jeda sangat singkat setelah repaint
			chromedp.Sleep(1*time.Second),

			// ------------------------------------------

			// 3. AMBIL SCREENSHOT
			chromedp.ActionFunc(func(ctx context.Context) error {
				fmt.Println("📸 Mengambil screenshot...")
				return nil
			}),
			chromedp.FullScreenshot(&buf, 100),
		)

		if err != nil {
			fmt.Printf("❌ Gagal di URL %d: %v\n", idx+1, err)
		} else {
			if err := os.WriteFile(outputFile, buf, 0644); err != nil {
				fmt.Println("❌ Gagal simpan file:", err)
			} else {
				fmt.Printf("✨ Berhasil disimpan: %s\n", outputFile)
			}
		}

		timeCancel()
		tabCancel()
	}

}
