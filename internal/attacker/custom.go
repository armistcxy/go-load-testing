package attacker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/go-faker/faker/v4"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Field struct {
	Name string
	Type string
}

type CustomTargeterBuilder struct {
	URL    string
	Method string
	Header map[string]string
	Fields []Field
}

func NewCustomTargetBuilder(fig *Figure) CustomTargeterBuilder {
	fields := make([]Field, 0)
	for fieldName, fieldType := range fig.Fields {
		fields = append(fields, Field{fieldName, fieldType})
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	return CustomTargeterBuilder{
		URL:    fig.URL,
		Method: fig.Method,
		Header: fig.Header,
		Fields: fields,
	}
}

func (ctb CustomTargeterBuilder) BuildCustomTargeter() (vegeta.Targeter, error) {
	// Prepare payload
	payload := make(map[string]interface{})

	for _, field := range ctb.Fields {
		if randomGenerator, exist := generators[field.Type]; exist {
			payload[field.Name] = randomGenerator()
		} else {
			return nil, fmt.Errorf("no equivalent random generators for type %s", field.Type)
		}
	}

	// Prepare header
	header := http.Header{}
	for key, val := range ctb.Header {
		header.Set(key, val)
		header.Add(key, val)
	}

	return func(target *vegeta.Target) error {
		if target == nil {
			return vegeta.ErrNilTarget
		}

		target.URL = ctb.URL
		target.Method = ctb.Method
		target.Header = header

		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		target.Body = body

		return nil
	}, nil
}

type RandomMethod func() interface{}

var generators = map[string]func() interface{}{
	// Geographical
	"Latitude":  func() interface{} { return faker.Latitude() },
	"Longitude": func() interface{} { return faker.Longitude() },
	"Address":   func() interface{} { return faker.GetRealAddress() },

	// Datetime
	"UnixTime":   func() interface{} { return faker.UnixTime() },
	"Date":       func() interface{} { return faker.Date() },
	"TimeString": func() interface{} { return faker.TimeString() },
	"MonthName":  func() interface{} { return faker.MonthName() },
	"YearString": func() interface{} { return faker.YearString() },
	"DayOfWeek":  func() interface{} { return faker.DayOfWeek() },
	"DayOfMonth": func() interface{} { return faker.DayOfMonth() },
	"Timestamp":  func() interface{} { return faker.Timestamp() },
	"Century":    func() interface{} { return faker.Century() },
	"Timezone":   func() interface{} { return faker.Timezone() },
	"Timeperiod": func() interface{} { return faker.Timeperiod() },

	// Internet
	"Email":      func() interface{} { return faker.Email() },
	"MacAddress": func() interface{} { return faker.MacAddress() },
	"DomainName": func() interface{} { return faker.DomainName() },
	"URL":        func() interface{} { return faker.URL() },
	"Username":   func() interface{} { return faker.Username() },
	"IPv4":       func() interface{} { return faker.IPv4() },
	"IPv6":       func() interface{} { return faker.IPv6() },
	"Password":   func() interface{} { return faker.Password() },

	// Words and Sentences
	"Word":      func() interface{} { return faker.Word() },
	"Sentence":  func() interface{} { return faker.Sentence() },
	"Paragraph": func() interface{} { return faker.Paragraph() },

	// Payment
	"CCType":             func() interface{} { return faker.CCType() },
	"CCNumber":           func() interface{} { return faker.CCNumber() },
	"Currency":           func() interface{} { return faker.Currency() },
	"AmountWithCurrency": func() interface{} { return faker.AmountWithCurrency() },

	// Person
	"TitleMale":       func() interface{} { return faker.TitleMale() },
	"TitleFemale":     func() interface{} { return faker.TitleFemale() },
	"FirstName":       func() interface{} { return faker.FirstName() },
	"FirstNameMale":   func() interface{} { return faker.FirstNameMale() },
	"FirstNameFemale": func() interface{} { return faker.FirstNameFemale() },
	"LastName":        func() interface{} { return faker.LastName() },
	"Name":            func() interface{} { return faker.Name() },

	// Phone
	"Phonenumber":         func() interface{} { return faker.Phonenumber() },
	"TollFreePhoneNumber": func() interface{} { return faker.TollFreePhoneNumber() },
	"E164PhoneNumber":     func() interface{} { return faker.E164PhoneNumber() },

	// UUID
	"UUIDHyphenated": func() interface{} { return faker.UUIDHyphenated() },
	"UUIDDigit":      func() interface{} { return faker.UUIDDigit() },
}
