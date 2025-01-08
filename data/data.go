package data

import (
	"app-cli/conf"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type PackTaskStatus string

const (
	PackTaskStatusPending      PackTaskStatus = "就绪"
	PackTaskStatusUnpacked     PackTaskStatus = "解压成功"
	PackTaskStatusFailed       PackTaskStatus = "打包失败"
	PackTaskStatusBuildSuccess PackTaskStatus = "编译成功"
	PackTaskStatusSuccess      PackTaskStatus = "打包成功"
)

type PackTask struct {
	ID            string `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	ProjectName   string         `json:"projectName"`
	Description   string         `json:"description"`
	Status        string         `json:"status"`
	DownloadCount int            `json:"downloadCount"`
}

func (PackTask) TableName() string {
	return "pack_task"
}

type Data struct {
	cfg *conf.Config
	db  *gorm.DB
}

func NewData(cfg *conf.Config, db *gorm.DB) *Data {
	return &Data{
		cfg: cfg,
		db:  db,
	}
}

func (d *Data) SavePackTask(t *PackTask) error {
	err := d.db.Create(t).Error
	if err != nil {
		return fmt.Errorf("create pack task failed ,err :%w ", err)
	}
	return nil
}

func (d *Data) UpdateTaskStatus(id string, status PackTaskStatus) error {

	return d.db.Model(&PackTask{}).Where("id=?", id).UpdateColumn("status", status).Error
}

func (d *Data) FindPackTask(pageNo, PageSize int) ([]PackTask, error) {
	slog.Info("find pack task", slog.Int("pageNo", pageNo), slog.Int("PageSize", PageSize))

	var tasks = make([]PackTask, 0, 10)
	err := d.db.Model(PackTask{}).Offset((pageNo - 1) * PageSize).Limit(PageSize).Order("id desc").Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func NewGormDB(cfg *conf.Config) *gorm.DB {

	mysqlConfig := mysql.Config{
		DSN: "root:mysql_5Y57Pi@tcp(192.168.1.246:3306)/pack?charset=utf8mb4&parseTime=True&loc=Local",
	}
	db, err := gorm.Open(mysql.New(mysqlConfig))
	if err != nil {
		panic(err)
	}
	return db.Debug()
}
