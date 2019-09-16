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
	"time"
)

var (
	ErrNotSupportedMethod = "Request method '%s' not supported"
	ErrMissingParameter   = "Missing parameter '%s'"
)

var (
	server *http.Server
)

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
	w.Write([]byte("delete"))
}
func get(w http.ResponseWriter, r *http.Request) {
	fileId := r.URL.Query().Get("fileId")
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
	http.HandleFunc("/file", handleFile)
	server = &http.Server{
		Addr:    fmt.Sprintf("%s%d", ":", conf.Conf.Http.Ports.HttpPort),
		Handler: http.DefaultServeMux,
		//WriteTimeout: time.Second * time.Duration(conf.Conf.Http.Timeout),
	}
	go func() {
		fmt.Printf("Http server listening at: %d\n", conf.Conf.Http.Ports.HttpPort)
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
		fmt.Println("Http server shutdown")
		return server.Shutdown(ctx)
	}
	return nil
}
