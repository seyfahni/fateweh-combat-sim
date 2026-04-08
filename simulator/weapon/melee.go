package weapon

import (
	"github.com/seyfahni/fateweh-combat-sim/simulator/dice"
	"github.com/seyfahni/fateweh-combat-sim/simulator/log"
)

type Melee struct {
	Dice     int
	Modifier int
}

func (w Melee) RollDamage(dice dice.D6) (int, log.Simulation) {
	details := make([]log.Simulation, 0, w.Dice)

	ones := 0
	critical := w.Dice > 0

	damage := w.Modifier
	for i := 0; i < w.Dice; i++ {
		result := dice.RollD6()
		details = append(details, log.MessageF("rolled %d with dice %d", result, i+1))
		if result < 6 {
			critical = false
			if result == 1 {
				ones++
				if ones >= 2 {
					return 0, log.MessageF("attacking with %dD6%+d, missed", w.Dice, w.Modifier).AndDetails(details...)
				}
			}
		}
		damage += result
	}
	if critical {
		details = append(details, log.Message("rolled critical, rerolling exploding dice"))
		for result := 6; result >= 6; {
			result = dice.RollD6()
			details = append(details, log.MessageF("rolled %d with exploding dice", result))
			damage += result
		}
	}

	return damage, log.MessageF("attacking with %dD6%+d, dealing %d damage", w.Dice, w.Modifier, damage).AndDetails(details...)
}
