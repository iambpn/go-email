package mimeParser

import (
	"strings"

	"github.com/jhillyerd/enmime"
)

type MimePart struct {
	Content          MessageContent
	Attachments      []AttachmentContent
	EmbeddedContents []AttachmentContent
}

type AttachmentContent struct {
	FileName    string
	ContentType string
	ContentId   string
	Data        []byte
}
type MessageContent struct {
	Text string
	Html string
}

func ParseEnmimeParts(part *enmime.Part) *MimePart {
	mimeParts := &MimePart{}

	switch part.ContentType {
	case "multipart/alternative":
		child := part.FirstChild
		for child != nil {
			if content, _, ok := ParseContentType(child); ok && content != nil {
				if content.Text != "" {
					mimeParts.Content.Text = content.Text
				}

				if content.Html != "" {
					mimeParts.Content.Html = content.Html
				}
			}
			child = child.NextSibling
		}
	case "multipart/report", "multipart/related", "multipart/mixed":
		child := part.FirstChild
		for child != nil {
			splitContentType := strings.Split(child.ContentType, "/")

			if splitContentType[0] == "multipart" {
				mParts := ParseEnmimeParts(child)
				mimeParts = MergeMimePart(mimeParts, mParts)
			} else {
				if content, attachment, ok := ParseContentType(child); ok {
					if content != nil && content.Text != "" {
						mimeParts.Content.Text = content.Text
					} else if content != nil && content.Html != "" {
						mimeParts.Content.Html = content.Html
					}

					if attachment != nil && attachment.ContentId != "" {
						// embedded
						mimeParts.EmbeddedContents = append(mimeParts.EmbeddedContents, *attachment)
					} else if attachment != nil {
						// attachment
						mimeParts.Attachments = append(mimeParts.Attachments, *attachment)
					}
				}
			}
			child = child.NextSibling
		}
	}

	return mimeParts
}

func ParseContentType(childPart *enmime.Part) (*MessageContent, *AttachmentContent, bool) {
	content := &MessageContent{}

	switch childPart.ContentType {
	case "text/plain":
		{
			content.Text = string(childPart.Content)
			return content, nil, true
		}
	case "text/html":
		{
			content.Html = string(childPart.Content)
			return content, nil, true
		}
	default:
		if childPart.Disposition != "attachment" {
			break
		}

		attachment := AttachmentContent{
			FileName:    childPart.FileName,
			ContentType: childPart.ContentType,
			Data:        childPart.Content,
		}

		if childPart.ContentID != "" {
			attachment.ContentId = childPart.ContentID
		}

		return nil, &attachment, true
	}

	return nil, nil, false
}

func MergeMimePart(original, part *MimePart) *MimePart {
	if original.Content.Html == "" {
		original.Content.Html = part.Content.Html
	}

	if original.Content.Text == "" {
		original.Content.Text = part.Content.Text
	}

	original.Attachments = append(original.Attachments, part.Attachments...)
	original.EmbeddedContents = append(original.EmbeddedContents, part.EmbeddedContents...)

	return original
}
