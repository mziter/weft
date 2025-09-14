package examples

import (
	"testing"

	"github.com/mziter/weft"
	"github.com/mziter/weft/wefttest"
)

// TestBankAccountTransfer demonstrates testing deadlock scenarios.
func TestBankAccountTransfer(t *testing.T) {
	t.Skip("Deadlock detection not yet implemented")

	wefttest.Explore(t, 100, func(s *weft.Scheduler) {
		account1 := NewBankAccount(100)
		account2 := NewBankAccount(100)

		// This can deadlock due to lock ordering
		s.Go(func(ctx weft.Context) {
			Transfer(account1, account2, 50)
		})

		s.Go(func(ctx weft.Context) {
			Transfer(account2, account1, 30)
		})

		s.Wait()

		// Check final balances
		total := account1.Balance() + account2.Balance()
		if total != 200 {
			t.Errorf("money disappeared: expected total=200, got %d", total)
		}
	})
}

// TestSafeBankAccountTransfer shows how the fixed version works.
func TestSafeBankAccountTransfer(t *testing.T) {
	wefttest.Explore(t, 100, func(s *weft.Scheduler) {
		account1 := NewBankAccount(100)
		account2 := NewBankAccount(100)

		success1, success2 := false, false

		// These transfers use proper lock ordering
		s.Go(func(ctx weft.Context) {
			success1 = SafeTransfer(account1, account2, 50)
		})

		s.Go(func(ctx weft.Context) {
			success2 = SafeTransfer(account2, account1, 30)
		})

		s.Wait()

		// Verify no money was lost
		total := account1.Balance() + account2.Balance()
		if total != 200 {
			t.Errorf("money disappeared: expected total=200, got %d", total)
		}

		// At least one transfer should succeed
		if !success1 && !success2 {
			t.Error("both transfers failed, expected at least one to succeed")
		}
	})
}