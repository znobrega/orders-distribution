package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Order struct {
	CustomerName    string
	Contact         string
	ShippingAddress string
	GrandTotal      float64
	DateOrderPlaced time.Time
	Items           []Item
}

type Item struct {
	Cost        float64
	ShippingFee float64
	TaxAmount   float64
	Product     Product
}

type Product struct {
	Name         string
	Category     string
	Weight       int
	Price        float64
	CreationDate time.Time
}

const (
	Range1and3Months         = "1-3 months"
	Range4and6Months         = "4-6 months"
	Range7and12Months        = "7-12 months"
	RangeGreaterThan12Months = ">12 months"

	dateLayout = "2006-01-02 15:04:05"
)

func main() {
	orders := buildMockOrders()

	initData, endData, err := getDateRange()
	if err != nil {
		log.Fatal(err)
	}

	monthsCounters := distributeOrders(orders, initData, endData)
	fmt.Println(monthsCounters)
}

func getDateRange() (time.Time, time.Time, error) {
	if len(os.Args) < 3 {
		return time.Time{}, time.Time{}, fmt.Errorf("missing date arguments")
	}

	initialDate := os.Args[1]
	initDate, err := time.Parse(dateLayout, initialDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	finalDate := os.Args[2]
	endDate, err := time.Parse(dateLayout, finalDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return initDate, endDate, nil
}

func distributeOrders(orders []Order, initData time.Time, endData time.Time) map[string]int {
	monthsCounters := make(map[string]int, 0)
	for _, order := range orders {
		if !order.DateOrderPlaced.After(initData) {
			continue
		}

		if !order.DateOrderPlaced.Before(endData) {
			continue
		}

		getItemTimeRange(order, initData, endData, monthsCounters)
	}
	return monthsCounters
}

func getItemTimeRange(order Order, initData time.Time, endData time.Time, monthsCounters map[string]int) {
	for _, item := range order.Items {
		if !item.Product.CreationDate.After(initData) {
			continue
		}

		if !item.Product.CreationDate.Before(endData) {
			continue
		}

		months := int((time.Now().Sub(item.Product.CreationDate).Hours() / 24) / 30)

		switch {
		case months >= 1 && months <= 3:
			monthsCounters[Range1and3Months]++
		case months >= 4 && months <= 6:
			monthsCounters[Range4and6Months]++
		case months >= 7 && months <= 12:
			monthsCounters[Range7and12Months]++
		case months > 12:
			monthsCounters[RangeGreaterThan12Months]++
		}
	}
}

func monthsToSubtract(months int) time.Duration {
	hourlyMonths := -720 * months
	return time.Duration(hourlyMonths) * time.Hour
}

func buildMockOrders() []Order {
	item := Item{
		Cost:        10,
		ShippingFee: 1,
		TaxAmount:   1,
		Product: Product{
			Name:         "product 1",
			Category:     "health",
			Weight:       1,
			Price:        10,
			CreationDate: time.Now().Add(-2),
		},
	}

	firstOrder := Order{
		CustomerName:    "carlos",
		Contact:         "nobregacarlos@gmail.com",
		ShippingAddress: "brazil",
		GrandTotal:      100,
		DateOrderPlaced: time.Now().Add(monthsToSubtract(1)),
		Items: []Item{
			item,
		},
	}

	var items []Item
	for i := 0; i < 30; i++ {
		newItem := item
		newItem.Product.CreationDate = time.Now().Add(monthsToSubtract(i))
		items = append(items, newItem)
	}

	firstOrder.Items = append(firstOrder.Items, items...)

	var orders []Order
	for i := 0; i < 10; i++ {
		newOrder := firstOrder
		newOrder.DateOrderPlaced = time.Now().Add(monthsToSubtract(i))
		orders = append(orders, newOrder)
	}

	return orders
}
