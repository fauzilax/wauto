package helper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// --- FUNGSI DOWNLOAD DARI GRAFANA LINUX ---
func DownloadGrafanaImage(url string, token string, filepath string) error {
	fmt.Println("⏳ Meminta server Linux merender dashboard (biasanya butuh 10-15 detik)...")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	// Set timeout 60 detik (rendering di Linux butuh waktu)
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server grafana menolak request. Status Code: %d", resp.StatusCode)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
