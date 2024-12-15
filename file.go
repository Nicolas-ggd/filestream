package fstream

import (
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// File struct is final face of uploaded file, it includes necessary field to use them after file is uploaded
type File struct {
	// Original uploaded file name
	FileName string
	// FileUniqueName is unique name
	FileUniqueName string
	// Uploaded file path
	FilePath string
	// Uploaded file extension
	FileExtension string
	// Uploaded file size
	FileSize string
}

// RFileRequest struct is used for http request, use this struct to bind uploaded file
type RFileRequest struct {
	// File is an interface to access the file part of a multipart message.
	File multipart.File
	// A FileHeader describes a file part of a multipart request.
	UploadFile *multipart.FileHeader
	// Maximum range of chunk uploads
	MaxRange int
	// Uploaded file size
	FileSize int
	// Upload directory
	UploadDirectory string
	// FileUniqueName is identifier to generate unique name for files
	FileUniqueName bool
}

// generateUuidUniqueName function generates unique string using UUID
func (r *File) generateUuidUniqueName(request *RFileRequest) {
	ext := filepath.Ext(request.UploadFile.Filename)

	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalln(err)
	}

	r.FileUniqueName = fmt.Sprintf("%s%s", id.String(), ext)
}

// RemoveUploadedFile function removes uploaded file from uploaded directory, it takes param and returns nothing:
//
// Takes:
//
//   - RFileRequest struct
//
// Use this function in your handler after file is uploaded
func RemoveUploadedFile(r *RFileRequest) {
	filePath := filepath.Join(r.UploadDirectory, r.UploadFile.Filename)

	e := os.Remove(filePath)
	if e != nil {
		log.Printf("error removing file: %v", e)
	}
}

// prettyByteSize function is used to concrete the file size
func prettyByteSize(b int) string {
	bf := float64(b)

	for _, unit := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
		if math.Abs(bf) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", bf, unit)
		}
		bf /= 1024.0
	}

	return fmt.Sprintf("%.1fYiB", bf)
}

// StoreChunk cares slice of chunks and returns final results and error
//
//   - File -  struct is final version about file information
//
//   - error - functions cares about occurred errors and returns it.
//
// Function creates new directory for chunks if it doesn't exist, if directory already exists it appends received chunks in current chunks and if entire file is uploaded then File struct is returned
func StoreChunk(r *RFileRequest) (*File, error) {
	var rFile *File

	// Create new directory for uploaded chunk
	filePath := filepath.Join(r.UploadDirectory + r.UploadFile.Filename)

	// Create new directory if it doesn't exist
	if _, err := os.Stat(r.UploadDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(r.UploadDirectory, 0777)
		if err != nil {
			return nil, fmt.Errorf("failed to create new temporary directory: %v", err)
		}
	}

	// Open the file for appending and creating if it doesn't exist
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Printf("error closing file: %v", cerr)
		}
	}()

	// Copy the chunk data to the file
	if _, err := io.Copy(f, r.File); err != nil {
		return nil, fmt.Errorf("failed to copying file: %v", err)
	}

	// If the entire file is uploaded, finalize entire process and return file information
	if r.MaxRange >= r.FileSize {
		fileInfo, err := f.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to stating file: %v", err)
		}

		// Calculate file size in bytes
		size := prettyByteSize(int(fileInfo.Size()))

		// Bind File struct and return
		rFile = &File{
			FileName:       r.UploadFile.Filename,
			FileUniqueName: r.UploadFile.Filename,
			FilePath:       filePath,
			FileExtension:  filepath.Ext(r.UploadFile.Filename),
			FileSize:       size,
		}
	}

	return rFile, nil
}

// IsAllowExtension function checks if file extension is allowed to upload, it takes following params
//
//   - fileExtensions - array of strings, which is looks like: []string{".jpg", ".jpeg"}, note that this is fileExtensions which is allowed to receive
//
//   - fileName - string, this parameter is file name which is like ".jpeg", ".jpg"
//
// Returns:
//
//   - bool - function returns false if extension isn't allowed to receive, it returns true if extension is allowed to receive
func IsAllowExtension(fileExtensions []string, fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))

	// range over the received extensions to check if file is ok to accept
	for _, allowed := range fileExtensions {
		if ext == allowed {
			return true
		}
	}

	return false
}
