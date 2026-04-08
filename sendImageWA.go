package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"wauto/config"
	"wauto/helper"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func mainSendImageWAGroup() {

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("❌ File .env tidak ditemukan!")
		return
	}

	// 1. Inisialisasi Database Sesi
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

	// 2. Cek Sesi Login (QR Code otomatis muncul jika belum login)
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
				fmt.Println("✅ Login Berhasil!")
			}
		}
	} else {
		// Jika sudah pernah login, langsung konek
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Jeda sebentar agar koneksi stabil
	time.Sleep(3 * time.Second)

	if client.IsConnected() {
		fmt.Println("✅ Terhubung ke WhatsApp")

		// --- CONFIG TARGET ---
		groupID := cfg.NOMORWAGROUP     // GANTI DENGAN JID GRUP ANDA
		filePath := "test_image_wa.png" // PASTIKAN FILE INI ADA DI FOLDER
		caption := "📊 *Laporan Dashboard Grafana Terkini*"

		// 3. Eksekusi Pengiriman Gambar
		helper.SendImageToGroup(client, groupID, filePath, caption)

		// 4. Jeda & Exit (Agar gambar terkirim sempurna sebelum aplikasi ditutup)
		fmt.Println("⏳ Menunggu proses upload & kirim selesai...")
		time.Sleep(7 * time.Second)

		client.Disconnect()
		fmt.Println("🚀 Selesai. Program ditutup.")
		os.Exit(0)
	} else {
		fmt.Println("❌ Gagal terhubung ke WhatsApp.")
	}
}

// Fungsi Inti untuk Mengirim Gambar ke Grup
func sendImageToGroup_(client *whatsmeow.Client, groupID string, filePath string, caption string) {
	// A. Baca file gambar
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ File %s tidak ditemukan: %v\n", filePath, err)
		return
	}

	// B. Upload file ke WhatsApp Server menggunakan MediaImage (yang benar)
	fmt.Println("☁️  Sedang mengunggah gambar ke server WhatsApp...")
	resp, err := client.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("❌ Gagal upload gambar: %v\n", err)
		return
	}

	// C. Susun JID Target
	recipient, err := types.ParseJID(groupID)
	if err != nil {
		fmt.Printf("❌ Format ID Grup Salah: %v\n", err)
		return
	}

	// D. Bangun struktur Pesan Gambar dengan Format Protobuf Terbaru
	msg := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Caption:       proto.String(caption),
			Mimetype:      proto.String("image/png"),
			URL:           &resp.URL, // HARUS BESAR SEMUA
			DirectPath:    &resp.DirectPath,
			MediaKey:      resp.MediaKey,
			FileEncSHA256: resp.FileEncSHA256, // HARUS BESAR SEMUA (SHA)
			FileSHA256:    resp.FileSHA256,    // HARUS BESAR SEMUA (SHA)
			FileLength:    proto.Uint64(uint64(len(data))),
		},
	}

	// E. Kirim Pesan
	_, err = client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		fmt.Printf("❌ Gagal mengirim gambar ke grup: %v\n", err)
	} else {
		fmt.Printf("✅ Berhasil! Gambar %s telah terkirim ke grup.\n", filePath)
	}
}
