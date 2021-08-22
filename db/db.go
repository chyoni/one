package db

import (
	"github.com/chiwon99881/one/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dbName      string = "onecoin.db"
	chainBucket string = "chainBucket"
	blockBucket string = "blockBucket"
	cursor      string = "aCursor"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		db = dbPointer
		utils.HandleErr(err)
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(chainBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blockBucket))
			utils.HandleErr(err)
			return nil
		})
		utils.HandleErr(err)
		return db
	}
	return db
}

func Close() {
	DB().Close()
}

func SaveChainDB(data []byte) {
	DB().Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(chainBucket))
		err := b.Put([]byte(cursor), data)
		utils.HandleErr(err)
		return nil
	})
}

func SaveBlockDB(key string, data []byte) {
	DB().Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blockBucket))
		err := b.Put([]byte(key), data)
		utils.HandleErr(err)
		return nil
	})
}
