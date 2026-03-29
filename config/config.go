package config

import "github.com/spf13/viper"

// Struktur untuk menampung konfigurasi (opsional tapi rapi)
type Config struct {
	MYTOKEN string `mapstructure:"MYTOKEN"`
	//TargetWA     string `mapstructure:"TARGET_WA"`
	GRAFANAURL1 string `mapstructure:"GRAFANAURL1"`
	GRAFANAURL2 string `mapstructure:"GRAFANAURL2"`
	GRAFANAURL3 string `mapstructure:"GRAFANAURL3"`
	GRAFANAURL4 string `mapstructure:"GRAFANAURL4"`
	GRAFANAURL5 string `mapstructure:"GRAFANAURL5"`
	GRAFANAURL6 string `mapstructure:"GRAFANAURL6"`
	GRAFANAURL7 string `mapstructure:"GRAFANAURL7"`
}

func LoadConfig() (config Config, err error) {
	// 1. Beritahu Viper lokasi dan nama filenya
	viper.AddConfigPath(".")    // Cari di folder saat ini
	viper.SetConfigName(".env") // Nama file
	viper.SetConfigType("env")  // Ekstensi file

	// 2. Baca Environment Variables dari sistem jika ada
	viper.AutomaticEnv()

	// 3. Baca file konfigurasi
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// 4. Masukkan (unmarshal) nilai ke struct Config
	err = viper.Unmarshal(&config)
	return
}
