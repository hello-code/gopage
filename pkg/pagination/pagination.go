package pagination

import (
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

// Page ...
type Page struct {
	Index       int           // current page number
	Size        int           // per page size
	Total       int           // total pages count
	Rows        int           // data rows
	Numbers     int           // how many page numbers will show in html page
	NumberRange []int         // a slice for create page numbers
	Request     *http.Request // for request search conditions
}

// NewPage create a pager.
// rows is total data rows.
func NewPage(r *http.Request, rows int) *Page {
	// set default page index and size
	index, size, numbers := SetPageAndSize(r)

	if rows < size {
		return &Page{
			Index:   1,
			Size:    size,
			Total:   1,
			Rows:    rows,
			Numbers: numbers,
		}
	}

	p := new(Page)
	p.Request = r
	p.Index = index
	p.Size = size
	p.Total = getTotalPageNumbers(rows, size)

	p.Rows = rows
	p.Numbers = numbers
	p.NumberRange = createPageNumberRange(index, p.Total, numbers)

	return p
}

// SetPageAndSize set default value if no value pass to server.
func SetPageAndSize(r *http.Request) (index, size, numbers int) {
	// page index
	var i string
	if i = r.FormValue("index"); i == "" {
		i = "1"
	}
	index, err := strconv.Atoi(i)
	if err != nil {
		log.Println("convert:", err.Error())
	}

	// page size
	var s string
	if s = r.FormValue("size"); s == "" {
		s = "5"
	}
	size, err = strconv.Atoi(s)
	if err != nil {
		log.Println("convert:", err.Error())
	}

	// page numbers
	var n string
	if n = r.FormValue("nums"); n == "" {
		n = "5"
	}
	numbers, err = strconv.Atoi(n)
	if err != nil {
		log.Println("convert:", err.Error())
	}

	return
}

// total pages
func getTotalPageNumbers(total int, size int) int {
	num := (int)(total / size)
	if total%size > 0 {
		num++
	}
	return num
}

// Create a slice based on how many page numbers your want to
// display on the page.
func createPageNumberRange(index, total, numbers int) []int {
	var pageNumbers []int
	switch {
	case index > total-numbers && total > numbers:
		start := total - numbers + 1
		pageNumbers = make([]int, numbers)
		for i := range pageNumbers {
			pageNumbers[i] = start + i
		}
	case index >= numbers && total > numbers:
		start := index - numbers + 1
		pageNumbers = make([]int, int(math.Min(float64(numbers), float64(index+1))))
		for i := range pageNumbers {
			pageNumbers[i] = start + i
		}
	default:
		pageNumbers = make([]int, int(math.Min(float64(numbers), float64(total))))
		for i := range pageNumbers {
			pageNumbers[i] = i + 1
		}
	}
	return pageNumbers
}

// StartEnd is currently show data from ... to rows
func (p *Page) StartEnd() (start, end int) {
	if p.Index > 1 {
		start = p.Index*p.Size - p.Size + 1
		end = p.Index * p.Size
		if end > p.Rows {
			end = p.Rows
		}
		return
	}
	return 1, p.Size
}

// CurrentPage is the page link when you clicked
func (p *Page) CurrentPage(page int) string {
	link, err := url.ParseRequestURI(p.Request.RequestURI)
	if err != nil {
		log.Println(err.Error())
	}
	values := link.Query()
	if page == 1 {
		values.Del("index")
	} else {
		values.Set("index", strconv.Itoa(page))
	}
	link.RawQuery = values.Encode()
	return link.String()
}

// FirstPage is the first page
func (p *Page) FirstPage() (link string) {
	return p.CurrentPage(1)
}

// LastPage is the last page
func (p *Page) LastPage() (link string) {
	return p.CurrentPage(p.Total)
}

// PrevPage is previous page
func (p *Page) PrevPage() (link string) {
	if p.Index > 1 {
		link = p.CurrentPage(p.Index - 1)
	}
	return link
}

// NextPage is next page
func (p *Page) NextPage() (link string) {
	if p.Index < p.Total {
		link = p.CurrentPage(p.Index + 1)
	}
	return link
}
