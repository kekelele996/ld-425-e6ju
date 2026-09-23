package config

// UploadConfig 文件上传配置。
type UploadConfig struct {
	Dir          string `env:"UPLOAD_DIR" envDefault:"/tmp/home-renovation/uploads"`
	MaxSizeMB    int64  `env:"MAX_UPLOAD_SIZE_MB" envDefault:"20"`
	PublicPrefix string `env:"UPLOAD_PUBLIC_PREFIX" envDefault:"/uploads"`
}
