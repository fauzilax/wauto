package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"wauto/config"

	_ "github.com/mattn/go-sqlite3" // Pastikan ini ada
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func sendToGroup(client *whatsmeow.Client, groupID string, message string) {
	// Format JID untuk grup adalah @g.us
	recipient, err := types.ParseJID(groupID)
	if err != nil {
		fmt.Printf("❌ Format ID Grup Salah: %v\n", err)
		return
	}

	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	_, err = client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		fmt.Printf("❌ Gagal kirim ke grup: %v\n", err)
	} else {
		fmt.Println("✅ Pesan berhasil terkirim ke Grup!")
	}
}

// Fungsi Helper untuk kirim teks
func sendText(client *whatsmeow.Client, target string, message string) {
	recipient, _ := types.ParseJID(target + "@s.whatsapp.net")
	msg := &waProto.Message{
		Conversation: proto.String(message),
	}
	_, err := client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		fmt.Println("❌ Gagal kirim pesan:", err)
	} else {
		fmt.Println("✅ Pesan terkirim ke", target)
	}
}

func mainSendWA() {

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("❌ File .env tidak ditemukan!")
		return
	}

	// 1. Inisialisasi Database Sesi (Agar tidak scan QR terus)
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

	// 2. Koneksi ke WhatsApp
	err = client.Connect()
	if err != nil {
		panic(err)
	}

	// Tunggu sebentar sampai koneksi benar-benar stabil
	time.Sleep(2 * time.Second)

	if client.IsConnected() {
		fmt.Println("✅ Terhubung ke WhatsApp")

		// 3. IDENTITAS GRUP (Ganti dengan ID grup Anda)
		groupID := cfg.NOMORWAGROUP
		pesan := "Halo semuanya! Saya Bot dari Go salam kenal semua."

		// 4. Kirim Pesan
		sendToGroup(client, groupID, pesan)

		// 5. JEDA SEBELUM EXIT (Sangat Penting)
		fmt.Println("⏳ Menunggu paket data terkirim...")
		time.Sleep(5 * time.Second)

		// 6. Putuskan Koneksi & Keluar
		client.Disconnect()
		fmt.Println("🚀 Selesai. Program ditutup.")
		os.Exit(0)
	} else {
		fmt.Println("❌ Gagal terhubung. Pastikan sudah scan QR sebelumnya.")
	}

	// // 1. Setup Database untuk simpan Sesi
	// dbLog := waLog.Stdout("Database", "DEBUG", true)
	// container, err := sqlstore.New(context.Background(), "sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	// if err != nil {
	// 	panic(err)
	// }

	// // 2. Ambil device pertama (atau buat baru jika belum ada)
	// deviceStore, err := container.GetFirstDevice(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	// clientLog := waLog.Stdout("Client", "DEBUG", true)
	// client := whatsmeow.NewClient(deviceStore, clientLog)

	// // 3. Cek Login / Scan QR
	// if client.Store.ID == nil {
	// 	qrChan, _ := client.GetQRChannel(context.Background())
	// 	err = client.Connect()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for evt := range qrChan {
	// 		if evt.Event == "code" {
	// 			fmt.Println("\n👉 SCAN QR DI BAWAH INI DENGAN WHATSAPP:")

	// 			// Parameter: (data_qr, level_error_correction, output_ke_stdout)
	// 			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)

	// 		} else if evt.Event == "success" {
	// 			fmt.Println("✅ Login Berhasil!")
	// 		}
	// 	}

	// } else {
	// 	err = client.Connect()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// // 4. Kirim Pesan Testing
	// targetNomor := "6285721911444" // Ganti dengan nomor tujuan (Gunakan kode negara)
	// sendText(client, targetNomor, "Halo! testing 2.")

	// //Menjaga program tetap jalan agar tidak langsung mati
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// <-c

	// client.Disconnect()
	os.Exit(1)
}
