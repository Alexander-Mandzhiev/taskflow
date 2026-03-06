package password

import (
	"testing"
)

func TestNewBcryptHasher(t *testing.T) {
	t.Run("zero cost uses default", func(t *testing.T) {
		h := NewBcryptHasher(0)
		if h == nil {
			t.Fatal("NewBcryptHasher(0) = nil")
		}
		if h.Cost <= 0 {
			t.Errorf("Cost = %d, want positive default", h.Cost)
		}
	})

	t.Run("negative cost uses default", func(t *testing.T) {
		h := NewBcryptHasher(-1)
		if h.Cost <= 0 {
			t.Errorf("Cost = %d, want positive default", h.Cost)
		}
	})

	t.Run("explicit cost", func(t *testing.T) {
		h := NewBcryptHasher(4)
		if h.Cost != 4 {
			t.Errorf("Cost = %d, want 4", h.Cost)
		}
	})
}

func TestBcryptHasher_Hash_Compare(t *testing.T) {
	// cost 4 для быстрых тестов
	h := NewBcryptHasher(4)

	t.Run("Hash then Compare succeeds", func(t *testing.T) {
		plain := "secret123"
		hashed, err := h.Hash(plain)
		if err != nil {
			t.Fatalf("Hash: %v", err)
		}
		if hashed == "" || hashed == plain {
			t.Errorf("Hash: expected bcrypt hash, got %q", hashed)
		}
		if err := h.Compare(hashed, plain); err != nil {
			t.Errorf("Compare(hashed, plain): %v", err)
		}
	})

	t.Run("Compare wrong password fails", func(t *testing.T) {
		hashed, err := h.Hash("correct")
		if err != nil {
			t.Fatalf("Hash: %v", err)
		}
		if err := h.Compare(hashed, "wrong"); err == nil {
			t.Error("Compare(hashed, wrong): expected error")
		}
	})

	t.Run("Compare invalid hash returns error", func(t *testing.T) {
		if err := h.Compare("not-a-bcrypt-hash", "any"); err == nil {
			t.Error("Compare(invalid hash, any): expected error")
		}
	})

	t.Run("Hash produces different salts", func(t *testing.T) {
		h1, _ := h.Hash("same")
		h2, _ := h.Hash("same")
		if h1 == h2 {
			t.Error("two hashes of same password should differ (salt)")
		}
		if err := h.Compare(h1, "same"); err != nil {
			t.Errorf("Compare first hash: %v", err)
		}
		if err := h.Compare(h2, "same"); err != nil {
			t.Errorf("Compare second hash: %v", err)
		}
	})
}
