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
	ErrIncorrectFileID = errors.New("invalid fileId")
	ErrNoFilesFound    = errors.New("no files were found")
)

type Result struct {
	Ok       bool   `json:"ok"`
	Msg      string `json:"msg"`
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
	fmt.Printf("UPLOAD:Path:%s\n", path)
	result = handleResult(fileId, part.FileName())
	return
}

func Get(fileId string) (*FileInfo, error) {
	path, err := getPath(fileId)
	if err != nil {
		return nil, err
	}
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
	fmt.Printf("GET:%s\n", filePath)
	return &FileInfo{file, fileInfo.Name()}, nil
}

func Delete(fileId string) *Result {
	res := &Result{Ok: false}
	path, err := getPath(fileId)
	if err != nil {
		res.Msg = err.Error()
		return res
	}
	//abs, err := filepath.Abs(filepath.Dir(fileInfo.File.Name()))

	fmt.Printf("DELETE:%s\n", path)
	err = os.RemoveAll(path)
	if err != nil {
		res.Msg = err.Error()
		return res
	}
	res.Ok = true
	return res
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

func getPath(fileId string) (string, error) {
	if len(fileId) != 44 {
		return "", ErrIncorrectFileID
	}
	uuid := fileId[:32]
	minute := utils.LetterToNumber(fileId[32:])
	return filepath.Join(conf.Conf.Data.DataDir, minute, uuid), nil
}

func handleResult(fileId, filename string) *Result {
	return &Result{Ok: true, FileName: filename, FileId: fileId}
}
