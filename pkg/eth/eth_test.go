package eth

import (
	"os"
	"strings"
	"testing"
)

const defaultRPCURL = "https://ethereum-rpc.publicnode.com"

func newClient(t *testing.T) *Client {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}
	url := os.Getenv("ETH_RPC_URL")
	if url == "" {
		url = defaultRPCURL
	}
	return NewClient(url)
}

func TestParseHexUint64(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    uint64
		wantErr bool
	}{
		{"with prefix", "0x10", 16, false},
		{"without prefix", "ff", 255, false},
		{"zero", "0x0", 0, false},
		{"empty", "0x", 0, true},
		{"invalid", "0xzz", 0, true},
		{"overflow", "0x1" + strings.Repeat("0", 16), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHexUint64(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestRandomFromSeed(t *testing.T) {
	const seed = "0xabcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"

	t.Run("deterministic", func(t *testing.T) {
		a, err := randomFromSeed(seed, 1, 100)
		if err != nil {
			t.Fatal(err)
		}
		b, err := randomFromSeed(seed, 1, 100)
		if err != nil {
			t.Fatal(err)
		}
		if a != b {
			t.Errorf("not deterministic: %d != %d", a, b)
		}
	})

	t.Run("within range", func(t *testing.T) {
		for i := int64(0); i < 50; i++ {
			// シードの末尾を変えて様々な値を試す。
			s := seed[:len(seed)-1] + string("0123456789abcdef"[i%16])
			got, err := randomFromSeed(s, 10, 20)
			if err != nil {
				t.Fatal(err)
			}
			if got < 10 || got > 20 {
				t.Errorf("got %d out of [10,20]", got)
			}
		}
	})

	t.Run("single value range", func(t *testing.T) {
		got, err := randomFromSeed(seed, 7, 7)
		if err != nil {
			t.Fatal(err)
		}
		if got != 7 {
			t.Errorf("got %d, want 7", got)
		}
	})

	t.Run("invalid seed", func(t *testing.T) {
		if _, err := randomFromSeed("0xzz", 1, 10); err == nil {
			t.Error("expected error for invalid seed")
		}
	})
}

func TestRandomStringFromSeed(t *testing.T) {
	const seed = "0xabcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789"

	t.Run("length and charset", func(t *testing.T) {
		s, err := randomStringFromSeed(seed, 8)
		if err != nil {
			t.Fatal(err)
		}
		if len(s) != 8 {
			t.Fatalf("len = %d, want 8", len(s))
		}
		for _, r := range s {
			if !strings.ContainsRune(charset, r) {
				t.Errorf("char %q not in charset", r)
			}
		}
	})

	t.Run("deterministic", func(t *testing.T) {
		a, _ := randomStringFromSeed(seed, 8)
		b, _ := randomStringFromSeed(seed, 8)
		if a != b {
			t.Errorf("not deterministic: %q != %q", a, b)
		}
	})

	t.Run("invalid seed", func(t *testing.T) {
		if _, err := randomStringFromSeed("0xzz", 8); err == nil {
			t.Error("expected error for invalid seed")
		}
	})
}

func TestBlockNumber(t *testing.T) {
	c := newClient(t)
	got, err := c.BlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	if got < 20_000_000 {
		t.Errorf("block number %d looks too low for mainnet", got)
	}
	t.Logf("block number = %d", got)
}

func TestLatestBlockHash(t *testing.T) {
	c := newClient(t)
	got, err := c.LatestBlockHash()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, "0x") || len(got) != 66 {
		t.Errorf("unexpected block hash format: %q", got)
	}
	t.Logf("latest block hash = %s", got)
}

func TestRandomInRange(t *testing.T) {
	c := newClient(t)
	got, err := c.RandomInRange(1, 100)
	if err != nil {
		t.Fatal(err)
	}
	if got < 1 || got > 100 {
		t.Errorf("got %d out of [1,100]", got)
	}
	t.Logf("random in [1,100] = %d", got)
}

func TestRandomInRange_InvalidRange(t *testing.T) {
	c := NewClient("http://unused.invalid")
	if _, err := c.RandomInRange(10, 1); err == nil {
		t.Error("expected error for min > max")
	}
}

func TestRandomString(t *testing.T) {
	c := newClient(t)
	s, err := c.RandomString(8)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 8 {
		t.Errorf("len = %d, want 8", len(s))
	}
	for _, r := range s {
		if !strings.ContainsRune(charset, r) {
			t.Errorf("char %q not in charset", r)
		}
	}
	t.Logf("random string = %s", s)
}

func TestRandomString_InvalidLength(t *testing.T) {
	c := NewClient("http://unused.invalid")
	if _, err := c.RandomString(0); err == nil {
		t.Error("expected error for length 0")
	}
}
