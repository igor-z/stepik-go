package main

// не меняйте импорты, они нужны для проверки
import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// account представляет счет
type account struct {
	balance   int
	overdraft int
}

func main() {
	var acc account
	var trans []int
	var err error

	fmt.Print("-> ")

	if acc, trans, err = parseInput(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(acc, trans)
}

// parseInput считывает счет и список транзакций из os.Stdin.
func parseInput() (account, []int, error) {
	accSrc, transSrc := readInput()
	var err error
	var acc account
	var trans []int

	if acc, err = parseAccount(accSrc); err != nil {
		return account{}, nil, err
	}

	if trans, err = parseTransactions(transSrc); err != nil {
		return account{}, nil, err
	}
	return acc, trans, nil
}

// readInput возвращает строку, которая описывает счет
// и срез строк, который описывает список транзакций.
// эту функцию можно не менять
func readInput() (string, []string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	accSrc := scanner.Text()
	var transSrc []string
	for scanner.Scan() {
		transSrc = append(transSrc, scanner.Text())
	}
	return accSrc, transSrc
}

// parseAccount парсит счет из строки
// в формате balance/overdraft.
func parseAccount(src string) (account, error) {
	var balance, overdraft int
	var err error

	parts := strings.Split(src, "/")
	if balance, err = strconv.Atoi(parts[0]); err != nil {
		return account{}, err
	}

	if overdraft, err = strconv.Atoi(parts[1]); err != nil {
		return account{}, err
	}

	if overdraft < 0 {
		return account{}, errors.New("expect overdraft >= 0")
	}
	if balance < -overdraft {
		return account{}, errors.New("balance cannot exceed overdraft")
	}
	return account{balance, overdraft}, nil
}

// parseTransactions парсит список транзакций из строки
// в формате [t1 t2 t3 ... tn].
func parseTransactions(src []string) ([]int, error) {
	trans := make([]int, len(src))
	for idx, s := range src {
		var t int
		var err error
		if t, err = strconv.Atoi(s); err != nil {
			return trans, err
		}
		trans[idx] = t
	}
	return trans, nil
}
