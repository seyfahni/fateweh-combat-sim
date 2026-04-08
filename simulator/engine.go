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

type World struct {
	Player Participant
	Enemy  Participant
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

type PlayerTurn struct {
	World
}

func (t *PlayerTurn) Next(rng dice.D6) (SimulationStep, log.Simulation) {
	turnLog := log.Message("player's turn")
	details := make([]log.Simulation, 0, len(t.Player.Actions))

	for _, action := range t.Player.Actions {
		details = append(details, action.Execute(rng, &t.Enemy))
		if t.Enemy.Health <= 0 {
			details = append(details, log.Message("killed enemy"))
			return nil, turnLog.AndDetails(details...)
		}
	}
	return &EnemyTurn{
		World: t.World,
	}, turnLog.AndDetails(details...)
}

type EnemyTurn struct {
	World
}

func (t *EnemyTurn) Next(rng dice.D6) (SimulationStep, log.Simulation) {
	turnLog := log.Message("enemy's turn")
	details := make([]log.Simulation, 0, len(t.Player.Actions))

	for _, action := range t.Enemy.Actions {
		details = append(details, action.Execute(rng, &t.Player))
		if t.Player.Health <= 0 {
			details = append(details, log.Message("killed player"))
			return nil, turnLog.AndDetails(details...)
		}
	}
	return &PlayerTurn{
		World: t.World,
	}, turnLog.AndDetails(details...)
}
