package simulator

import (
	"github.com/seyfahni/fateweh-combat-sim/simulator/dice"
	"github.com/seyfahni/fateweh-combat-sim/simulator/log"
)

type Participant struct {
	Name    string
	Actions []Action
	Health  int
}

type Action interface {
	Execute(rng dice.D6, target *Participant) log.Simulation
}

type AttackAction struct {
	Weapon
}

func (a *AttackAction) Execute(rng dice.D6, target *Participant) log.Simulation {
	damage, damageLog := a.RollDamage(rng)
	previousHealth := target.Health
	target.Health = max(target.Health-damage, 0)

	return log.MessageF("attacking %s", target.Name).AndDetails(
		damageLog,
		log.MessageF("dealt %d damage, reducing health from %d to %d", damage, previousHealth, target.Health),
	)
}

type Weapon interface {
	RollDamage(dice.D6) (int, log.Simulation)
}

func Simulate(rng dice.D6, state SimulationStep, maxSteps int) log.Simulation {
	history := make(log.Group, 0, 16)
	for steps := 0; state != nil && steps < maxSteps; steps++ {
		var step log.Simulation
		state, step = state.Next(rng)
		history = append(history, step)
	}
	if state == nil {
		history = append(history, log.Message("simulation ended"))
	} else {
		history = append(history, log.Message("max steps reached"))
	}
	return history
}

type SimulationStep interface {
	Next(rng dice.D6) (SimulationStep, log.Simulation)
}

type Turn struct {
	Self   Participant
	Target Participant
}

func (t *Turn) Next(rng dice.D6) (SimulationStep, log.Simulation) {
	turnLog := log.MessageF("%s's turn", t.Self.Name)
	details := make([]log.Simulation, 0, len(t.Self.Actions))

	for _, action := range t.Self.Actions {
		details = append(details, action.Execute(rng, &t.Target))
		if t.Target.Health <= 0 {
			details = append(details, log.MessageF("killed %s", t.Target.Name))
			return nil, turnLog.AndDetails(details...)
		}
	}
	return &Turn{
		Self:   t.Target,
		Target: t.Self,
	}, turnLog.AndDetails(details...)
}
