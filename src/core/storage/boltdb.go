package storage

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"encoding/binary"
	"github.com/boltdb/bolt"
	"time"
)

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
		_, err = tx.CreateBucketIfNotExists([]byte(base.BoltDBConfig.Boltdb.TableName))
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

func (s *BoltDBStorageHandler) MarkFile(filePath string, modTime time.Time) bool {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		uploadFileMarkBucket := tx.Bucket([]byte(base.BoltDBConfig.Boltdb.TableName))
		if uploadFileMarkBucket != nil {
			value := make([]byte, 8)
			binary.LittleEndian.PutUint64(value, uint64(modTime.UnixMilli()))
			return uploadFileMarkBucket.Put([]byte(filePath), value)
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

func (s *BoltDBStorageHandler) CheckFile(filePath string, modTime time.Time) bool {
	value := []byte("")
	err := s.DB.View(func(tx *bolt.Tx) error {
		uploadFileMarkBucket := tx.Bucket([]byte(base.BoltDBConfig.Boltdb.TableName))
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
	// 如果key不存在，说明该文件未被标记，需要上传
	if value == nil {
		return true
	}
	i := int64(binary.LittleEndian.Uint64(value))
	// 如果本地文件更新时间大于记录中的文件更新时间，说明该文件需要同步
	if modTime.UnixMilli() > i {
		return true
	}
	return false
}
