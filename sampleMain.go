// package main

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"go.mau.fi/whatsmeow"
// 	"go.mau.fi/whatsmeow/store/sqlstore"
// 	"go.mau.fi/whatsmeow/types/events"
// 	waLog "go.mau.fi/whatsmeow/util/log"
// )

// func main() {

// 	// 1. Setup Logging & Database
// 	dbLog := waLog.Stdout("Database", "DEBUG", true)
// 	// BENAR:
// 	container, err := sqlstore.New(context.Background(), "sqlite3", "file:whatsapp_session.db?_foreign_keys=on", dbLog)
// 	if err != nil {
// 		panic(err)
// 	}

// 	deviceStore, err := container.GetFirstDevice(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}

// 	clientLog := waLog.Stdout("Client", "DEBUG", true)
// 	client := whatsmeow.NewClient(deviceStore, clientLog)

// 	// 2. Handler (Opsional jika hanya ingin mengirim, bukan menerima)
// 	client.AddEventHandler(func(evt interface{}) {
// 		switch evt.(type) { // Hapus 'v :='
// 		case *events.Connected:
// 			fmt.Println("WhatsApp Terhubung!")
// 		}
// 	})

// 	// 3. Login / Scan QR
// 	if client.Store.ID == nil {
// 		qrChan, _ := client.GetQRChannel(context.Background())
// 		err = client.Connect()
// 		if err != nil {
// 			panic(err)
// 		}
// 		for evt := range qrChan {
// 			if evt.Event == "code" {
// 				// Tampilkan QR di terminal atau buat file PNG QR
// 				fmt.Println("Scan QR ini:", evt.Code)
// 			} else {
// 				fmt.Println("Hasil Login:", evt.Event)
// 			}
// 		}
// 	} else {
// 		err = client.Connect()
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	// Tunggu sebentar sampai koneksi stabil
// 	fmt.Println("Siap mengirim pesan...")

// 	// 4. Contoh Fungsi Pengiriman Screenshot (nanti dipanggil di Cron)
// 	// sendGrafanaScreenshot(client, "6281234567890@s.whatsapp.net", "screenshot.png")

//		// Menjaga aplikasi tetap jalan
//		c := make(chan os.Signal, 1)
//		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//		<-c
//		client.Disconnect()
//	}
package main
