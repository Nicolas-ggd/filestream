package fstream

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"image"
	"image/jpeg"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

func TestIsAllowedExtension(t *testing.T) {
	testCases := []struct {
		fileExtension []string
		fileName      string
		expected      bool
		expectedName  string
	}{
		{
			fileExtension: []string{".jpeg", ".jpg"},
			fileName:      "test.jpeg",
			expected:      true,
			expectedName:  "Success",
		},
		{
			fileExtension: []string{".png", ".webp"},
			fileName:      "test.jpg",
			expected:      false,
			expectedName:  "Failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedName, func(t *testing.T) {
			res := IsAllowExtension(tc.fileExtension, tc.fileName)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestRemoveUploadedFile(t *testing.T) {
	// create test directory
	err := os.MkdirAll("test", 0777)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("test")

	// simulate file upload
	testFileName := "testfile.txt"
	testFilePath := filepath.Join("test", testFileName)
	_, err = os.Create(testFilePath)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		uploadDir    string
		fileName     string
		expected     error
		expectedName string
	}{
		{
			uploadDir:    "test",
			fileName:     testFileName,
			expected:     nil,
			expectedName: "Success",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedName, func(t *testing.T) {
			err = RemoveUploadedFile(tc.uploadDir, tc.fileName)
			assert.Equal(t, tc.expected, err)

			if _, statErr := os.Stat(testFilePath); !os.IsNotExist(statErr) {
				t.Errorf("File %s was not removed", testFilePath)
			}
		})
	}
}

func TestStoreChunk(t *testing.T) {
	testCases := []struct {
		name             string
		fileContent      []byte
		maxRange         int
		fileUniqueName   bool
		expectError      bool
		expectedFileSize int
	}{
		{
			name:             "Successful file upload with unique name",
			fileContent:      []byte("This is a test chunk"),
			maxRange:         19, // Full file uploaded
			fileUniqueName:   true,
			expectError:      false,
			expectedFileSize: 19,
		},
		{
			name:             "Successful file upload without unique name",
			fileContent:      []byte("Another test chunk"),
			maxRange:         20, // Full file uploaded
			fileUniqueName:   false,
			expectError:      false,
			expectedFileSize: 20,
		},
		{
			name:             "Partial upload",
			fileContent:      []byte("Partial chunk"),
			maxRange:         7, // Partial file uploaded
			fileUniqueName:   false,
			expectError:      false,
			expectedFileSize: 0, // File should not be finalized
		},
		{
			name:             "Error due to maxRange exceeding file size",
			fileContent:      []byte("Invalid max range"),
			maxRange:         50, // maxRange > FileSize
			fileUniqueName:   false,
			expectError:      true,
			expectedFileSize: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "")
			if err != nil {
				assert.Error(t, err)
			}
			defer os.Remove(file.Name())

			_, err = file.Write(tc.fileContent)
			if err != nil {
				assert.Error(t, err)
			}
			defer file.Close()

			multipartFile, err := os.Open(file.Name())
			if err != nil {
				assert.Error(t, err)
			}
			defer multipartFile.Close()

			fileHeader := &multipart.FileHeader{
				Filename: filepath.Base(file.Name()),
				Size:     int64(len(tc.fileContent)),
			}

			r := &RFileRequest{
				File:            multipartFile,
				UploadFile:      fileHeader,
				MaxRange:        tc.maxRange,
				FileSize:        len(tc.fileContent),
				UploadDirectory: t.TempDir(), // Temporary directory for test
				FileUniqueName:  tc.fileUniqueName,
			}

			resFile, err := StoreChunk(r)
			if err != nil {
				assert.Error(t, err)
			}

			expectedFilePath := filepath.Join(r.UploadDirectory + fileHeader.Filename)
			fmt.Printf("Expected Path: %s\n", expectedFilePath)

			assert.FileExists(t, expectedFilePath)

			actualContent, err := os.ReadFile(expectedFilePath)
			require.NoError(t, err)
			assert.Equal(t, tc.fileContent, actualContent)

			if resFile != nil {
				assert.Equal(t, fileHeader.Filename, resFile.FileName)
				assert.Equal(t, expectedFilePath, resFile.FilePath)
				assert.Equal(t, filepath.Ext(fileHeader.Filename), resFile.FileExtension)
			}
		})
	}
}

func TestRemoveExifMetadata(t *testing.T) {
	testCases := []struct {
		name        string
		setupFunc   func() (string, error)
		expectError bool
	}{
		{
			name: "Exif metadata removed successfully",
			setupFunc: func() (string, error) {
				file, err := os.CreateTemp("", "*.jpg")
				if err != nil {
					return "", err
				}
				defer file.Close()

				// create dummy image
				img := image.NewRGBA(image.Rect(0, 0, 100, 100))
				if err := jpeg.Encode(file, img, nil); err != nil {
					return "", err
				}
				return file.Name(), nil
			},
			expectError: false,
		},
		{
			name: "failed to open image",
			setupFunc: func() (string, error) {
				return "invalid_path.jpg", nil
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputPath, err := tc.setupFunc()
			if err != nil {
				t.Fatal(err)
			}

			err = RemoveExifMetadata(inputPath)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Validate output file
				file, err := os.Open(inputPath)
				assert.NoError(t, err)
				defer file.Close()

				_, _, err = image.Decode(file)
				assert.NoError(t, err)
			}

			if _, err = os.Stat(inputPath); err == nil {
				os.Remove(inputPath)
			}
		})
	}
}
