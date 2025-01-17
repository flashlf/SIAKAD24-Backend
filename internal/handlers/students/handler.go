package students

import (
	"encoding/json"
	"errors"
	"lumen/go-siakad/api"
	"lumen/go-siakad/tools"
	"net/http"
	"strconv"

	gorillaSchema "github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Fungsi umum untuk mengambil data siswa
func GetStudents(w http.ResponseWriter, r *http.Request, queryParams api.StudentsListParam, limit int, offset int) ([]*tools.StudentDetails, error) {
	// Inisialisasi decoder dan parsing query params

	var decoder = gorillaSchema.NewDecoder()
	if err := decoder.Decode(&queryParams, r.URL.Query()); err != nil {
		log.Error(err)
		return nil, err
	}
	// Inisialisasi database
	database, err := tools.Init()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Menentukan query berdasarkan apakah kita ingin mencari by ID atau list dengan paginasi
	var students []*tools.StudentDetails

	if queryParams.NISN != "" {
		// Jika NISN diberikan, cari berdasarkan NISN
		sanitizeNISN, err := strconv.ParseUint(queryParams.NISN, 10, 32)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		student, err := database.GetStudentByNISN(sanitizeNISN)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	} else {
		// Jika tidak ada NISN, ambil list dengan paginasi
		students, err = database.GetStudents(queryParams, limit, offset)
		if err != nil {
			return nil, err
		}
	}
	return students, nil
}

func LoadByID(w http.ResponseWriter, r *http.Request) {
	var params = api.StudentsListParam{}

	students, err := GetStudents(w, r, params, 0, 0)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warn(tools.NotFoundError)
		api.CustomErrorHandler(w, err, http.StatusNotFound)
		return
	}

	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}

	// Jika students adalah slice of pointers ([]*tools.StudentDetails)
	var studentValues []tools.StudentDetails

	// Menyalin data dari slice of pointers ke slice of values
	for _, studentPtr := range students {
		studentValues = append(studentValues, *studentPtr) // Dereference pointer dan append ke slice
	}

	responseMessage := "Data Not Found"

	if len(studentValues) > 0 {
		responseMessage = "Data Found"
	}

	var response = api.StudentsListResponse{
		Code:    http.StatusOK,
		Message: responseMessage,
		Data:    studentValues,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func LoadList(w http.ResponseWriter, r *http.Request) {
	queryParams := api.StudentsListParam{}

	// Misalnya kita ambil query params `limit` dan `offset` dari URL
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	students, err := GetStudents(w, r, queryParams, limit, offset)
	if err != nil {
		api.InternalErrorHandler(w, err)
		return
	}
	var studentList []tools.StudentDetails

	for _, studentPtr := range students { // Iterasi melalui slice pointer
		studentList = append(studentList, tools.StudentDetails{
			Nisn:     studentPtr.Nisn,
			Jurusan:  studentPtr.Jurusan,
			Nama:     studentPtr.Nama,
			Alamat:   studentPtr.Alamat,
			TglLahir: studentPtr.TglLahir,
			FileFoto: studentPtr.FileFoto,
		})
	}

	responseMessage := "Data Not Found"

	if len(studentList) > 0 {
		responseMessage = "Data Found"
	}

	response := api.StudentsListResponse{
		Code:    http.StatusOK,
		Message: responseMessage,
		Data:    studentList,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
