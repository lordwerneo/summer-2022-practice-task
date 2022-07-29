package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

type Trains []Train

type Train struct {
	TrainID            int       `json:"trainId"`
	DepartureStationID int       `json:"departureStationID"`
	ArrivalStationID   int       `json:"arrivalStationId"`
	Price              float32   `json:"price"`
	ArrivalTime        time.Time `json:"arrivalTime"`
	DepartureTime      time.Time `json:"departureTime"`
}

func (t Train) String() string {
	return fmt.Sprintf("TrainID: %v, DepartureStationID: %v, ArrivalStationID: %v, Price: %v, ArrivalTime: %v,"+
		" DepartureTime: %v", t.TrainID, t.DepartureStationID, t.ArrivalStationID, t.Price,
		t.ArrivalTime.Format("15:04:05"), t.DepartureTime.Format("15:04:05"),
	)
}

func (t *Train) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	tempId, _ := v["trainId"].(float64)
	t.TrainID = int(tempId)

	tempDepId, _ := v["departureStationId"].(float64)
	t.DepartureStationID = int(tempDepId)

	tempArrId, _ := v["arrivalStationId"].(float64)
	t.ArrivalStationID = int(tempArrId)

	tempPrice, _ := v["price"].(float64)
	t.Price = float32(tempPrice)

	layout := "15:04:05"

	tempTime, _ := v["arrivalTime"].(string)
	t.ArrivalTime, _ = time.Parse(layout, tempTime)

	tempTime, _ = v["departureTime"].(string)
	t.DepartureTime, _ = time.Parse(layout, tempTime)

	return nil
}

var (
	UnsupportedCriteria      = errors.New("unsupported criteria")
	EmptyStation             = errors.New("empty station")
	EmptyDepartureStation    = errors.New("empty departure station")
	EmptyArrivalStation      = errors.New("empty arrival station")
	BadStationInput          = errors.New("bad station input")
	BadDepartureStationInput = errors.New("bad departure station input")
	BadArrivalStationInput   = errors.New("bad arrival station input")
)

func main() {
	var (
		departureStation string
		arrivalStation   string
		criteria         string
		result           []Train
	)
	//	... запит даних від користувача
	fmt.Println("Enter departure station ID")
	departureStation = readUserInput()
	fmt.Println("Enter arrival station ID")
	arrivalStation = readUserInput()
	fmt.Println("Enter criteria")
	criteria = readUserInput()

	//result, err := FindTrains(departureStation, arrivalStation, criteria))
	result, err := FindTrains(departureStation, arrivalStation, criteria)

	//	... обробка помилки
	if err != nil {
		err = fmt.Errorf("entered incorrect parameters: %w", err)
		fmt.Println(err)
		return
	}
	//	... друк result
	for _, v := range result {
		fmt.Println(v)
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	// ... код
	if err := checkCriteria(criteria); err != nil {
		return nil, err
	}

	departure, err := checkStation(departureStation)
	if err != nil {
		if errors.Is(err, BadStationInput) {
			return nil, BadDepartureStationInput
		}
		return nil, EmptyDepartureStation
	}

	arrival, err := checkStation(arrivalStation)
	if err != nil {
		if errors.Is(err, BadStationInput) {
			return nil, BadArrivalStationInput
		}
		return nil, EmptyArrivalStation
	}

	var trains Trains
	trains = importData()
	trains = selectAndSortTrains(trains, arrival, departure, criteria)

	return trains, nil // маєте повернути правильні значення
}

func readUserInput() (userInput string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userInput = scanner.Text()
	return userInput
}

func importData() (importedData Trains) {
	filename := "data.json"
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, &importedData)

	return importedData
}

func checkCriteria(criteria string) error {
	allowedCriteria := map[string]struct{}{
		"price":          {},
		"arrival-time":   {},
		"departure-time": {}}
	if _, ok := allowedCriteria[criteria]; !ok {
		return UnsupportedCriteria
	}
	return nil
}

func checkStation(station string) (int, error) {
	if station == "" {
		return 0, EmptyStation
	}
	result, err := strconv.Atoi(station)
	if err != nil {
		return 0, BadStationInput
	}
	if result < 1 {
		return 0, BadStationInput
	}
	return result, nil
}

func sortByPrice(trains Trains) Trains {
	sort.Slice(trains, func(i, j int) bool {
		if trains[i].Price == trains[j].Price {
			return trains[i].TrainID < trains[j].TrainID
		}
		return trains[i].Price < trains[j].Price
	})
	return trains
}

func sortByTime(trains Trains, criteria string) Trains {
	sort.Slice(trains, func(i, j int) bool {
		if criteria == "departure-time" {
			if trains[i].DepartureTime.Equal(trains[j].DepartureTime) {
				return trains[i].TrainID < trains[j].TrainID
			}
			return trains[i].DepartureTime.Before(trains[j].DepartureTime)
		} else {
			if trains[i].ArrivalTime.Equal(trains[j].ArrivalTime) {
				return trains[i].TrainID < trains[j].TrainID
			}
			return trains[i].ArrivalTime.Before(trains[j].ArrivalTime)
		}
	})

	return trains
}

func selectAndSortTrains(trains Trains, arrival int, departure int, criteria string) (sortedTrains Trains) {
	for _, v := range trains {
		if v.DepartureStationID == departure && v.ArrivalStationID == arrival {
			sortedTrains = append(sortedTrains, v)
		}
	}

	if len(sortedTrains) > 1 {
		switch criteria {
		case "price":
			sortedTrains = sortByPrice(sortedTrains)
		default:
			sortedTrains = sortByTime(sortedTrains, criteria)
		}
	}
	return limitTrains(sortedTrains)
}

func limitTrains(trains Trains) (newTrains Trains) {
	for _, v := range trains {
		newTrains = append(newTrains, v)
		if len(newTrains) == 3 {
			return newTrains
		}
	}
	return newTrains
}
