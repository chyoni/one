package interfaces

type DBOperation interface {
	GetExistChain() []byte
}
