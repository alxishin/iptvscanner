package tools

import (
	"encoding/json"
	"fmt"
	"github.com/pawanpaudel93/go-m3u-parser/m3uparser"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

func Iptvscanner() {
	// userAgent and timeout is optional. default timeout is 5 seconds and userAgent is latest chrome version 86.
	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36"
	timeout := 5 // in seconds

	parser := m3uparser.M3uParser{UserAgent: userAgent, Timeout: timeout}
	// file path can also be used /home/pawan/Downloads/ru.m3u
	parser.ParseM3u("https://smolnp.github.io/IPTVru/IPTVru.m3u", true, true)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	var result = make([]m3uparser.Channel, len(parser.GetStreamsSlice()))
	//var result []m3uparser.Channel
	var wg sync.WaitGroup

	for i, stream := range parser.GetStreamsSlice() {
		wg.Add(1)
		stream := stream
		go func(i int) {
			defer wg.Done()
			tmp := worker(client, stream)
			if tmp != nil {
				//result = append(result, tmp)
				result[i] = tmp
			}
		}(i)
	}
	wg.Wait()
	result = slices.DeleteFunc(
		result,
		func(thing m3uparser.Channel) bool {
			return thing == nil
		},
	)

	result = slices.Clip(result)
	ToFile("ru.m3u", result)
}

func worker(client http.Client, stream m3uparser.Channel) m3uparser.Channel {
	request, _ := http.NewRequest("GET", stream["url"].(string), nil)
	response, err := client.Do(request)

	if err == nil && response.StatusCode == 200 {
		headers := []string{
			"application/vnd.apple.mpegurl",
			"application/octet-stream",
			"audio/mpeg",
			"audio/aacp",
			"application/x-mpegurl",
			"application/vnd.apple.mpegurl; charset=utf-8"}
		//s[strings.ToLower(response.Header.Get("Content-Type"))] = strings.ToLower(response.Header.Get("Content-Type"))
		if slices.Contains(headers, strings.ToLower(response.Header.Get("Content-Type"))) {
			return stream
		}
	}
	return nil
}

func ToFile(fileName string, streamsInfo []m3uparser.Channel) {
	var format string
	if len(strings.Split(fileName, ".")) > 1 {
		format = strings.ToLower(strings.Split(fileName, ".")[1])
	}

	if format == "json" {
		json, _ := json.MarshalIndent(streamsInfo, "", "    ")
		json = []byte(strings.ReplaceAll(string(json), `: ""`, ": null"))
		if !strings.Contains(fileName, "json") {
			fileName = fileName + ".json"
		}
		writeData(fileName, json)
	} else if format == "m3u" {
		content := []string{"#EXTM3U"}
		for _, stream := range streamsInfo {
			line := "#EXTINF:-1"
			if tvg, ok := stream["tvg"]; ok {
				tvg := tvg.(map[string]string)
				for key, val := range tvg {
					if val != "" {
						line += fmt.Sprintf(` tvg-%s="%s"`, key, val)
					}
				}
			}
			if logo, ok := stream["logo"]; ok && logo != "" {
				line += fmt.Sprintf(` tvg-logo="%s"`, logo)
			}
			if country, ok := stream["country"]; ok {
				country := country.(map[string]string)
				if code, ok := country["code"]; ok && code != "" {
					line += fmt.Sprintf(` tvg-country="%s"`, code)
				}
			}
			if language, ok := stream["language"]; ok && language != "" {
				line += fmt.Sprintf(` tvg-language="%s"`, language)
			}
			if category, ok := stream["category"]; ok && category != "" {
				line += fmt.Sprintf(` group-title="%s"`, category)
			}
			if title, ok := stream["title"]; ok && title != "" {
				line += fmt.Sprintf(`,%s`, title)
			}
			content = append(content, line)
			content = append(content, stream["url"].(string))
		}
		writeData(fileName, []byte(strings.Join(content, "\n")))
	} else {
		log.Println("File extension not present/supported !!!")
	}
}

func writeData(fileName string, data []byte) {
	os.WriteFile(fileName, data, 0666)
}
