package ext

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

// Initialize faker
func init() {
	gofakeit.Seed(time.Now().UnixNano())
}

// Gen utility struct for generating mock data
type Gen struct{}

// GenerateMobileNumber generates a random mobile number
func (g *Gen) GenerateMobileNumber(countryCode int) string {
	prefixes := []int{62, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99}
	numberPrefix := prefixes[rand.Intn(len(prefixes))]
	number := rand.Intn(89999999) + 10000000
	return fmt.Sprintf("+%d-%d%d", countryCode, numberPrefix, number)
}

// QueryParser parses a query string into a map
func (g *Gen) QueryParser(query string) map[string]string {
	data := make(map[string]string)
	re := regexp.MustCompile(`(\w+)\((.*)\)`)
	match := re.FindStringSubmatch(query)
	if len(match) < 3 {
		return data
	}
	data["type"] = match[1]
	payload := strings.Split(strings.ReplaceAll(match[2], "\\", ""), "$")
	for _, item := range payload {
		kv := strings.Split(item, "=")
		if len(kv) == 2 {
			key := regexp.MustCompile(`[^a-zA-Z_]`).ReplaceAllString(kv[0], "")
			data[key] = kv[1]
		}
	}
	return data
}

// GenStatic generates static data based on a query
func (g *Gen) GenStatic(query string) interface{} {
	if _, err := strconv.Atoi(query); err == nil {
		return query
	}
	if strings.Contains(query, "list") {
		return g.GenList(strings.Replace(query, "list-", "", 1))
	}

	data := g.QueryParser(query)
	dataType := data["type"]

	// Type-specific generators
	switch dataType {
	case "name":
		return gofakeit.Name()
	case "email":
		domain := data["domain"]
		if domain == "" {
			domain = "gmail.com"
		}
		return gofakeit.Email()
	case "password":
		length, _ := strconv.Atoi(data["len"])
		if length == 0 {
			length = 8
		}
		return gofakeit.Password(true, true, true, true, false, length)
	case "text", "str":
		length, _ := strconv.Atoi(data["len"])
		if length == 0 {
			length = 5
		}
		return gofakeit.LetterN(uint(length))
	case "int":
		length, _ := strconv.Atoi(data["len"])
		if length == 0 {
			length = 3
		}
		min := intPow(10, length-1)
		max := intPow(10, length) - 1
		return rand.Intn(max-min+1) + min
	case "time":
		return time.Now().Format("15:04:05")
	case "date":
		// Customize date generation logic if needed
		return time.Now().Format("2006-01-02")
	case "address":
		return gofakeit.Address().Address
	case "company":
		return gofakeit.Company()
	case "phone":
		code, _ := strconv.Atoi(data["code"])
		if code == 0 {
			code = 91
		}
		return g.GenerateMobileNumber(code)
	case "bool":
		return rand.Intn(2) == 1
	case "float":
		return rand.Float64() * 100
	case "age":
		min, _ := strconv.Atoi(data["min"])
		if min == 0 {
			min = 1
		}
		max, _ := strconv.Atoi(data["max"])
		if max == 0 {
			max = 100
		}
		return rand.Intn(max-min+1) + min
	case "description":
		words, _ := strconv.Atoi(data["words"])
		if words == 0 {
			words = 4
		}
		return gofakeit.Sentence(words)
	case "image":
		width, _ := strconv.Atoi(data["width"])
		if width == 0 {
			width = 200
		}
		height, _ := strconv.Atoi(data["height"])
		if height == 0 {
			height = 200
		}
		return fmt.Sprintf("https://via.placeholder.com/%dx%d", width, height)
	default:
		return "Invalid type"
	}
}

// GenList generates a list of mock data based on the query
func (g *Gen) GenList(query string) []interface{} {
	data := g.QueryParser(query)
	dataType := data["type"]
	amount, _ := strconv.Atoi(data["amount"])
	if amount == 0 {
		amount = 3
	}

	result := make([]interface{}, amount)
	switch dataType {
	case "int":
		min, _ := strconv.Atoi(data["min"])
		max, _ := strconv.Atoi(data["max"])
		for i := 0; i < amount; i++ {
			result[i] = rand.Intn(max-min+1) + min
		}
	case "str":
		for i := 0; i < amount; i++ {
			result[i] = gofakeit.LetterN(10)
		}
	case "name":
		for i := 0; i < amount; i++ {
			result[i] = gofakeit.Name()
		}
	case "email":
		domain := data["domain"]
		if domain == "" {
			domain = "gmail.com"
		}
		for i := 0; i < amount; i++ {
			result[i] = gofakeit.Email()
		}
	}
	return result
}

// GenDict recursively generates a map based on a schema
func (g *Gen) GenDict(data map[string]interface{}) map[string]interface{} {
	copyData := make(map[string]interface{}, len(data))
	for k, v := range data {
		switch value := v.(type) {
		case map[string]interface{}:
			amount, ok := value["_$amount"].(int)
			if ok {
				list := make([]map[string]interface{}, amount)
				for i := 0; i < amount; i++ {
					list[i] = g.GenDict(value)
				}
				copyData[k] = list
			} else {
				copyData[k] = g.GenDict(value)
			}
		default:
			copyData[k] = g.GenStatic(fmt.Sprintf("%v", value))
		}
	}
	return copyData
}

// GenerateObject generates objects based on a schema
func (g *Gen) GenerateObject(schema map[string]interface{}, amount int) []map[string]interface{} {
	if amount <= 1 {
		return []map[string]interface{}{g.GenDict(schema)}
	}

	result := make([]map[string]interface{}, amount)
	for i := 0; i < amount; i++ {
		result[i] = g.GenDict(schema)
	}
	return result
}

// Helper function to calculate integer power
func intPow(base, exp int) int {
	result := 1
	for exp > 0 {
		result *= base
		exp--
	}
	return result
}


// func main() {
// 	gen := Gen{}

// 	// Example schema
// 	schema := map[string]interface{}{
// 		"name": "name",
// 		"email": map[string]interface{}{
// 			"type":      "email",
// 			"domain":    "example.com",
// 			"_$amount":  2,
// 		},
// 		"age": "int(len=2)",
// 	}

// 	// Generate objects
// 	result := gen.GenerateObject(schema, 2)
// 	fmt.Println(result)
// }
