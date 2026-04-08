package helper

import (
	"context"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// --- FUNGSI KIRIM WHATSAPP (Anti-Error URL/SHA256) ---
func SendImageToGroup(client *whatsmeow.Client, groupID string, filePath string, caption string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ File gambar tidak ditemukan: %v\n", err)
		return
	}

	fmt.Println("☁️  Sedang mengunggah gambar ke server WhatsApp...")
	// Menggunakan whatsmeow.MediaImage yang benar
	resp, err := client.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("❌ Gagal upload gambar: %v\n", err)
		return
	}

	recipient, err := types.ParseJID(groupID)
	if err != nil {
		fmt.Printf("❌ Format ID Grup Salah: %v\n", err)
		return
	}

	// Format pesannya sudah menggunakan standar Protobuf terbaru (HURUF BESAR)
	msg := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Caption:       proto.String(caption),
			Mimetype:      proto.String("image/png"),
			URL:           &resp.URL,
			DirectPath:    &resp.DirectPath,
			MediaKey:      resp.MediaKey,
			FileEncSHA256: resp.FileEncSHA256,
			FileSHA256:    resp.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
		},
	}

	_, err = client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		fmt.Printf("❌ Gagal mengirim ke grup: %v\n", err)
	} else {
		fmt.Println("✅ Berhasil! Gambar telah terkirim ke grup.")
	}
}
