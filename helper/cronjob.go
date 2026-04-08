package helper

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func CronjobSimple() {
	c := cron.New()
	// Contoh: Mengirim setiap hari jam 8 pagi
	c.AddFunc("0 8 * * *", func() {
		log.Println("Memulai proses otomatisasi...")
		// Panggil fungsi screenshot & kirim WA di sini
	})
	c.Start()
}
func Cronjob() {
	// Membuat scheduler baru
	c := cron.New()

	// Menambahkan fungsi dengan standar cron: setiap 10 menit
	// Format: Menit Jam HariBulan Bulan HariMinggu
	_, err := c.AddFunc("*/1 * * * *", func() {
		log.Println("--- Memulai proses otomatisasi (Setiap 1 Menit) ---")

		// Menambahkan Timeout agar jika proses screenshot macet,
		// tidak membebani memori selamanya.
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		runAutomationTask(ctx)
	})

	if err != nil {
		log.Fatalf("Gagal menambahkan cron: %v", err)
	}

	// Memulai penjadwal di background
	c.Start()
	log.Println("Cronjob berhasil dijalankan...")
}

func runAutomationTask(ctx context.Context) {
	// Simulasi proses
	log.Println("Start Executing")

	// Simpan Function yang ingin dijalankan dibawah ini
	//SendTextWA("Pesan ini otomatis dari system !! Laporan Grafana hari ini")

	// Gunakan select untuk menangani timeout atau proses selesai
	select {
	case <-time.After(5 * time.Second): // Simulasi proses 5 detik
		log.Println("Proccess Complete !!")
	case <-ctx.Done():
		log.Println("Error: Proses dihentikan karena mencapai batas waktu (timeout)")
	}
}
