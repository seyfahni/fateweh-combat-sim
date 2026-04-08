package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/seyfahni/fateweh-combat-sim/simulator"
	"github.com/seyfahni/fateweh-combat-sim/simulator/dice"
	"github.com/seyfahni/fateweh-combat-sim/simulator/weapon"
)

type ConsolePrinter struct{}

func (p ConsolePrinter) Print(line string) error {
	_, err := fmt.Println(line)
	return err
}

func main() {
	world := simulator.World{
		Player: simulator.Participant{
			Name: "player",
			Actions: []simulator.Action{
				&simulator.AttackAction{
					Weapon: weapon.Melee{
						Dice:     1,
						Modifier: 3,
					},
				},
			},
			Health: 25,
		},
		Enemy: simulator.Participant{
			Name: "enemy",
			Actions: []simulator.Action{
				&simulator.AttackAction{
					Weapon: weapon.Melee{
						Dice:     2,
						Modifier: 1,
					},
				},
			},
			Health: 23,
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	pouch := &dice.Pouch{Random: rng}

	result := simulator.Simulate(pouch, &simulator.PlayerTurn{
		World: world,
	}, 100)

	err := result.PrintTo(ConsolePrinter{})
	if err != nil {
		panic(err)
	}
}
