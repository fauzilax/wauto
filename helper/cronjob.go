package helper

import (
	"log"

	"github.com/robfig/cron/v3"
)

func cronjob() {
	c := cron.New()
	// Contoh: Mengirim setiap hari jam 8 pagi
	c.AddFunc("0 8 * * *", func() {
		log.Println("Memulai proses otomatisasi...")
		// Panggil fungsi screenshot & kirim WA di sini
	})
	c.Start()
}
