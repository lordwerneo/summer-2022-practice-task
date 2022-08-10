package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	//Errors Description
	emptyDeptSt = "empty departure station"
	badDeptSt   = "bad departure station input"
	emptyArrSt  = "empty arrival station"
	badArrSt    = "bad arrival station input"
	unsCriteria = "unsupported criteria"
)

func TestFindTrainsSuccess(t *testing.T) {
	assert := assert.New(t)
	testsOK := map[string]struct {
		arrStation string
		depStation string
		criteria   string
		want       Trains
		wantErr    error
	}{
		"Price": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "price",
			want: Trains{{1177, 1902, 1929, 164.65, time.Date(0, time.January, 1, 10, 25, 00, 0, time.UTC), time.Date(0, time.January, 1, 16, 36, 00, 0, time.UTC)},
				{1178, 1902, 1929, 164.65, time.Date(0, time.January, 1, 10, 25, 0, 0, time.UTC), time.Date(0, time.January, 1, 16, 36, 00, 0, time.UTC)},
				{1141, 1902, 1929, 176.77, time.Date(0, time.January, 1, 12, 15, 00, 0, time.UTC), time.Date(0, time.January, 1, 16, 48, 00, 0, time.UTC)}},
			wantErr: nil,
		},
		"arrivalTime": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "arrival-time",
			want: Trains{{978, 1902, 1929, 258.53, time.Date(0, time.January, 1, 04, 15, 00, 0, time.UTC), time.Date(0, time.January, 1, 13, 10, 00, 0, time.UTC)},
				{1316, 1902, 1929, 209.73, time.Date(0, time.January, 1, 05, 55, 00, 0, time.UTC), time.Date(0, time.January, 1, 13, 52, 00, 0, time.UTC)},
				{2201, 1902, 1929, 280, time.Date(0, time.January, 1, 06, 15, 00, 0, time.UTC), time.Date(0, time.January, 1, 14, 55, 00, 0, time.UTC)}},
			wantErr: nil,
		},
		"departureTime": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "departure-time",
			want: Trains{{1386, 1902, 1929, 220.49, time.Date(0, time.January, 01, 8, 30, 00, 0, time.UTC), time.Date(0, time.January, 1, 13, 03, 00, 0, time.UTC)},
				{978, 1902, 1929, 258.53, time.Date(0, time.January, 1, 04, 15, 0, 0, time.UTC), time.Date(0, time.January, 1, 13, 10, 00, 0, time.UTC)},
				{1316, 1902, 1929, 209.73, time.Date(0, time.January, 1, 05, 55, 00, 0, time.UTC), time.Date(0, time.January, 1, 13, 52, 00, 0, time.UTC)}},
			wantErr: nil,
		},
	}

	for name, tc := range testsOK {
		t.Run(name, func(t *testing.T) {
			got, gotErr := FindTrains(tc.depStation, tc.arrStation, tc.criteria)
			if assert.NoError(gotErr) {
				assert.Len(got, 3)
				assert.Equal(tc.want, got)
			}
		})
	}
}

func TestFindTrainsNil(t *testing.T) {
	assert := assert.New(t)
	testsNIL := map[string]struct {
		arrStation string
		depStation string
		criteria   string
		want       Trains
		wantErr    error
	}{
		"UnknownArrStation": {
			depStation: "1902",
			arrStation: "777",
			criteria:   "price",
			want:       nil,
			wantErr:    nil,
		},
		"UnknownDeptStation": {
			depStation: "777",
			arrStation: "1929",
			criteria:   "price",
			want:       nil,
			wantErr:    nil,
		},
	}

	for name, tc := range testsNIL {
		t.Run(name, func(t *testing.T) {
			got, gotErr := FindTrains(tc.depStation, tc.arrStation, tc.criteria)
			assert.Nil(got)
			assert.Nil(gotErr)
		})
	}
}
func TestFindTrainsNegative(t *testing.T) {
	assert := assert.New(t)
	testsNotOK := map[string]struct {
		arrStation string
		depStation string
		criteria   string
		want       Trains
		wantErr    string
	}{
		"InvalidDepStation": {
			depStation: "w",
			arrStation: "1929",
			criteria:   "price",
			want:       nil,
			wantErr:    badDeptSt,
		},
		"InvalidArrStation": {
			depStation: "1902",
			arrStation: "19[[",
			criteria:   "price",
			want:       nil,
			wantErr:    badArrSt,
		},
		"InvalidArrStationAndCriteria": {
			depStation: "1902",
			arrStation: "19[[",
			criteria:   "priceds",
			want:       nil,
			wantErr:    badArrSt,
		},
		"UnsCriteria": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "duck",
			want:       nil,
			wantErr:    unsCriteria,
		},
		"EmptyCriteria": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "",
			want:       nil,
			wantErr:    unsCriteria,
		},
		"EmptyDeptStation": {
			depStation: "",
			arrStation: "1929",
			criteria:   "price",
			want:       nil,
			wantErr:    emptyDeptSt,
		},
		"EmptyArrStation": {
			depStation: "1902",
			arrStation: "",
			criteria:   "price",
			want:       nil,
			wantErr:    emptyArrSt,
		},
		"EmptyArrStationSpace": {
			depStation: "1902",
			arrStation: " ",
			criteria:   "price",
			want:       nil,
			wantErr:    badArrSt,
		},
	}

	for name, tc := range testsNotOK {
		t.Run(name, func(t *testing.T) {
			got, gotErr := FindTrains(tc.depStation, tc.arrStation, tc.criteria)
			if assert.Error(gotErr) && assert.Nil(got) {
				assert.EqualError(gotErr, tc.wantErr)
			}
		})
	}
}

