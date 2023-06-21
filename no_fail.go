package main

import (
	"errors"
	"fmt"
	"os"
)

// label - уникальное наименование
type label string

// command - команда, которую можно выполнять в игре
type command label

// список доступных команд
var (
	eat  = command("eat")
	take = command("take")
	talk = command("talk to")
)

// полный список объектов в игре
var (
	apple = thing{"apple", map[command]string{
		eat:  "mmm, delicious!",
		take: "you have an apple now",
	}}
	bob = thing{"bob", map[command]string{
		talk: "Bob says hello",
	}}
	coin = thing{"coin", map[command]string{
		take: "you have a coin now",
	}}
	mirror = thing{"mirror", map[command]string{
		take: "you have a mirror now",
		talk: "mirror does not answer",
	}}
	mushroom = thing{"mushroom", map[command]string{
		eat:  "tastes funny",
		take: "you have a mushroom now",
	}}
)

// thing - объект, который существует в игре
type thing struct {
	name    label
	actions map[command]string
}

// supports() возвращает true, если объект
// поддерживает команду action
func (t thing) supports(action command) bool {
	_, ok := t.actions[action]
	return ok
}

// step описывает шаг игры: сочетание команды и объекта
type step struct {
	cmd command
	obj thing
}

// isValid() возвращает true, если объект
// совместим с командой
func (s step) isValid() bool {
	return s.obj.supports(s.cmd)
}

// начало решения

// invalidStepError - ошибка, которая возникает,
// когда команда шага не совместима с объектом
type invalidStepError struct {
	st step
}

func (t invalidStepError) Error() string {
	return fmt.Sprintf("cannot %s", t.st)
}

// notEnoughObjectsError - ошибка, которая возникает,
// когда в игре закончились объекты определенного типа
type notEnoughObjectsError struct {
	st step
}

func (t notEnoughObjectsError) Error() string {
	return fmt.Sprintf("there are no %ss left", t.st.obj)
}

// commandLimitExceededError - ошибка, которая возникает,
// когда игрок превысил лимит на выполнение команды
type commandLimitExceededError struct {
	cmd command
}

func (t commandLimitExceededError) Error() string {
	return fmt.Sprintf("you don't want to eat anymore")
}

// objectLimitExceededError - ошибка, которая возникает,
// когда игрок превысил лимит на количество объектов
// определенного типа в инвентаре
type objectLimitExceededError struct {
	obj   thing
	limit int
}

func (t objectLimitExceededError) Error() string {
	return fmt.Sprintf("you already have a %s", t.obj)
}

// gameOverError - ошибка, которая произошла в игре
type gameOverError struct {
	// количество шагов, успешно выполненных
	// до того, как произошла ошибка
	err    error
	nSteps int
	// ...
}

func (t gameOverError) Error() string {
	return fmt.Sprintf("%s", t.err)
}

func (t gameOverError) Unwrap() error {
	return t.err
}

// player - игрок
type player struct {
	// количество съеденного
	nEaten int
	// количество диалогов
	nDialogs int
	// инвентарь
	inventory []thing
}

// has() возвращает true, если у игрока
// в инвентаре есть предмет obj
func (p *player) has(obj thing) bool {
	for _, got := range p.inventory {
		if got.name == obj.name {
			return true
		}
	}
	return false
}

// do() выполняет команду cmd над объектом obj
// от имени игрока
func (p *player) do(cmd command, obj thing) error {
	// действуем в соответствии с командой
	switch cmd {
	case eat:
		if p.nEaten > 1 {
			return commandLimitExceededError{eat}
		}
		p.nEaten++
	case take:
		if p.has(obj) {
			return objectLimitExceededError{obj, 1}
		}
		p.inventory = append(p.inventory, obj)
	case talk:
		if p.nDialogs > 0 {
			return commandLimitExceededError{talk}
		}
		p.nDialogs++
	}
	return nil
}

// newPlayer создает нового игрока
func newPlayer() *player {
	return &player{inventory: []thing{}}
}

// game описывает игру
type game struct {
	// игрок
	player *player
	// объекты игрового мира
	things map[label]int
	// количество успешно выполненных шагов
	nSteps int
}

// has() проверяет, остались ли в игровом мире указанные предметы
func (g *game) has(obj thing) bool {
	count := g.things[obj.name]
	return count > 0
}

// execute() выполняет шаг step
func (g *game) execute(st step) error {
	// проверяем совместимость команды и объекта
	if !st.isValid() {
		return gameOverError{invalidStepError{st}, g.nSteps}
	}

	// когда игрок берет или съедает предмет,
	// тот пропадает из игрового мира
	if st.cmd == take || st.cmd == eat {
		if !g.has(st.obj) {
			return gameOverError{notEnoughObjectsError{st}, g.nSteps}
		}
		g.things[st.obj.name]--
	}

	// выполняем команду от имени игрока
	if err := g.player.do(st.cmd, st.obj); err != nil {
		return gameOverError{err, g.nSteps}
	}

	g.nSteps++
	return nil
}

// newGame() создает новую игру
func newGame() *game {
	p := newPlayer()
	things := map[label]int{
		apple.name:    2,
		coin.name:     3,
		mirror.name:   1,
		mushroom.name: 1,
	}
	return &game{p, things, 0}
}

// giveAdvice() возвращает совет, который
// поможет игроку избежать ошибки err в будущем
func giveAdvice(err error) string {
	invalidStep := invalidStepError{}
	if errors.As(err, &invalidStep) {
		return fmt.Sprintf("things like '%s %s' are impossible", invalidStep.st.cmd, invalidStep.st.obj.name)
	}

	notEnoughObjects := notEnoughObjectsError{}
	if errors.As(err, &notEnoughObjects) {
		return fmt.Sprintf("be careful with scarce %ss", notEnoughObjects.st.obj.name)
	}

	commandLimit := commandLimitExceededError{}
	if errors.As(err, &commandLimit) {
		if commandLimit.cmd == eat {
			return fmt.Sprintf("eat less")
		}
		if commandLimit.cmd == talk {
			return fmt.Sprintf("talk to less")
		}
	}

	objectLimit := objectLimitExceededError{}
	if errors.As(err, &objectLimit) {
		return fmt.Sprintf("don't be greedy, %d %s is enough", objectLimit.limit, objectLimit.obj.name)
	}
	return ""
}

// конец решения

func main() {
	gm := newGame()
	steps := []step{
		//{eat, apple},
		//{talk, bob},
		//{take, coin},
		//{eat, mushroom},
		{eat, mirror},
		{eat, coin},
		{talk, bob},
		{talk, bob},
	}

	for _, st := range steps {
		if err := tryStep(gm, st); err != nil {
			fmt.Println(err)
			fmt.Println(giveAdvice(err))
			os.Exit(1)
		}
	}
	fmt.Println("You win!")
}

// tryStep() выполняет шаг игры и печатает результат
func tryStep(gm *game, st step) error {
	fmt.Printf("trying to %s %s... ", st.cmd, st.obj.name)
	if err := gm.execute(st); err != nil {
		fmt.Println("FAIL")
		return err
	}
	fmt.Println("OK")
	return nil
}
