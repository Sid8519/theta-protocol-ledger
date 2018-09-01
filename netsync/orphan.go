package netsync

import (
	"container/list"

	"github.com/thetatoken/ukulele/blockchain"
	"github.com/thetatoken/ukulele/common"
)

const (
	maxOrphanBlockPoolSize = 64
	maxOrphanCCPoolSize    = 64
)

type OrphanBlockPool struct {
	blocks          *list.List
	hashToBlock     map[string]*list.Element
	prevHashToBlock map[string]*list.Element
}

func NewOrphanBlockPool() *OrphanBlockPool {
	return &OrphanBlockPool{
		blocks:          list.New(),
		hashToBlock:     make(map[string]*list.Element),
		prevHashToBlock: make(map[string]*list.Element),
	}
}

func (bp *OrphanBlockPool) Contains(block *blockchain.Block) bool {
	_, ok := bp.hashToBlock[block.Hash.String()]
	return ok
}

func (bp *OrphanBlockPool) Add(block *blockchain.Block) {
	if bp.blocks.Len() >= maxOrphanBlockPoolSize {
		bp.RemoveOldest()
	}

	if bp.Contains(block) {
		return
	}

	el := bp.blocks.PushBack(block)
	bp.hashToBlock[block.Hash.String()] = el
	bp.prevHashToBlock[block.ParentHash.String()] = el
}

func (bp *OrphanBlockPool) Remove(block *blockchain.Block) {
	el, ok := bp.hashToBlock[block.Hash.String()]
	if !ok {
		// block is not in pool.
		return
	}
	bp.blocks.Remove(el)
	delete(bp.hashToBlock, block.Hash.String())
	delete(bp.prevHashToBlock, block.ParentHash.String())
}

func (bp *OrphanBlockPool) RemoveOldest() {
	el := bp.blocks.Front()
	if el == nil {
		return
	}
	block := el.Value.(*blockchain.Block)
	bp.Remove(block)
}

func (bp *OrphanBlockPool) TryGetNextBlock(hash common.Bytes) *blockchain.Block {
	el, ok := bp.prevHashToBlock[hash.String()]
	if !ok {
		return nil
	}
	block := el.Value.(*blockchain.Block)
	bp.Remove(block)
	return block
}

type OrphanCCPool struct {
	ccs      *list.List
	hashToCC map[string]*list.Element
}

func NewOrphanCCPool() *OrphanCCPool {
	return &OrphanCCPool{
		ccs:      list.New(),
		hashToCC: make(map[string]*list.Element),
	}
}

func (cp *OrphanCCPool) Contains(cc *blockchain.CommitCertificate) bool {
	_, ok := cp.hashToCC[cc.BlockHash.String()]
	return ok
}

func (cp *OrphanCCPool) Add(cc *blockchain.CommitCertificate) {
	if cp.ccs.Len() >= maxOrphanCCPoolSize {
		cp.RemoveOldest()
	}

	if cp.Contains(cc) {
		return
	}

	el := cp.ccs.PushBack(cc)
	cp.hashToCC[cc.BlockHash.String()] = el
}

func (cp *OrphanCCPool) Remove(cc *blockchain.CommitCertificate) {
	el, ok := cp.hashToCC[cc.BlockHash.String()]
	if !ok {
		// block is not in pool.
		return
	}
	cp.ccs.Remove(el)
	delete(cp.hashToCC, cc.BlockHash.String())
}

func (cp *OrphanCCPool) RemoveOldest() {
	el := cp.ccs.Front()
	if el == nil {
		return
	}
	cc := el.Value.(*blockchain.CommitCertificate)
	cp.Remove(cc)
}

func (cp *OrphanCCPool) TryGetCCByBlockHash(hash common.Bytes) *blockchain.CommitCertificate {
	el, ok := cp.hashToCC[hash.String()]
	if !ok {
		return nil
	}
	cc := el.Value.(*blockchain.CommitCertificate)
	cp.Remove(cc)
	return cc
}
