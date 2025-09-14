package examples

import (
	"unsafe"

	"github.com/mziter/weft"
)

// BankAccount represents a simple bank account with thread-safe operations.
// This is a classic example for testing concurrent access patterns.
type BankAccount struct {
	mu      weft.RWMutex
	balance int
}

// NewBankAccount creates a new account with the given initial balance.
func NewBankAccount(initialBalance int) *BankAccount {
	return &BankAccount{
		balance: initialBalance,
	}
}

// Deposit adds money to the account.
func (b *BankAccount) Deposit(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.balance += amount
}

// Withdraw removes money from the account if sufficient funds exist.
func (b *BankAccount) Withdraw(amount int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.balance >= amount {
		b.balance -= amount
		return true
	}
	return false
}

// Balance returns the current balance.
func (b *BankAccount) Balance() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.balance
}

// Transfer moves money between two accounts.
// This can cause deadlock if not carefully ordered.
func Transfer(from, to *BankAccount, amount int) bool {
	// This implementation can deadlock if another Transfer
	// is happening in the opposite direction simultaneously
	from.mu.Lock()
	defer from.mu.Unlock()

	if from.balance < amount {
		return false
	}

	to.mu.Lock()
	defer to.mu.Unlock()

	from.balance -= amount
	to.balance += amount
	return true
}

// SafeTransfer uses lock ordering to prevent deadlock.
func SafeTransfer(from, to *BankAccount, amount int) bool {
	// Order locks by memory address to prevent deadlock
	first, second := from, to
	if uintptr(unsafe.Pointer(from)) > uintptr(unsafe.Pointer(to)) {
		first, second = to, from
	}

	first.mu.Lock()
	defer first.mu.Unlock()
	second.mu.Lock()
	defer second.mu.Unlock()

	if from.balance < amount {
		return false
	}

	from.balance -= amount
	to.balance += amount
	return true
}