package imapwrapper

import (
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	mimeParser "github.com/iambpn/go-email/pkg/mime_parser"
	"github.com/jhillyerd/enmime"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type ImapWrapper struct {
	Host     string
	Port     string
	Username string
	Password string
	client   *client.Client
}

type IwMailbox struct {
	Name         string
	Delimiter    string
	MessageCount uint32 // All Message count
	UnseenCount  uint32
}

type IwPreviewMessage struct {
	Uid          uint32
	Flags        []string
	Envelope     *imap.Envelope
	MessageSize  uint32
	InternalDate time.Time
	PreviewText  string
}

type IwMessage struct {
	IwPreviewMessage
	MimeContent mimeParser.MimePart
	Header      map[string][]string
}

func (iw *ImapWrapper) Connect() error {
	c, err := client.DialTLS(fmt.Sprintf("%s:%s", iw.Host, iw.Port), nil)

	if err != nil {
		return err
	}

	iw.client = c
	err = iw.login()

	if err != nil {
		return err
	}

	return nil
}

func (iw *ImapWrapper) login() error {
	if err := iw.client.Login(iw.Username, iw.Password); err != nil {
		iw.client.Logout()
		return err
	}

	return nil
}

func (iw *ImapWrapper) Logout() {
	iw.client.Logout()
}

func (iw *ImapWrapper) GetMailBoxes() ([]IwMailbox, error) {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	go func() {
		done <- iw.client.List("", "*", mailboxes)
	}()

	iwMailboxes := []IwMailbox{}
	for mbox := range mailboxes {
		status, err := iw.client.Select(mbox.Name, true)

		if err != nil {
			return nil, err
		}

		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{
			imap.SeenFlag,
		}
		ids, err := iw.client.Search(criteria)

		if err != nil {
			return nil, err
		}

		iwMailboxes = append(iwMailboxes, IwMailbox{
			Delimiter:    mbox.Delimiter,
			Name:         mbox.Name,
			MessageCount: status.Messages,
			UnseenCount:  uint32(len(ids)),
		})
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return iwMailboxes, nil
}

type PreviewMessageResponse struct {
	Data         *[]IwPreviewMessage
	TotalMessage uint32
}

/*
- Retrieve IMAP preview messages in descending order
*/
func (iw *ImapWrapper) GetPreviewMessages(mailboxName string, page, limit int) (*PreviewMessageResponse, error) {
	mStatus, err := iw.client.Select(mailboxName, true)

	if err != nil {
		return nil, err
	}

	from := mStatus.Messages - uint32(page*limit)
	to := from + uint32(limit)

	if mStatus.Messages < uint32(page*limit) {
		max := uint32(page * limit)
		extra := max - mStatus.Messages

		to = mStatus.Messages
		if extra > uint32(limit) {
			from = mStatus.Messages
		} else {
			from = (uint32(page) - 1) * uint32(limit)
		}
	}

	if mStatus.Messages == from && mStatus.Messages == to {
		return &PreviewMessageResponse{
			Data:         &[]IwPreviewMessage{},
			TotalMessage: mStatus.Messages,
		}, nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	messagesCh := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- iw.client.Fetch(seqSet, []imap.FetchItem{imap.FetchFlags, imap.FetchInternalDate, imap.FetchRFC822Size, imap.FetchEnvelope, imap.FetchUid}, messagesCh)
	}()

	iwMessages := orderedmap.New[uint32, *IwPreviewMessage]()
	availableSeqSet := new(imap.SeqSet)

	for message := range messagesCh {
		availableSeqSet.AddNum(message.SeqNum)
		iwMessages.Set(message.SeqNum, &IwPreviewMessage{
			Uid:          message.Uid,
			Flags:        message.Flags,
			Envelope:     message.Envelope,
			MessageSize:  message.Size,
			InternalDate: message.InternalDate,
		})
	}

	if err := <-done; err != nil {
		return nil, err
	}

	// fetch message preview
	fetchSection := imap.BodySectionName{
		BodyPartName: imap.BodyPartName{
			Specifier: imap.PartSpecifier("TEXT"),
		},
		Peek: true,
	}

	messagesCh = make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- iw.client.Fetch(availableSeqSet, []imap.FetchItem{imap.FetchBodyStructure, fetchSection.FetchItem()}, messagesCh)
	}()

	for message := range messagesCh {
		mimeMsg, err := iw.fixPartialMimeMessage(message.GetBody(&fetchSection), message.BodyStructure.MIMEType+"/"+message.BodyStructure.MIMESubType, message.BodyStructure.Params["boundary"])

		if err != nil {
			return nil, err
		}

		parser := enmime.NewParser()
		part, err := parser.ReadParts(strings.NewReader(mimeMsg))

		if err != nil {
			return nil, err
		}

		mimeParts := mimeParser.ParseEnmimeParts(part)

		if iwMessage, ok := iwMessages.Get(message.SeqNum); ok {
			minLen := math.Min(float64(len(mimeParts.Content.Text)), float64(100))
			iwMessage.PreviewText = mimeParts.Content.Text[:int(minLen)]
		}
	}

	if err := <-done; err != nil {
		return nil, err
	}

	iwMessagesArr := []IwPreviewMessage{}

	for pair := iwMessages.Newest(); pair != nil; pair = pair.Prev() {
		iwMessagesArr = append(iwMessagesArr, *pair.Value)
	}

	return &PreviewMessageResponse{Data: &iwMessagesArr, TotalMessage: mStatus.Messages}, nil
}

func (iw *ImapWrapper) GetMessage(mailboxName string, uid uint32) (*IwMessage, error) {
	_, err := iw.client.Select(mailboxName, false)

	if err != nil {
		return nil, err
	}

	uidSeqSet := new(imap.SeqSet)
	uidSeqSet.AddNum(uid)

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- iw.client.UidFetch(uidSeqSet, []imap.FetchItem{imap.FetchRFC822, imap.FetchFull}, messages)
	}()

	iwMessage := IwMessage{}
	for message := range messages {
		iwMessage.Flags = message.Flags
		iwMessage.Envelope = message.Envelope
		iwMessage.InternalDate = message.InternalDate
		iwMessage.MessageSize = message.Size
		iwMessage.Uid = message.Uid

		eParser := enmime.NewParser()
		bodySection, err := imap.ParseBodySectionName(imap.FetchRFC822)

		if err != nil {
			return nil, err
		}

		part, err := eParser.ReadParts(message.GetBody(bodySection))
		iwMessage.Header = part.Header

		if err != nil {
			return nil, err
		}

		mimePart := mimeParser.ParseEnmimeParts(part)

		minLen := math.Min(float64(len(mimePart.Content.Text)), 100)
		iwMessage.PreviewText = mimePart.Content.Text[:int(minLen)]
		iwMessage.MimeContent = *mimePart
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return &iwMessage, nil
}

func (iw *ImapWrapper) UpdateMessage(mailbox string, uid uint32, flagsToAdd, flagsToRemove []string) error {
	_, err := iw.client.Select(mailbox, false)

	if err != nil {
		return err
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uid)

	// Fetch the message's current flags
	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- iw.client.UidFetch(seqSet, []imap.FetchItem{imap.FetchFlags}, messages)
	}()

	var msg *imap.Message = <-messages

	if err := <-done; err != nil {
		return err
	}

	// Prepare the new set of flags
	updatedFlags := []interface{}{} // must be of interface{} array
	for _, existingFlag := range msg.Flags {
		flagMatch := false
		for _, flagToRemove := range flagsToRemove {
			if existingFlag == flagToRemove {
				flagMatch = true
				break
			}
		}
		if !flagMatch {
			updatedFlags = append(updatedFlags, existingFlag)
		}
	}

	for _, flag := range flagsToAdd {
		updatedFlags = append(updatedFlags, flag)
	}

	if len(updatedFlags) <= 0 {
		updatedFlags = nil
	}

	if err := iw.client.UidStore(seqSet, imap.FormatFlagsOp(imap.SetFlags, true), updatedFlags, nil); err != nil {
		return err
	}

	return nil
}

func (iw *ImapWrapper) fixPartialMimeMessage(data io.Reader, contentType string, boundary string) (string, error) {
	strMsg, err := iw.readerToString(data)

	if err != nil {
		return "", err
	}

	strMsg = fmt.Sprintf("MIME-Version: 1.0\nContent-Type: %s; boundary=\"%s\"\n\n%s", contentType, boundary, strMsg)

	return strMsg, nil
}

func (iw *ImapWrapper) readerToString(literals io.Reader) (string, error) {
	if b, err := io.ReadAll(literals); err == nil {
		return string(b), nil
	} else {
		return "", err
	}
}
