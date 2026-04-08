package helper

import (
	"context"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// --- FUNGSI BARU: ANALISIS GAMBAR DENGAN GEMINI ---
func analyzeImageWithGemini(ctx context.Context, apiKey, imagePath string) (string, error) {
	// A. Inisialisasi Client Gemini
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	// B. Gunakan model Gemini 1.5 Flash (Cepat & Hemat untuk analisis visual dasar)
	//model := client.GenerativeModel("gemini-1.5-flash")
	model := client.GenerativeModel("models/gemini-flash-latest")

	// C. Baca file gambar menjadi bytes
	imgData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	// D. Tulis PROMPT (Perintah) untuk AI
	// Semakin spesifik prompt Anda, semakin bagus hasilnya.
	prompt := genai.Text("Gemini analisis gambar dashboard monitoring Grafana ini. Sebutkan jika ada lonjakan grafik yang tidak wajar, error (biasanya warna merah), atau status tidak normal. Berikan ringkasan kondisi sistem saat ini secara singkat dan jelas dalam Bahasa Indonesia.")

	// E. Kirim gambar dan prompt ke Gemini
	resp, err := model.GenerateContent(ctx,
		prompt,
		//genai.ImageData("image/png", imgData),
		genai.ImageData("png", imgData),
	)
	if err != nil {
		return "", err
	}

	// F. Ambil teks hasil generate
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		// Mengasumsikan respons berupa teks tunggal
		if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(text), nil
		}
	}

	return "Gemini tidak memberikan analisis.", nil
}
