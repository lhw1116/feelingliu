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
	conn := modles.RedisPool
	defer func() {
		if e := conn.Close(); e != nil {
			return
		}
	}()

	options := newOptions(opts...)
	if options.Timeout {
		_, err := conn.Set(key, value, time.Duration(modles.RedisInfo.CacheTime)).Result()
		return err
	}
	_, err := conn.Set(key, value, 0).Result()
	return err
}

func GetKey(key string) (data interface{}, err error) {
	conn := modles.RedisPool
	defer func() {
		if err := conn.Close(); err != nil {
			return
		}
	}()

	data, err = conn.Get(key).Result()
	return
}

func DelKey(key string) error {
	conn := modles.RedisPool
	defer func() {
		if e := conn.Close(); e != nil {
			return
		}
	}()

	_, err := conn.Del(key).Result()
	return err
}

func INCRKey(key string) error {
	conn := modles.RedisPool
	defer func() {
		if e := conn.Close(); e != nil {
			return
		}
	}()

	_, err := conn.Incr(key).Result()
	return err
}
