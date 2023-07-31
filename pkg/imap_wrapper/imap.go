package imapwrapper

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type ImapWrapper struct {
	host     string
	port     string
	username string
	password string
	Client   *client.Client
}

type IwMailbox struct {
	Name         string
	Delimiter    string
	MessageCount uint32 // All Message count
	UnseenCount  uint32
}

type IwFetchAllMessage struct {
	Flags         []string
	Envelope      *imap.Envelope
	MessageSize   uint32
	InternalDate  time.Time
	BodyStructure *imap.BodyStructure
	// Body          *imap.
}

func (iw *ImapWrapper) Connect() error {
	c, err := client.DialTLS(fmt.Sprintf("%s:%s", iw.host, iw.port), nil)

	if err != nil {
		return err
	}

	err = iw.login()

	if err != nil {
		return err
	}

	iw.Client = c
	return nil
}

func (iw *ImapWrapper) login() error {
	if err := iw.Client.Login(iw.username, iw.password); err != nil {
		iw.Client.Logout()
		return err
	}

	return nil
}

func (iw *ImapWrapper) Logout() {
	iw.Client.Logout()
}

func (iw *ImapWrapper) GetMailBoxes() ([]IwMailbox, error) {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	go func() {
		done <- iw.Client.List("", "*", mailboxes)
	}()

	iwMailboxes := []IwMailbox{}
	for mbox := range mailboxes {
		status, err := iw.Client.Select(mbox.Name, true)

		if err != nil {
			return nil, err
		}

		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{
			imap.SeenFlag,
		}
		ids, err := iw.Client.Search(criteria)

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

/*
- Retrieve IMAP messages in descending order
*/
func (iw *ImapWrapper) Messages(mailboxName string, page, limit int) ([]IwFetchAllMessage, error) {
	mStatus, err := iw.Client.Select(mailboxName, true)

	if err != nil {
		return nil, err
	}

	to := uint32(1)
	from := uint32(1)

	if mStatus.Messages > uint32(page*limit) {
		from = mStatus.Messages - uint32(page*limit)
		to = from + uint32(limit)
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	messagesCh := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- iw.Client.Fetch(seqSet, []imap.FetchItem{imap.FetchFlags, imap.FetchInternalDate, imap.FetchRFC822Size, imap.FetchEnvelope, imap.FetchBodyStructure, "BODY.PEEK[TEXT]"}, messagesCh)
	}()

	iwMessages := []IwFetchAllMessage{}

	for message := range messagesCh {
		// message.Bo
		iwMessages = append(iwMessages, IwFetchAllMessage{
			Flags:        message.Flags,
			Envelope:     message.Envelope,
			MessageSize:  message.Size,
			InternalDate: message.InternalDate,
		})
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return iwMessages, nil
}

//TODO: add preview message to the iwMessage
// TODO: fix properties to include in iwFetchALlMessage Struct
