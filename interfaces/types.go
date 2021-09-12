package interfaces

type DBOperation interface {
	GetExistChain() []byte
	SaveChainDB(data []byte)
}
