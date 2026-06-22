package service

import (
	"fmt"
	"math/rand"
)

var nameAdjectives = []string{
	"awesome", "super", "funky", "happy", "sneaky", "mighty", "cosmic",
	"groovy", "spicy", "jolly", "fuzzy", "zesty", "bouncy", "dapper",
	"sleepy", "turbo", "wobbly", "snazzy", "cheeky", "noble",
	"glowing", "rusty", "velvet", "electric", "midnight", "golden",
	"silly", "brave", "quirky", "dazzling", "lucky", "humble",
	"feisty", "mellow", "rowdy", "plucky", "swift", "gentle",
	"radiant", "sparkly", "grumpy", "breezy", "crispy", "loyal",
	"witty", "fierce", "shiny", "chunky", "nimble", "epic", "uma",
}

var nameNouns = []string{
	"burrito", "pizza", "taco", "noodle", "waffle", "pickle", "donut",
	"mango", "panda", "otter", "penguin", "raccoon", "walrus", "narwhal",
	"muffin", "pancake", "dumpling", "biscuit", "nacho", "cupcake",
	"avocado", "bagel", "pretzel", "lobster", "hedgehog", "koala",
	"platypus", "wombat", "ferret", "meerkat", "quokka", "axolotl",
	"croissant", "ravioli", "samosa", "churro", "burger", "sushi",
	"falcon", "badger", "gecko", "lemur", "moose", "puffin",
	"mochi", "ramen", "popcorn", "brownie", "mushroom", "pineapple", "musume",
}

// generateTicketName returns a fun, human-friendly name like "awesome burrito".
func generateTicketName() string {
	adjective := nameAdjectives[rand.Intn(len(nameAdjectives))]
	noun := nameNouns[rand.Intn(len(nameNouns))]
	return fmt.Sprintf("%s %s", adjective, noun)
}
