package main

import (
	"embed"
	"log"
	"path"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/runtime"
)

//go:embed public/*
var f embed.FS

func main() {
	dirs, err := f.ReadDir("public")
	if err != nil {
		log.Fatalf("Error: failed to read directory: %+v", err)
	}

	files := []runtime.ServeFile{}

	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}

		filename := dir.Name()
		extension := path.Ext(filename)

		data, err := f.ReadFile("public/" + filename)
		if err != nil {
			log.Fatalf("Error: failed to get file data for %s: %+v", filename, err)
		}

		switch extension {
		case ".css":
			log.Printf("Serving file: %s", filename)
			files = append(files, runtime.ServeFile{
				Path:     filename,
				Data:     data,
				MimeType: "text/css",
			})
		case ".ttf":
			log.Printf("Serving file: %s", filename)
			files = append(files, runtime.ServeFile{
				Path:     filename,
				Data:     data,
				MimeType: "font/ttf",
			})
		default:
			log.Printf("Warn: unrecognised file type '%s'", filename)
		}
	}

	runtime.RunFileServer(files)
}
