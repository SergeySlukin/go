package main

import (
	"regexp"
	"strings"
	"strconv"
	"errors"
)

var uniqueUrls = make(map[string]uint64)

var urlRegexp = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
var httpInfo = regexp.MustCompile(`HTTP\/[0-9](\.[0-9])?\"\s[0-9]*\s[0-9]*`)
var searchEngineRegexp = regexp.MustCompile(`\)\s[a-zA-Z0-9]*`)

/**
 Обрабатываем строки
 */
func stringHandler(s string) {
	res, ok := parseUrl(s)
	if !ok {
		panic(errors.New("error parse url"))
	}
	statusCode, bytes, _ := parseHttpInfo(s)

	searchEngine, ok := parseSearchEngine(s)
	if !ok {
		searchEngine = "direct"
	}
	addStatistics(res, statusCode, bytes, searchEngine)
}

/**
Парсим урл
 */
func parseUrl(s string) (string, bool) {
	findUrl := urlRegexp.FindString(s)
	if len(findUrl) > 0 {
		return findUrl, true
	}
	return "", false
}

/**
Парсим статус и кол-во байт
 */
func parseHttpInfo(s string) (uint16, uint64, error) {
	bytes := 0
	respInfo := httpInfo.FindString(s)
	respArray := strings.Split(respInfo, " ")
	statusCode, err := strconv.Atoi(respArray[1])
	if err != nil {
		panic(err)
	}

	if len(respArray[2]) > 0 {
		bytes, err = strconv.Atoi(respArray[2])
		if err != nil {
			panic(err)
		}
	}
	return uint16(statusCode), uint64(bytes), nil
}

/**
Парсим поисковик
 */
func parseSearchEngine(s string) (string, bool) {
	searchEngine := searchEngineRegexp.FindString(s)
	if len(searchEngine) > 0 {
		return searchEngine[2:], true
	}
	return "", false
}

/**
Добавляем статистику
 */
func addStatistics(url string, status uint16, bytes uint64, searchEngine string) {
	mu.Lock()
	defer mu.Unlock()
	uniqueUrls[url]++
	statistics.StatusCodes[status]++
	statistics.Traffic += bytes
	statistics.Crawlers[searchEngine]++
}
