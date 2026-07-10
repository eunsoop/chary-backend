package model

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewWorkspace(t *testing.T) {
	ws, err := NewWorkspace("Enterprise Workspace")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ws.Name != "Enterprise Workspace" {
		t.Errorf("expected name 'Enterprise Workspace', got '%s'", ws.Name)
	}

	if _, err := uuid.Parse(ws.ID); err != nil {
		t.Errorf("expected valid ID UUID, got error: %v", err)
	}
}

func TestNewGlobalScopeDeterministicUUIDv5(t *testing.T) {
	wsID1 := uuid.New().String()
	wsID2 := uuid.New().String()
	domain := "acme.com"

	gs1, err := NewGlobalScope(wsID1, domain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gs1.ID == uuid.Nil.String() {
		t.Errorf("global scope ID cannot be Nil UUID")
	}

	// Verify determinism: same workspace and same domain must yield identical UUID
	gs1Dup, err := NewGlobalScope(wsID1, domain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gs1.ID != gs1Dup.ID {
		t.Errorf("expected deterministic UUIDv5, got %s and %s", gs1.ID, gs1Dup.ID)
	}

	// Verify isolation: same domain but different workspace must yield different UUID
	gs2, err := NewGlobalScope(wsID2, domain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gs1.ID == gs2.ID {
		t.Errorf("expected different workspace to yield different UUIDv5, but both got %s", gs1.ID)
	}

	// Verify isolation: same workspace but different domain must yield different UUID
	gs3, err := NewGlobalScope(wsID1, "other.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gs1.ID == gs3.ID {
		t.Errorf("expected different domain to yield different UUIDv5, but both got %s", gs1.ID)
	}
}

func TestNewTeamScope(t *testing.T) {
	wsID := uuid.New().String()
	ts, err := NewTeamScope(wsID, "Engineering Task Force")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ts.Name != "Engineering Task Force" {
		t.Errorf("expected name 'Engineering Task Force', got '%s'", ts.Name)
	}

	if ts.WorkspaceID != wsID {
		t.Errorf("expected workspace ID '%s', got '%s'", wsID, ts.WorkspaceID)
	}
}

func TestCategoryAndChannel(t *testing.T) {
	scopeID := uuid.New().String()
	cat, err := NewCategory(scopeID, ScopeGlobal, "Engineering", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cat.ScopeID != scopeID || cat.ScopeType != ScopeGlobal || cat.Name != "Engineering" || cat.Order != 10 {
		t.Errorf("category fields mismatch")
	}

	ch, err := NewChannel(cat.ID, "alerts", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ch.CategoryID != cat.ID || ch.Name != "alerts" || !ch.IsPrivate {
		t.Errorf("channel fields mismatch")
	}
}

func TestScopeKeyHandling(t *testing.T) {
	wsID := uuid.New().String()
	teamID := uuid.New().String()
	catID := uuid.New().String()
	chID := uuid.New().String()

	// Global Key Test
	globalKey, err := BuildGlobalScopeKey(wsID, catID, chID)
	if err != nil {
		t.Fatalf("unexpected global build error: %v", err)
	}

	expectedGlobalKey := "workspace:" + wsID + ":global.category:" + catID + ".channel:" + chID
	if globalKey != expectedGlobalKey {
		t.Errorf("expected global key '%s', got '%s'", expectedGlobalKey, globalKey)
	}

	parsedGlobal, err := ParseScopeKey(globalKey)
	if err != nil {
		t.Fatalf("unexpected global parse error: %v", err)
	}

	if parsedGlobal.WorkspaceID != wsID || !parsedGlobal.IsGlobal || parsedGlobal.CategoryID != catID || parsedGlobal.ChannelID != chID || parsedGlobal.TeamID != "" {
		t.Errorf("parsed global scope key mismatch")
	}

	// Team Key Test
	teamKey, err := BuildTeamScopeKey(wsID, teamID, catID, chID)
	if err != nil {
		t.Fatalf("unexpected team build error: %v", err)
	}

	expectedTeamKey := "workspace:" + wsID + ":team:" + teamID + ".category:" + catID + ".channel:" + chID
	if teamKey != expectedTeamKey {
		t.Errorf("expected team key '%s', got '%s'", expectedTeamKey, teamKey)
	}

	parsedTeam, err := ParseScopeKey(teamKey)
	if err != nil {
		t.Fatalf("unexpected team parse error: %v", err)
	}

	if parsedTeam.WorkspaceID != wsID || parsedTeam.IsGlobal || parsedTeam.TeamID != teamID || parsedTeam.CategoryID != catID || parsedTeam.ChannelID != chID {
		t.Errorf("parsed team scope key mismatch")
	}

	// Invalid Keys Tests
	invalidKeys := []string{
		"workspace:invalid:global.category:" + catID + ".channel:" + chID,
		"workspace:" + wsID + ":global.category:invalid.channel:" + chID,
		"workspace:" + wsID + ":global.category:" + catID + ".channel:invalid",
		"workspace:" + wsID + ":team:" + teamID + ".category:" + catID,
		"workspace:" + wsID + ":team:invalid.category:" + catID + ".channel:" + chID,
		"workspace:" + wsID + ":something:invalid",
		"invalidprefix:" + wsID + ":global.category:" + catID + ".channel:" + chID,
	}

	for _, k := range invalidKeys {
		if _, err := ParseScopeKey(k); err == nil {
			t.Errorf("expected parsing error for key '%s', but succeeded", k)
		}
	}
}

func TestIntegrationsAndWebhooks(t *testing.T) {
	chID := uuid.New().String()

	// Test Integration Creation
	intg, err := NewIntegration(chID, "Slack Incoming Sync", IntegrationWebhook)
	if err != nil {
		t.Fatalf("unexpected integration creation error: %v", err)
	}

	if intg.ChannelID != chID || intg.Name != "Slack Incoming Sync" || intg.Type != IntegrationWebhook || !intg.IsEnabled {
		t.Errorf("integration fields mismatch")
	}

	// Test Incoming Webhook Creation
	incomingWh, err := NewIncomingWebhook(intg.ID)
	if err != nil {
		t.Fatalf("unexpected incoming webhook creation error: %v", err)
	}

	if incomingWh.IntegrationID != intg.ID || incomingWh.Type != WebhookIncoming || len(incomingWh.SecretToken) != 64 {
		t.Errorf("incoming webhook fields mismatch: token length %d", len(incomingWh.SecretToken))
	}

	// Test Outgoing Webhook Creation
	targetURL := "https://api.external.com/webhooks"
	outgoingWh, err := NewOutgoingWebhook(intg.ID, targetURL)
	if err != nil {
		t.Fatalf("unexpected outgoing webhook creation error: %v", err)
	}

	if outgoingWh.IntegrationID != intg.ID || outgoingWh.Type != WebhookOutgoing || outgoingWh.TargetURL != targetURL {
		t.Errorf("outgoing webhook fields mismatch")
	}
}

func TestRolesTagsAndEmojis(t *testing.T) {
	scopeID := uuid.New().String()

	// Test Role Creation (Global Scope)
	role, err := NewRole(scopeID, ScopeGlobal, "Workspace Admin", "#FF0000", `{"can_delete_messages": true}`, 1)
	if err != nil {
		t.Fatalf("unexpected role creation error: %v", err)
	}
	if role.ScopeID != scopeID || role.ScopeType != ScopeGlobal || role.Name != "Workspace Admin" || role.Color != "#FF0000" || role.Order != 1 {
		t.Errorf("role fields mismatch")
	}

	// Test Custom Emoji Creation (Team Scope)
	emoji, err := NewCustomEmoji(scopeID, ScopeTeam, "party_parrot", "https://assets.com/parrot.gif")
	if err != nil {
		t.Fatalf("unexpected emoji creation error: %v", err)
	}
	if emoji.ScopeID != scopeID || emoji.ScopeType != ScopeTeam || emoji.Name != "party_parrot" || emoji.ImageURL != "https://assets.com/parrot.gif" {
		t.Errorf("custom emoji fields mismatch")
	}

	// Test Tag Creation (Team Scope specific)
	tag, err := NewTag(scopeID, "Urgent", "#FF5500")
	if err != nil {
		t.Fatalf("unexpected tag creation error: %v", err)
	}
	if tag.TeamScopeID != scopeID || tag.Name != "Urgent" || tag.Color != "#FF5500" {
		t.Errorf("tag fields mismatch")
	}
}


