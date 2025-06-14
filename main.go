package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var flagsjson, flagojson string

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func Atoi(s string) int64 {
	var (
		n uint64
		i int
		v byte
	)
	for ; i < len(s); i++ {
		d := s[i]
		// fmt.Println(d)
		if '0' <= d && d <= '9' {
			v = d - '0'
		} else if 'a' <= d && d <= 'z' {
			v = d - 'a' + 10
		} else if 'A' <= d && d <= 'Z' {
			v = d - 'A' + 10
		} else {
			n = 0
			break
		}
		n *= uint64(10)
		n += uint64(v)
	}
	return int64(n)
}

func convertDateToISO(data string) (string, int) {
	// fmt.Println("date conversion in progess")
	var count int = 0
	var re = regexp.MustCompile(`(?m)\{"\$date":\{"\$numberLong":"(?:[^\\"]|\\\\|\\")*"\}\}`)
	for _, match := range re.FindAllString(data, -1) {
		// fmt.Println(match, "date found at index", i)
		res := strings.Split(match, `"`)
		res[5] = strings.ReplaceAll(res[5], `-`, ``)

		tm := Atoi(res[5])
		t := time.Unix(0, tm*int64(time.Millisecond))
		data = strings.ReplaceAll(string(data), match, `ISODate(`+t.Format("2006-01-02T06:01:17.171Z")+`)`)
		count++
		// fmt.Print(res[5])
	}
	// fmt.Print(count)
	return data, count
}
func numberLongToInt(data string) (string, int) {
	// fmt.Println("date conversion in progess")
	var count int = 0

	var renum = regexp.MustCompile(`\{"\$numberLong":"(?:[^"]|"")*"\}`)
	for _, match := range renum.FindAllString(data, -1) {
		// fmt.Println(match, "$longint found at index", i)

		res := strings.Split(match, `"`)
		// fmt.Println(res[3])
		data = strings.Replace(string(data), match, res[3], -1)
		count++
	}
	// fmt.Print(count)
	return data, count
}
func numberIntToInt(data string) (string, int) {
	var count int = 0

	var renum = regexp.MustCompile(`\{"\$numberInt":"(?:[^"]|"")*"\}`)
	for _, match := range renum.FindAllString(data, -1) {
		// fmt.Println(match, "$numint found at index", i)

		res := strings.Split(match, `"`)
		data = strings.Replace(data, match, res[3], -1)
		count++
		// fmt.Println(res[3])
	}
	// fmt.Print(count)
	return data, count
}

func numberDoubleToDouble(data string) (string, int) {
	var count int = 0

	var renum = regexp.MustCompile(`\{"\$numberDouble":"(?:[^"]|"")*"\}`)
	for _, match := range renum.FindAllString(data, -1) {
		// fmt.Println(match, "$numdouble found at index", i)

		res := strings.Split(match, `"`)
		data = strings.Replace(data, match, res[3], -1)
		count++
		// fmt.Println(res[3])
	}
	// fmt.Print(count)
	return data, count
}
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func checkExtention(path string) bool {
	//.json
	if len(path) < 6 {
		return false
	}
	return strings.ToLower(path[len(path)-5:]) == ".json"
}
func init() {
	flag.StringVar(&flagsjson, "source", "export.json", "mongo exported json file")
	flag.StringVar(&flagojson, "output", "output.json", "output json file")
	flag.Parse()
	fmt.Println("------init")
}

func main() {
	// CheckArgs("<sourceJson>", "<outputJson>")
	// sourceJson, outputJson := os.Args[1], os.Args[2]
	if flagsjson == "" || flagojson == "" {
		fmt.Println("Please provide mongo exported json file!!   --source xyz.json ")
		return
	}
	isExist, err := exists(flagsjson)
	if err != nil {
		panic(err)
	}
	if !isExist || !checkExtention(flagsjson) {
		fmt.Println("Please provide valid mongo exported json file!!   --source xyz.json ")
		return
	}
	// checkExtention(flagsjson)
	fmt.Println(flagsjson + "------" + flagojson)
	var data string = ""
	f, err := os.OpenFile(flagsjson, os.O_RDONLY, os.ModePerm)
	check(err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		s = strings.ReplaceAll(s, ": ", ":")
		data += s

	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}
	var dateCount, numberLong, numberInt, numberDouble int
	data, dateCount = convertDateToISO(data)
	data, numberLong = numberLongToInt(data)
	data, numberInt = numberIntToInt(data)
	data, numberDouble = numberDoubleToDouble(data)
	total := dateCount + numberInt + numberLong + numberDouble
	fmt.Println("Total changes: ", total)
	// ioutil.WriteFile(flagojson, []byte(data), os.ModePerm)
	err = os.WriteFile(flagojson, []byte(data), os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}
}
