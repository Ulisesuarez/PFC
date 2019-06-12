package util

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"bitbucket.org/inehealth/idonia-common/common"
)

//DefaultElementPerPage if the limit per page is not defined, it will be used
const DefaultElementPerPage = 10

//PageKeyQS is the default key in the querystring for the value of the number page requested
const PageKeyQS = "page"

//CountKeyQS is the default key in the querystring for the value of the number of elements requested per page
const CountKeyQS = "count"

//SortKeyQS is the default key in the querystring for the sort
const SortKeyQS = "sort[]"

//GetPaginatorFromRq Extract the paginator from the querystring
func GetPaginatorFromRq(r *http.Request, i interface{}) (reqPage int, reqCount int, orders []Order) {

	reqPage, err := strconv.Atoi(r.URL.Query().Get(PageKeyQS))
	if err != nil || reqPage < 1 {
		reqPage = 1
	}

	reqCount, err = strconv.Atoi(r.URL.Query().Get(CountKeyQS))
	if err != nil || reqCount < 1 {
		reqCount = DefaultElementPerPage
	}

	queryString, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Printf("Error trying to parse the querystring: %s", err.Error())
	}

	// {url}?sort[]=-field1&sort[]=field2&sort[]=-field3...
	if sort := []string(queryString[SortKeyQS]); len(sort) > 0 {

		//extract the attr binding between json and model representation
		dataStruct := common.ExtractStructBinding(i)
		for _, field := range sort {
			direction := Ascendant
			field = strings.TrimSpace(field)
			if strings.HasPrefix(field, "-") {
				field = field[1:] //removes the prefix
				direction = Descendant
			}
			//Skips the field if not exists
			if val, ok := dataStruct[field]; ok {
				orders = append(orders, Order{Field: strings.ToLower(val), Type: direction})
			}

		}
	}

	return reqPage, reqCount, orders
}

//GetTotalPages calculates the number of the total page for the result
func GetTotalPages(totalItems int, reqCount int) int {

	//reqCount never should be < 1...
	if totalItems < 1 || reqCount < 1 {
		return 0
	}

	increment := 0
	if math.Mod(float64(totalItems), float64(reqCount)) > 0 {
		increment = 1
	}

	return int(math.Floor(float64(totalItems)/float64(reqCount))) + increment

}

//CalculateOffset calculates the offset for the query
func CalculateOffset(page int, itemsPerPage int) int {
	offset := (page - 1) * itemsPerPage

	if offset < 0 {
		return 0
	}

	return offset
}

//GenerateSQLByOrders generates the sql from a Order[]
func GenerateSQLByOrders(orders []Order) string {

	if len(orders) == 0 {
		return ""
	}

	firstElement := true
	var sqlBuffer bytes.Buffer
	for _, order := range orders {

		//if the field is not set we must continue with the next...
		if len(order.Field) == 0 {
			continue
		}

		if firstElement {
			sqlBuffer.WriteString(" ORDER BY ")
			firstElement = false
		} else {
			sqlBuffer.WriteString(" , ")
		}

		sqlBuffer.WriteString(order.Field + " " + order.Type)

	}

	return sqlBuffer.String()
}
