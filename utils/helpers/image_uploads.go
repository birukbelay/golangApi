package helpers

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/birukbelay/item/utils/global"
)

//const MaxUploadSize = 2 * 1024 * 1024 // 2 mb
const uploadPath = "../../public/assets/images"

func UploadFile(r *http.Request, updating bool, updatingFileName string, uniquePath string) (string, error, string, int) {



	//LogTrace("here ", "==========")
	// parse and validate file and post parameters
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		//Tr(global.InvalidFile, err)
		return "", err, global.InvalidFile, http.StatusBadRequest
	}



	defer file.Close()
	// Get and print out file size
	fileSize := fileHeader.Size

	//LogTrace("Upload-> file size ", fileSize)
	//fmt.Printf("File size (bytes): %v\n", fileSize)
	// validate file size
	if fileSize > global.MaxUploadSize {

		return "", err, global.FileTooBig, http.StatusBadRequest
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err, global.InvalidFileType, http.StatusBadRequest
	}

	// check file type,   only needs the first 512 bytes
	detectedFileType := http.DetectContentType(fileBytes)
	switch detectedFileType {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png":
		//case "application/pdf":
		break
	default:
		return "", err, "INVALID_FILE_TYPE", http.StatusBadRequest
	}


	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		return "", err, "CANT_READ_FILE_TYPE", http.StatusInternalServerError
	}
	//LogTrace("fileEndings", fileEndings)



	var imgName string
	if !updating{
		fileName := randToken(12)
		img := strings.TrimSuffix(fileHeader.Filename, fileEndings[0])
		imgName = img + "-" + fileName +fileEndings[0]
	}else{
		imgName = updatingFileName
	}

	wd, err := os.Getwd()
	fmt.Println("wd..",wd)


	if err != nil {
		return "", err, "CANT_READ_FILE_TYPE", http.StatusInternalServerError
	}

	newPath := filepath.Join(wd, "public", "assets", "images", uniquePath)
	fullPath := filepath.Join(newPath, imgName)

	//LogTrace("path", newPath)

	//fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)
	// write file
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		err := os.MkdirAll(newPath, 0777)
		if err != nil {
			return "", err,"CANT_CREATE_PATH", http.StatusInternalServerError
		}
	}
	newFile, err := os.Create(fullPath)
	if err != nil {
		LogTrace("os.Create",err)
		return "", err, "CANT_WRITE_FILE", http.StatusInternalServerError
	}

	defer newFile.Close() // idempotent, okay to call twice
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {

		return "", err, global.CantWriteFile, http.StatusInternalServerError
	}
	return imgName, nil, global.Success, http.StatusCreated

}

func randToken(len int) string {
	b := make([]byte, len)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
