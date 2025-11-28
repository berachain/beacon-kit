// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blobreactor

import (
	"sync"
	"time"

	"github.com/berachain/beacon-kit/log"
	"github.com/cometbft/cometbft/p2p"
)

const (
	// defaultMaxReputationScore is the starting reputation score for all peers
	defaultMaxReputationScore = 10
	// defaultDisconnectThreshold is the reputation score below or at which peers are banned from blob protocol
	defaultDisconnectThreshold = 0
	// defaultBadBehaviorPenalty is the reputation points deducted per protocol violation
	defaultBadBehaviorPenalty = 2
	// defaultGoodBehaviorReward is the reputation points added for successfully providing valid blobs
	defaultGoodBehaviorReward = 1
	// defaultBanPeriod is how long peers are banned when reputation drops below threshold
	defaultBanPeriod = 2 * time.Hour
	// defaultStaleReputationTimeout determines how long to keep inactive peer reputations before deleting them
	defaultStaleReputationTimeout = 24 * time.Hour
)

// ReputationConfig contains reputation system parameters
type ReputationConfig struct {
	MaxReputationScore  int           // Starting reputation score for all peers
	DisconnectThreshold int           // Reputation score below which peers are banned from blob protocol
	BadBehaviorPenalty  int           // Reputation points deducted per protocol violation
	GoodBehaviorReward  int           // Reputation points added for successfully providing valid blobs
	BanPeriod           time.Duration // How long peers are banned when reputation drops below threshold
}

// DefaultReputationConfig returns default reputation configuration
func DefaultReputationConfig() ReputationConfig {
	return ReputationConfig{
		MaxReputationScore:  defaultMaxReputationScore,
		DisconnectThreshold: defaultDisconnectThreshold,
		BadBehaviorPenalty:  defaultBadBehaviorPenalty,
		GoodBehaviorReward:  defaultGoodBehaviorReward,
		BanPeriod:           defaultBanPeriod,
	}
}

// WithDefaults returns a new ReputationConfig with zero values replaced by defaults
func (c ReputationConfig) WithDefaults() ReputationConfig {
	defaults := DefaultReputationConfig()

	if c.MaxReputationScore == 0 {
		c.MaxReputationScore = defaults.MaxReputationScore
	}
	if c.DisconnectThreshold == 0 {
		c.DisconnectThreshold = defaults.DisconnectThreshold
	}
	if c.BadBehaviorPenalty == 0 {
		c.BadBehaviorPenalty = defaults.BadBehaviorPenalty
	}
	if c.GoodBehaviorReward == 0 {
		c.GoodBehaviorReward = defaults.GoodBehaviorReward
	}
	if c.BanPeriod == 0 {
		c.BanPeriod = defaults.BanPeriod
	}

	return c
}

type peerReputation struct {
	totalScore      int
	disconnectCount int
	bannedUntil     time.Time
	lastSeen        time.Time
}

// ReputationManager tracks peer behavior and maintains reputation scores that persist
// across disconnections.
//
// Peers start with a maximum score. Scores decrease on bad behaviour and increase for
// good behavior. When a peer's score drops below the configured threshold, they are
// ignored at the reactor level and banned for the configured period.
//
// Peers remain connected to CometBFT for consensus but are excluded from blob protocol
// operations. Reputation entries persist in memory across disconnections, allowing peers
// to resume with their previous score. A background cleanup process removes stale entries
// after the configured timeout period.
type ReputationManager struct {
	reputations map[p2p.ID]*peerReputation
	mu          sync.RWMutex
	config      ReputationConfig
	logger      log.Logger
}

func NewReputationManager(logger log.Logger, config ReputationConfig) *ReputationManager {
	return &ReputationManager{
		reputations: make(map[p2p.ID]*peerReputation),
		config:      config,
		logger:      logger,
	}
}

func (rm *ReputationManager) RecordBadBehavior(peerID p2p.ID, err error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rep, exists := rm.reputations[peerID]
	if !exists {
		rep = &peerReputation{
			totalScore: rm.config.MaxReputationScore,
		}
		rm.reputations[peerID] = rep
	}

	rep.lastSeen = time.Now()
	rep.totalScore = max(0, rep.totalScore-rm.config.BadBehaviorPenalty)

	rm.logger.Warn("Peer reputation violation",
		"peer", peerID,
		"penalty", rm.config.BadBehaviorPenalty,
		"new_score", rep.totalScore,
		"error", err)

	if rep.totalScore <= rm.config.DisconnectThreshold {
		rep.disconnectCount++
		rep.bannedUntil = time.Now().Add(rm.config.BanPeriod)
		rm.logger.Warn("Peer reputation below threshold, banning from blob protocol",
			"peer", peerID,
			"score", rep.totalScore,
			"ban_count", rep.disconnectCount,
			"banned_until", rep.bannedUntil)
	}
}

func (rm *ReputationManager) RecordGoodBehavior(peerID p2p.ID) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rep, exists := rm.reputations[peerID]
	if !exists {
		rep = &peerReputation{
			totalScore: rm.config.MaxReputationScore,
		}
		rm.reputations[peerID] = rep
	}

	rep.lastSeen = time.Now()
	rep.totalScore = min(rep.totalScore+rm.config.GoodBehaviorReward, rm.config.MaxReputationScore)
}

func (rm *ReputationManager) ShouldAcceptPeer(peerID p2p.ID) bool {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rep, exists := rm.reputations[peerID]
	if !exists {
		return true
	}

	rep.lastSeen = time.Now()

	// If not banned, accept
	if rep.bannedUntil.IsZero() {
		return true
	}

	// If ban has expired, clear it, reset score, and accept
	if time.Now().After(rep.bannedUntil) {
		rep.bannedUntil = time.Time{}
		rep.totalScore = rm.config.MaxReputationScore
		rm.logger.Info("Peer ban expired, resetting reputation", "peer", peerID, "new_score", rep.totalScore)
		return true
	}

	// Still banned
	return false
}

// CleanupStaleReputations removes reputation entries for peers that haven't been seen recently.
// Banned peers are kept until their ban expires.
func (rm *ReputationManager) CleanupStaleReputations() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	now := time.Now()
	for peerID, rep := range rm.reputations {
		// Keep banned peers until ban expires
		if !rep.bannedUntil.IsZero() && now.Before(rep.bannedUntil) {
			continue
		}

		// Remove stale entries
		if now.Sub(rep.lastSeen) > defaultStaleReputationTimeout {
			delete(rm.reputations, peerID)
		}
	}
}
