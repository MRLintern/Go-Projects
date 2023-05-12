
//--web application which allows users to enter a list of numbers and then performs some statistical calculations--//

package main

//get packages for application

import (

	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

//HTML template for UI

const (

	pageTop = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head><title>Statistics</title>
<body><h3>Statistics</h3>
<p>Computes basic statistics for a given list of numbers</p>`

    form = `<form action="/" method="POST">

<label for="numbers">Numbers (comma or space-separated):</label><br />
<input type="text" name="numbers" size="30"><br />
<input type="submit" value="Calculate">
</form>`

    pageBottom = `</body></html>`
    anError = `<p class="error">%s</p>`
)

//struct for the statistics to be calculated
type statistics struct {

	numbers []float64 //list of numbers entered by user
	average float64
	median float64
}

func main() {

	//this function takes a path and a reference to a function
	//to call when the path is requested; the / path (web app homepage)
	http.HandleFunc("/", homePage)

	//starts up a web server at the given TCP network address; localhost and port number 9001
	if err := http.ListenAndServe(":9001", nil); err != nil {

		log.Fatal("failed to start server", err)
	}
}

//function is called whenever the statistics web site is visited.
func homePage(writer http.ResponseWriter, request *http.Request) {

    err := request.ParseForm() // Must be called before writing response

    fmt.Fprint(writer, pageTop, form)

    if err != nil {

        fmt.Fprintf(writer, anError, err)

    } else {

        if numbers, message, ok := processRequest(request); ok {

            stats := getStats(numbers)

            fmt.Fprint(writer, formatStats(stats))

        } else if message != "" {

            fmt.Fprintf(writer, anError, message)
        }
    }

    fmt.Fprint(writer, pageBottom)
}

//function to read the forms data from the request value
func processRequest(request *http.Request) ([]float64, string, bool) {

    var numbers []float64

    if slice, found := request.Form["numbers"]; found && len(slice) > 0 {

        text := strings.Replace(slice[0], ",", " ", -1)

        for _, field := range strings.Fields(text) {

        	//convert string to float64; can be float32 
            if x, err := strconv.ParseFloat(field, 64); err != nil {

                return numbers, "'" + field + "' is invalid", false

            } else {

                numbers = append(numbers, x)
            }
        }
    }

    if len(numbers) == 0 {

        return numbers, "", false // no data first time form is shown
    }

    return numbers, "", true
}

//function allows us display the results online
func formatStats(stats statistics) string {

    return fmt.Sprintf(`<table border="1">

<tr><th colspan="2">Results</th></tr>
<tr><td>Numbers</td><td>%v</td></tr>
<tr><td>Count</td><td>%d</td></tr>
<tr><td>Average</td><td>%f</td></tr>
<tr><td>Median</td><td>%f</td></tr>
</table>`, stats.numbers, len(stats.numbers), stats.average, stats.median)

}

//accepts a slice (dynamic array in C) of numbers.
//the function will then populate the stats result
func getStats(numbers []float64) (stats statistics) {

    stats.numbers = numbers

    //calculating the median (middle number); use sort package; 
    //numbers need to sorted in ascending order
    sort.Float64s(stats.numbers)

    stats.average = sum(numbers) / float64(len(numbers))

    stats.median = median(numbers)

    return stats
}

//function to iterate over the numbers produce the sum of these numbers
func sum(numbers []float64) (total float64) {

    for _, x := range numbers {

        total += x
    }

    return total
}

//function to calculate the median statistic
func median(numbers []float64) float64 {

    middle := len(numbers) / 2

    result := numbers[middle]

    if len(numbers)%2 == 0 {

        result = (result + numbers[middle-1]) / 2
    }

    return result
}


