package tools

import (
	"feelingliu/modles"
	"time"
)

type Options struct {
	Timeout bool // 是否设置过期时间
}

var defaultOptions = Options{
	Timeout: false,
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	// 初始化默认值
	opt := defaultOptions

	for _, o := range opts {
		o(&opt) // 依次调用opts函数列表中的函数，为服务选项（opt变量）赋值
	}

	return opt
}

func SetTimeout(timeout bool) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

func SetKey(key string, value interface{}, opts ...Option) error {


	options := newOptions(opts...)
	if options.Timeout {
		err := modles.RedisPool.Set(key, value, time.Duration(modles.RedisInfo.CacheTime) * time.Second).Err()
		return err
	}
	err := modles.RedisPool.Set(key, value, time.Duration(0) * time.Second).Err()
	return err
}

func GetKey(key string) (data interface{}, err error) {
	data, err = modles.RedisPool.Get(key).Result()
	return
}

func DelKey(key string) error {
	err := modles.RedisPool.Del(key).Err()
	return err
}

func INCRKey(key string) error {

	err := modles.RedisPool.Incr(key).Err()
	return err
}
