package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	maxTrainCount = 3
	minTrainCount = 1
	timeFormat    = "15:04:05"
)

var (
	unsupportedCriteria      = errors.New("unsupported criteria")
	emptyStation             = errors.New("empty station")
	emptyDepartureStation    = errors.New("empty departure station")
	emptyArrivalStation      = errors.New("empty arrival station")
	badStationInput          = errors.New("bad station input")
	badDepartureStationInput = errors.New("bad departure station input")
	badArrivalStationInput   = errors.New("bad arrival station input")
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
		t.ArrivalTime.Format(timeFormat), t.DepartureTime.Format(timeFormat),
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
	layout := timeFormat
	tempTime, _ := v["arrivalTime"].(string)
	t.ArrivalTime, _ = time.Parse(layout, tempTime)
	tempTime, _ = v["departureTime"].(string)
	t.DepartureTime, _ = time.Parse(layout, tempTime)

	return nil
}

func (t Trains) sortTrains(criteria string) {
	sort.SliceStable(t, func(i, j int) bool {
		switch criteria {
		case "price":
			return t[i].Price < t[j].Price
		case "departure-time":
			return t[i].DepartureTime.Before(t[j].DepartureTime)
		default:
			return t[i].ArrivalTime.Before(t[j].ArrivalTime)
		}
	})
}

func main() {
	fmt.Println("Enter departure station ID")
	departureStation := readUserInput()
	fmt.Println("Enter arrival station ID")
	arrivalStation := readUserInput()
	fmt.Println("Enter criteria")
	criteria := readUserInput()
	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		err = fmt.Errorf("stoped working: %w", err)
		fmt.Println(err)
		return
	}

	for _, v := range result {
		fmt.Println(v)
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	departure, err := checkStation(departureStation)
	if err != nil {
		if errors.Is(err, badStationInput) {
			return nil, badDepartureStationInput
		}
		return nil, emptyDepartureStation
	}

	arrival, err := checkStation(arrivalStation)
	if err != nil {
		if errors.Is(err, badStationInput) {
			return nil, badArrivalStationInput
		}
		return nil, emptyArrivalStation
	}

	if err = checkCriteria(criteria); err != nil {
		return nil, err
	}

	trains, err := importData()
	if err != nil {
		err = fmt.Errorf("something went wrong while importing data from .json: %w", err)
		return nil, err
	}

	trains = selectAndSortTrains(trains, arrival, departure, criteria)
	if len(trains) < 1 {
		return nil, nil
	}

	return trains, nil
}

func readUserInput() (userInput string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}

func importData() (importedData Trains, err error) {
	filename := "data.json"
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &importedData)
	if err != nil {
		return nil, err
	}

	return importedData, nil
}

func checkCriteria(criteria string) error {
	allowedCriteria := map[string]struct{}{
		"price":          {},
		"arrival-time":   {},
		"departure-time": {}}
	if _, ok := allowedCriteria[criteria]; !ok {
		return unsupportedCriteria
	}

	return nil
}

func checkStation(station string) (int, error) {
	if station == "" {
		return 0, emptyStation
	}
	result, err := strconv.Atoi(station)
	if err != nil {
		return 0, badStationInput
	}
	if result < 1 {
		return 0, badStationInput
	}

	return result, nil
}

func selectAndSortTrains(trains Trains, arrival int, departure int, criteria string) (sortedTrains Trains) {
	for _, v := range trains {
		if v.DepartureStationID == departure && v.ArrivalStationID == arrival {
			sortedTrains = append(sortedTrains, v)
		}
	}

	if len(sortedTrains) > minTrainCount {
		sortedTrains.sortTrains(criteria)
	}

	if len(sortedTrains) > maxTrainCount {
		return sortedTrains[:3]
	}

	return sortedTrains
}
