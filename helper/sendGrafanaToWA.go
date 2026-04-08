package helper

import (
	"context"
	"fmt"
	"os"
	"time"
	"wauto/config"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// --- KONFIGURASI UTAMA ---

const (
	// Ganti IP/Domain sesuai server Linux Anda, pastikan ada /render/

	ImageFile   = "grafana_report_today.png"
	CaptionText = "📊 *Laporan Dashboard Grafana Terkini*\n Automatic send from goBot system."
)

func SendGrafToWA() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("❌ File .env tidak ditemukan!")
		return
	}
	urls := []string{
		cfg.GRAFANAURL1,
		cfg.GRAFANAURL2,
		cfg.GRAFANAURL3,
	}
	GrafanaToken := cfg.MYTOKEN
	WhatsAppGroup := cfg.NOMORWAGROUP // ID Grup target

	for _, url := range urls {
		// Inisialisasi Grafana URL
		GrafanaURL := url

		fmt.Println("🤖 Bot Wauto Dimulai...")

		// 1. DOWNLOAD GAMBAR DARI SERVER LINUX
		err = DownloadGrafanaImage(GrafanaURL, GrafanaToken, ImageFile)
		if err != nil {
			fmt.Printf("❌ Gagal mengambil gambar dari server: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Gambar berhasil diunduh dari server Linux!")

		// 2. ANALISIS GAMBAR DENGAN GEMINI AI (Fitur Baru!)
		ctx := context.Background()

		fmt.Println("🧠 Sedang meminta Gemini AI menganalisis gambar...")
		aiAnalysis, err := analyzeImageWithGemini(ctx, cfg.GEMINIAPIKEY, ImageFile)
		if err != nil {
			fmt.Printf("⚠️ Gagal analisis AI (tetap lanjut kirim gambar): %v\n", err)
			aiAnalysis = "*(Gagal mendapatkan analisis AI)*"
		} else {
			fmt.Println("✅ Gemini AI berhasil menganalisis gambar.")
		}

		// 3. PERSIAPAN WHATSAPP BOT
		dbLog := waLog.Stdout("Database", "ERROR", true)
		container, err := sqlstore.New(context.Background(), "sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
		if err != nil {
			panic(err)
		}

		deviceStore, err := container.GetFirstDevice(context.Background())
		if err != nil {
			panic(err)
		}

		clientLog := waLog.Stdout("Client", "ERROR", true)
		client := whatsmeow.NewClient(deviceStore, clientLog)

		// 4. LOGIN WHATSAPP (QR Code jika belum sesi)
		if client.Store.ID == nil {
			qrChan, _ := client.GetQRChannel(context.Background())
			err = client.Connect()
			if err != nil {
				panic(err)
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					fmt.Println("\n👉 SCAN QR DI BAWAH INI DENGAN WHATSAPP:")
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				} else if evt.Event == "success" {
					fmt.Println("✅ Login WhatsApp Berhasil!")
				}
			}
		} else {
			err = client.Connect()
			if err != nil {
				panic(err)
			}
		}

		time.Sleep(3 * time.Second)

		// 5. KIRIM GAMBAR KE GRUP
		if client.IsConnected() {
			formattedCaption := fmt.Sprintf("📊 *Laporan Dashboard Grafana*\n\n🤖 *Analisis AI Gemini:* \n%s", aiAnalysis)

			fmt.Println("✅ Terhubung ke WhatsApp server.")
			SendImageToGroup(client, WhatsAppGroup, ImageFile, formattedCaption)

			fmt.Println("⏳ Menunggu proses pengiriman selesai...")
			time.Sleep(5 * time.Second)
			client.Disconnect()
			fmt.Println("🚀 Selesai. Program ditutup.")
		} else {
			fmt.Println("❌ Gagal terhubung ke WhatsApp.")
		}

		time.Sleep(5 * time.Second)
	}
}
