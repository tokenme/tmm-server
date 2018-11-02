package common

type Config struct {
	Domain               string            `default:"tmm.tokenmama.io"`
	AppName              string            `default:"tmm"`
	BaseUrl              string            `default:"https://tmm.tokenmama.io"`
	CDNUrl               string            `default:"https://cdn.tmm.io/"`
	QRCodeUrl            string            `default:"qr.tmm.io"`
	Port                 int               `default:"8008"`
	Geth                 string            `default:"geth.xibao100.com"`
	GethWSS              string            `required:"true"`
	ShareUrl             string            `required:"true"`
	Template             string            `required:"true"`
	LogPath              string            `required:"true"`
	TokenProfilePath     string            `required:"true"`
	PhonedataPath        string            `required:"true"`
	TokenSalt            string            `required:"true"`
	LinkSalt             string            `required:"true"`
	TMMAgentWallet       WalletConfig      `required:"true"`
	TMMPoolWallet        WalletConfig      `required:"true"`
	TMMTokenAddress      string            `required:"true"`
	TMMEscrowAddress     string            `required:"true"`
	MinTMMExchange       uint              `required:"true"`
	MinTMMRedeem         uint              `required:"true"`
	DefaultAppTaskTS     int64             `required:"true"`
	DefaultShareTaskTS   int64             `required:"true"`
	DefaultDeviceBalance uint64            `required:"true"`
	SentryDSN            string            `required:"true"`
	MySQL                MySQLConfig       `required:"true"`
	Redis                RedisConfig       `required:"true"`
	SQS                  SQSConfig         `required:"true"`
	Qiniu                QiniuConfig       `required:"true"`
	Aliyun               AliyunConfig      `required:"true"`
	ReCaptcha            ReCaptchaConfig   `required:"true"`
	Wechat               WechatConfig      `required:"true"`
	IOSXinge             XingeConfig       `required:"true"`
	AndroidXinge         XingeConfig       `required:"true"`
	TwilioToken          string            `required:"true"`
	EthplorerAPIKey      string            `required:"true"`
	EtherscanAPIKey      string            `required:"true"`
	Slack                SlackConfig       `required:"true"`
	CoinbaseAPI          CoinbaseAPIConfig `required:"true"`
	GeoIP                string            `required:"true"`
	ProxyApiKey          string            `required:"true"`
	GrowthRate           float64           `required:"true"`
	MinGrowthTS          int               `required:"true"`
	InviteBonus          uint              `required:"true"`
	InviterBonus         uint              `required:"true"`
	InviteBonusRate      float64           `required:"true"`
	MaxBindDevice        int               `required:"true"`
	Contact              ContactConfig     `required:"true"`
	AppReleaseVersion    AppVersion        `required:"true"`
	Debug                bool              `default:"false"`
	EnableWeb            bool              `default:"false"`
	EnableGC             bool              `default:"false"`
	EnableTx             bool              `default:"false"`
	EnableOrderBook      bool              `default:"false"`
}

type MySQLConfig struct {
	Host   string `required:"true"`
	User   string `required:"true"`
	Passwd string `required:"true"`
	DB     string `default:"tokenme"`
}

type RedisConfig struct {
	Master string `required:"true"`
	Slave  string
}

type WalletConfig struct {
	Salt string `required:"true"`
	Data string `required:"true"`
	Key  string `required:"true"`
}

type SQSConfig struct {
	Region    string `default:"ap-northeast-1"`
	AccountId string `required:"true"`
	AK        string `required:"true"`
	Secret    string `required:"true"`
	TxQueue   string `default:"ucoin-tx"`
	GasQueue  string `default:"ucoin-gas"`
	Token     string `default:""`
}

type QiniuConfig struct {
	AK         string `required:"true"`
	Secret     string `required:"true"`
	Bucket     string `required:"true"`
	AvatarPath string `required:"true"`
	LogoPath   string `required:"true"`
	ImagePath  string `required:"true"`
	Domain     string `required:"true"`
	Pipeline   string `required:"true"`
	NotifyURL  string `required:"true"`
}

type CoinbaseAPIConfig struct {
	Key    string `required:"true"`
	Secret string `required:"true"`
}

type ContactConfig struct {
	Telegram string `required:"true"`
	Wechat   string `required:"true"`
	Website  string `required:"true"`
}

type AliyunConfig struct {
	RegionId string `required:"true"`
	AK       string `required:"true"`
	AS       string `required:"true"`
}

type ReCaptchaConfig struct {
	Key      string `required:"true"`
	Secret   string `required:"true"`
	Hostname string `required:"true"`
}

type WechatConfig struct {
	AppId   string `required:"true"`
	MchId   string `required:"true"`
	Key     string `required:"true"`
	CertCrt string `required:"true"`
	CertKey string `required:"true"`
}

type SlackConfig struct {
	Token           string `required:"true"`
	FeedbackChannel string `required:"true"`
	OpsChannel      string `required:"true"`
	CaptchaChannel  string `required:"true"`
}

type AppVersion struct {
	IOS     string `required:"true"`
	Android string `required:"true"`
}

type XingeConfig struct {
	AppId     string `required:"true"`
	SecretKey string `required:"true"`
}
