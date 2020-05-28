package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AuraBlast ...
func AuraBlast(c *match.Card) {

	c.Name = "Aura Blast"
	c.Civ = civ.Nature
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			card.Player.MapContainer(match.BATTLEZONE, func(creature *match.Card) {
				creature.AddCondition(cnd.PowerAttacker, 2000, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was given power attacker +2000", creature.Name))
			})

		}

	})

}

// BrainSerum ...
func BrainSerum(c *match.Card) {

	c.Name = "Brain Serum"
	c.Civ = civ.Water
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Spell, fx.ShieldTrigger, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			card.Player.DrawCards(2)

		}

	})

}

// BurningPower ...
func BurningPower(c *match.Card) {

	c.Name = "Burning Power"
	c.Civ = civ.Fire
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Search(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Select 1 creature from your battlezone that will gain \"Power Attacker +2000\"", 1, 1, false)

			for _, creature := range creatures {

				creature.AddCondition(cnd.PowerAttacker, 2000, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was given power attacker +2000", creature.Name))

			}

		}

	})

}

// ChaosStrike ...
func ChaosStrike(c *match.Card) {

	c.Name = "Chaos Strike"
	c.Civ = civ.Fire
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmICasted(card, ctx) {

			creatures := match.Search(card.Player, ctx.Match, ctx.Match.Opponent(card.Player), match.BATTLEZONE, "Select 1 of your opponent's creatures that will be tapped", 1, 1, false)

			for _, creature := range creatures {

				creature.Tapped = true
				ctx.Match.Chat("Server", fmt.Sprintf("%s was tapped by %s", creature.Name, card.Name))

			}

		}

	})

}