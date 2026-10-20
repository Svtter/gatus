package api

import (
	"encoding/xml"
	"fmt"
	"sort"
	"time"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/storage/store"
	"github.com/TwiN/gatus/v5/storage/store/common/paging"
	fiber "github.com/gofiber/fiber/v2"
)

const historyFeedItemLimit = 50

// HistoryFeedHandler serves RSS/Atom feeds of endpoint status events at /history.rss
// and /history.atom, matching the paths used by Atlassian Statuspage
type HistoryFeedHandler struct {
	cfg *config.Config
}

func (handler HistoryFeedHandler) GetRSS(c *fiber.Ctx) error {
	baseURL := c.Protocol() + "://" + c.Hostname()
	document := rssDocument{
		Version: "2.0",
		Channel: rssChannel{
			Title:       handler.cfg.UI.Title,
			Link:        baseURL,
			Description: handler.cfg.UI.Description,
		},
	}
	for _, event := range fetchHistoryFeedEvents() {
		document.Channel.Items = append(document.Channel.Items, rssItem{
			Title:       event.title(),
			Link:        baseURL + "/endpoints/" + event.key,
			PubDate:     event.timestamp.Format(time.RFC1123Z),
			GUID:        rssGUID{IsPermaLink: "false", Value: event.guid()},
			Description: event.description(),
		})
	}
	output, err := xml.Marshal(document)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	c.Set("Content-Type", "application/rss+xml; charset=utf-8")
	return c.Status(fiber.StatusOK).SendString(xml.Header + string(output))
}

func (handler HistoryFeedHandler) GetAtom(c *fiber.Ctx) error {
	baseURL := c.Protocol() + "://" + c.Hostname()
	feed := atomFeed{
		Title: handler.cfg.UI.Title,
		ID:    baseURL,
		Links: []atomLink{{Href: baseURL}},
	}
	events := fetchHistoryFeedEvents()
	if len(events) > 0 {
		feed.Updated = events[0].timestamp.UTC().Format(time.RFC3339)
	} else {
		feed.Updated = time.Now().UTC().Format(time.RFC3339)
	}
	for _, event := range events {
		feed.Entries = append(feed.Entries, atomEntry{
			Title:   event.title(),
			ID:      event.guid(),
			Updated: event.timestamp.UTC().Format(time.RFC3339),
			Link:    atomLink{Href: baseURL + "/endpoints/" + event.key},
			Summary: event.description(),
		})
	}
	output, err := xml.Marshal(feed)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	c.Set("Content-Type", "application/atom+xml; charset=utf-8")
	return c.Status(fiber.StatusOK).SendString(xml.Header + string(output))
}

func fetchHistoryFeedEvents() []historyFeedEvent {
	statuses, err := store.Get().GetAllEndpointStatuses(paging.NewEndpointStatusParams().WithEvents(1, historyFeedItemLimit))
	if err != nil {
		return nil
	}
	var events []historyFeedEvent
	for _, status := range statuses {
		for _, event := range status.Events {
			if event.Type == endpoint.EventStart {
				continue
			}
			events = append(events, historyFeedEvent{
				key:       status.Key,
				name:      status.Name,
				group:     status.Group,
				eventType: event.Type,
				timestamp: event.Timestamp,
			})
		}
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].timestamp.After(events[j].timestamp)
	})
	if len(events) > historyFeedItemLimit {
		events = events[:historyFeedItemLimit]
	}
	return events
}

type historyFeedEvent struct {
	key       string
	name      string
	group     string
	eventType endpoint.EventType
	timestamp time.Time
}

func (e historyFeedEvent) guid() string {
	return fmt.Sprintf("urn:gatus:event:%s:%d:%s", e.key, e.timestamp.Unix(), e.eventType)
}

func (e historyFeedEvent) description() string {
	switch e.eventType {
	case endpoint.EventHealthy:
		return "服务恢复正常 (Healthy)"
	case endpoint.EventUnhealthy:
		return "服务出现异常 (Unhealthy)"
	default:
		return "开始监控 (Monitoring started)"
	}
}

func (e historyFeedEvent) title() string {
	name := e.name
	if e.group != "" {
		name = "[" + e.group + "] " + name
	}
	return name + "：" + e.description()
}

type rssDocument struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	XMLName     xml.Name  `xml:"channel"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	PubDate     string  `xml:"pubDate"`
	GUID        rssGUID `xml:"guid"`
	Description string  `xml:"description"`
}

type rssGUID struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type atomFeed struct {
	XMLName xml.Name    `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string      `xml:"title"`
	ID      string      `xml:"id"`
	Updated string      `xml:"updated"`
	Links   []atomLink  `xml:"link"`
	Entries []atomEntry `xml:"entry"`
}

type atomLink struct {
	Rel  string `xml:"rel,attr,omitempty"`
	Href string `xml:"href,attr"`
}

type atomEntry struct {
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Link    atomLink `xml:"link"`
	Summary string   `xml:"summary"`
}
