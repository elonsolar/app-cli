package conf

const (
	UploadDir   = "D:\\upload"
	DownloadDir = "D:\\download"

	UnpackDir = "D:\\projects\\Android-SDK@4.29.82201_20241008\\UniPlugin-Hello-AS\\app\\src\\main\\assets\\apps"

	AppId    = "__UNI__6B07A5B"
	BuildDir = "D:\\projects\\Android-SDK@4.29.82201_20241008\\UniPlugin-Hello-AS"

	OutputDir = "D:\\projects\\Android-SDK@4.29.82201_20241008\\UniPlugin-Hello-AS\\app\\build\\outputs\\apk\\release"
)

type Config struct {
	Port        int                `json:"port"`
	Dsn         string             `json:"dsn"`
	UploadDir   string             `json:"upload_dir"`
	DownloadDir string             `json:"download_dir"`
	Project     map[string]Project `json:"project"`
}

type Project struct {
	Name      string `json:"name"`
	UnpackDir string `json:"unpack_dir"`
	BuildDir  string `json:"build_dir"`
	OutputDir string `json:"output_dir"`
}
