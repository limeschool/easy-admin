package cert

import (
	"errors"
	"github.com/limeschool/easy-admin/server/config"
	"io"
	"os"
	"sync"
)

type cert struct {
	mu sync.RWMutex
	m  map[string][]byte
}

type Cert interface {
	Get(name string) ([]byte, error)
	GetCert(name string) []byte
}

func New(cs []config.Cert) Cert {
	ct := cert{
		m:  make(map[string][]byte),
		mu: sync.RWMutex{},
	}

	ct.mu.Lock()
	defer ct.mu.Unlock()

	for _, item := range cs {
		file, err := os.Open(item.Path)
		if err != nil {
			panic("cert初始化失败:" + err.Error())
		}
		val, err := io.ReadAll(file)
		if err != nil {
			panic("读取cert证书失败:" + err.Error())
		}
		ct.m[item.Name] = val
	}
	return &ct
}

// Get
//
//	@Description: 获取指定名称的证书，不存在则报错
//	@receiver o
//	@param name
//	@return []byte
//	@return error
func (o *cert) Get(name string) ([]byte, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if o.m[name] == nil {
		return nil, errors.New("not exist cert")
	}
	return o.m[name], nil
}

// GetCert
//
//	@Description: 获取指定名称的证书，不存在则返回nil
//	@receiver o
//	@param name
//	@return []byte
func (o *cert) GetCert(name string) []byte {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.m[name]
}
