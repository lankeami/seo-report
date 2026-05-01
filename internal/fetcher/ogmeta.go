package fetcher

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// FetchOGMeta fetches Open Graph metadata from the given URL.
// It never returns an error; failures result in an empty OGMeta.
func FetchOGMeta(ctx context.Context, rawURL string) OGMeta {
	if rawURL == "" {
		return OGMeta{}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return OGMeta{}
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; seo-report-bot/1.0; +https://github.com/jaychinthrajah/seo-report)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return OGMeta{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return OGMeta{}
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		return OGMeta{}
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return OGMeta{}
	}

	return extractOGMeta(doc)
}

func extractOGMeta(doc *html.Node) OGMeta {
	var meta OGMeta
	var metaDesc string // fallback: <meta name="description">
	var titleText string

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "meta":
				prop := attrVal(n, "property")
				name := attrVal(n, "name")
				content := attrVal(n, "content")
				switch prop {
				case "og:title":
					meta.OGTitle = content
				case "og:description":
					meta.OGDescription = content
				case "og:image":
					meta.OGImage = content
				}
				if strings.EqualFold(name, "description") && metaDesc == "" {
					metaDesc = content
				}
			case "title":
				if titleText == "" {
					titleText = nodeText(n)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	if meta.OGDescription == "" {
		meta.OGDescription = metaDesc
	}
	if meta.OGTitle == "" {
		meta.OGTitle = titleText
	}

	return meta
}

func attrVal(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func nodeText(n *html.Node) string {
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			sb.WriteString(c.Data)
		}
	}
	return strings.TrimSpace(sb.String())
}
