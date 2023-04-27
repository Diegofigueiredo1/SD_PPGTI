package remotelist

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

type ListData struct {
	ListID string `json:"list_id"`
	Values []int  `json:"values"`
}

type RemoteList struct {
	mu        sync.Mutex
	items     map[string][]int
	dataFile  string
	dataLoaded bool
}

func (l *RemoteList) Append(args []interface{}, reply *bool) error {
	if len(args) != 2 {
		return errors.New("número inválido de argumentos")
	}

	listID, ok := args[0].(string)
	if !ok {
		return errors.New("primeiro argumento deve ser uma string")
	}

	value, ok := args[1].(int)
	if !ok {
		return errors.New("segundo argumento deve ser um int")
	}

	log.Printf("Valor recebido para adicionar na lista '%s': %d\n", listID, value)

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.items[listID]; !ok {
		l.items[listID] = []int{}
	}

	l.items[listID] = append(l.items[listID], value)
	fmt.Println(l.items)

	err := l.saveListDataToFile()
	if err != nil {
		return err
	}

	*reply = true
	return nil
}

func (l *RemoteList) Remove(listID string, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.items[listID]) > 0 {
		*reply = l.items[listID][len(l.items[listID])-1]
		l.items[listID] = l.items[listID][:len(l.items[listID])-1]
		fmt.Println(l.items)

		err := l.saveListDataToFile()
		if err != nil {
			return err
		}
	} else {
		return errors.New("lista vazia")
	}

	return nil
}

func (l *RemoteList) Get(args []interface{}, reply *int) error {
	if len(args) != 2 {
		return errors.New("número inválido de argumentos")
	}

	listID, ok := args[0].(string)
	if !ok {
		return errors.New("primeiro argumento deve ser uma string")
	}

	index, ok := args[1].(int)
	if !ok {
		return errors.New("segundo argumento deve ser um int")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.items[listID]; !ok {
		return errors.New("lista não encontrada")
	}

	if index < 0 || index >= len(l.items[listID]) {
		return errors.New("índice fora do intervalo")
	}

	*reply = l.items[listID][index]
	return nil
}

func (l *RemoteList) Size(listID string, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	list, ok := l.items[listID]
	if !ok {
		return errors.New("lista não encontrada")
	}

	*reply = len(list)
	return nil
}

func (l *RemoteList) saveListDataToFile() error {
	file, err := os.Create(l.dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	listData := make([]ListData, 0, len(l.items))
	for listID, values := range l.items {
		listData = append(listData, ListData{
			ListID: listID,
			Values: values,
		})
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(listData)
	if err != nil {
		return err
	}

	return nil
}

func (l *RemoteList) loadListDataFromFile() error {
	file, err := os.Open(l.dataFile)
	if err != nil {
		// Se o arquivo não existir, não há dados para carregar
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var listData []ListData
	err = decoder.Decode(&listData)
	if err != nil {
		return err
	}

	l.items = make(map[string][]int)
	for _, data := range listData {
		l.items[data.ListID] = data.Values
	}

	return nil
}

func NewRemoteList(dataFile string) (*RemoteList, error) {
	list := &RemoteList{
		items:    make(map[string][]int),
		dataFile: dataFile,
	}

	err := list.loadListDataFromFile()
	if err != nil {
		return nil, err
	}

	list.dataLoaded = true
	return list, nil
}