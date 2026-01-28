package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// database model
type Produk struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}


// masukin ke database
var produk = []Produk{
	{
		ID:          1,
		Name:        "Baju",
		Description: "Benda untuk dipakai di badan",
	},
	{
		ID:          2,
		Name:        "Celana",
		Description: "benda untuk dipakai di badan bagian bawah",
	},
	{
		ID:          3,
		Name:        "Jaket",
		Description: "benda untuk dipakai di badan bagian atas dan menutupi tangan",
	},
}

func getProdukById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/") // mengembalikan angka setelah path
	id, err := strconv.Atoi(idStr) // Atoi mengembalikan 2 value jadi butuh 2 variable
	if err != nil {  // bukan nil berarti ada error, jika ada error
		http.Error(w, "invalid Produk ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk { // loop for range, abaikan index karena hanya untuk copy
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json") // menulis respon header
			json.NewEncoder(w).Encode(p) // jadikan bentuk json
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func updateProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid Produk ID", http.StatusBadRequest)
		return
	}

	var updateProduk Produk

	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "invalid Request", http.StatusBadRequest)
		return
	}

	for i := range produk { // ini for loop untuk merubah karena pake i - index
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}

}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid Produk ID", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)

			// index diset ulang
			for j := range produk {
				produk[j].ID = j + 1
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Produk berhasil dihapus",
			})
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

func main() {
	// GET api produk detail
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukById(w, r)
		} else if r.Method == "PUT" {
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}
	})

	// GET api-produk

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			var produkBaru Produk

			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "invalid Request", http.StatusBadRequest)
				return
			}

			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(produkBaru)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "server is running",
		})
	})
	fmt.Println("server running di port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
