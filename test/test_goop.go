package main

import (
	"fmt"
	"log"

	goop "github.com/ez0000001000000/Goop"
)

func main() {
	fmt.Println("Testing Goop package with OpenAI.com...")

	// Test HTTP GET
	resp, err := goop.Get("https://openai.com")
	if err != nil {
		log.Fatal("Failed to fetch:", err)
	}

	// Test HTML parsing
	doc := goop.HTMLParse(resp)
	if doc.Error != nil {
		log.Fatal("Failed to parse HTML:", doc.Error)
	}

	fmt.Println("=== OpenAI.com Page Structure ===")

	// Test finding all major sections
	nav := doc.Find("nav")
	if nav.Error == nil {
		fmt.Println(" Found navigation element")
		navLinks := nav.FindAll("a")
		fmt.Printf("  Found %d navigation links\n", len(navLinks))
		for i, link := range navLinks {
			if i < 5 { // Show first 5
				text := link.Text()
				href := link.Attrs()["href"]
				fmt.Printf("    %d. %s -> %s\n", i+1, text, href)
			}
		}
	} else {
		fmt.Println(" No navigation found")
	}

	// Test finding headings
	headings := doc.FindAll("h1", "h2", "h3")
	fmt.Printf(" Found %d headings (h1, h2, h3)\n", len(headings))

	// Test finding buttons/CTAs
	buttons := doc.FindAll("button")
	fmt.Printf(" Found %d buttons\n", len(buttons))

	// Test finding links
	links := doc.FindAll("a")
	fmt.Printf(" Found %d total links\n", len(links))

	// Test finding images
	images := doc.FindAll("img")
	fmt.Printf(" Found %d images\n", len(images))

	// Test finding forms
	forms := doc.FindAll("form")
	fmt.Printf(" Found %d forms\n", len(forms))

	// Test finding scripts
	scripts := doc.FindAll("script")
	fmt.Printf(" Found %d scripts\n", len(scripts))

	// Test finding meta tags
	metaTags := doc.FindAll("meta")
	fmt.Printf(" Found %d meta tags\n", len(metaTags))

	// Test page title
	title := doc.Find("title")
	if title.Error == nil {
		fmt.Printf(" Page title: %s\n", title.Text())
	} else {
		fmt.Println(" No title found")
	}

	// Test body content
	body := doc.Find("body")
	if body.Error == nil {
		bodyText := body.FullText()
		if len(bodyText) > 200 {
			fmt.Printf(" Body content (first 200 chars): %s...\n", bodyText[:200])
		} else {
			fmt.Printf(" Body content: %s\n", bodyText)
		}
	}

	// Test HTML structure
	fmt.Println("\n=== HTML Structure Sample ===")
	if len(links) > 0 {
		fmt.Printf("First link HTML: %s\n", links[0].HTML())
	}
	if len(images) > 0 {
		fmt.Printf("First image HTML: %s\n", images[0].HTML())
	}

	fmt.Printf("\n Goop package successfully scraped OpenAI.com!\n")
	fmt.Printf(" Summary: %d links, %d images, %d buttons, %d forms\n",
		len(links), len(images), len(buttons), len(forms))
}
