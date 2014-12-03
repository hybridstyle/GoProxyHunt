package proxyhunt

type IP struct {
	Addr string	`json:"addr"`
	Ip string	`json:"ip"`
	Port int	`json:"port"`
	Ctime int64	`json:"-"`
	Utime int64	`json:"-"`
}
