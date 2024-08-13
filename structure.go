// structure.go
package taskwrappr

type Block struct {
    Actions    []*Action
	Executed   bool
    Memory     *MemoryMap
	LastResult *Variable
}

func NewBlock(parentMemory *MemoryMap) *Block {
    mem := NewMemoryMap(parentMemory)
    return &Block{
        Actions: []*Action{},
        Memory:  mem,
    }
}