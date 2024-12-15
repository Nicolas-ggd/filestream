[![Go Report Card](https://goreportcard.com/badge/github.com/Nicolas-ggd/filestream)](https://goreportcard.com/report/github.com/Nicolas-ggd/filestream)
![Go Version](https://img.shields.io/github/go-mod/go-version/Nicolas-ggd/filestream)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FNicolas-ggd%2Ffilestream.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2FNicolas-ggd%2Ffilestream?ref=badge_shield&issueType=license)
![License](https://img.shields.io/github/license/Nicolas-ggd/filestream)
![Issues Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)
![Version](https://img.shields.io/github/v/release/Nicolas-ggd/filestream)

# FileStream: Event-Driven File Upload Engine

FileStream is an open-source project, which is aim to gain more experience in golang, the purpose of this project is to build robust file upload systems that are efficient, resumable, and ready for real-time integrations.

## Purpose of this project

This project was born from a desire to learn Go deeply while building something practical and useful. My aim is to enhance my skills and gain real world experience by working on an open-source project that can also attract contributions from others who share the same passion. This project is learning journey for me and developers also which decide to collaborating and create a reusable engine.

## Roadmap

Here are some exciting features in the pipeline:
- Chunked uploads - Upload large files in smaller, manageable chunks.
- WebSocket Notifications - Real time updates for uploading process.
- Event driven - Publish events via NATS or other event system.
- Storage - In starting level it's good to support Minio for example.
- Customizable options - Customizable options is like remove metadata or using virus scan.
- Resumable upload - Resume interrupted uploads without starting over.

## How to contribute

Here’s how you can get involved:
1. Report issues: If you find any bug or issue, please open a [issue](https://github.com/Nicolas-ggd/filestream/issues)
2. Fork and code: Check out the [open issues](https://github.com/Nicolas-ggd/filestream/pulls) and submit a pull request.

## Project setup

1. Clone the repository:
    ```
    git clone https://github.com/Nicolas-ggd/filestream
    ```
2. Install dependencies:
    ```
   make dependencies
   ```
3. Run tests:
    ```
   make test
   ```
   
## Usage
```go
import(
    "github.com/gin-gonic/gin"
    "github.com/Nicolas-ggd/filestream"
)

func Upload(c *gin.Context) {
   // Retrieve the uploaded file
   uploadFile, err := c.FormFile("file")
   if err != nil {
     c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file"})
     return
   }

   // Open the uploaded file
   file, err := uploadFile.Open()
   if err != nil {
     fmt.Println("Error opening uploaded file:", err)
     c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
     return
   }
   
   defer file.Close()


   // Construct the FileRequest
   fileReq := fstream.RFileRequest{
       File:            file,
       UploadFile:      uploadFile,
       MaxRange:        rangeMax,
       FileSize:        fileSize,
       UploadDirectory: "uploads/",
       FileUniqueName:  true, // remember that if you set true here, you will receive unique name file
   }


   // Call StoreChunk to handle the uploaded chunk
   prFile, err := fstream.StoreChunk(&fileReq)
   if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": "failed to store chunk"})
      return
   }
   
   // You can write your own if statement to check whatever you want
   if isLast {
    // You can perform your own logic here, before return 200 status, 
    // it's better to remove chunks which is uploaded and doesn't use anymore
     fstream.RemoveUploadedFile(&fileReq)
   }

   c.JSON(http.StatusOK, gin.H{"message": "file chunk processed"})
}
```

`fstream` offers extension check, all you need is that to provide which extension do you want to allow

```go
   import(
     "github.com/Nicolas-ggd/filestream"
   )

    // Declare extension slice, slice contains all file extension which is ok for you to allow in your system
    var allowExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}

    // This function is very simple, but save your time to compare extension and current file to known if it's valid for you
    allowed := fstream.IsAllowExtension(allowExtensions, "filename.png")
	
    if allowed {
     // your logic goes here...	
    }
```
   
## License
FileStream is open-source software licensed under the MIT License.
