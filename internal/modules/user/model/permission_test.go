package model

import (
	"testing"
)

func TestPermissionInheritanceAndEvaluation(t *testing.T) {
	wsID := "01909a3c-d3c2-70b9-8ce7-459de4db6b78"
	teamID := "01909a3c-d3c2-70b9-8ce7-459de4db6b79"
	catID := "01909a3c-d3c2-70b9-8ce7-459de4db6b7a"
	chID := "01909a3c-d3c2-70b9-8ce7-459de4db6b7b"

	globalScopeKey := "workspace:" + wsID + ":global"
	teamScopeKey := "workspace:" + wsID + ":team:" + teamID
	teamCategoryKey := teamScopeKey + ".category:" + catID
	teamChannelKey := teamCategoryKey + ".channel:" + chID

	user, err := NewUser("jane_doe", "jane@doe.com", "Jane Doe", "", false)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// 1. Initial State: No permissions anywhere
	if user.CanInScope(PermMessageWrite, teamChannelKey) {
		t.Error("expected access to be denied initially")
	}

	// 2. Direct Channel-level permission grant
	// PermMessageWrite is Block: 2, Mask: 1 << 0
	// Initialize blocks: len 3 to cover block 2
	user.ScopeMap[teamChannelKey] = []uint64{0, 0, 1 << 0}
	if !user.CanInScope(PermMessageWrite, teamChannelKey) {
		t.Error("expected access granted at channel level")
	}

	// MessageWrite should not inherit upward (should not have category/team write)
	if user.CanInScope(PermMessageWrite, teamCategoryKey) {
		t.Error("expected access denied at parent category level since grant is channel-specific")
	}

	// 3. Inheritance Test: Grant permission at Category level, verify Channel inherits it
	delete(user.ScopeMap, teamChannelKey) // clear direct channel permission
	user.ScopeMap[teamCategoryKey] = []uint64{0, 0, 1 << 0}

	if !user.CanInScope(PermMessageWrite, teamChannelKey) {
		t.Error("expected channel to inherit PermMessageWrite from category")
	}

	// 4. Inheritance Test: Grant permission at Team Scope level, verify Channel inherits it
	delete(user.ScopeMap, teamCategoryKey) // clear category permission
	user.ScopeMap[teamScopeKey] = []uint64{0, 0, 1 << 0}

	if !user.CanInScope(PermMessageWrite, teamChannelKey) {
		t.Error("expected channel to inherit PermMessageWrite from team scope")
	}

	// 5. Global Override Test: Grant global role management at Global workspace scope
	// This acts as a master super-admin override for any resources under this workspace
	delete(user.ScopeMap, teamScopeKey) // clear team scope permission
	user.ScopeMap[globalScopeKey] = []uint64{1 << 0, 0, 0} // PermRoleManage is Block 0, Mask 1 << 0

	// Verify the user can perform messages, category management, etc. anywhere in the workspace
	if !user.CanInScope(PermMessageWrite, teamChannelKey) {
		t.Error("expected global super-admin override to grant PermMessageWrite in team channel")
	}
	if !user.CanInScope(PermCategoryManage, teamCategoryKey) {
		t.Error("expected global super-admin override to grant PermCategoryManage")
	}

	// Verify it does NOT override permissions on a different workspace
	otherWorkspaceChannel := "workspace:00000000-0000-0000-0000-000000000000:global.category:cat.channel:ch"
	if user.CanInScope(PermMessageWrite, otherWorkspaceChannel) {
		t.Error("expected super-admin override of workspace A to NOT grant permissions in workspace B")
	}
}

func TestPermissionBoundsSafety(t *testing.T) {
	wsID := "01909a3c-d3c2-70b9-8ce7-459de4db6b78"
	channelKey := "workspace:" + wsID + ":global.category:cat.channel:ch"

	user, _ := NewUser("bounds_tester", "test@bounds.com", "Tester", "", false)

	// User has a shorter bitmask array (length 1), but we are checking a block 2 permission
	user.ScopeMap[channelKey] = []uint64{1} // Block 0 has bit, but blocks 1-3 do not exist in slice

	// This should return false and NOT panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("permission evaluation panicked: %v", r)
		}
	}()

	if user.CanInScope(PermMessageWrite, channelKey) {
		t.Error("expected false since block 2 is out of slice bounds")
	}
}
