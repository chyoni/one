package db

import (
	"github.com/chiwon99881/one/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dbName        string = "onecoin.db"
	chainBucket   string = "chainBucket"
	blockBucket   string = "blockBucket"
	mempoolBucket string = "mempoolBucket"
	mempoolData   string = "mempoolData"
	cursor        string = "aCursor"
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
			_, err = t.CreateBucketIfNotExists([]byte(mempoolBucket))
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

func GetExistChain() []byte {
	var newestHash []byte
	DB().View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(chainBucket))
		newestHash = b.Get([]byte(cursor))
		return nil
	})
	return newestHash
}

func GetExistMempool() []byte {
	var memData []byte
	DB().View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(mempoolBucket))
		memData = b.Get([]byte(mempoolData))
		return nil
	})
	return memData
}

func PushOnMempool(data []byte) {
	DB().Update(func(t *bolt.Tx) error {
		err := t.DeleteBucket([]byte(mempoolBucket))
		utils.HandleErr(err)
		b, err := t.CreateBucketIfNotExists([]byte(mempoolBucket))
		utils.HandleErr(err)
		err = b.Put([]byte(mempoolData), data)
		utils.HandleErr(err)
		return nil
	})
}

func FindBlock(hash string) []byte {
	var blockAsBytes []byte
	DB().View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blockBucket))
		blockAsBytes = b.Get([]byte(hash))
		return nil
	})
	return blockAsBytes
}
