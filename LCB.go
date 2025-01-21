package LCB

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	// "log"
	"net/http"
	"reflect"
	"sync"
)

// Пример использования
// state := NewNewState[*T]()
// state.Append(key, &element)
// element := *state.GetElementByKey(key)

type NewState[T any] struct {
	Mu sync.Mutex
	State map[int64]T
	Count int64
	WarningZone []int64
}

func NewNewState[T any]() *NewState[T] {
	return &NewState[T]{
		State: make(map[int64]T),
		Mu: sync.Mutex{},
		Count: -99999999999999,
		WarningZone: []int64{-99999999999999, -90000000000001},
	}
}

func (ns *NewState[T]) GetElement(key int64) T {
	ns.Mu.Lock()
	element := ns.State[key]
	ns.Mu.Unlock()

	return element
}

func (ns *NewState[T]) SetElement(key int64, element T) {
	if ns.WarningZone[0] < key && key < ns.WarningZone[1] {
		fmt.Print("WARNING, the value of key ran out of SAFE zone, Warning Zone: -99999999999999 : -90000000000001")
	}
	ns.Mu.Lock()
	ns.State[key] = element
	ns.Mu.Unlock()
}

func (ns *NewState[T]) AddElement(element T) {
	if ns.WarningZone[0] > ns.Count && ns.Count > ns.WarningZone[1] {
		fmt.Print("WARNING, the value of ns.Count ran out of Warning zone, Warning Zone: -99999999999999 : -90000000000001")
	}
	ns.Count += 1
	ns.Mu.Lock()
	ns.State[ns.Count] = element
	ns.Mu.Unlock()
}

func (ns *NewState[T]) DeleteElement(key int64) {
	ns.Mu.Lock()
	delete(ns.State, key)
	ns.Mu.Unlock()
}

func (ns *NewState[T]) GetKeyByNameFieldAndVal(nameField string, valueTarget any) (int64, error) {
	var result int64
	var er error
	var flag bool

	ns.Mu.Lock()
	for key, valueT := range ns.State {
		value, err := ns.getFieldValue(valueT, nameField)
		if err != nil {
			er = err
			break
		}

		if value == valueTarget {
			result = key
			flag = true
			break
		}
	}
	ns.Mu.Unlock()

	if !flag {
		er = fmt.Errorf("значение не найдено")
	}

	return result, er
}

func (ns *NewState[T]) getFieldValue(obj T, fieldName string) (any, error) {
	v := reflect.ValueOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ожидалась структура, но получен %s", v.Kind())
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("поле %s не найдено", fieldName)
	}

	return field.Interface(), nil
}

func (b *Bot) AddHandler(filter func(update Update) bool, callback func(update Update)) {
	b.handlers = append(b.handlers, Handler{Filter: filter, Callback: callback})
}

func (b *Bot) Start() {
	go b.pollUpdates()
	go b.processUpdates()
}

func (b *Bot) pollUpdates() {
	defer close(b.updatesChan)
	for {
		updates, err := b.getUpdates(b.lastUpdateId)
		if err != nil {
			fmt.Println("Error getting updates:", err)
			continue
		}

		for _, update := range updates {
			if b.lastUpdateId <= update.UpdateID {
				b.lastUpdateId = update.UpdateID + 1
			}
			b.updatesChan <- update
		}
	}
}

func (b *Bot) processUpdates() {
	for update := range b.updatesChan {
		for _, handler := range b.handlers {
			if handler.Filter(update) {
				go handler.Callback(update)
			}
		}
	}
}

func (b *Bot) getUpdates(offset int64) ([]Update, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", b.Token, offset)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var updates TelegramResponse
	err = json.Unmarshal(body, &updates)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if !updates.Ok {
		return nil, fmt.Errorf("telegram API returned an error33")
	}

	if len(updates.Result) > 0 && b.logs {
		var formattedJSON bytes.Buffer
		if err := json.Indent(&formattedJSON, body, "", "  "); err != nil {
			fmt.Println("Ошибка форматирования JSON: %v", err)
			return nil, fmt.Errorf("Ошибка форматирования JSON")
		}
		fmt.Println(formattedJSON.String())
	}

	for _, update := range updates.Result {
		if update.Message != nil && b.logs {
			fmt.Printf("Сообщение от %s: %s\n", update.Message.From.FirstName, update.Message.Text)
		}
	}

	return updates.Result, nil
}