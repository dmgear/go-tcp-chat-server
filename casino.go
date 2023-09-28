package main

import (
	"fmt"
	"math/rand"
	"time"
)

var DENOMINATIONS = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
var SUITS = []string{"Hearts", "Diamonds", "Clubs", "Spades"}

type Gopher struct {
	state bool
	players []string
	
}

type Card struct {
	Denomination string
	Suit         string
}

type Deck struct {
	cards []Card
}

func (d *Deck) generate_deck() {
	for _, denom := range DENOMINATIONS {
		for _, suit := range SUITS {
			card := Card{
				Denomination: denom,
				Suit:         suit,
			}
			d.cards = append(d.cards, card)
		}
	}
}

func (d *Deck) printDeck() {
	for _, card := range d.cards {
		fmt.Printf("Card: %s of %s\n", card.Denomination, card.Suit)
	}
}

func (d *Deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) draw(c *Client) {
	if len(d.cards) == 0 {
		fmt.Println("Game over!")
	}
	c.hand = append(c.hand, d.cards[len(d.cards)-1])
	d.cards = d.cards[:len(d.cards)-1]
}

func (d *Deck) gopherDealHand(r *Room, c *Client) {
	// check if being called in casino room, then check how many clients exist in the room to determine how many cards need to be dealt to each client
	if r.name == "Casino" && len(r.Members) <= 4 {
		// iterate over the members map and deal 7 cards to each client
		for i := 0; i < 7; i++ {
			for _, _ = range r.Members {// <-- note the blank identifiers used because we wont actually be using the variables for anything
				d.draw(c)
			}
		}
	} else if r.name == "Casino" && len(r.Members) > 4 {
		for i := 0; i < 5; i++ {
			for _, _ = range r.Members {
				d.draw(c)
			}
		}
	}
}

func (c *Client) hasCard(username string, target string) {
	for _, client := range clients {
		if client.username == username {
			for _, card := range c.hand {
				if card.Denomination == target {
					c.hand = append(c.hand, card)
					return
				}
			}
		}
	}
}

func (g *Gopher) startGame(r *Room) *Gopher {
	g.state = true
	g.players = make([]string, 0)

	for _, user := range r.Members {
		g.players = append(g.players, user)
	}
	return g
}