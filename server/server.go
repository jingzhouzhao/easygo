package server

import (
	"context"
	"easygo/conf"
	"easygo/storage"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrNotSupportedMethod = "Request method '%s' not supported"
	ErrMissingParameter   = "Missing parameter '%s'"
)

var (
	ErrNotFound = errors.New("Not Found")
)

var (
	server *http.Server
)

func router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello EasyGo!"))
		return
	}
	if strings.HasPrefix(path, "/files/") {
		handleFile(w, r)
		return
	}
	respErr(w, http.StatusNotFound, ErrNotFound)
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		post(w, r)
	case http.MethodGet:
		get(w, r)
	case http.MethodDelete:
		delete(w, r)
	default:
		respErr(w, http.StatusBadRequest, errors.New(fmt.Sprintf(ErrNotSupportedMethod, r.Method)))
		return
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		respErr(w, http.StatusBadRequest, err)
		return
	}

	var results []*storage.Result
	for {
		part, err := mr.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}
			respErr(w, http.StatusInternalServerError, err)
			return
		}

		filename := part.FileName()
		if filename == "" {
			continue
		}

		result, err := storage.Save(part)
		if err != nil {
			respErr(w, http.StatusInternalServerError, err)
			return
		}
		results = append(results, result)

	}
	resp(w, results)
}

func delete(w http.ResponseWriter, r *http.Request) {
	fileId := getFileId(r.URL.Path)
	if fileId == "" {
		respErr(w, http.StatusBadRequest, errors.New(fmt.Sprintf(ErrMissingParameter, "fileId")))
		return
	}
	res := storage.Delete(fileId)
	resp(w, res)
}
func get(w http.ResponseWriter, r *http.Request) {
	fileId := getFileId(r.URL.Path)
	fmt.Println("fileid:", fileId)
	if fileId == "" {
		respErr(w, http.StatusBadRequest, errors.New(fmt.Sprintf(ErrMissingParameter, "fileId")))
		return
	}
	fileInfo, err := storage.Get(fileId)
	if err != nil {
		respErr(w, http.StatusInternalServerError, err)
		return
	}
	defer fileInfo.File.Close()
	w.Header().Set("Pragma", "No-cache")
	w.Header().Set("Cache-Control", "No-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-disposition", "attachment;filename="+fileInfo.FileName)
	w.WriteHeader(http.StatusOK)
	buf := make([]byte, 4096)
	for {
		n, err := fileInfo.File.Read(buf)
		if err != nil {
			if err == io.EOF {
				w.Write(buf[0:n])
				break
			}
			respErr(w, http.StatusInternalServerError, err)
			return
		}
		w.Write(buf[0:n])
	}

}

func getFileId(path string) string {
	idx := strings.LastIndex(path, "/")
	return path[idx+1:]
}

func resp(w http.ResponseWriter, results interface{}) {
	rs, err := json.Marshal(results)
	if err != nil {
		respErr(w, http.StatusInternalServerError, err)
		return
	}
	w.Write(rs)
}

func respErr(w http.ResponseWriter, errCode int, err error) {
	fmt.Errorf("errorCode:%d,error:%s", errCode, err.Error())
	w.WriteHeader(errCode)
	w.Write([]byte(err.Error()))
}

func Start() {
	http.HandleFunc("/", router)
	server = &http.Server{
		Addr:    fmt.Sprintf("%s%d", ":", conf.Conf.Http.Ports.HttpPort),
		Handler: http.DefaultServeMux,
		//WriteTimeout: time.Second * time.Duration(conf.Conf.Http.Timeout),
	}
	go func() {
		fmt.Printf("EasyGo server listening at: %d\n", conf.Conf.Http.Ports.HttpPort)
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				fmt.Println("Http server closed")
				return
			}
			panic(err)
		}
	}()
}

func Listen() {

}

func Close() error {
	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		fmt.Println("EasyGo server shutdown")
		return server.Shutdown(ctx)
	}
	return nil
}
