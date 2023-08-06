package mimeParser

import (
	"testing"

	"github.com/jhillyerd/enmime"
)

var part = enmime.NewPart("multipart/alternative")

func TestParseEnmimeParts_multipartAlternative(t *testing.T) {
	var text = "text email"
	var html = "<h1>html email</h1>"
	textPart := enmime.NewPart("text/plain")
	textPart.Content = []byte(text)
	htmlPart := enmime.NewPart("text/html")
	htmlPart.Content = []byte(html)
	part.AddChild(textPart)
	part.AddChild(htmlPart)
	mimePart := ParseEnmimeParts(part)

	if mimePart.Content.Text != text {
		t.Error("Text did not matched", mimePart.Content.Text)
	}

	if mimePart.Content.Html != html {
		t.Error("Html did not matched", mimePart.Content.Html)
	}

}
