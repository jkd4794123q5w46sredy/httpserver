package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	pgx "github.com/jackc/pgx/v5"
)

type Unit struct {
	UnitN   int `json:"unit_number"`
	ID      int `json:"id" db:"id"`
	Signal1 int `json:"signal1" db:"signal1"`
	Signal2 int `json:"signal2" db:"signal2"`
	Signal3 int `json:"signal3" db:"signal3"`
	Signal4 int `json:"signal4" db:"signal4"`
}

func establishDBConn() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:12345@pgDB:5432/units?sslmode=disable")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	return conn, nil
}

func mainFetchData(w http.ResponseWriter, _ *http.Request) {
	dbConnection, err := establishDBConn()
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer dbConnection.Close(context.Background())

	var unit1Week []Unit
	var unit2Week []Unit
	var unit3Week []Unit
	var unit4Week []Unit

	// Fetching data for Unit#1
	rows1, err := dbConnection.Query(context.Background(), "SELECT * FROM unit1 ORDER BY id")
	if err != nil {
		fmt.Fprint(w, "The error during db fetch occurred: ", err)
		return
	}
	defer rows1.Close()
	for rows1.Next() {
		var unit Unit
		err := rows1.Scan(&unit.ID, &unit.Signal1, &unit.Signal2, &unit.Signal3, &unit.Signal4)
		if err != nil {
			fmt.Fprint(w, "The error during db fetch occurred: ", err)
			return
		}
		unit.UnitN = 1
		unit1Week = append(unit1Week, unit)
	}

	fmt.Fprintf(w, "\n=== Unit #1 ===\n")
	for _, unit := range unit1Week {
		fmt.Fprintf(w, "Day: %d\nSignal1: %d, Signal2: %d, Signal3: %d, Signal4: %d\n\n",
			unit.ID, unit.Signal1, unit.Signal2, unit.Signal3, unit.Signal4)
	}

	// Fetching data for Unit#2
	rows2, err := dbConnection.Query(context.Background(), "SELECT * FROM unit2 ORDER BY id")
	if err != nil {
		fmt.Fprint(w, "The error during db fetch occurred: ", err)
		return
	}
	defer rows2.Close()
	for rows2.Next() {
		var unit Unit
		err := rows2.Scan(&unit.ID, &unit.Signal1, &unit.Signal2, &unit.Signal3)
		if err != nil {
			fmt.Fprint(w, "The error during db fetch occurred: ", err)
			return
		}
		unit.UnitN = 2
		unit2Week = append(unit2Week, unit)
	}

	fmt.Fprintf(w, "\n=== Unit #2 ===\n")
	for _, unit := range unit2Week {
		fmt.Fprintf(w, "Day: %d\nSignal1: %d, Signal2: %d, Signal3: %d\n\n",
			unit.ID, unit.Signal1, unit.Signal2, unit.Signal3)
	}

	// Fetching data for Unit#3
	rows3, err := dbConnection.Query(context.Background(), "SELECT * FROM unit3 ORDER BY id")
	if err != nil {
		fmt.Fprint(w, "The error during db fetch occurred: ", err)
		return
	}
	defer rows3.Close()
	for rows3.Next() {
		var unit Unit
		err := rows3.Scan(&unit.ID, &unit.Signal1, &unit.Signal2, &unit.Signal3)
		if err != nil {
			fmt.Fprint(w, "The error during db fetch occurred: ", err)
			return
		}
		unit.UnitN = 3
		unit3Week = append(unit3Week, unit)
	}

	fmt.Fprintf(w, "\n=== Unit #3 ===\n")
	for _, unit := range unit3Week {
		fmt.Fprintf(w, "Day: %d\nSignal1: %d, Signal2: %d, Signal3: %d\n\n",
			unit.ID, unit.Signal1, unit.Signal2, unit.Signal3)
	}

	// Fetching data for Unit#4
	rows4, err := dbConnection.Query(context.Background(), "SELECT * FROM unit4 ORDER BY id")
	if err != nil {
		fmt.Fprint(w, "The error during db fetch occurred: ", err)
		return
	}
	defer rows4.Close()
	for rows4.Next() {
		var unit Unit
		err := rows4.Scan(&unit.ID, &unit.Signal1, &unit.Signal2, &unit.Signal3)
		if err != nil {
			fmt.Fprint(w, "The error during db fetch occurred: ", err)
			return
		}
		unit.UnitN = 4
		unit4Week = append(unit4Week, unit)
	}

	fmt.Fprintf(w, "\n=== Unit #4 ===\n")
	for _, unit := range unit4Week {
		fmt.Fprintf(w, "Day: %d\nSignal1: %d, Signal2: %d, Signal3: %d\n\n",
			unit.ID, unit.Signal1, unit.Signal2, unit.Signal3)
	}
}

