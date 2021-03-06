package proof

import (
	"errors"
	"fmt"
	"time"

	"getmelange.com/zooko/message"
	"github.com/melange-app/nmcd/wire"
)

type blockHeader struct {
	// We can throw out all the auxiliary data once
	// we have verified the block and put it in the chain.
	MerkleRoot wire.ShaHash
	Timestamp  time.Time
	Hash       wire.ShaHash
	Height     int
	// *wire.BlockHeader

	Previous *blockHeader
	IsBottom bool
}

type chainBase struct {
	wire.ShaHash
	Height int
}

type blockchain struct {
	Base *chainBase
	Top  *blockHeader

	Direct   map[string]*blockHeader
	Orphaned map[string]*blockHeader

	SyncedTime time.Time

	Height int
}

func (b *blockchain) verifyTransaction(m message.NamecoinTransaction) bool {
	blck, ok := b.Direct[m.BlockId]
	if !ok {
		fmt.Println(m.BlockId, "not in blockchain.")
		return false
	}

	var hashes []wire.ShaHash
	for _, v := range m.VerificationHashes {
		temp, err := wire.NewShaHash(v)
		if err != nil {
			fmt.Println("Unable to create Verification Sha hash", err)
			return false
		}

		hashes = append(hashes, *temp)
	}

	txHash, err := wire.NewShaHashFromStr(m.TxId)
	if err != nil {
		fmt.Println("Cannot make TX Sha Hash", err)
		return false
	}

	_, err = VerifyMerkleBranch(wire.MerkleBranch{
		BranchSideMask: m.Branch,
		BranchHash:     hashes,
	}, &blck.MerkleRoot, txHash)
	if err != nil {
		fmt.Println("Error verifying Merkle Branch", err)
	}

	return err == nil
}

func (b *blockchain) createChainHeader(h *wire.BlockHeader) (*blockHeader, error) {
	hash, err := h.BlockSha()
	return &blockHeader{
		MerkleRoot: h.MerkleRoot,
		Timestamp:  h.Timestamp,
		Hash:       hash,
	}, err
}

// will only accept headers if they are the "next in the chain"
func (b *blockchain) addHeader(h *wire.BlockHeader) error {
	block, err := b.createChainHeader(h)
	if err != nil {
		return err
	}

	if _, ok := b.Direct[block.Hash.String()]; ok {
		return errors.New("We already have this block...")
	}

	// Add the direct block mapping.
	b.Direct[block.Hash.String()] = block

	// Is the PreviousBlock Hash already in the Chain?
	if check, ok := b.Direct[h.PrevBlock.String()]; ok {
		block.Height = check.Height + 1
		block.Previous = check

		if block.Height > b.Height {
			b.Top = block
			b.Height = block.Height
			// fmt.Println("ACCEPTED New Top Header", block.Hash.String()[:20], "at height", block.Height)
		} else {
			// fmt.Println("ACCEPTED Orphaned Header", block.Hash.String()[:20], "at height", block.Height)
		}

		if b.SyncedTime.Before(h.Timestamp) {
			b.SyncedTime = h.Timestamp
		}
	} else if b.Top == nil && b.Base.IsEqual(&h.PrevBlock) {
		// This block is at the top of the chain.
		block.IsBottom = true
		block.Height = b.Base.Height + 1

		b.Height = block.Height
		b.Top = block

		b.SyncedTime = h.Timestamp
		// fmt.Println("ACCEPTED Base Header", block.Hash.String()[:20], "at height", block.Height)
	} else {
		// This block is orphaned
		b.Orphaned[block.Hash.String()] = block
		// fmt.Println("ACCEPTED Orphaned Header", block.Hash.String()[:20], "at height", block.Height)
	}

	return nil
}

type blockchainManager struct {
	acceptChannel chan interface{}
	syncedChannel chan struct{}

	isSynced bool
	chain    *blockchain
}

func (b *blockchainManager) waitForSync() {
	if b.isSynced {
		return
	}

	fmt.Println("Waiting for blockchain sync.")
	_, done := <-b.syncedChannel
	fmt.Println("Blockchain has finished syncing.", !done)
}

func (b *blockchainManager) VerifyTransactions(m ...message.NamecoinTransaction) bool {
	b.waitForSync()

	for _, v := range m {
		results := make(chan bool)
		b.acceptChannel <- transactionVerification{
			NamecoinTransaction: v,
			resultChan:          results,
		}

		if answer := <-results; !answer {
			return false
		}
	}

	return true
}

func (b *blockchainManager) TopHeader() (*blockHeader, int32) {
	return b.chain.Top, int32(b.chain.Height)
}

func CreateBlockchainManager() *blockchainManager {
	b := &blockchainManager{
		acceptChannel: make(chan interface{}),
		syncedChannel: make(chan struct{}),
		chain: &blockchain{
			Height: TopResolverHeight,
			Base: &chainBase{
				ShaHash: *ResolverLocatorHashes[0],
				Height:  TopResolverHeight,
			},
			Direct: make(map[string]*blockHeader),
		},
	}
	go b.acceptanceLoop()

	return b
}

type transactionVerification struct {
	message.NamecoinTransaction
	resultChan chan bool
}

func (b *blockchainManager) acceptanceLoop() {
	for {
		obj := <-b.acceptChannel

		switch t := obj.(type) {
		case *wire.BlockHeader:
			b.addHeader(t)
			if b.chain.SyncedTime.After(time.Now().Add(-1*time.Hour)) && !b.isSynced {
				close(b.syncedChannel)
				b.isSynced = true
			}
		case transactionVerification:
			result := b.chain.verifyTransaction(t.NamecoinTransaction)
			t.resultChan <- result
		}
	}
}

func (b *blockchainManager) addHeader(t *wire.BlockHeader) {
	err := b.chain.addHeader(t)
	if err != nil {
		fmt.Println("Unable to add header to chain", err)
	}
}

func (b *blockchainManager) AddHeader(hdrs ...*wire.BlockHeader) {
	// fmt.Println(hdrs[0].PrevBlock)
	for _, v := range hdrs {
		b.acceptChannel <- v
	}
}
