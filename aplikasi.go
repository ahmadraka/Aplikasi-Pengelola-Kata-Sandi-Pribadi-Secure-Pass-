package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
)

type account struct {
	nama_layanan  string
	email         string
	kata_sandi    string
	tgl_pembaruan string
	urutanInput   int
}

const namaAplikasi = "SecurePass"

var daftarAkun []account
var counterInput int = 0
var Input = bufio.NewScanner(os.Stdin)

func bacaInput() string {
	Input.Scan()
	return strings.TrimSpace(Input.Text())
}

func waktuSekarang() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func header() {
	fmt.Println("Aplikasi Pengelola Kata Sandi Pribadi (SecurePass)")
}

func klasifikasiKataSandi(password string) string {
	var panjang int = len(password)
	if panjang == 0 {
		return "Kosong"
	}

	besar := false
	kecil := false
	angka := false
	simbol := false

	for i := 0; i < len(password); i++ {
		c := rune(password[i])
		switch {
		case unicode.IsUpper(c):
			besar = true
		case unicode.IsLower(c):
			kecil = true
		case unicode.IsDigit(c):
			angka = true
		default:
			simbol = true
		}
	}

	if panjang < 8 {
		return "Lemah"
	}
	if besar && kecil && angka && simbol {
		return "Kuat"
	}
	return "Sedang"
}

func tambahAkun() {
	fmt.Println("Tambah Akun Baru")

	fmt.Print("Nama Layanan : ")
	nama := bacaInput()
	if nama == "" {
		fmt.Println("Input tidak valid!")
		return
	}

	fmt.Print("Email        : ")
	email := bacaInput()

	fmt.Print("Kata Sandi   : ")
	sandi := bacaInput()

	counterInput++

	akun := account{
		nama_layanan:  nama,
		email:         email,
		kata_sandi:    sandi,
		tgl_pembaruan: waktuSekarang(),
		urutanInput:   counterInput,
	}

	daftarAkun = append(daftarAkun, akun)
	fmt.Println("Akun berhasil ditambahkan!")
}

func ubahAkun() {
	if len(daftarAkun) == 0 {
		fmt.Println("Belum ada akun tersimpan.")
		return
	}

	tampilkanDaftarAkun()

	fmt.Print("\nNomor akun yang ingin diubah: ")
	var nomor int
	fmt.Sscanf(bacaInput(), "%d", &nomor)

	if nomor < 1 || nomor > len(daftarAkun) {
		fmt.Println("Input tidak valid!")
		return
	}

	idx := nomor - 1
	fmt.Println("(Biarkan kosong untuk mempertahankan nilai lama)")

	fmt.Printf("Nama Layanan [%s]: ", daftarAkun[idx].nama_layanan)
	if v := bacaInput(); v != "" {
		daftarAkun[idx].nama_layanan = v
	}

	fmt.Printf("Email [%s]: ", daftarAkun[idx].email)
	if v := bacaInput(); v != "" {
		daftarAkun[idx].email = v
	}

	fmt.Printf("Kata Sandi [%s]: ", daftarAkun[idx].kata_sandi)
	if v := bacaInput(); v != "" {
		daftarAkun[idx].kata_sandi = v
	}

	daftarAkun[idx].tgl_pembaruan = waktuSekarang()
	fmt.Println("Akun berhasil diubah!")
}

func hapusAkun() {
	if len(daftarAkun) == 0 {
		fmt.Println("Belum ada akun tersimpan.")
		return
	}

	tampilkanDaftarAkun()

	fmt.Print("\nNomor akun yang ingin dihapus: ")
	var nomor int
	fmt.Sscanf(bacaInput(), "%d", &nomor)

	if nomor < 1 || nomor > len(daftarAkun) {
		fmt.Println("Nomor akun tidak valid.")
		return
	}

	namaHapus := daftarAkun[nomor-1].nama_layanan
	daftarAkun = append(daftarAkun[:nomor-1], daftarAkun[nomor:]...)
	fmt.Printf("Akun '%s' berhasil dihapus!\n", namaHapus)
}

func tampilkanDaftarAkun() {
	if len(daftarAkun) == 0 {
		fmt.Println("Belum ada akun tersimpan.")
		return
	}

	fmt.Println("Daftar Akun")
	fmt.Printf("%-4s %-20s %-28s %-14s %-8s %-20s\n",
		"No.", "Nama Layanan", "Email", "Kata Sandi", "Kekuatan", "Terakhir Diperbarui")

	for i := 0; i < len(daftarAkun); i++ {
		akun := daftarAkun[i]
		kekuatan := klasifikasiKataSandi(akun.kata_sandi)
		fmt.Printf("%-4d %-20s %-28s %-14s %-8s %-20s\n",
			i+1, akun.nama_layanan, akun.email,
			akun.kata_sandi, kekuatan, akun.tgl_pembaruan)
	}
}

