package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/seyfahni/fateweh-combat-sim/simulator"
	"github.com/seyfahni/fateweh-combat-sim/simulator/dice"
	"github.com/seyfahni/fateweh-combat-sim/simulator/weapon"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := root.Execute(); err != nil {
		println(err)
		os.Exit(1)
	}
}

func init() {
	root.AddCommand(oneVsOne)
	root.PersistentFlags().Int64("seed", 0, "randomness seed")

	oneVsOne.Flags().Int("player-health", 15, "player health points")
	oneVsOne.Flags().String("player-weapon-type", "unarmed", "player weapon type (unarmed, melee)")
	oneVsOne.Flags().Int("player-damage-modifier", 0, "player damage modifier")
	oneVsOne.Flags().Int("player-weapon-dice", 1, "player weapon dices")

	oneVsOne.Flags().Int("enemy-health", 15, "enemy health points")
	oneVsOne.Flags().String("enemy-weapon-type", "unarmed", "enemy weapon type (unarmed, melee)")
	oneVsOne.Flags().Int("enemy-damage-modifier", 0, "enemy damage modifier")
	oneVsOne.Flags().Int("enemy-weapon-dice", 1, "enemy weapon dices")

	oneVsOne.Flags().Int("max-steps", 100, "maximum number of steps to simulate")
}

var root = &cobra.Command{
	Use:   "fateweh-combat-sim",
	Short: "Simulate FreeFate & Fernweh mashup combat simulator",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var oneVsOne = &cobra.Command{
	Use:   "one-vs-one",
	Short: "Simulate a 1v1 combat situation",
	RunE: func(cmd *cobra.Command, args []string) error {
		seed, err := cmd.Flags().GetInt64("seed")
		if err != nil {
			return err
		}
		if seed == 0 {
			seed = time.Now().UnixNano()
		}

		maxSteps, err := cmd.Flags().GetInt("max-steps")
		if err != nil {
			return err
		}

		player, err := configureParticipant(cmd, "player")
		if err != nil {
			return err
		}

		enemy, err := configureParticipant(cmd, "enemy")
		if err != nil {
			return err
		}

		world := simulator.World{
			Player: player,
			Enemy:  enemy,
		}

		rng := rand.New(rand.NewSource(seed))
		pouch := &dice.Pouch{Random: rng}

		result := simulator.Simulate(pouch, &simulator.PlayerTurn{
			World: world,
		}, maxSteps)

		err = result.PrintTo(ConsolePrinter{})
		if err != nil {
			return err
		}

		return nil
	},
}

func configureParticipant(cmd *cobra.Command, name string) (simulator.Participant, error) {
	h, err := cmd.Flags().GetInt(name + "-health")
	if err != nil {
		return simulator.Participant{}, err
	}

	w, err := configureWeapon(cmd, name)
	if err != nil {
		return simulator.Participant{}, err
	}

	return simulator.Participant{
		Name: name,
		Actions: []simulator.Action{
			&simulator.AttackAction{
				Weapon: w,
			},
		},
		Health: h,
	}, nil
}

var ErrUnrecognisedWeapon = errors.New("unrecognised weapon")

func configureWeapon(cmd *cobra.Command, name string) (simulator.Weapon, error) {
	weaponType, err := cmd.Flags().GetString(name + "-weapon-type")
	if err != nil {
		return nil, err
	}

	switch weaponType {
	case "unarmed":
		m, err := cmd.Flags().GetInt(name + "-damage-modifier")
		if err != nil {
			return nil, err
		}
		return weapon.Unarmed{
			Modifier: m,
		}, nil
	case "melee":
		m, err := cmd.Flags().GetInt(name + "-damage-modifier")
		if err != nil {
			return nil, err
		}
		d, err := cmd.Flags().GetInt(name + "-weapon-dice")
		if err != nil {
			return nil, err
		}
		return weapon.Melee{
			Dice:     d,
			Modifier: m,
		}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnrecognisedWeapon, weaponType)
	}
}
