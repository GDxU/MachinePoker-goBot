package main

import (
    "log"
    "github.com/loganjspears/joker/hand"
)

type PokerPlayer interface {
    Play(game *Game) int
}

type defaultPlayer struct { }

func (p *defaultPlayer) Play(game *Game) int {
	// consider all cards when calculating odds
	all := append(game.Community, game.Self.Cards...)
	allCards := Cards(all)

	// convert to joker hand and calculate ranking
	allHand := hand.New(allCards)
	logger.Println("** myHand:", allHand)

	switch game.State {
	case "pre-flop":
		return calculatePreflopBet(game, allHand)
	case "flop":
		return calculateBet(game, allHand)
	case "turn":
		return calculateBet(game, allHand)
	case "river":
		return calculateBet(game, allHand)
	default:
		log.Fatal("Undefined game state:", game.State)
		return -1
	}
}

func calculatePreflopBet(game *Game, allHand *hand.Hand) int {
	if allHand.Ranking() == hand.Pair {
		return raise(game)
	}
	return call(game)
}

func calculateBet(game *Game, allHand *hand.Hand) int {
	if safeguard(game, allHand) {
		myHand := hand.New(Cards(game.Self.Cards))
		if myHand.Ranking() == hand.Pair && allHand.Ranking() >= hand.TwoPair {
			return raise(game)
		} else if allHand.Ranking() >= hand.Pair || game.Self.Wagered > 50 {
			return call(game)
		}
	}
	return fold(game)
}

func raise(game *Game) int {
	if game.Betting.CanRaise {
		logger.Println("-> raising:", game.Betting.Raise)
		return game.Betting.Raise
	}
	return call(game)
}

func call(game *Game) int {
	logger.Println("-> calling:", game.Betting.Call)
	return game.Betting.Call
}

func fold(game *Game) int {
	logger.Println("-> folding")
	return 0
}

func safeguard(game *Game, allHand *hand.Hand) bool {
	if game.Betting.Call < 100 {
		return true
	}
	myHand := hand.New(Cards(game.Self.Cards))
	if myHand.Ranking() == hand.Pair && allHand.Ranking() >= hand.ThreeOfAKind {
		return true
	}
	return false
}
