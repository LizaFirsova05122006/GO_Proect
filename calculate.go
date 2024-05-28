package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	idCounter int
	mu        sync.Mutex
)

func calculateResult(expression string) int {
	result := 0
	arg1 := int(expression[0] - '0')
	arg2 := int(expression[2] - '0')
	operation := string(expression[1])

	switch operation {
	case "+":
		result = arg1 + arg2
	case "-":
		result = arg1 - arg2
	case "*":
		result = arg1 * arg2
	case "/":
		result = arg1 / arg2
	}

	return result
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	calculate := r.URL.Query().Get("exp")
	if calculate == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// Проверяем выражение
	if !isValidExpression(calculate) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\n422 - невалидные данные\n}")
		return
	}

	// Открываем файл calculate.csv
	file, err := os.OpenFile("calculate.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Генерируем уникальный ID
	ID := "1"
	if len(rows) > 0 { // Если есть данные в файле
		lastRow := rows[len(rows)-1]
		lastID := lastRow[0]
		incrementedID := incrementID(lastID)
		ID = incrementedID
	}

	// Добавляем в CSV файл
	writer := csv.NewWriter(file)

	err = writer.Write([]string{ID, calculate})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Flush()

	// Проверяем выражение еще раз
	if !isValidExpression(calculate) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\n422 - невалидные данные\n}")
		return
	}

	// Записываем новую запись в файл CSV
	err = writer.Write([]string{ID, calculate})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rezults := calculateResult(calculate)
	resultsFile, err := os.OpenFile("results.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resultsFile.Close()

	// Создаем writer для файла results.csv
	resultsWriter := csv.NewWriter(resultsFile)

	// Записываем результат в файл
	err = resultsWriter.Write([]string{ID, strconv.Itoa(rezults)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultsWriter.Flush()

	// Возвращаем ответ с ID
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "{\n\"ID\": %s\n}", ID)
}

func incrementID(lastID string) string {
	incrementedID, _ := strconv.Atoi(lastID)
	incrementedID++
	return strconv.Itoa(incrementedID)
}

func isValidExpression(exp string) bool {
	validChars := "+-*/0123456789"
	for _, char := range exp {
		if !strings.Contains(validChars, string(char)) {
			return false
		}
	}
	return true
}

func expressionsHandler(w http.ResponseWriter, r *http.Request) {
	// Открываем файл calculate.csv
	file, err := os.Open("calculate.csv")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Читаем данные из calculate.csv
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Создаем карту для хранения результатов
	results := make(map[string]string)

	// Открываем файл results.csv
	resultsFile, err := os.Open("results.csv")
	if err == nil {
		defer resultsFile.Close()
		resultsReader := csv.NewReader(resultsFile)
		existingRecords, err := resultsReader.ReadAll()
		if err == nil {
			for _, record := range existingRecords {
				results[record[0]] = record[1]
			}
		}
	}

	expressionsList := make([]map[string]string, 0)
	for _, record := range records {
		expressionID := record[0]
		status := "not calculated"
		result := ""

		// Проверяем если ID уже присутствует в results.csv
		if value, ok := results[expressionID]; ok {
			status = "calculated"
			result = value
		}

		expressionsList = append(expressionsList, map[string]string{
			"id":     expressionID,
			"status": status,
			"result": result,
		})
	}

	jsonResponse, err := json.Marshal(expressionsList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func expressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	expressionID := r.URL.Path[len("/api/v1/expressions/"):]

	file, err := os.Open("calculate.csv")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Создаем карту для хранения результатов
	results := make(map[string]string)

	// Открываем файл results.csv
	resultsFile, err := os.Open("results.csv")
	if err == nil {
		defer resultsFile.Close()
		resultsReader := csv.NewReader(resultsFile)
		existingRecords, err := resultsReader.ReadAll()
		if err == nil {
			for _, record := range existingRecords {
				results[record[0]] = record[1]
			}
		}
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, record := range records {
		if record[0] == expressionID {
			expression := record[1]
			status := "not calculated"
			result := ""

			// Проверяем если ID уже присутствует в results.csv
			if value, ok := results[expressionID]; ok {
				status = "calculated"
				result = value
			}

			// Perform calculation if needed

			jsonResponse, err := json.Marshal(map[string]interface{}{
				"id":         expressionID,
				"status":     status,
				"result":     result,
				"expression": expression,
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResponse)
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "\n{\n404 - выражение не найдено\n}")
}

func internalTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/api/v1/internal/task/"):]

	file, err := os.Open("calculate.csv")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var taskData map[string]string
	for _, record := range records {
		if record[0] == taskID {
			arg1 := string(record[1][0])
			arg2 := string(record[1][2])
			operation := string(record[1][1])

			switch operation {
			case "+":
				operation = "add"
			case "-":
				operation = "subtract"
			case "*":
				operation = "multiply"
			case "/":
				operation = "divide"
			}

			operationTime := "" // You can calculate the operation time here if needed

			taskData = map[string]string{
				"id":             taskID,
				"arg1":           arg1,
				"arg2":           arg2,
				"operation":      operation,
				"operation_time": operationTime,
			}
			break
		}
	}

	if taskData == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(map[string]map[string]string{"task": taskData})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		expression := r.FormValue("expression")
		// Далее можно обработать выражение, например, вычислить его значение
		fmt.Fprintf(w, "Вы ввели выражение: %s", expression)
	} else {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	http.HandleFunc("/api/v1/expressions", expressionsHandler)
	http.HandleFunc("/api/v1/expressions/", expressionByIDHandler)
	http.HandleFunc("/api/v1/internal/task/", internalTaskHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "home.html")
	})
	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request) {
		expression := r.FormValue("username")
		if !isValidExpression(expression) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "{\n422 - невалидные данные\n}")
			return
		}

		// Открываем файл calculate.csv
		file, err := os.OpenFile("calculate.csv", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		rows, err := reader.ReadAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Генерируем уникальный ID
		ID := "1"
		if len(rows) > 0 { // Если есть данные в файле
			lastRow := rows[len(rows)-1]
			lastID := lastRow[0]
			incrementedID := incrementID(lastID)
			ID = incrementedID
		}

		// Добавляем в CSV файл
		writer := csv.NewWriter(file)

		err = writer.Write([]string{ID, expression})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Flush()
		fmt.Fprintf(w, "Ваши данные записаны успешно, можете при помощи стрелочки вернутся обратно и вввести еще выражение")
		fmt.Fprintf(w, "\nID вашего выражения : %s", ID)
		fmt.Fprintf(w, "\nТакже вы можете перейти на: \n1) http://localhost:8080/api/v1/expressions")
		fmt.Fprintf(w, "\n2) http://localhost:8080/api/v1/expressions/<ID>\n3) http://localhost:8080/api/v1/internal/task/")
	})
	http.ListenAndServe(":8080", nil)
}
