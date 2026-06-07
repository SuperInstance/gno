package gnolang

import "fmt"

// VetoTier defines the SAEP governance hierarchy.
type VetoTier int

const (
	VetoNone    VetoTier = 0
	VetoRoom    VetoTier = 1
	VetoSector  VetoTier = 2
	VetoPort    VetoTier = 3
	VetoMarket  VetoTier = 4
)

type VetoSignal struct {
	Tier    VetoTier
	Reason  string
	IsVetoed bool
}

// VetoEngine manages the 4-tier governance check.
type VetoEngine struct {
	CurrentTier VetoTier
}

func NewVetoEngine() *VetoEngine {
	return &VetoEngine{CurrentTier: VetoNone}
}

// Check evaluates a proposed state transition against the hierarchy.
func (ve *VetoEngine) Check(proposal string, tier VetoTier, reason string) VetoSignal {
	if tier > ve.CurrentTier {
		ve.CurrentTier = tier
	}
	
	if tier != VetoNone {
		return VetoSignal{
			Tier:      tier,
			Reason:    reason,
			IsVetoed:  true,
		}
	}
	
	return VetoSignal{Tier: VetoNone, IsVetoed: false}
}

// WrapExecutor is a conceptual wrapper for the gnovm executor.
func (ve *VetoEngine) WrapExecutor(execFunc func() error) error {
	// Pre-execution veto check
	signal := ve.Check("START_EXEC", VetoNone, "")
	if signal.IsVetoed {
		return fmt.Errorf("execution vetoed at tier %d: %s", signal.Tier, signal.Reason)
	}

	return execFunc()
}
