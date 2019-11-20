package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

type Manga struct {
	Title        string
	Alternative  string
	Released     string
	Author       string
	Type         string
	TotalChapter string
	Status       string
	Rating       string
	Image        string
	Description  string
}
type Chapter struct {
	Title  string
	Number string
	Images []string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	var manga []Manga
	res, err := http.Get(os.Getenv("LIST"))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	var mangaLink []string
	doc.Find(".series").Each(func(i int, sel *goquery.Selection) {
		link, _ := sel.Attr("href")
		mangaLink = append(mangaLink, link)
		// fmt.Println(link)
	})
	mangaLength := len(mangaLink)

	for i := 0; i < mangaLength; i++ {

		fmt.Println(mangaLink[i])
		mangas := GetManga(mangaLink[i])

		manga = append(manga, mangas)

	}
}

func GetManga(url string) Manga {
	var manga Manga
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".infox").Children().Each(func(i int, sel *goquery.Selection) {
		if i == 0 {
			manga.Title = strings.TrimRight(sel.Text(), "Indonesia")
			manga.Title = strings.TrimRight(manga.Title, " Bahasa")
		}

		if i == 1 {
			manga.Alternative = sel.Text()
		}

		var a []string
		// var b string
		if i == 2 {
			scanner := bufio.NewScanner(strings.NewReader(sel.Text()))
			for scanner.Scan() {
				a = append(a, scanner.Text())
			}
			var info = make(map[string]interface{})
			for _, value := range a {
				header := value[:strings.IndexByte(value, ':')]
				val := value[strings.LastIndex(value, ":"):]
				val = strings.Replace(val, ": ", "", 1)
				info[header] = val
			}
			manga.Released = info["Released"].(string)
			manga.Author = info["Author"].(string)
			manga.Type = info["Type"].(string)
			manga.TotalChapter = info["Total Chapter"].(string)
			manga.TotalChapter = info["Status"].(string)

		}
	})

	doc.Find(".desc").Children().EachWithBreak(func(i int, sel *goquery.Selection) bool {
		if i == 1 {
			manga.Description = sel.Text()
			manga.Description = strings.Replace(manga.Description, "komikcast", "kami", -1)
			return false
		}
		return true
	})

	doc.Find(".rating").Children().EachWithBreak(func(i int, sel *goquery.Selection) bool {
		manga.Rating = strings.Replace(sel.Text(), "Rating ", "", 1)
		return false
	})

	doc.Find(".thumb").Children().EachWithBreak(func(i int, sel *goquery.Selection) bool {
		manga.Image, _ = sel.Attr("src")
		return false
	})
	var chapters []Chapter
	doc.Find(".leftoff").Children().Each(func(i int, sel *goquery.Selection) {
		link, _ := sel.Attr("href")
		chapter := GetChapter(link)
		chapters = append(chapters, chapter)
		fmt.Println(len(chapter.Images))
		fmt.Println()
	})
	return manga
}

func GetChapter(url string) Chapter {
	var chapter Chapter

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var images []string

	doc.Find("#readerarea").Children().Each(func(i int, sel *goquery.Selection) {
		c, _ := sel.Attr("src")
		if i > 2 {
			sel.Find("img").Each(func(i int, sel *goquery.Selection) {
				a, _ := sel.Attr("src")
				images = append(images, string(a))
			})
			a, _ := sel.Attr("src")
			if a != "" {
				images = append(images, string(a))
			} else {
				sel.Find("img").Children().Each(func(i int, sel *goquery.Selection) {
					sel.Children().Each(func(i int, sel *goquery.Selection) {
						a, _ := sel.Attr("src")
						images = append(images, string(a))
					})
				})
			}
		}
		if c != "" {
			sel.Children().Each(func(i int, sel *goquery.Selection) {
				a, _ := sel.Attr("src")
				if a != "" {
					images = append(images, string(a))
				}
			})
		}
	})

	e := doc.Find("title").Text()
	e = strings.TrimRight(e, "Komikcast")
	e = strings.Replace(e, "Bahasa Indonesia", "", 1)
	es := len(e)
	e = e[:es-4]
	chapter.Title = e
	chapter.Images = images
	fmt.Println(chapter.Title)
	fmt.Println(chapter)

	return chapter
}
