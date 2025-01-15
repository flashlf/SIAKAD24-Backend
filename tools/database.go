package tools

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type StudentDetails struct {
	Nisn     uint32 `json:"nisn" gorm:"primaryKey;column:NISN;<-:create"`
	Jurusan  string `json:"jurusan" gorm:"column:kd_mata_diklat"`
	Nama     string `json:"nama" gorm:"column:Nama_siswa"`
	Alamat   string `json:"alamat" gorm:"column:Alamat_siswa"`
	TglLahir string `json:"tgl_lahir" gorm:"column:Tgl_lahir"`
	FileFoto string `json:"file_foto" gorm:"column:Foto_siswa"`
}

func (StudentDetails) TableName() string {
	return "siswa"
}

type LecturerDetails struct {
	KodeGuru       string `json:"kd_guru" gorm:"primaryKey;column:kd_guru;<-:create"`
	KodeKompetensi string `json:"kd_kompetensi" gorm:"column:Kode_KK"`
	Nama           string `json:"nama" gorm:"column:nm_guru"`
	NIP            string `json:"nip" gorm:"column:NIP"`
	Alamat         string `json:"alamat" gorm:"alamat_guru"`
	Telp           string `json:"telp" gorm:"telp_guru"`
}

func (LecturerDetails) TableName() string {
	return "guru"
}

var NotFoundError = errors.New("Data Not Found")

type Database struct {
	DB *gorm.DB
}

type DatabaseInterface interface {
	GetStudentByNISN(nisn int) (*StudentDetails, error)
	GetStudents(paramInterface interface{}, limit int, offset int) (*[]StudentDetails, error)
	Init() error
}

func Init() (*Database, error) {
	dsn := "root@/sekolah?timeout=90s"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("failed to initialize database connection")
	}

	return &Database{DB: db}, nil
}

func (d *Database) GetStudentByNISN(nisn uint64) (*StudentDetails, error) {
	var data = StudentDetails{}

	sqlDB, err := d.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to access database")
	}

	defer sqlDB.Close()

	err = d.DB.Where("NISN = ?", nisn).First(&data).Order("NISN").Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (d *Database) GetStudents(params interface{}, limit int, offset int) ([]*StudentDetails, error) {
	var listStudents = []*StudentDetails{}

	sqlDB, err := d.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to access database")
	}

	defer sqlDB.Close()

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

func GetList[T any, U any](d *gorm.DB, params U, limit int, offset int) ([]*T, error) {
	var list = []*T{}

	sqlDB, err := d.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to access database")
	}

	defer sqlDB.Close()

	query := d.Where(params)
	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}
