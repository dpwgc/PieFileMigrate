package storage

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"PieFileMigrate/src/util"
	"github.com/boltdb/bolt"
)

const tableName = "upload_file_mark"

func NewBoltDBStorageHandler() Handler {
	base.InitBoltDBConfig()
	db, err := initBoltDBStorage()
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	base.LogHandler.Println(constant.LogInfoTag, "本地数据库(BoltDB)加载成功")
	return &BoltDBStorageHandler{
		DB: db,
	}
}

type BoltDBStorageHandler struct {
	DB *bolt.DB
}

func initBoltDBStorage() (*bolt.DB, error) {
	// 创建或者打开数据库
	db, err := bolt.Open(base.BoltDBConfig.Boltdb.Db, 0600, nil)
	if err != nil {
		return nil, err
	}

	// 创建表
	err = db.Update(func(tx *bolt.Tx) error {
		// 创建upload_file_mark表
		_, err = tx.CreateBucketIfNotExists([]byte(tableName))
		if err != nil {
			return err
		}
		return nil
	})
	//更新失败
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (s *BoltDBStorageHandler) MarkFile(filePath string) bool {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		uploadFileMarkBucket := tx.Bucket([]byte(tableName))
		if uploadFileMarkBucket != nil {
			return uploadFileMarkBucket.Put([]byte(filePath), []byte(util.GetLocalDateTime()))
		}
		return nil
	})
	//写入数据库失败
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, err)
		return false
	}
	return true
}

func (s *BoltDBStorageHandler) CheckFile(filePath string) bool {
	value := []byte("")
	err := s.DB.View(func(tx *bolt.Tx) error {
		uploadFileMarkBucket := tx.Bucket([]byte(tableName))
		if uploadFileMarkBucket != nil {
			value = uploadFileMarkBucket.Get([]byte(filePath))
		}
		return nil
	})
	//查询数据库失败
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, err)
		return false
	}
	return value != nil
}