func sequentialSearch(keyword string) []int {
	var hasil []int
	kw := strings.ToLower(keyword)

	for i := 0; i < len(daftarAkun); i++ {
		if strings.Contains(strings.ToLower(daftarAkun[i].nama_layanan), kw) {
			hasil = append(hasil, i)
		}
	}
	return hasil
}

func binarySearch(keyword string) int {
	sorted := make([]account, len(daftarAkun))
	for i := 0; i < len(daftarAkun); i++ {
		sorted[i] = daftarAkun[i]
	}

	for i := 1; i < len(sorted); i++ {
		key := sorted[i]
		j := i - 1
		for j >= 0 && strings.ToLower(sorted[j].nama_layanan) > strings.ToLower(key.nama_layanan) {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}

	kw := strings.ToLower(keyword)
	kiri, kanan := 0, len(sorted)-1

	for kiri <= kanan {
		tengah := (kiri + kanan) / 2
		namaTengah := strings.ToLower(sorted[tengah].nama_layanan)

		if namaTengah == kw {
			for i := 0; i < len(daftarAkun); i++ {
				if strings.ToLower(daftarAkun[i].nama_layanan) == kw {
					return i
				}
			}
			return -1
		} else if namaTengah < kw {
			kiri = tengah + 1
		} else {
			kanan = tengah - 1
		}
	}
	return -1
}

func cariAkun() {
	if len(daftarAkun) == 0 {
		fmt.Println("Belum ada akun tersimpan.")
		return
	}

	fmt.Println("Cari Akun")
	fmt.Println("1. Sequential Search")
	fmt.Println("2. Binary Search")
	fmt.Print("Pilih metode: ")
	pilihan := bacaInput()

	fmt.Print("Masukkan nama layanan: ")
	keyword := bacaInput()

	switch pilihan {
	case "1":
		hasil := sequentialSearch(keyword)
		if len(hasil) == 0 {
			fmt.Printf("Layanan '%s' tidak ditemukan.\n", keyword)
			return
		}
		fmt.Printf("Hasil untuk '%s': %d akun ditemukan\n", keyword, len(hasil))
		fmt.Printf("%-4s %-20s %-28s %-14s %-8s\n", "No.", "Nama Layanan", "Email", "Kata Sandi", "Kekuatan")
		for i := 0; i < len(hasil); i++ {
			idx := hasil[i]
			akun := daftarAkun[idx]
			fmt.Printf("%-4d %-20s %-28s %-14s %-8s\n",
				idx+1, akun.nama_layanan, akun.email,
				akun.kata_sandi, klasifikasiKataSandi(akun.kata_sandi))
		}

	case "2":
		idx := binarySearch(keyword)
		if idx == -1 {
			fmt.Printf("Layanan '%s' tidak ditemukan.\n", keyword)
			return
		}
		akun := daftarAkun[idx]
		fmt.Printf("Akun ditemukan di posisi #%d:\n", idx+1)
		fmt.Printf("  Nama Layanan : %s\n", akun.nama_layanan)
		fmt.Printf("  Email        : %s\n", akun.email)
		fmt.Printf("  Kata Sandi   : %s\n", akun.kata_sandi)
		fmt.Printf("  Kekuatan     : %s\n", klasifikasiKataSandi(akun.kata_sandi))
		fmt.Printf("  Diperbarui   : %s\n", akun.tgl_pembaruan)

	default:
		fmt.Println("Pilihan tidak valid.")
	}
}

func lebihKecil(a, b account, berdasarkanNama bool, ascending bool) bool {
	if berdasarkanNama {
		if ascending {
			return strings.ToLower(a.nama_layanan) < strings.ToLower(b.nama_layanan)
		}
		return strings.ToLower(a.nama_layanan) > strings.ToLower(b.nama_layanan)
	}
	if ascending {
		return a.urutanInput < b.urutanInput
	}
	return a.urutanInput > b.urutanInput
}

func selectionSort(berdasarkanNama bool, ascending bool) {
	n := len(daftarAkun)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if lebihKecil(daftarAkun[j], daftarAkun[minIdx], berdasarkanNama, ascending) {
				minIdx = j
			}
		}
		daftarAkun[i], daftarAkun[minIdx] = daftarAkun[minIdx], daftarAkun[i]
	}
}

func insertionSort(berdasarkanNama bool, ascending bool) {
	n := len(daftarAkun)
	for i := 1; i < n; i++ {
		key := daftarAkun[i]
		j := i - 1
		for j >= 0 && lebihKecil(key, daftarAkun[j], berdasarkanNama, ascending) {
			daftarAkun[j+1] = daftarAkun[j]
			j--
		}
		daftarAkun[j+1] = key
	}
}

