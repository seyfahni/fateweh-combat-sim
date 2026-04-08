package weapon

import (
	"github.com/seyfahni/fateweh-combat-sim/simulator/dice"
	"github.com/seyfahni/fateweh-combat-sim/simulator/log"
)

type Unarmed struct {
	Modifier int
}

func (u Unarmed) RollDamage(dice dice.D6) (int, log.Simulation) {
	details := make([]log.Simulation, 0, 1)

	damage := u.Modifier
	result := dice.RollD6()
	details = append(details, log.MessageF("rolled %d with dice", result))
	damage += result
	for result >= 6 {
		result = dice.RollD6()
		details = append(details, log.MessageF("rolled %d with exploding dice", result))
		damage += result
	}

	damage = damage / 2
	details = append(details, log.Message("damage halved due to unarmed attack"))

	return damage, log.MessageF("attacking unarmed with modifier %+d, dealing %d damage", u.Modifier, damage).AndDetails(details...)
}
