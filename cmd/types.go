package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	WalletId string // unique id
	Balance  int64  // balance
	state    int32  // state 0 lock 1 unlock
}

func NewWallet() *Wallet {
	return &Wallet{
		WalletId: uuid.NewString(),
		Balance:  0,
		state:    0,
	}
}

var GlobalWalletCollect WalletCollect

type WalletCollect struct {
	N         int
	buf       []*Wallet
	mu        sync.Mutex
	indexer   map[string]int
	lockArray []*MyLock
}

func (w *WalletCollect) Add(wallet *Wallet) {
	w.mu.Lock()
	w.buf = append(w.buf, wallet)
	w.indexer[wallet.WalletId] = w.N
	w.N++
	w.mu.Unlock()

	slog.Info("wallet add success", "wallet_id", wallet.WalletId)
}

func (w *WalletCollect) Get(wallentId string) *Wallet {
	index, ok := w.indexer[wallentId]
	if !ok {
		return nil
	}

	return w.buf[index]
}

func (w *WalletCollect) getTransferLock(walletFromA, walletFromB string) *MyLock {
	var key string
	if walletFromA < walletFromB {
		key = walletFromA + "/" + walletFromB
	} else {
		key = walletFromB + "/" + walletFromA
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	for _, l := range w.lockArray {
		if l.key == key {
			return l
		}
	}

	mylock := &MyLock{
		key:  key,
		lock: sync.Mutex{},
	}
	w.lockArray = append(w.lockArray, mylock)

	return mylock
}

func (w *WalletCollect) Transfer(walletFromA, walletFromB string, amount int64) error {
	indexA, ok := w.indexer[walletFromA]
	if !ok {
		return fmt.Errorf("wallent not exists: %s", walletFromA)
	}
	wallentA := w.buf[indexA]

	indexB, ok := w.indexer[walletFromB]
	if !ok {
		return fmt.Errorf("wallent not exists: %s", walletFromB)
	}
	wallentB := w.buf[indexB]

	timeoutCtx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*100)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("system busy")
		default:
			if atomic.LoadInt32(&wallentA.state) == 0 && atomic.LoadInt32(&wallentB.state) == 0 {
				// get lock
				l := w.getTransferLock(walletFromA, walletFromB)

				// try to lock
				l.lock.Lock()
				wallentA.state = 1
				wallentB.state = 1

				if wallentA.Balance < amount {
					l.lock.Unlock()
					return fmt.Errorf("insufficient balance")
				}
				wallentA.Balance -= amount
				wallentB.Balance += amount
				l.lock.Unlock()

				atomic.StoreInt32(&wallentA.state, 0)
				atomic.StoreInt32(&wallentB.state, 0)

				return nil
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func init() {
	GlobalWalletCollect.buf = make([]*Wallet, 0, 100000)
	GlobalWalletCollect.indexer = make(map[string]int)
	GlobalWalletCollect.Add(&Wallet{WalletId: "A", Balance: 1000})
	GlobalWalletCollect.Add(&Wallet{WalletId: "B", Balance: 1000})
	GlobalWalletCollect.Add(&Wallet{WalletId: "C", Balance: 1000})
}