func mainInsertData(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "The error during request's body reading occurred: %v", err)
		return
	}

	if len(bodyBytes) == 0 {
		fmt.Fprintf(w, "The client sent POST request with empty body")
		return
	}

	var unit Unit
	err = json.Unmarshal(bodyBytes, &unit)
	if err != nil {
		fmt.Fprintf(w, "The error during JSON unmarshaling occurred: %v", err)
		return
	}

	dbConnection, err := establishDBConn()
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer dbConnection.Close(context.Background())

	switch unit.UnitN {
	case 1:
		_, err = dbConnection.Exec(context.Background(), "INSERT INTO unit1 (signal1, signal2, signal3, signal4) VALUES ($1, $2, $3, $4)",
			unit.Signal1, unit.Signal2, unit.Signal3, unit.Signal4)
		fmt.Fprint(w, "INSERT operation successful: ", err)
	case 2:
		_, err = dbConnection.Exec(context.Background(), "INSERT INTO unit2 (signal1, signal2, signal3) VALUES ($1, $2, $3)",
			unit.Signal1, unit.Signal2, unit.Signal3)
		fmt.Fprint(w, "INSERT operation successful: ", err)
	case 3:
		_, err = dbConnection.Exec(context.Background(), "INSERT INTO unit3 (signal1, signal2, signal3) VALUES ($1, $2, $3)",
			unit.Signal1, unit.Signal2, unit.Signal3)
		fmt.Fprint(w, "INSERT operation successful: ", err)
	case 4:
		_, err = dbConnection.Exec(context.Background(), "INSERT INTO unit4 (signal1, signal2, signal3) VALUES ($1, $2, $3)",
			unit.Signal1, unit.Signal2, unit.Signal3)
		fmt.Fprint(w, "INSERT operation successful: ", err)
	default:
		fmt.Fprintf(w, "Invalid unit number: %d", unit.UnitN)
		return
	}

	if err != nil {
		fmt.Fprintf(w, "INSERT operation failed: %v", err)
		return
	}
	fmt.Fprintf(w, "INSERT operation successful")
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	dbConnection, err := establishDBConn()
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer dbConnection.Close(context.Background())

	var unit1Week []Unit
	var unit2Week []Unit
	var unit3Week []Unit
	var unit4Week []Unit

	// Fetching data for Unit#1
	rows1, err := dbConnection.Query(context.Background(), "SELECT * FROM unit1 ORDER BY id")
	if err != nil {
		fmt.Fprint(w, "The error during db fetch occurred: ", err)
		return
	}
	defer rows1.Close()
	for rows1.Next() {
		var unit Unit
		err := rows1.Scan(&unit.ID, &unit.Signal1, &unit.Signal2, &unit.Signal3, &unit.Signal4)
		if err != nil {
			fmt.Fprint(w, "The error during db fetch occurred: ", err)
			return
		}
		unit.UnitN = 1
		unit1Week = append(unit1Week, unit)
	}

	//fmt.Fprint(w, "\nUnit#1")
	u1sig1Sum := 0
	u1sig2Sum := 0
	u1sig3Sum := 0
	u1sig4Sum := 0
	for _, unit := range unit1Week {
		u1sig1Sum += unit.Signal1
		u1sig2Sum += unit.Signal2
		u1sig3Sum += unit.Signal3
		u1sig4Sum += unit.Signal4
	}

	// Outputting statistics for Unit#1
	fmt.Fprint(w, "\nUnit#1")
	fmt.Fprintf(w, "\nSig#1. Level: %d", u1sig1Sum)
	fmt.Fprintf(w, "\nSig#2. Level: %d", u1sig2Sum)
	fmt.Fprintf(w, "\nSig#3. Level: %d", u1sig3Sum)
	fmt.Fprintf(w, "\nSig#4. Level: %d", u1sig4Sum)

	// Signal 1
	if u1sig1Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#1 Signal#1 code green")
	} else if u1sig1Sum/5 >= 4 && u1sig1Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#1 Signal#1 code orange")
	} else if u1sig1Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#1 Signal#1 code red")
	}

	// Signal 2
	if u1sig2Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#1 Signal#2 code green")
	} else if u1sig2Sum/5 >= 4 && u1sig2Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#1 Signal#2 code orange")
	} else if u1sig2Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#1 Signal#2 code red")
	}

	// Signal 3
	if u1sig3Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#1 Signal#3 code green")
	} else if u1sig3Sum/5 >= 4 && u1sig3Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#1 Signal#3 code orange")
	} else if u1sig3Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#1 Signal#3 code red")
	}

	// Signal 4
	if u1sig4Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#1 Signal#4 code green")
	} else if u1sig4Sum/5 >= 4 && u1sig4Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#1 Signal#4 code orange")
	} else if u1sig4Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#1 Signal#4 code red")
	}

	// Fetching data for Unit#2
	rows2, err := dbConnection.Query(context.Background(), "SELECT * FROM unit2 ORDER BY id")
	if err != nil {
		fmt.Fprint(w, "The error during db fetch occurred: ", err)
		return
	}
	defer rows2.Close()
	for rows2.Next() {
		var unit Unit
		err := rows2.Scan(&unit.ID, &unit.Signal1, &unit.Signal2, &unit.Signal3)
		if err != nil {
			fmt.Fprint(w, "The error during db fetch occurred: ", err)
			return
		}
		unit.UnitN = 2
		unit2Week = append(unit2Week, unit)
	}

	//Calculate statistics for Unit#2
	u2sig1Sum := 0
	u2sig2Sum := 0
	u2sig3Sum := 0
	for _, unit := range unit2Week {
		u2sig1Sum += unit.Signal1
		u2sig2Sum += unit.Signal2
		u2sig3Sum += unit.Signal3
	}

	// Outputting statistics for Unit#2
	fmt.Fprint(w, "\nUnit#2")
	fmt.Fprintf(w, "\nSig#1. Level: %d", u2sig1Sum)
	fmt.Fprintf(w, "\nSig#2. Level: %d", u2sig2Sum)
	fmt.Fprintf(w, "\nSig#3. Level: %d", u2sig3Sum)

	// Signal 1
	if u2sig1Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#2 Signal#1 code green")
	} else if u2sig1Sum/5 >= 4 && u2sig1Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#2 Signal#1 code orange")
	} else if u2sig1Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#2 Signal#1 code red")
	}

	// Signal 2
	if u2sig2Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#2 Signal#2 code green")
	} else if u2sig2Sum/5 >= 4 && u2sig2Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#2 Signal#2 code orange")
	} else if u2sig2Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#2 Signal#2 code red")
	}

	// Signal 3
	if u2sig3Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#2 Signal#3 code green")
	} else if u2sig3Sum/5 >= 4 && u2sig3Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#2 Signal#3 code orange")
	} else if u2sig3Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#2 Signal#3 code red")
	}

	// Fetching data for Unit#3
	rows3, err := dbConnection.Query(context.Background(), "SELECT * FROM unit3 ORDER BY id")
	if err != nil {
		fmt.Fprint(w, "The error during db fetch occurred: ", err)
		return
	}
	defer rows3.Close()
	for rows3.Next() {
		var unit Unit
		err := rows3.Scan(&unit.ID, &unit.Signal1, &unit.Signal2, &unit.Signal3)
		if err != nil {
			fmt.Fprint(w, "The error during db fetch occurred: ", err)
			return
		}
		unit.UnitN = 3
		unit3Week = append(unit3Week, unit)
	}

	/*fmt.Fprint(w, "\nUnit#3")
	for _, unit := range unit3Week {
		fmt.Fprintf(w, "Unit #%d\nDay: %d\nSignal1: %d, Signal2: %d, Signal3: %d\n", unit.UnitN, unit.ID, unit.Signal1, unit.Signal2, unit.Signal3)
	}*/

	//Calculate statistics for Unit#3
	u3sig1Sum := 0
	u3sig2Sum := 0
	u3sig3Sum := 0
	for _, unit := range unit3Week {
		u3sig1Sum += unit.Signal1
		u3sig2Sum += unit.Signal2
		u3sig3Sum += unit.Signal3
	}

	// Outputting statistics for Unit#3
	fmt.Fprint(w, "\nUnit#3")
	fmt.Fprintf(w, "\nSig#1. Level: %d", u3sig1Sum)
	fmt.Fprintf(w, "\nSig#2. Level: %d", u3sig2Sum)
	fmt.Fprintf(w, "\nSig#3. Level: %d", u3sig3Sum)

	// Signal 1
	if u3sig1Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#3 Signal#1 code green")
	} else if u3sig1Sum/5 >= 4 && u3sig1Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#3 Signal#1 code orange")
	} else if u3sig1Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#3 Signal#1 code red")
	}

	// Signal 2
	if u3sig2Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#3 Signal#2 code green")
	} else if u3sig2Sum/5 >= 4 && u3sig2Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#3 Signal#2 code orange")
	} else if u3sig2Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#3 Signal#2 code red")
	}

	// Signal 3
	if u3sig3Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#3 Signal#3 code green")
	} else if u3sig3Sum/5 >= 4 && u3sig3Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#3 Signal#3 code orange")
	} else if u3sig3Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#3 Signal#3 code red")
	}

	// Fetching data for Unit#4
	rows4, err := dbConnection.Query(context.Background(), "SELECT * FROM unit4 ORDER BY id")
	if err != nil {
		fmt.Fprint(w, "The error during db fetch occurred: ", err)
		return
	}
	defer rows4.Close()
	for rows4.Next() {
		var unit Unit
		err := rows4.Scan(&unit.ID, &unit.Signal1, &unit.Signal2, &unit.Signal3)
		if err != nil {
			fmt.Fprint(w, "The error during db fetch occurred: ", err)
			return
		}
		unit.UnitN = 4
		unit4Week = append(unit4Week, unit)
	}

	/*fmt.Fprint(w, "\nUnit#4")
	for _, unit := range unit4Week {
		fmt.Fprintf(w, "Unit #%d\nDay: %d\nSignal1: %d, Signal2: %d, Signal3: %d\n", unit.UnitN, unit.ID, unit.Signal1, unit.Signal2, unit.Signal3)
	}*/

	//Calculate statistics for Unit#4
	u4sig1Sum := 0
	u4sig2Sum := 0
	u4sig3Sum := 0
	for _, unit := range unit4Week {
		u4sig1Sum += unit.Signal1
		u4sig2Sum += unit.Signal2
		u4sig3Sum += unit.Signal3
	}

	// Outputting statistics for Unit#4
	fmt.Fprint(w, "\nUnit#4")
	fmt.Fprintf(w, "\nSig#1. Level: %d", u4sig1Sum)
	fmt.Fprintf(w, "\nSig#2. Level: %d", u4sig2Sum)
	fmt.Fprintf(w, "\nSig#3. Level: %d", u4sig3Sum)

	// Signal 1
	if u4sig1Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#4 Signal#1 code green")
	} else if u4sig1Sum/5 >= 4 && u4sig1Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#4 Signal#1 code orange")
	} else if u4sig1Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#4 Signal#1 code red")
	}

	// Signal 2
	if u4sig2Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#4 Signal#2 code green")
	} else if u4sig2Sum/5 >= 4 && u4sig2Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#4 Signal#2 code orange")
	} else if u4sig2Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#4 Signal#2 code red")
	}

	// Signal 3
	if u4sig3Sum/5 < 4 {
		fmt.Fprintf(w, "\nUnit#4 Signal#3 code green")
	} else if u4sig3Sum/5 >= 4 && u4sig3Sum/5 < 7 {
		fmt.Fprintf(w, "\nUnit#4 Signal#3 code orange")
	} else if u4sig3Sum/5 >= 7 {
		fmt.Fprintf(w, "\nUnit#4 Signal#3 code red")
	}
}

func mainhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		mainFetchData(w, r)
	}

	if r.Method == http.MethodPost {
		mainInsertData(w, r)
	}
}

func main() {
	http.HandleFunc("/", mainhandler)
	http.HandleFunc("/run_stats", statsHandler)

	fmt.Println("HTTP server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
