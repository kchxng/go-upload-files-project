package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title File Upload API
// @version 1.0
// @description API for uploading multiple files
// @host localhost:4000
// @BasePath /

// @Summary Upload multiple files
// @Description Upload multiple files to the server
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Files to upload"
// @Success 200 {string} string "Files uploaded successfully"
// @Failure 400 {string} string "Failed to parse form"
// @Failure 500 {string} string "Failed to save file"
// @Router /upload [post]
func uploadFiles(c *fiber.Ctx) error {
	// Parse the multipart form
	form, err := c.MultipartForm()
	fmt.Printf("file: %v\n", form)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	// Get all files from the "files" field
	files := form.File["files"]
	if len(files) == 0 {
		file := &multipart.FileHeader{
			Filename: "no_upload_file.txt",
			Size:     int64(10 * 1024 * 1024), // 10 MB
			Header:   make(map[string][]string),
		}
		log.Printf("Custom no file upload: %v", file)
		// return file, nil
		if file.Filename == "no_upload_file.txt" {
			log.Println("Custom no upload file")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file custom that mean no file uploaded"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No files uploaded"})
	} else {
		fileCount, err := UploadFiles(c, files)
		if err != nil {
			log.Printf("err: %v\n", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("%d files uploaded successfully", *fileCount),
		})
	}
}

func UploadFiles(c *fiber.Ctx, files []*multipart.FileHeader) (*int, error) {
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory")
	}
	for i, file := range files {
		// filePath := filepath.Join(uploadDir, file.Filename)
		log.Println(file.Filename)
		fileExt := GetFileExtension(file.Filename)
		timestamp := time.Now().UnixMicro()
		newFileName := fmt.Sprintf("order(%v)_%d%s", i, timestamp, fileExt)
		filePath := fmt.Sprintf("%s/%s", uploadDir, newFileName)
		if err := c.SaveFile(file, filePath); err != nil {
			return nil, fmt.Errorf("failed to save file: %s", file.Filename)
		}

	}
	// msg := fmt.Sprintf("%d files uploaded successfully", len(files))
	filesCount := len(files)
	log.Printf("%d files uploaded successfully\n", filesCount)
	return &filesCount, nil
}

func GetFileExtension(fileName string) string {
	fileName = strings.ToLower(fileName) // Ensure case insensitivity
	for i := len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '.' {
			fmt.Println(fileName[i:])
			return fileName[i:]
		}
	}
	fmt.Println(fileName)
	return ""
}

func main() {
	app := fiber.New()

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// File upload route
	app.Post("/upload", uploadFiles)
	app.Static("/file-attach", "./uploads")

	log.Fatal(app.Listen(":4000"))
}
