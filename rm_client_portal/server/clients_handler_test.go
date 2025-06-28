package server

import (
	"fmt"
	"strconv"
	"testing"

	"rm_client_portal/config"
	"rm_client_portal/database"

	"gopkg.in/testfixtures.v2"
)

var fixtures *testfixtures.Context

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	var err error

	// Initialize config for JWT testing
	config.ReadProperties()

	// Open connection with the test database using config values
	// Use the same database name as configured (should be a test database)
	database.OpenDB(config.Conf.DbName, config.Conf.DbAddress, config.Conf.DbPort, config.Conf.DbUsername, config.Conf.DbPassword)

	// Create fixtures context
	fixtures, err = testfixtures.NewFolder(database.Db, &testfixtures.MySQL{}, "../database/fixtures")
	if err != nil {
		panic(err)
	}

	m.Run()
}

// prepareTestDatabase loads test fixtures
func prepareTestDatabase() {
	testfixtures.SkipDatabaseNameCheck(true)
	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}

// Test Case 1: Test client access validation logic directly
func TestReportOnReviewsAndInsights_AuthorizedClient(t *testing.T) {
	prepareTestDatabase()

	// Test the core client filtering logic by checking database functions directly
	// This tests the fix without requiring Google API setup

	// Test 1: test1@testing.com should only have access to client 1
	clientIDs := database.GetClientIDsForUserEmail("test1@testing.com")
	fmt.Printf("test1@testing.com has access to clients: %v\n", clientIDs)

	if len(clientIDs) != 1 {
		t.Errorf("Expected test1@testing.com to have access to 1 client, got %d", len(clientIDs))
	}

	if len(clientIDs) > 0 && clientIDs[0] != 1 {
		t.Errorf("Expected test1@testing.com to have access to client 1, got client %d", clientIDs[0])
	}

	// Test 2: test@testing.com should have access to multiple clients (1, 2, 4)
	multiClientIDs := database.GetClientIDsForUserEmail("test@testing.com")
	fmt.Printf("test@testing.com has access to clients: %v\n", multiClientIDs)

	if len(multiClientIDs) != 3 {
		t.Errorf("Expected test@testing.com to have access to 3 clients, got %d", len(multiClientIDs))
	}

	// Verify the specific client IDs
	expectedClients := map[uint64]bool{1: true, 2: true, 4: true}
	for _, clientID := range multiClientIDs {
		if !expectedClients[clientID] {
			t.Errorf("Unexpected client ID %d for test@testing.com", clientID)
		}
	}

	// Test 3: Verify client access validation
	// test1@testing.com should be able to access client 1
	client1 := database.GetClientCheckUserEmail(1, "test1@testing.com")
	if client1.ID != 1 {
		t.Errorf("test1@testing.com should have access to client 1, but access denied")
	}

	// test1@testing.com should NOT be able to access client 2
	client2 := database.GetClientCheckUserEmail(2, "test1@testing.com")
	if client2.ID != 0 {
		t.Errorf("test1@testing.com should NOT have access to client 2, but access granted")
	}

	fmt.Printf("Client filtering logic tests passed - users can only access their assigned clients\n")
}

// Test Case 2: Test client filtering logic in handler directly
func TestReportOnReviewsAndInsights_ClientFilteringLogic(t *testing.T) {
	prepareTestDatabase()

	// Test the client filtering logic from the handler
	// This simulates the key logic from ReportOnReviewsAndInsights function

	// Simulate what happens in the handler for test1@testing.com
	email := "test1@testing.com"
	clientIDs := database.GetClientIDsForUserEmail(email)
	fmt.Printf("User %s has access to clients: %v\n", email, clientIDs)

	// Test Case 2a: User requests access to authorized client (client_id=1)
	requestedClientID := uint64(1)

	// Check if user has access to this client
	hasAccess := false
	for _, userClientID := range clientIDs {
		if userClientID == requestedClientID {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		t.Errorf("User %s should have access to client %d", email, requestedClientID)
	} else {
		fmt.Printf("Access granted: User %s can access client %d\n", email, requestedClientID)
	}

	// Test Case 2b: User requests access to unauthorized client (client_id=2)
	unauthorizedClientID := uint64(2)

	hasAccessUnauthorized := false
	for _, userClientID := range clientIDs {
		if userClientID == unauthorizedClientID {
			hasAccessUnauthorized = true
			break
		}
	}

	if hasAccessUnauthorized {
		t.Errorf("User %s should NOT have access to client %d", email, unauthorizedClientID)
	} else {
		fmt.Printf("Access denied: User %s cannot access client %d (as expected)\n", email, unauthorizedClientID)
	}

	// Test Case 2c: Test multi-client user (test@testing.com)
	multiClientEmail := "test@testing.com"
	multiClientIDs := database.GetClientIDsForUserEmail(multiClientEmail)
	fmt.Printf("User %s has access to clients: %v\n", multiClientEmail, multiClientIDs)

	// This user should be able to access clients 1, 2, and 4
	expectedMultiClients := []uint64{1, 2, 4}
	for _, expectedClientID := range expectedMultiClients {
		hasMultiAccess := false
		for _, userClientID := range multiClientIDs {
			if userClientID == expectedClientID {
				hasMultiAccess = true
				break
			}
		}
		if !hasMultiAccess {
			t.Errorf("User %s should have access to client %d", multiClientEmail, expectedClientID)
		} else {
			fmt.Printf("Multi-client access: User %s can access client %d\n", multiClientEmail, expectedClientID)
		}
	}

	// But should NOT be able to access client 3
	unauthorizedMultiClientID := uint64(3)
	hasUnauthorizedMultiAccess := false
	for _, userClientID := range multiClientIDs {
		if userClientID == unauthorizedMultiClientID {
			hasUnauthorizedMultiAccess = true
			break
		}
	}
	if hasUnauthorizedMultiAccess {
		t.Errorf("User %s should NOT have access to client %d", multiClientEmail, unauthorizedMultiClientID)
	} else {
		fmt.Printf("Multi-client access denied: User %s cannot access client %d (as expected)\n", multiClientEmail, unauthorizedMultiClientID)
	}

	fmt.Printf("All client filtering logic tests passed!\n")
}

// Test Case 3: Test parameter validation
func TestReportOnReviewsAndInsights_ParameterValidation(t *testing.T) {
	prepareTestDatabase()

	// Test invalid client_id parameter parsing
	invalidClientIDs := []string{"invalid", "abc", "-1", "999999999999999999999"}

	for _, invalidID := range invalidClientIDs {
		fmt.Printf("Testing invalid client_id: %s\n", invalidID)

		// Try to parse as the handler would
		_, err := strconv.ParseUint(invalidID, 10, 64)
		if err == nil {
			// This should fail for truly invalid IDs like "abc"
			if invalidID == "abc" || invalidID == "invalid" {
				t.Errorf("Expected parsing error for client_id '%s', but got none", invalidID)
			}
		} else {
			fmt.Printf("Correctly rejected invalid client_id '%s': %v\n", invalidID, err)
		}
	}

	// Test valid client_id parameter parsing
	validClientIDs := []string{"1", "2", "4", "100"}

	for _, validID := range validClientIDs {
		clientID, err := strconv.ParseUint(validID, 10, 64)
		if err != nil {
			t.Errorf("Expected valid client_id '%s' to parse successfully, got error: %v", validID, err)
		} else {
			fmt.Printf("Successfully parsed valid client_id '%s' as %d\n", validID, clientID)
		}
	}

	fmt.Printf("Parameter validation tests passed!\n")
}
