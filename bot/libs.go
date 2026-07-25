package bot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/rand"
	"gopkg.in/macaron.v1"
)

func generateCode() string {
	code := ""
	num := strconv.Itoa(generateRandNum(99))

	words, err := urlToLines(WordsUrl)
	if err != nil {
		loge(err)
	}

	code = words[generateRandNum(2047)]

	return code + num
}

func generateRandNum(max int) int {
	rand.Seed(uint64(time.Now().UnixNano()))
	min := 0
	rn := rand.Intn(max-min+1) + min
	return rn
}

func urlToLines(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return linesFromReader(resp.Body)
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func getTgId(ctx *macaron.Context) int64 {
	tgids := ctx.Params("telegramid")
	tgid, err := strconv.Atoi(tgids)
	if err != nil {
		loge(err)
	}
	return int64(tgid)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func getGroup() int64 {
	if conf.Dev {
		return GroupDev
	}
	return Group
}

func formatNumber(n float64) string {
	return formatNumberWithPrecision(n, 9)
}

func formatNumberWithPrecision(n float64, precision int) string {
	format := fmt.Sprintf("%%.%df", precision)
	str := fmt.Sprintf(format, n)

	// Split into integer and decimal parts
	parts := strings.Split(str, ".")
	if len(parts) == 0 {
		return str
	}

	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = "." + parts[1]
	}

	// Add commas to integer part
	if len(integerPart) > 3 {
		var result strings.Builder
		for i, char := range integerPart {
			if i > 0 && (len(integerPart)-i)%3 == 0 {
				result.WriteRune(',')
			}
			result.WriteRune(char)
		}
		integerPart = result.String()
	}

	return integerPart + decimalPart
}