func urutkanAkun() {
	if len(daftarAkun) == 0 {
		fmt.Println("Belum ada akun tersimpan.")
		return
	}

	fmt.Println("Urutkan Akun")
	fmt.Println("Berdasarkan:")
	fmt.Println("  1. Nama Layanan")
	fmt.Println("  2. Waktu Input")
	fmt.Print("Pilih kriteria: ")
	pilihUrut := bacaInput()

	if pilihUrut != "1" && pilihUrut != "2" {
		fmt.Println("Pilihan tidak valid.")
		return
	}

	fmt.Println("Metode Pengurutan:")
	fmt.Println("  1. Selection Sort")
	fmt.Println("  2. Insertion Sort")
	fmt.Print("Pilih metode: ")
	pilihMetode := bacaInput()

	berdasarkanNama := pilihUrut == "1"

	fmt.Println("Arah Pengurutan:")
	fmt.Println("  1. Ascending (terlama)")
	fmt.Println("  2. Descending (terbaru)")
	fmt.Print("Pilih arah: ")
	pilihArah := bacaInput()

	if pilihArah != "1" && pilihArah != "2" {
		fmt.Println("Pilihan tidak valid.")
		return
	}

	ascending := pilihArah == "1"

	labelUrut := "Waktu Input"
	if berdasarkanNama {
		labelUrut = "Nama Layanan"
	}
	labelArah := "Descending"
	if ascending {
		labelArah = "Ascending"
	}

	switch pilihMetode {
	case "1":
		selectionSort(berdasarkanNama, ascending)
		fmt.Printf("Data diurutkan berdasarkan %s %s menggunakan Selection Sort!\n", labelUrut, labelArah)
	case "2":
		insertionSort(berdasarkanNama, ascending)
		fmt.Printf("Data diurutkan berdasarkan %s %s menggunakan Insertion Sort!\n", labelUrut, labelArah)
	default:
		fmt.Println("Pilihan tidak valid.")
		return
	}

	tampilkanDaftarAkun()
}

func tampilkanStatistik() {
	fmt.Println("\nSTATISTIK")

	totalAkun := len(daftarAkun)
	fmt.Printf("Total Akun Tersimpan : %d\n\n", totalAkun)

	if totalAkun == 0 {
		fmt.Println("Belum ada akun untuk ditampilkan statistiknya.")
		return
	}

	var (
		jumlahKosong int
		jumlahLemah  int
		jumlahSedang int
		jumlahKuat   int
	)

	for i := 0; i < len(daftarAkun); i++ {
		switch klasifikasiKataSandi(daftarAkun[i].kata_sandi) {
		case "Kosong":
			jumlahKosong++
		case "Lemah":
			jumlahLemah++
		case "Sedang":
			jumlahSedang++
		case "Kuat":
			jumlahKuat++
		}
	}

	fmt.Println("\nKlasifikasi Kekuatan Kata Sandi:")
	fmt.Printf("  Kuat   : %d akun\n", jumlahKuat)
	fmt.Printf("  Sedang : %d akun\n", jumlahSedang)
	fmt.Printf("  Lemah  : %d akun\n", jumlahLemah)
	if jumlahKosong > 0 {
		fmt.Printf("  Kosong : %d akun\n", jumlahKosong)
	}

	fmt.Println("\nKriteria Kekuatan:")
	fmt.Println("  Lemah  : < 8 karakter")
	fmt.Println("  Sedang : >= 8 karakter, belum semua kriteria terpenuhi")
	fmt.Println("  Kuat   : >= 8 karakter + huruf besar + huruf kecil + angka + simbol")
}

func menu() {
	header()

	for {
		fmt.Println("\nMENU UTAMA")
		fmt.Println("1. Tambah Akun")
		fmt.Println("2. Ubah Akun")
		fmt.Println("3. Hapus Akun")
		fmt.Println("4. Tampilkan Semua Akun")
		fmt.Println("5. Cari Akun")
		fmt.Println("6. Urutkan Akun")
		fmt.Println("7. Statistik")
		fmt.Println("0. Keluar")
		fmt.Print("Pilih menu: ")
		pilihan := bacaInput()

		switch pilihan {
		case "1":
			tambahAkun()
		case "2":
			ubahAkun()
		case "3":
			hapusAkun()
		case "4":
			tampilkanDaftarAkun()
		case "5":
			cariAkun()
		case "6":
			urutkanAkun()
		case "7":
			tampilkanStatistik()
		case "0":
			fmt.Println("Terima kasih telah menggunakan Aplikasi SecurePass!")
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid, coba lagi.")
		}
	}
}

func main() {
	menu()
}
