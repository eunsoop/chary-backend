package model

import (
	"strings"
)

type Permission struct {
	Block int
	Mask  uint64
	Name  string
}

// Block 0: Core Workspace Administration
var (
	PermRoleManage  = Permission{Block: 0, Mask: 1 << 0, Name: "PermRoleManage"}
	PermEmojiManage = Permission{Block: 0, Mask: 1 << 1, Name: "PermEmojiManage"}
	PermTagManage   = Permission{Block: 0, Mask: 1 << 2, Name: "PermTagManage"} // Exclusive to Team scope
)

// Block 1: Container & Channel Management
var (
	PermCategoryManage    = Permission{Block: 1, Mask: 1 << 0, Name: "PermCategoryManage"}
	PermChannelCreate     = Permission{Block: 1, Mask: 1 << 1, Name: "PermChannelCreate"}
	PermChannelDelete     = Permission{Block: 1, Mask: 1 << 2, Name: "PermChannelDelete"}
	PermIntegrationManage = Permission{Block: 1, Mask: 1 << 3, Name: "PermIntegrationManage"}
)

// Block 2: Symmetrical Message & Collaboration
var (
	PermMessageWrite     = Permission{Block: 2, Mask: 1 << 0, Name: "PermMessageWrite"}
	PermMessageDeleteAny = Permission{Block: 2, Mask: 1 << 1, Name: "PermMessageDeleteAny"}
	PermExpressionAdd    = Permission{Block: 2, Mask: 1 << 2, Name: "PermExpressionAdd"}
)

// Block 3: Compliance & Security Audit
var (
	PermDataExport   = Permission{Block: 3, Mask: 1 << 0, Name: "PermDataExport"}
	PermAuditLogView = Permission{Block: 3, Mask: 1 << 1, Name: "PermAuditLogView"}
)

// CanInScope checks if the user has the requested permission in the target scope.
// It first evaluates global overrides and then traverses up the scope key hierarchy.
func (u *User) CanInScope(perm Permission, targetScope string) bool {
	// 1. Evaluate top-level Global Master Bit mask override first
	if u.hasGlobalOverride(perm, targetScope) {
		return true
	}

	// 2. Traversal up the resource tree via string slicing (Zero-Allocation)
	for currentScope := targetScope; currentScope != ""; currentScope = getParentScope(currentScope) {
		if bits, exists := u.ScopeMap[currentScope]; exists {
			// Bounds check to avoid slice index panic on dynamic/smaller bitmask arrays
			if perm.Block >= 0 && perm.Block < len(bits) {
				if (bits[perm.Block] & perm.Mask) != 0 {
					return true // Access Granted
				}
			}
		}
	}
	return false // Access Denied
}

// hasGlobalOverride checks if the user has the permission or super-admin role at the workspace root level.
func (u *User) hasGlobalOverride(perm Permission, targetScope string) bool {
	if !strings.HasPrefix(targetScope, "workspace:") {
		return false
	}

	parts := strings.SplitN(targetScope, ":", 3)
	if len(parts) < 2 {
		return false
	}
	workspaceID := parts[1]

	// Slice off any trailing sub-elements of workspaceID to isolate the workspace UUID
	if idx := strings.Index(workspaceID, "."); idx != -1 {
		workspaceID = workspaceID[:idx]
	}
	if idx := strings.Index(workspaceID, ":"); idx != -1 {
		workspaceID = workspaceID[:idx]
	}

	globalScopeKey := "workspace:" + workspaceID + ":global"
	if bits, exists := u.ScopeMap[globalScopeKey]; exists {
		// Override 1: User explicitly has target permission at global scope level
		if perm.Block >= 0 && perm.Block < len(bits) {
			if (bits[perm.Block] & perm.Mask) != 0 {
				return true
			}
		}

		// Override 2: User is a Global Administrator (has PermRoleManage at global scope)
		if PermRoleManage.Block >= 0 && PermRoleManage.Block < len(bits) {
			if (bits[PermRoleManage.Block] & PermRoleManage.Mask) != 0 {
				return true
			}
		}
	}

	return false
}

// getParentScope traverses up one level in the dot-separated scope path.
func getParentScope(scope string) string {
	idx := strings.LastIndex(scope, ".")
	if idx == -1 {
		return ""
	}
	return scope[:idx]
}
