package lecturers

import (
	"encoding/json"
	"errors"
	"lumen/go-siakad/api"
	"lumen/go-siakad/tools"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gorillaSchema "github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LecturerDetails struct {
	KodeGuru       string `json:"-" gorm:"primaryKey;column:kd_guru;<-:create"`
	KodeKompetensi string `json:"-" gorm:"column:Kode_KK"`
	NIP            string `json:"nip" gorm:"column:NIP"`
	Nama           string `json:"nama" gorm:"column:nm_guru"`
	Alamat         string `json:"alamat" gorm:"column:alamat_guru"`
	Telp           string `json:"telp" gorm:"column:telp_guru"`
}

func (LecturerDetails) TableName() string {
	return "guru"
}

func GetLecturers(c *gin.Context) {
	decoder := gorillaSchema.NewDecoder()
	interfaceParam := api.TeacherParam{}
	if err := decoder.Decode(&interfaceParam, c.Request.URL.Query()); err != nil {
		log.Error(err)
		api.CustomErrorHandler(c.Writer, err, http.StatusNotFound)
		return
	}

	db, err := tools.Init()
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(c.Writer, err)
		return
	}

	var lecturer []*LecturerDetails

	// Misalnya kita ambil query params `limit` dan `offset` dari URL
	limit, err := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}
	offset, err := strconv.Atoi(c.Request.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	lecturer, err = tools.GetList[LecturerDetails, api.TeacherParam](db.DB, interfaceParam, limit, offset)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warn(tools.NotFoundError)
		api.CustomErrorHandler(c.Writer, err, http.StatusNotFound)
		return
	}

	if err != nil {
		api.InternalErrorHandler(c.Writer, err)
		return
	}

	responseMessage := "Data Not Found"

	if len(lecturer) > 0 {
		responseMessage = "Data Found"
	}

	var response = api.GenericListCountResponse[*LecturerDetails]{
		Code:    http.StatusOK,
		Message: responseMessage,
		Records: len(lecturer),
		Data:    lecturer,
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.Writer).Encode(response)
}
