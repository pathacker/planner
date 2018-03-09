// Planner reads the named files, default $HOME/lib/calendar,
// and writes to standard output in calendar order any lines
// containing matching dates for today and tomomrrow.
// The '-n days' flag changes the number of days compared.
// No special processing is done for weekends.
//
// Recognized date formats are "4/26", "Apr 26", "26 April".
// Only the first three runes of the month name are matched.
//
// All comparisions are case insensitive.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/user"
	"regexp"
	"sort"
	"strings"
	"time"
)

type regexpItem struct {
	day int
	ra  []*regexp.Regexp
}

type result struct {
	n int
	s string
}
type byDay []result

func (a byDay) Len() int           { return len(a) }
func (a byDay) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDay) Less(i, j int) bool { return a[i].n < a[j].n }

func checkDates(fn string, ria []regexpItem) {
	var results []result
	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	a := strings.Split(string(buf), "\n")
	for _, s := range a {
		for _, ri := range ria {
			for _, r := range ri.ra {
				if r.MatchString(strings.ToLower(s)) {
					results = append(results, result{ri.day, s})
				}
			}
		}
	}
	sort.Sort(byDay(results))
	for _, r := range results {
		fmt.Printf("%s\n", r.s)
	}
}

func buildRegexpArray(n int) []regexpItem {
	var ria []regexpItem
	var sa [3]string
	t := time.Now()
	for i := 1; i <= n; i++ {
		var ra []*regexp.Regexp
		mn := strings.ToLower(t.Month().String())[0:3]
		sa[0] = fmt.Sprintf(`(^|[^a-z])%s[a-z]* %d($|[^\d])`, mn, t.Day())
		sa[1] = fmt.Sprintf(`(^|[^\d])%d %s[a-z]*`, t.Day(), mn)
		sa[2] = fmt.Sprintf(`(^|[^\d]|\d/)%d/%d[^\d]`, t.Month(), t.Day())
		for _, s := range sa {
			ra = append(ra, regexp.MustCompile(s))
		}
		ri := regexpItem{
			day: i,
			ra:  ra,
		}
		ria = append(ria, ri)
		t = t.Add(time.Hour * 24)
	}
	return ria
}

func main() {
	// calendar [-n days] [file...]
	var np = flag.Int("n", 2, "number of days to print") // today and tomorrow
	flag.Parse()
	ria := buildRegexpArray(*np)
	if flag.NArg() == 0 {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		checkDates(fmt.Sprintf("%s/lib/calendar", u.HomeDir), ria)
	} else {
		a := flag.Args()
		for _, s := range a {
			checkDates(s, ria)
		}
	}
}
