package storage

import (
	"easygo/conf"
	"easygo/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrIncorrectFileID = errors.New("invalid FileId")
	ErrNoFilesFound    = errors.New("no files were found")
)

type Result struct {
	FileName string `json:"fileName"`
	FileId   string `json:"fileId"`
}

type FileInfo struct {
	File     *os.File
	FileName string
}

func Save(part *multipart.Part) (result *Result, err error) {
	path, fileId, err := handlePath(part.FileName())
	if err != nil {
		fmt.Println("GetPath Error:", err)
		return
	}
	dst, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("OpenFile Error:", err)
		return
	}
	defer dst.Close()
	var done bool = false
	var n int
	for {
		buffer := make([]byte, 4096)
		n, err = part.Read(buffer)
		if err != nil {
			if err == io.EOF {
				done = true
			} else {
				fmt.Println("Part read error:", err)
				return
			}
		}
		n, err = dst.Write(buffer[0:n])
		if err != nil {
			fmt.Println("Part write error:", err)
			return
		}
		if done {
			break
		}
	}
	fmt.Printf("Path:%s,FileName:%s\n", path, part.FileName())
	result = handleResult(fileId, part.FileName())
	return
}

func Get(fileId string) (*FileInfo, error) {
	minute, uuid, _ := splitFileId(fileId)
	fmt.Printf("%s,%s\n", uuid, minute)
	path := filepath.Join(conf.Conf.Data.DataDir, minute, uuid)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, ErrNoFilesFound
	}
	fileInfo := files[0]
	filePath := filepath.Join(path, fileInfo.Name())
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return &FileInfo{file, fileInfo.Name()}, nil
}

func handlePath(filename string) (string, string, error) {
	uuid := strings.Replace(uuid.New().String(), "-", "", -1)
	nowMinute := time.Now().Format("200601021504")
	path := filepath.Join(conf.Conf.Data.DataDir, nowMinute, uuid)
	if !utils.FileExists(path) {
		if err := os.MkdirAll(path, 644); err != nil {
			return "", "", err
		}
	}
	return filepath.Join(path, filename), handleFileId(uuid, nowMinute), nil
}

func handleFileId(uuid, minute string) string {
	var fileId []string
	fileId = append(fileId, uuid)
	fileId = append(fileId, utils.NumberToLetter(minute))
	return strings.Join(fileId, "")
}

func splitFileId(fileId string) (minute string, uuid string, err error) {
	if len(fileId) != 44 {
		err = ErrIncorrectFileID
		return
	}
	uuid = fileId[:32]
	minute = utils.LetterToNumber(fileId[32:])
	return
}

func handleResult(fileId, filename string) *Result {
	return &Result{FileName: filename, FileId: fileId}
}
