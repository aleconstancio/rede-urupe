package gateway

import (
	"context"
	"time"

	"github.com/aleconstancio/talos/v2/engine/core/agent"
	"nucleo-engine/internal/data/sqlite"
	"nucleo-engine/internal/pkg/timeutil"
)

// MazePulseDetector implements agent.PulseDetector for Discord channels.
type MazePulseDetector struct {
	Repo *sqlite.Repository
}

func (d *MazePulseDetector) ShouldProcess(ctx context.Context, input *agent.PipelineInput) (bool, agent.EventType, string, error) {
	// Reactive events always process
	if input.EventType == "reactive" || string(input.EventType) == string(agent.EventReactive) {
		return true, agent.EventReactive, "direct_mention", nil
	}

	// Check passive mode
	isPassive := d.Repo.GetConfig("passive_mode", "true") == "true"
	if isPassive {
		return false, agent.EventPassive, "passive_mode", nil
	}

	// Get gateway state
	state, err := d.Repo.GetGatewayState(input.ChannelID)
	if err != nil {
		return false, agent.EventSkipped, "state_error", err
	}

	now := timeutil.Now()

	// 3-minute cooldown
	if !state.LastPulseTimestamp.IsZero() && now.Sub(state.LastPulseTimestamp) < 3*time.Minute {
		return false, agent.EventSkipped, "cooldown", nil
	}

	// Volume pulse: 20 messages since last pulse
	msgCount, _ := d.Repo.CountMessagesSince(input.ChannelID, state.LastPulseTimestamp)
	if msgCount >= 20 {
		return true, agent.EventAmbient, "volume_pulse", nil
	}

	// Time pulse: 1 hour since last intervention
	if !state.LastPulseTimestamp.IsZero() && now.Sub(state.LastPulseTimestamp) > 1*time.Hour {
		return true, agent.EventAmbient, "time_pulse", nil
	}

	return false, agent.EventSkipped, "below_threshold", nil
}
