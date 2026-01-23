package kvk

import (
	"context"
	"os"
	"testing"
)

// TEST_API_KEY is the key provided in the public documentation for the KVK Test Environment
const TEST_API_KEY = "l7xx1f2691f2520d487b902f4e0b57a0b197"

func TestIntegration(t *testing.T) {
	// Skip if explicitly requested (e.g. in CI without net access, though this is a unit test file, we treat it as integration)
	if os.Getenv("SKIP_INTEGRATION") != "" {
		t.Skip("Skipping integration test")
	}

	cfg := &Config{
		BaseURL: "https://api.kvk.nl/test/api",
		APIKey:  TEST_API_KEY,
	}
	client := NewClient(cfg)

	ctx := context.Background()

	t.Run("Search", func(t *testing.T) {
		// Searching for "kvk" or similar common term in test env
		// The test environment data is limited, "B.V." or specific names might be needed.
		// Detailed docs say: "Maak je eigen zoek query door zoektermen toe te voegen"
		// Let's try a broad search or a known test case if documented.
		// The docs don't list specific test cases easily without digging, so we try a generic search.
		resp, err := client.Search(ctx, "test")
		if err != nil {
			t.Fatalf("Search failed: %v", err)
		}

		t.Logf("Search returned %d results", resp.TotaalAantal)
		if len(resp.Resultaten) > 0 {
			t.Logf("First result: %+v", resp.Resultaten[0])
		}
	})

	t.Run("GetBasicProfile", func(t *testing.T) {
		// We need a valid KVK number from the test environment.
		// Often 90000000+ numbers are used for testing or we can pick one from search if available.
		// Let's first search to get a valid number, then fetch profile.
		// This makes the test self-sustaining.

		searchResp, err := client.Search(ctx, "bv") // "bv" is likely to return results
		if err != nil {
			t.Fatalf("Setup Search failed: %v", err)
		}

		if len(searchResp.Resultaten) == 0 {
			t.Skip("No search results found to test GetBasicProfile")
		}

		targetKvk := searchResp.Resultaten[0].KvkNummer
		t.Logf("Testing GetBasicProfile with KVK: %s", targetKvk)

		profile, err := client.GetBasicProfile(ctx, targetKvk)
		if err != nil {
			t.Fatalf("GetBasicProfile failed: %v", err)
		}

		if profile.KvkNummer != targetKvk {
			t.Errorf("Expected KVK %s, got %s", targetKvk, profile.KvkNummer)
		}
		t.Logf("Got profile: %+v", profile)
	})
}
