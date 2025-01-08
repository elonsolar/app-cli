package handler

import (
	"app-cli/conf"
	"app-cli/data"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func ErrorResponse(msg string) Response {
	return Response{
		Code: 500,
		Msg:  msg,
		Data: nil,
	}
}

func OkResponse(data any) Response {
	return Response{
		Code: 200,
		Msg:  "ok",
		Data: data,
	}
}

type Handler struct {
	data *data.Data
	cfg  *conf.Config
}

func NewHandler(cfg *conf.Config, data *data.Data) *Handler {
	return &Handler{
		data: data,
		cfg:  cfg,
	}
}

func (h *Handler) CreatePackTask(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(200, ErrorResponse(err.Error()))
		return
	}

	fileList, exist := form.File["file"]
	if !exist || len(fileList) != 1 {
		c.JSON(200, ErrorResponse("上传一个文件"))
		return
	}

	projectName, _ := getFormValue(form.Value, "projectName")
	version, _ := getFormValue(form.Value, "version")
	description, _ := getFormValue(form.Value, "description")
	slog.Info("version", slog.String("version", version))

	// todo save to local
	var fileName = fileList[0].Filename
	if !strings.HasSuffix(fileName, ".zip") {
		c.JSON(200, ErrorResponse("只支持zip"))
		return
	}

	err = c.SaveUploadedFile(fileList[0], filepath.Join(h.cfg.UploadDir, fileName))
	if err != nil {
		c.JSON(200, ErrorResponse(err.Error()))
		return
	}

	executor, err := NewPackExecutor(h.cfg, projectName, h.data)
	if err != nil {
		c.JSON(200, ErrorResponse(err.Error()))
		return
	}

	go executor.execute(fileName, description)

}

func (h *Handler) ListPackTask(c *gin.Context) {
	tasks, err := h.data.FindPackTask(1, 10)
	if err != nil {
		c.JSON(200, ErrorResponse(err.Error()))
		return
	}
	c.JSON(200, OkResponse(tasks))
}

func (h *Handler) Download(c *gin.Context) {
	var param struct {
		Id string `form:"id" json:"id"`
	}
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(200, ErrorResponse(err.Error()))
		return
	}

	var fileName = fmt.Sprintf("app-release-%s.apk", param.Id)

	c.Header("fileName", fileName)
	c.Header("msg", "文件下载成功")
	c.File(filepath.Join(h.cfg.DownloadDir, fileName))
	// var fileName = filepath.Join("./", fmt.Sprintf("app-release-%s.apk", param.Id))

	// fmt.Println(fileName == "app-release-1864929139635523584.apk")
	// c.Header("fileName", fileName)
	// c.Header("msg", "文件下载成功")
	// // c.File(filepath.Join(h.cfg.DownloadDir, fileName))
	// c.File(filepath.Join("./", fileName))
}

func getFormValue(value map[string][]string, key string) (string, bool) {

	vl, exist := value[key]
	if !exist {
		return "", false
	}

	if len(vl) == 0 {
		return "", false
	}
	return vl[0], true
}
