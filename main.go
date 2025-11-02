package main

import (
	"strings"
)

type Room struct {
	name       string
	items      map[string]bool
	exits      map[string]*Room
	onLook     func() string
	doorOpened bool
}

type Player struct {
	currentRoom *Room
	inventory   map[string]bool
	hasBackpack bool
}

var world map[string]*Room
var player Player

func initGame() {
	world = make(map[string]*Room)
	player.inventory = make(map[string]bool)
	player.hasBackpack = false

	kitchen := &Room{
		name:  "кухня",
		items: map[string]bool{"чай": true},
		exits: make(map[string]*Room),
	}
	// onLook для перемещения в кухню
	kitchen.onLook = func() string {
		return "кухня, ничего интересного. можно пройти - коридор"
	}
	world[kitchen.name] = kitchen

	room := &Room{
		name: "комната",
		items: map[string]bool{
			"рюкзак":    true,
			"ключи":     true,
			"конспекты": true,
		},
		exits: make(map[string]*Room),
	}
	room.onLook = func() string {
		return "ты в своей комнате. можно пройти - коридор"
	}
	world[room.name] = room

	corridor := &Room{
		name:  "коридор",
		items: map[string]bool{},
		exits: make(map[string]*Room),
	}
	corridor.onLook = func() string {
		return "ничего интересного. можно пройти - кухня, комната, улица"
	}
	world[corridor.name] = corridor

	street := &Room{
		name:  "улица",
		items: map[string]bool{},
		exits: make(map[string]*Room),
	}
	street.onLook = func() string {
		return "на улице весна. можно пройти - домой"
	}
	world[street.name] = street

	/* Настроим пути */
	kitchen.exits["коридор"] = corridor
	room.exits["коридор"] = corridor
	corridor.exits["кухня"] = kitchen
	corridor.exits["комната"] = room
	corridor.exits["улица"] = street

	// Начальная локация
	player.currentRoom = kitchen
}

func handleCommand(cmd string) string {
	parts := strings.Split(cmd, " ")
	command := parts[0]
	args := parts[1:]

	commands := map[string]func([]string) string{
		"осмотреться": lookCommand,
		"идти":        goCommand,
		"взять":       takeCommand,
		"надеть":      wearCommand,
		"применить":   useCommand,
	}

	if cmdFunc, ok := commands[command]; ok {
		return cmdFunc(args)
	}
	return "неизвестная команда"
}

// Описание функций
func lookCommand(args []string) string {
	if player.currentRoom.name == "кухня" {
		if !player.hasBackpack {
			return "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"
		}
		return "ты находишься на кухне, на столе: чай, надо идти в универ. можно пройти - коридор"
	}

	if player.currentRoom.name == "комната" {
		if len(player.currentRoom.items) == 0 {
			return "пустая комната. можно пройти - коридор"
		}
		parts := []string{}
		if player.currentRoom.items["ключи"] || player.currentRoom.items["конспекты"] {
			items := []string{}
			if player.currentRoom.items["ключи"] {
				items = append(items, "ключи")
			}
			if player.currentRoom.items["конспекты"] {
				items = append(items, "конспекты")
			}
			parts = append(parts, "на столе: "+strings.Join(items, ", "))
		}
		if player.currentRoom.items["рюкзак"] {
			parts = append(parts, "на стуле: рюкзак")
		}
		return strings.Join(parts, ", ") + ". можно пройти - коридор"
	}

	if player.currentRoom.onLook != nil {
		return player.currentRoom.onLook()
	}
	return "ничего интересного"
}

func goCommand(args []string) string {
	if len(args) < 1 {
		return "куда идти? Команда: идти 'название комнаты'"
	}
	dest := args[0]

	nextRoom, ok := player.currentRoom.exits[dest]
	if !ok {
		return "нет пути в " + dest
	}

	// Проверка двери на улицу
	if dest == "улица" && !player.currentRoom.doorOpened {
		return "дверь закрыта"
	}

	player.currentRoom = nextRoom
	return nextRoom.onLook()
}

func takeCommand(args []string) string {
	if len(args) < 1 {
		return "что брать? Команда: взять 'название предмета'"
	}
	item := args[0]

	if !player.hasBackpack {
		return "некуда класть"
	}

	if player.currentRoom.items[item] {
		player.inventory[item] = true
		delete(player.currentRoom.items, item)
		return "предмет добавлен в инвентарь: " + item
	}
	return "нет такого"
}

func wearCommand(args []string) string {
	if len(args) < 1 {
		return "что надеть? Команда: надеть 'название предмета'"
	}
	item := args[0]

	if item == "рюкзак" && player.currentRoom.items[item] {
		player.hasBackpack = true
		delete(player.currentRoom.items, item)
		return "вы надели: рюкзак"
	}

	return "нечего надеть"
}

func useCommand(args []string) string {
	if len(args) < 2 {
		return "неверная команда"
	}
	item := args[0]
	target := args[1]

	if !player.inventory[item] {
		return "нет предмета в инвентаре - " + item
	}

	if item == "ключи" && target == "дверь" {
		if player.currentRoom.name != "коридор" {
			return "не к чему применить"
		}
		if player.currentRoom.doorOpened {
			return "дверь уже открыта"
		}
		player.currentRoom.doorOpened = true
		return "дверь открыта"
	}

	return "не к чему применить"
}
