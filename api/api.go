package api

import (
	"encoding/json"
	"lumen/go-siakad/tools"
	"net/http"
)

type StudentsListParam struct {
	NISN    string `schema:"NISN" gorm:"column:NISN"`
	Jurusan string `schema:"jurusan" gorm:"column:kd_mata_diklat"`
	Limit   int    `schema:"limit" gorm:"-"`
	Offset  int    `schema:"offset" gorm:"-"`
}

type StudentsResponse struct {
	Nisn     uint32 `json:"nisn"`
	Jurusan  string `json:"jurusan"`
	Nama     string `json:"nama"`
	Alamat   string `json:"alamat"`
	TglLahir string `json:"tanggal_lahir"`
	FileFoto string `json:"file_foto"`
}

type StudentsListResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"info"`
	Data    []tools.StudentDetails `json:"data"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"info"`
}

func writeError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(resp)
}

var (
	CustomErrorHandler = func(w http.ResponseWriter, err error, code int) {
		writeError(w, err.Error(), code)
	}
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An Unexpected Error Occured.", http.StatusInternalServerError)
	}
)