func TestFindTrainsByViktor(t *testing.T) {
	testsTable := map[string]struct {
		arrStation    string
		depStation    string
		criteria      string
		expected      Trains
		expectedError error
	}{
		"successful_price": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "price",
			expected: Trains{
				{TrainID: 1177, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 164.65, ArrivalTime: time.Date(0, time.January, 1, 10, 25, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 16, 36, 0, 0, time.UTC)},
				{TrainID: 1178, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 164.65, ArrivalTime: time.Date(0, time.January, 1, 10, 25, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 16, 36, 0, 0, time.UTC)},
				{TrainID: 1141, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 176.77, ArrivalTime: time.Date(0, time.January, 1, 12, 15, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 16, 48, 0, 0, time.UTC)},
			},
			expectedError: nil,
		},
		"successful_arrival": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "arrival-time",
			expected: Trains{
				{TrainID: 978, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 258.53, ArrivalTime: time.Date(0, time.January, 1, 4, 15, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 13, 10, 0, 0, time.UTC)},
				{TrainID: 1316, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 209.73, ArrivalTime: time.Date(0, time.January, 1, 5, 55, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 13, 52, 0, 0, time.UTC)},
				{TrainID: 2201, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 280, ArrivalTime: time.Date(0, time.January, 1, 6, 15, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 14, 55, 0, 0, time.UTC)},
			},
			expectedError: nil,
		},
		"successful_departure": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "departure-time",
			expected: Trains{
				{TrainID: 1386, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 220.49, ArrivalTime: time.Date(0, time.January, 1, 8, 30, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 13, 3, 0, 0, time.UTC)},
				{TrainID: 978, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 258.53, ArrivalTime: time.Date(0, time.January, 1, 4, 15, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 13, 10, 0, 0, time.UTC)},
				{TrainID: 1316, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 209.73, ArrivalTime: time.Date(0, time.January, 1, 5, 55, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 13, 52, 0, 0, time.UTC)},
			},
			expectedError: nil,
		},
		"wrong_criteria": {
			depStation:    "1902",
			arrStation:    "1929",
			criteria:      "awef",
			expected:      nil,
			expectedError: errors.New("unsupported criteria"),
		},
		"absent_depStationId": {
			depStation:    "",
			arrStation:    "1929",
			criteria:      "departure",
			expected:      nil,
			expectedError: errors.New("empty departure station"),
		},
		"absent_arrStation": {
			depStation:    "1902",
			arrStation:    "",
			criteria:      "departure",
			expected:      nil,
			expectedError: errors.New("empty arrival station"),
		},
		"wrong_depStation": {
			depStation:    "12",
			arrStation:    "1929",
			criteria:      "price",
			expected:      nil,
			expectedError: nil,
		},
		"wrong_arrStation": {
			depStation:    "1902",
			arrStation:    "11",
			criteria:      "price",
			expected:      nil,
			expectedError: nil,
		},
		"bad_arrStation_input": {
			depStation:    "1902",
			arrStation:    "serg",
			criteria:      "price",
			expected:      nil,
			expectedError: errors.New("bad arrival station input"),
		},
		"bad_depStation_input": {
			depStation:    "serg",
			arrStation:    "1922",
			criteria:      "price",
			expected:      nil,
			expectedError: errors.New("bad departure station input"),
		},
	}

	for name, testCase := range testsTable {
		t.Run(name, func(tt *testing.T) {
			result, err := FindTrains(testCase.depStation, testCase.arrStation, testCase.criteria)
			if testCase.expectedError != nil {
				assert.EqualError(tt, err, testCase.expectedError.Error())
			} else {
				assert.NoError(tt, err)
			}
			assert.Equal(tt, testCase.expected, result)
		})
	}
}
