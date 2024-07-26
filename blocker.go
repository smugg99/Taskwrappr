// blocker.go
package taskwrappr

type Block struct {
    Actions    []*Action
	Executed   bool
    Memory     *MemoryMap
	LastResult *Variable
}

func NewBlock() *Block {
    return &Block{
        Actions: []*Action{},
    }
}