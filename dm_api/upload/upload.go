package upload

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type SharedFile struct {
	Path        string `json:"path"`
	MimeType    string `json:"mimeType"`
	Extension   string `json:"extension"`
	Size        int64  `json:"size"`
	DisplayName string `json:"displayName"`
}

// func invoiceHandler(w http.ResponseWriter, r *http.Request) {
// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer file.Close()

// 	// Сохранение файла на сервере
// 	filePath := "/var/dm/invoices/" + header.Filename
// 	out, err := os.Create(filePath)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer out.Close()

// 	_, err = io.Copy(out, file)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	var sharedFiles []SharedFile

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&sharedFiles); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, file := range sharedFiles {
		src, err := os.Open(file.Path)
		if err != nil {
			log.Println(err)
			continue
		}
		defer src.Close()

		destPath := filepath.Join("/opt/dm/uploads", "1.bin")
		dest, err := os.Create(destPath)
		if err != nil {
			log.Println(err)
			continue
		}
		defer dest.Close()

		if _, err := io.Copy(dest, src); err != nil {
			log.Println(err)
			continue
		}

		log.Printf("File %s uploaded to %s\n", file.DisplayName, destPath)
	}

	w.WriteHeader(http.StatusOK)
}
