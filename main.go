package main

import (
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"brange.net/quicknote/config"
	"brange.net/quicknote/db"
)

var cfg config.Config

func init() {
	rand.Seed(time.Now().UnixNano())
	cfg = config.LoadConfig()
	db.InitDB(cfg)
}

func generateLink() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	link := make([]byte, 10)
	for i := range link {
		link[i] = charset[rand.Intn(len(charset))]
	}
	return string(link)
}

// Handlers
func createHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/create.html"))
	tmpl.Execute(w, nil)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	content := r.FormValue("content")
	link := generateLink()

	const maxLength = 200
	if len(content) > maxLength {
		http.Error(w, "Error: Message exceeds the maximum length of 200 characters.", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err == nil {
		defer file.Close()

		dir := "uploads/" + link + "/"
		if err := os.MkdirAll(dir, 0755); err != nil {
			http.Error(w, "Error creating upload directory", http.StatusInternalServerError)
			return
		}

		out, err := os.Create(filepath.Join(dir, header.Filename))
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			http.Error(w, "Error writing file", http.StatusInternalServerError)
			return
		}

		createdPath := filepath.Join(dir, ".created")
		t := time.Now().UTC().Format(time.RFC3339)
		_ = os.WriteFile(createdPath, []byte(t+"\n"), 0644)
	}

	if err := db.SaveMessage(content, link); err != nil {
		http.Error(w, "Error saving message", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/link.html"))
	_ = tmpl.Execute(w, map[string]string{"Link": "localhost:8080/view/" + link})
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	link := r.URL.Path[len("/view/"):]
	content, err := db.GetMessage(link)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	defer db.DeleteMessage(link)

	var fileName string
	var isImage bool
	files, err := os.ReadDir("uploads/" + link)
	if err == nil && len(files) > 0 {
		for _, file := range files {
			if file.Name() != ".created" {
				fileName = file.Name()
				ext := filepath.Ext(fileName)
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
					isImage = true
				}
				break
			}
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/view.html"))
	err = tmpl.Execute(w, map[string]interface{}{
		"Content":    content,
		"FileExists": fileName != "",
		"Link":       link,
		"FileName":   fileName,
		"IsImage":    isImage,
	})
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/download/"):]

	fullPath := "uploads/" + filePath

	file, err := os.Open(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileName := filepath.Base(fullPath)

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")

	io.Copy(w, file)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/download/", downloadHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	http.ListenAndServe(":8080", nil)
}
