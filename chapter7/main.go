package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

var help string = `PIG
=================================================================================================================
The rules of the game are as follows:
* At the beginning, a coin will be flipped to determine if the human or computer goes first.
* Each player will roll dice until the end of their turn, adding the sum of the two dice to their total score.
* If a player does not roll any 1s, that player's will add the sum of the two dice to their score, and continue.
* If a player rolls one 1, that player's turn is over, and they forfeit the score they earned for that turn.
* If a player rolls two 1s, that player's turn is over, and they forfeit their total score.
* Either player can choose to end their turn at any time.
* The first player to reach a total score of 100 wins.
`
var randomizer *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

var highScores map[string]int

func init() {
	if jsonFile, errOpen := os.Open("highscores.json"); errOpen == nil {
		defer jsonFile.Close()
		jsonDecoder := json.NewDecoder(jsonFile)
		if errDecode := jsonDecoder.Decode(&highScores); errDecode != nil {
			panic(errDecode)
		}
	} else {
		highScores = make(map[string]int)
	}
}

func main() {
	fmt.Println(help)

	var playerName string
	if errAskPlayerName := survey.AskOne(
		&survey.Input{
			Message: "Please enter your name:",
		},
		&playerName,
	); errAskPlayerName != nil {
		panic(errAskPlayerName)
	}

	var humanTotalScore, computerTotalScore int
	var humanTurnScore, computerTurnScore int
	turn := randomizer.Intn(2)

	for humanTotalScore+humanTurnScore < 100 && computerTotalScore+computerTurnScore < 100 {
		if turn%2 == 0 {
			fmt.Println("\033[1m\033[32mYour turn!\033[0m")

			var roll bool
			if errAskToRoll := survey.AskOne(
				&survey.Confirm{
					Message: "Would you like to roll the dice?",
					Help:    help,
					Default: true,
				},
				&roll,
			); errAskToRoll != nil {
				panic(errAskToRoll)
			}

			randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
			if roll {
				firstRoll := randomizer.Intn(6) + 1
				secondRoll := randomizer.Intn(6) + 1
				fmt.Printf("\033[32mYou rolled a %d and a %d.\033[0m\n", firstRoll, secondRoll)
				if firstRoll == 1 || secondRoll == 1 {
					if firstRoll == 1 && secondRoll == 1 {
						humanTotalScore = 0
					}
					humanTurnScore = 0
					turn++
				} else {
					humanTurnScore += firstRoll + secondRoll
				}
			} else {
				fmt.Println("\033[1m\033[32mYou decided not to roll.\033[0m")
				humanTotalScore += humanTurnScore
				humanTurnScore = 0
				turn++
			}
		} else {
			fmt.Println("\033[1m\033[31mTheir turn!\033[0m")

			roll := randomizer.Intn(100) >= 50

			if roll {
				firstRoll := randomizer.Intn(6) + 1
				secondRoll := randomizer.Intn(6) + 1
				fmt.Printf("\033[31mThey rolled a %d and a %d.\033[0m\n", firstRoll, secondRoll)
				if firstRoll == 1 || secondRoll == 1 {
					if firstRoll == 1 && secondRoll == 1 {
						computerTotalScore = 0
					}
					computerTurnScore = 0
					turn++
				} else {
					computerTurnScore += firstRoll + secondRoll
				}
			} else {
				fmt.Println("\033[31mThey decided not to roll.\033[0m")
				computerTotalScore += computerTurnScore
				computerTurnScore = 0
				turn++
			}
		}
		fmt.Printf("\033[1m\033[33mYou: %d Them: %d\033[0m\n", humanTotalScore+humanTurnScore, computerTotalScore+computerTurnScore)
	}
	humanTotalScore += humanTurnScore
	computerTotalScore += humanTurnScore

	if humanTotalScore > computerTotalScore {
		fmt.Printf("You won in %d turns!\n", turn+1)

		gameScore := turn + 1
		highScore, postedHighScore := highScores[playerName]
		if !postedHighScore || gameScore < highScore {
			fmt.Printf("New high score for %s!\n", playerName)
			highScores[playerName] = gameScore
		}

		jsonFile, errCreate := os.Create("highscores.json")
		if errCreate != nil {
			panic(errCreate)
		}
		defer jsonFile.Close()
		jsonEncoder := json.NewEncoder(jsonFile)
		if errEncode := jsonEncoder.Encode(&highScores); errEncode != nil {
			panic(errEncode)
		}
	} else {
		fmt.Printf("You lose!")
	}
}
