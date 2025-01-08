package handler

import (
	"app-cli/conf"
	"app-cli/data"
	"archive/zip"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"fmt"

	"github.com/go-pay/util/snowflake"
)

type PackExecutor struct {
	projectName string
	project     conf.Project
	cfg         *conf.Config
	data        *data.Data
}

func NewPackExecutor(cfg *conf.Config, projectName string, data *data.Data) (*PackExecutor, error) {

	project, exist := cfg.Project[projectName]
	if !exist {
		return nil, fmt.Errorf("项目不存在 %s", projectName)
	}

	return &PackExecutor{
		projectName: projectName,
		project:     project,
		cfg:         cfg,
		data:        data,
	}, nil
}

func (e *PackExecutor) execute(srcFile string, description string) {

	node, err := snowflake.NewNode(1)
	if err != nil {
		slog.Error("new node failed", "err", err)
		return
	}

	// save record
	taskId := node.Generate().String()

	go e.data.SavePackTask(&data.PackTask{
		ID:          taskId,
		ProjectName: e.projectName,
		Description: description,
		Status:      string(data.PackTaskStatusPending),
	})

	err = e.unpackSrcFileToDest(srcFile)
	if err != nil {
		slog.Error("save record failed", "err", err)
		return
	}

	go e.data.UpdateTaskStatus(taskId, data.PackTaskStatusUnpacked)

	err = e.buildGradleRelease()
	if err != nil {
		slog.Error("build gradle failed", "err", err)
		return
	}

	go e.data.UpdateTaskStatus(taskId, data.PackTaskStatusBuildSuccess)

	// clean src file
	os.Remove(srcFile)

	err = e.moveReleaseForDownload(taskId)
	if err != nil {
		slog.Error("move release file failed", "err", err)
		return
	}
	go e.data.UpdateTaskStatus(taskId, data.PackTaskStatusSuccess)
}

func (e *PackExecutor) unpackSrcFileToDest(sourceFile string) error {
	archive, err := zip.OpenReader(filepath.Join(e.cfg.UploadDir, sourceFile))
	if err != nil {
		return fmt.Errorf("open zip file failed err : %w", err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		fpath := filepath.Join(e.project.UnpackDir, f.Name)

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(fpath, f.Mode())
			if err != nil {
				return fmt.Errorf("create dir failed err : %w", err)
			}
			continue
		}

		// 创建文件所在的目录
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("create dir failed err: %w", err)
		}

		// 解压文件
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open zip file failed err : %w", err)
		}
		defer rc.Close()

		fdest, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("create target file failed err: %w", err)
		}
		defer fdest.Close()

		_, err = io.Copy(fdest, rc)
		if err != nil {
			return fmt.Errorf("copy file failed err: %w", err)
		}
	}
	return nil
}

func (e *PackExecutor) buildGradleRelease() error {
	err := os.Chdir(e.project.BuildDir)
	if err != nil {
		return fmt.Errorf("change dir failed err : %w", err)
	}

	cmdName := "cmd"
	args := []string{"/C", "gradlew.bat assembleRelease"}

	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 运行命令
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("command failed: %v", err)
	}
	return nil
}

func (e *PackExecutor) moveReleaseForDownload(taskId string) error {

	src, err := os.Open(filepath.Join(e.project.OutputDir, "app-release.apk"))
	if err != nil {
		return fmt.Errorf("open src file failed, err :%w ", err)
	}

	dst, err := os.Create(filepath.Join(e.cfg.DownloadDir, fmt.Sprintf("app-release-%s.apk", taskId)))
	if err != nil {
		return fmt.Errorf("create dst file failed, err :%w ", err)
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("copy file failed, err :%w ", err)
	}

	return nil
}
