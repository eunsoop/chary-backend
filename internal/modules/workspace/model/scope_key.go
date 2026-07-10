package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type ScopeKey struct {
	WorkspaceID string
	IsGlobal    bool
	TeamID      string
	CategoryID  string
	ChannelID   string
}

func BuildGlobalScopeKey(workspaceID, categoryID, channelID string) (string, error) {
	if _, err := uuid.Parse(workspaceID); err != nil {
		return "", fmt.Errorf("invalid workspace ID: %w", err)
	}
	if _, err := uuid.Parse(categoryID); err != nil {
		return "", fmt.Errorf("invalid category ID: %w", err)
	}
	if _, err := uuid.Parse(channelID); err != nil {
		return "", fmt.Errorf("invalid channel ID: %w", err)
	}
	return fmt.Sprintf("workspace:%s:global.category:%s.channel:%s", workspaceID, categoryID, channelID), nil
}

func BuildTeamScopeKey(workspaceID, teamID, categoryID, channelID string) (string, error) {
	if _, err := uuid.Parse(workspaceID); err != nil {
		return "", fmt.Errorf("invalid workspace ID: %w", err)
	}
	if _, err := uuid.Parse(teamID); err != nil {
		return "", fmt.Errorf("invalid team ID: %w", err)
	}
	if _, err := uuid.Parse(categoryID); err != nil {
		return "", fmt.Errorf("invalid category ID: %w", err)
	}
	if _, err := uuid.Parse(channelID); err != nil {
		return "", fmt.Errorf("invalid channel ID: %w", err)
	}
	return fmt.Sprintf("workspace:%s:team:%s.category:%s.channel:%s", workspaceID, teamID, categoryID, channelID), nil
}

func ParseScopeKey(key string) (*ScopeKey, error) {
	if !strings.HasPrefix(key, "workspace:") {
		return nil, errors.New("invalid scope key format: must start with 'workspace:'")
	}

	parts := strings.SplitN(key, ":", 3)
	if len(parts) < 3 {
		return nil, errors.New("invalid scope key format: missing workspace ID or scope type")
	}

	workspaceID := parts[1]
	if _, err := uuid.Parse(workspaceID); err != nil {
		return nil, fmt.Errorf("invalid workspace ID: %w", err)
	}

	remaining := parts[2]

	var isGlobal bool
	var teamID string
	var categoryID string
	var channelID string

	if strings.HasPrefix(remaining, "global.category:") {
		isGlobal = true
		subParts := strings.Split(remaining, ".")
		if len(subParts) != 3 {
			return nil, errors.New("invalid global scope key format: expected 3 dot-separated segments")
		}

		catPart := subParts[1]
		if !strings.HasPrefix(catPart, "category:") {
			return nil, errors.New("invalid category prefix in scope key")
		}
		categoryID = strings.TrimPrefix(catPart, "category:")

		chPart := subParts[2]
		if !strings.HasPrefix(chPart, "channel:") {
			return nil, errors.New("invalid channel prefix in scope key")
		}
		channelID = strings.TrimPrefix(chPart, "channel:")

	} else if strings.HasPrefix(remaining, "team:") {
		isGlobal = false
		teamAndRemaining := strings.SplitN(remaining, ".", 2)
		if len(teamAndRemaining) != 2 {
			return nil, errors.New("invalid team scope key format: missing dot separator after team segment")
		}

		teamPart := teamAndRemaining[0]
		if !strings.HasPrefix(teamPart, "team:") {
			return nil, errors.New("invalid team prefix in scope key")
		}
		teamID = strings.TrimPrefix(teamPart, "team:")

		subParts := strings.Split(teamAndRemaining[1], ".")
		if len(subParts) != 2 {
			return nil, errors.New("invalid team scope key sub-components: expected category and channel segments")
		}

		catPart := subParts[0]
		if !strings.HasPrefix(catPart, "category:") {
			return nil, errors.New("invalid category prefix in scope key")
		}
		categoryID = strings.TrimPrefix(catPart, "category:")

		chPart := subParts[1]
		if !strings.HasPrefix(chPart, "channel:") {
			return nil, errors.New("invalid channel prefix in scope key")
		}
		channelID = strings.TrimPrefix(chPart, "channel:")
	} else {
		return nil, errors.New("invalid scope key type: must be 'global' or 'team'")
	}

	if _, err := uuid.Parse(categoryID); err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}
	if _, err := uuid.Parse(channelID); err != nil {
		return nil, fmt.Errorf("invalid channel ID: %w", err)
	}
	if !isGlobal {
		if _, err := uuid.Parse(teamID); err != nil {
			return nil, fmt.Errorf("invalid team ID: %w", err)
		}
	}

	return &ScopeKey{
		WorkspaceID: workspaceID,
		IsGlobal:    isGlobal,
		TeamID:      teamID,
		CategoryID:  categoryID,
		ChannelID:   channelID,
	}, nil
}
