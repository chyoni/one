package interfaces

type DBOperation interface {
	GetExistChain() []byte
	SaveChainDB(data []byte)
	FindBlock(hash string) []byte
	SaveBlockDB(key string, data []byte)
}
