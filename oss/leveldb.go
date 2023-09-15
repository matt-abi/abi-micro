package oss

import (
	"fmt"
	"time"

	"github.com/golang/leveldb"
	"github.com/golang/leveldb/db"
)

type LevelDBConfig struct {
	Dir string `json:"dir"`
}

type levelDB struct {
	db *leveldb.DB
}

func NewLevelDB(cfg *LevelDBConfig) (OSS, error) {

	db, err := leveldb.Open(cfg.Dir, nil)
	if err != nil {
		return nil, err
	}

	return &levelDB{db: db}, nil
}

func (S *levelDB) Get(key string) ([]byte, error) {
	return S.db.Get([]byte(key), nil)
}

func (S *levelDB) GetURL(key string) string {
	return ""
}

func (S *levelDB) GetSignURL(key string, expires time.Duration) (string, error) {
	return "", fmt.Errorf("GetSignURL not supported")
}

func (S *levelDB) Put(key string, data []byte, header map[string]string) error {
	return S.db.Set([]byte(key), data, &db.WriteOptions{Sync: true})
}

func (S *levelDB) PutSignURL(key string, expires time.Duration, header map[string]string) (string, error) {
	return "", fmt.Errorf("PutSignURL not supported")
}

func (S *levelDB) PostSignURL(key string, expires time.Duration, maxSize int64, header map[string]string) (string, map[string]string, error) {
	return "", nil, fmt.Errorf("PostSignURL not supported")
}

func (S *levelDB) Del(key string) error {
	return S.db.Delete([]byte(key), nil)
}

func (S *levelDB) Has(key string) (bool, error) {
	i := S.db.Find([]byte(key), nil)
	defer i.Close()
	return i.Next(), nil
}

func (S *levelDB) Recycle() {
	if S.db != nil {
		S.db.Close()
		S.db = nil
	}
}
