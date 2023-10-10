package jira

import "fmt"

type Description struct {
	Version int       `json:"version"`
	Type    string    `json:"type"`
	Content []Content `json:"content"`
}

func (d Description) String() string {
	var description string

	for _, content := range d.Content {
		description += content.String(false)
	}

	return description
}

type Content struct {
	Type    string    `json:"type"`
	Text    string    `"json:"text"`
	Attrs   Attrs     `"json:"attrs`
	Marks   []Marks   `json:"marks"`
	Content []Content `json:"content"`
}

type Marks struct {
	Type  string `json:"type"`
	Attrs Attrs  `json:"Attrs"`
}

type Attrs struct {
	Level      int     `json:"level"`
	Url        string  `json:"url"`
	Href       string  `json:"href"`
	Color      string  `json:"color"`
	Layout     string  `json:"layout"`
	Id         string  `json:"id"`
	Type       string  `json:"type"`
	Collection string  `json:"collection"`
	Width      float32 `json:"width"`
	Height     float32 `json:"height"`
}

func (c Content) String(excludeBlockMarks bool) string {
	var innerContent string
	format := "%s"

	if excludeBlockMarks == false {
		if c.Type == "paragraph" {
			format = "<p>%s</p>"
		}

		if c.Type == "code" {
			format = "<code>%s</code>"
		}
	}

	if c.Type == "bulletList" {
		format = "<ul>%s</ul>"
		excludeBlockMarks = true
	}

	if c.Type == "orderedList" {
		format = "<ol>%s</ol>"
		excludeBlockMarks = true
	}

	if c.Type == "listItem" {
		format = "<li>%s</li>"
	}

	if c.Type == "hardBreak" {
		format = "<br/>"
	}

	if c.Type == "text" {
		for _, mark := range c.Marks {
			if mark.Type == "strong" {
				format = "<strong>%s</strong>"
			}

			if mark.Type == "textColor" {
				textColor := mark.Attrs.Color

				format = fmt.Sprintf("<span style=\"color:%s\">%s</span>", textColor, "%s")
			}

			if mark.Type == "em" {
				format = "<em>%s</em>"
			}

			if mark.Type == "underline" {
				format = "<u>%s</u>"
			}

			if mark.Type == "strike" {
				format = "<s>%s</s>"
			}

			if mark.Type == "subsup" {
				markType := mark.Attrs.Type

				if markType == "sup" {
					format = "<sup>%s</sup>"
				}

				if markType == "sub" {
					format = "<sub>%s</sub>"
				}
			}

			if mark.Type == "code" {
				format = "<code>%s</code>"
			}
		}
	}

	if len(c.Content) > 0 {
		for _, c := range c.Content {
			innerContent += c.String(excludeBlockMarks)
		}
	}

	innerContent += c.Text

	return fmt.Sprintf(format, innerContent)
}
