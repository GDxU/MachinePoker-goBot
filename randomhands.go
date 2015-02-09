package main

import (
	"fmt"
    "log"
	"net/http"
	"math/rand"
	"github.com/loganjspears/joker/hand"
)

var BotName = "GOd of Gamblers"

func main() {
	http.HandleFunc("/bot/gog", botHandler)
	http.ListenAndServe("0.0.0.0:8081", nil)
}

func botHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
        case "GET":
            fmt.Fprintf(w, "{\"info\": { \"name\": \"%s\" } }", BotName)
        case "POST":
            game := ReadGame(r.Body)
            DisplayGame(game)

            var bet int
            if game.State != "complete" {
                bet = play(game)
            }
            fmt.Fprintf(w, "{\"bet\": \"%d\"}", bet)
        default:
            log.Fatal("Method unsupported")
    }
}

func play(game *Game) int {
	var ret int

    // consider all cards to calculate odds
    all := append(game.Community, game.Self.Cards...)
	myCards := Cards(all)

    // convert to joker hand and calculate ranking
    myHand := hand.New(myCards)
    fmt.Println("\n** Hand: ")
    fmt.Println(myHand)

    // TODO: printed value of rank is wrong, subract 1
    fmt.Printf("ranking: %s\n", myHand.Ranking()-1)

    // strategy
	if game.State == "pre-flop" {
		if myHand.Ranking() == hand.Pair {
			ret = raise(game)
		} else {
			ret = rand.Intn(2) * game.Betting.Call
            fmt.Println("-> returning:", ret)
		}
	} else {
        if myHand.Ranking() >= hand.ThreeOfAKind {
            ret = raise(game)
        } else if myHand.Ranking() == hand.Pair {
            fmt.Println("-> calling:", game.Betting.Call)
            ret = game.Betting.Call
        }
	}
	return ret
}

func raise(game *Game) int {
    fmt.Println("-> canRaise:", game.Betting.CanRaise)
	if game.Betting.CanRaise {
        fmt.Println("-> raising:", game.Betting.Raise)
		return game.Betting.Raise
	} else {
        fmt.Println("-> calling:", game.Betting.Call)
		return game.Betting.Call
	}
}
