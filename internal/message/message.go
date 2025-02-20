package message

import (
	"errors"
	"fmt"
	"strings"
)

type AuthMessage struct {
	Username string
	Password string
}

type SelectQuery struct {
	Columns   []string
	TableName string
	// Where     string
}

// username:password
func (m AuthMessage) ToBytes() []byte {
	return []byte(m.Username + ":" + m.Password)
}

func ParseAuthMessage(data string) (AuthMessage, error) {
	// username:password, min length 3 (a:b)
	if len(data) < 3 {
		return AuthMessage{}, errors.New("invalid data format")
	}

	split := strings.Split(data, ":")
	if len(split) != 2 {
		return AuthMessage{}, errors.New("invalid data format")
	}

	return AuthMessage{
		Username: split[0],
		Password: split[1],
	}, nil
}

func ParseSelectQuery(parts *[]string) (SelectQuery, error) {
	if len(*parts) < 4 { // at least 4 parts: SELECT * FROM table
		return SelectQuery{}, errors.New("invalid query, use at least: SELECT * FROM table")
	}

	isDistinct := false

	if strings.ToUpper((*parts)[1]) == "DISTINCT" {
		if len(*parts) < 5 { // at least 5 parts: SELECT DISTINCT * FROM table
			return SelectQuery{}, errors.New("invalid distinct query, use at least: SELECT DISTINCT * FROM table")
		}

		isDistinct = true
	}

	columns := []string{}

	startIndex := 1
	if isDistinct {
		startIndex = 2
	}

	tableNameIndex := len(*parts) - 1
	for i := startIndex; i < len(*parts); i++ {
		if strings.ToUpper((*parts)[i]) == "FROM" {
			tableNameIndex = i + 1
			break
		}

		columns = append(columns, strings.TrimRight((*parts)[i], ",")) // remove trailing comma
	}

	if len(columns) == 0 {
		return SelectQuery{}, errors.New("no columns specified for selection")
	}

	table := (*parts)[tableNameIndex]
	if table == "" {
		return SelectQuery{}, errors.New("no table specified for selection")
	}

	// TO DO: handle where clause

	if len(*parts) > (tableNameIndex + 1) {
		return SelectQuery{}, errors.New("invalid query length")
		// if strings.ToUpper(parts[tableNameIndex + 1]) == "WHERE" {
		//     if len(parts) < (tableNameIndex + 4) {
		//         return errors.New("invalid query, use at least: SELECT * FROM table WHERE column = value")
		//     }
		//     // handle where clause
		//     fmt.Println("Where clause: ", parts[tableNameIndex + 2], parts[tableNameIndex + 3], parts[tableNameIndex + 4])
		// }
	}

	selectQuery := SelectQuery{
		Columns:   columns,
		TableName: table,
	}

	return selectQuery, nil
}

func MatchQuery(data *string) error {
	trimmedData := strings.TrimSpace(*data)
	parts := strings.Fields(trimmedData)

	fmt.Println("Parts: ", parts)

	queryType := strings.ToUpper(parts[0])

	switch queryType {
	case "SELECT":
		selectQuery, err := ParseSelectQuery(&parts)
		if err != nil {
			return err
		}

		fmt.Println("Select Query: ", selectQuery)

		return nil
	default:
		return errors.New("unsupported query type")
	}
}
