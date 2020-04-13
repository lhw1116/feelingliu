package modles

import(
	"time"
)

type App struct {
	TimeFormat     string `json:"time_format"`
	JwtSecret      string `json:"jwt_secret"`
	TokenTimeout   int64  `json:"token_timeout"`
	StaticBasePath string `json:"static_base_path"`
	UploadBasePath string `json:"upload_base_path"`
	ImageRelPath   string `json:"image_rel_path"`
	AvatarRelPath  string `json:"avatar_rel_path"`
	ApiBaseUrl     string `json:"api_base_url"`
}

type Server struct {
	RunMode      string        `json:"run_mode"`
	ServerAddr   string        `json:"server_addr"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

type DataBase struct {
	Mode     string `json:"mode"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

type Redis struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	Password  string `json:"password"`
	DB        int    `json:"db"`
	CacheTime int    `json:"cache_time"`
}
