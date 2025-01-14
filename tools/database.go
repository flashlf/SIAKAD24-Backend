package tools

import (
	"errors"
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type StudentDetails struct {
	Nisn     uint32 `json:"nisn" gorm:"primaryKey;column:NISN"`
	Jurusan  string `json:"jurusan" gorm:"column:kd_mata_diklat"`
	Nama     string `json:"nama" gorm:"column:Nama_siswa"`
	Alamat   string `json:"alamat" gorm:"column:Alamat_siswa"`
	TglLahir string `json:"tgl_lahir" gorm:"column:Tgl_lahir"`
	FileFoto string `json:"file_foto" gorm:"column:Foto_siswa"`
}

func (StudentDetails) TableName() string {
	return "siswa"
}

var NotFoundError = errors.New("Data Not Found")

type Database struct {
	DB *gorm.DB
}

type DatabaseInterface interface {
	GetStudentByNISN(nisn int) (*StudentDetails, error)
	GetStudents(limit int, offset int) (*[]StudentDetails, error)
	SetupDatabase() error
}

func Init() (*Database, error) {
	dsn := "root@/sekolah?timeout=90s"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (d *Database) GetStudentByNISN(nisn uint64) (*StudentDetails, error) {
	var data = StudentDetails{}
	err := d.DB.Where("NISN = ?", nisn).First(&data).Order("NISN").Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (d *Database) GetStudents(params interface{}, limit int, offset int) ([]*StudentDetails, error) {
	var listStudents = []*StudentDetails{}
	fmt.Printf("Struct: %+v\n", params)
	t := reflect.TypeOf(params)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("Field: %s, Tag: %s\n", field.Name, field.Tag.Get("gorm"))
	}
	query := d.DB.Where(params)
	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&listStudents).Order("NISN").Error; err != nil {
		return nil, err
	}

	return listStudents, nil
}
