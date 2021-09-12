package db

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chiwon99881/one/utils"
	bolt "go.etcd.io/bbolt"
)

type DataBase struct{}

func (DataBase) GetExistChain() []byte {
	return getExistChain()
}

func (DataBase) SaveChainDB(data []byte) {
	saveChainDB(data)
}

const (
	chainBucket   string = "chainBucket"
	blockBucket   string = "blockBucket"
	mempoolBucket string = "mempoolBucket"
	mempoolData   string = "mempoolData"
	cursor        string = "aCursor"
)

var db *bolt.DB

func getPort() (string, error) {
	var port string
	for _, flag := range os.Args {
		if strings.Contains(flag, "-port") {
			data := strings.Split(flag, "=")
			port = data[1]
			return port, nil
		}
	}
	return "", errors.New("port is undefined")
}

func DB() *bolt.DB {
	if db == nil {
		port, err := getPort()
		if err != nil {
			utils.HandleErr(err)
		}
		dbPointer, err := bolt.Open(fmt.Sprintf("onecoin.%s.db", port), 0600, nil)
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

func saveChainDB(data []byte) {
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

func getExistChain() []byte {
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

func CreateAfterDeleteDB() {
	DB().Update(func(t *bolt.Tx) error {
		err := t.DeleteBucket([]byte(chainBucket))
		utils.HandleErr(err)
		err = t.DeleteBucket([]byte(blockBucket))
		utils.HandleErr(err)
		_, err = t.CreateBucket([]byte(chainBucket))
		utils.HandleErr(err)
		_, err = t.CreateBucket([]byte(blockBucket))
		utils.HandleErr(err)
		return nil
	})
}
