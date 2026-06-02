package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
)

type Video struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func SearchTikTok(query string, limit int) ([]Video, error) {
	if err := playwright.Run(); err != nil {
		return nil, err
	}

	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return nil, err
	}

	url := "https://www.tiktok.com/search?q=" + query

	_, err = page.Goto(url)
	if err != nil {
		return nil, err
	}

	time.Sleep(10 * time.Second)

	js := `
	() => {
		const items = [];
		document.querySelectorAll("a").forEach(el => {
			const link = el.href;
			if(link && link.includes("/video/")){
				items.push({
					title: el.innerText || "TikTok Video",
					url: link
				});
			}
		});
		return items;
	}
	`

	result, err := page.Evaluate(js)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(result)

	var videos []Video
	json.Unmarshal(data, &videos)

	if len(videos) > limit {
		videos = videos[:limit]
	}

	return videos, nil
}

func main() {
	query := "funny"

	if len(os.Args) > 1 {
		query = os.Args[1]
	}

	results, err := SearchTikTok(query, 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonData, _ := json.Marshal(results)
	fmt.Println(string(jsonData))
}