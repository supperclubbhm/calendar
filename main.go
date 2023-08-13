package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

var circledDates = map[time.Month][]int{
	time.September: {9},
	time.November:  {3},
}

func generateCalendar(year int, month time.Month) template.HTML {
	var builder strings.Builder

	fmt.Fprintf(&builder, `<h2>%s %d</h2><table class="calendar">
		<tr>
			<th>Sun</th><th>Mon</th><th>Tue</th><th>Wed</th><th>Thu</th><th>Fri</th><th>Sat</th>
		</tr>
		<tr>`, month, year)

	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	weekday := int(firstDay.Weekday())

	for day := 1; day <= daysInMonth; day++ {
		if day == 1 {
			fmt.Fprintf(&builder, strings.Repeat("<td></td>", weekday))
		}
		circle := ""
		if _, ok := circledDates[month]; ok {
			for _, d := range circledDates[month] {
				if day == d {
					circle = "*"
					break
				}
			}
		}
		fmt.Fprintf(&builder, "<td>%s%d</td>", circle, day)
		if (day+weekday)%7 == 0 || day == daysInMonth {
			fmt.Fprintf(&builder, "</tr>")
		}
	}

	fmt.Fprintf(&builder, `</table>`)
	return template.HTML(builder.String())
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var calMonths [][]template.HTML
	currentYear := 2023
	currentMonth := time.August
	for i := 0; i < 6; i++ {
		var months []template.HTML
		year := currentYear
		month := currentMonth + time.Month(i)
		if month > time.December {
			month -= 12
			year++
		}
		months = append(months, generateCalendar(year, month))
		calMonths = append(calMonths, months)
	}

	data := struct {
		Years [][]template.HTML
	}{
		Years: calMonths,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)

	port := 8080
	fmt.Printf("Server started on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
