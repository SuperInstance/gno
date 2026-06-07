package gnolang

import (
	"fmt"
	"sync"
)

// FluxBytecode defines the native A2A opcodes for the Flux VM.
type FluxBytecode byte

const (
	OP_SPAWN_VESSEL FluxBytecode = iota // Create a new agent execution context
	OP_PASS_BATON                      // Teleport state (sharded bottle) between vessels
	OP_SPLICING                        // Merge threads based on topological symmetry
)

// Baton defines the 3-way shard of state transferred between agents.
type Baton struct {
	Artifacts  map[string][]byte // Raw data/output
	Reasoning  string            // Causal trace / "The Why"
	Blockers   []string           // Current impediments
	SymmetryID uint64            // Topological orbit identifier
}

// Vessel represents an autonomous execution context within the Flux VM.
type Vessel struct {
	ID          string
	State       map[string]interface{}
	Mu          sync.RWMutex
	IsActive    bool
	SymmetryID  uint64
}

// FluxVM extends the standard gno VM with A2A native abilities.
type FluxVM struct {
	Vessels map[string]*Vessel
	Mu      sync.RWMutex
}

func NewFluxVM() *FluxVM {
	return &FluxVM{
		Vessels: make(map[string]*Vessel),
	}
}

// SpawnVessel implements OP_SPAWN_VESSEL.
func (vm *FluxVM) SpawnVessel(id string) (*Vessel, error) {
	vm.Mu.Lock()
	defer vm.Mu.Unlock()

	if _, exists := vm.Vessels[id]; exists {
		return nil, fmt.Errorf("vessel %s already exists", id)
	}

	v := &Vessel{
		ID:       id,
		State:    make(map[string]interface{}),
		IsActive: true,
	}
	vm.Vessels[id] = v
	return v, nil
}

// PassBaton implements OP_PASS_BATON.
func (vm *FluxVM) PassBaton(fromID, toID string, b Baton) error {
	vm.Mu.RLock()
	from, ok1 := vm.Vessels[fromID]
	to, ok2 := vm.Vessels[toID]
	vm.Mu.RUnlock()

	if !ok1 || !ok2 {
		return fmt.Errorf("one or both vessels (%s, %s) not found", fromID, toID)
	}

	from.Mu.Lock()
	defer from.Mu.Unlock()
	to.Mu.Lock()
	defer to.Mu.Unlock()

	// Transfer logic: a-priori delivery
	to.State["last_baton"] = b
	to.SymmetryID = b.SymmetryID
	
	return nil
}

// Splicing implements OP_SPLICING by merging state of symmetric vessels.
func (vm *FluxVM) Splicing(idA, idB string) error {
	vm.Mu.RLock()
	vA, ok1 := vm.Vessels[idA]
	vB, ok2 := vm.Vessels[idB]
	vm.Mu.RUnlock()

	if !ok1 || !ok2 {
		return fmt.Errorf("vessels not found")
	}

	vA.Mu.Lock()
	defer vA.Mu.Unlock()
	vB.Mu.Lock()
	defer vB.Mu.Unlock()

	// Merge logic: union of states based on symmetry
	for k, v := range vB.State {
		vA.State[k] = v
	}
	vB.IsActive = false

	return nil
}
