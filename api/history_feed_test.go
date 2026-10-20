package api

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/ui"
	"github.com/TwiN/gatus/v5/storage/store"
	fiber "github.com/gofiber/fiber/v2"
)

func TestHistoryFeedHandler_GetRSS(t *testing.T) {
	healthy := &endpoint.Endpoint{Group: "feed-rss", Name: "rss-endpoint"}
	if err := store.Get().InsertEndpointResult(healthy, &endpoint.Result{Success: true, Duration: time.Second, Timestamp: time.Now().Add(-10 * time.Minute)}); err != nil {
		t.Fatal("expected no error, got", err.Error())
	}
	if err := store.Get().InsertEndpointResult(healthy, &endpoint.Result{Success: false, Errors: []string{"connection refused"}, Duration: time.Second, Timestamp: time.Now().Add(-5 * time.Minute)}); err != nil {
		t.Fatal("expected no error, got", err.Error())
	}
	unhealthy := &endpoint.Endpoint{Group: "feed-rss", Name: "healthy-endpoint"}
	if err := store.Get().InsertEndpointResult(unhealthy, &endpoint.Result{Success: true, Duration: time.Second, Timestamp: time.Now().Add(-9 * time.Minute)}); err != nil {
		t.Fatal("expected no error, got", err.Error())
	}
	app := fiber.New()
	app.Get("/history.rss", HistoryFeedHandler{cfg: &config.Config{UI: ui.GetDefaultConfig()}}.GetRSS)
	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/history.rss", http.NoBody))
	if err != nil {
		t.Fatal("expected no error, got", err.Error())
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
	if contentType := response.Header.Get("Content-Type"); !strings.Contains(contentType, "application/rss+xml") {
		t.Fatalf("expected Content-Type to contain application/rss+xml, got %s", contentType)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("expected no error, got", err.Error())
	}
	output := string(body)
	if !strings.HasPrefix(output, xml.Header) {
		t.Fatal("expected output to start with xml declaration")
	}
	if !strings.Contains(output, `<rss version="2.0">`) {
		t.Fatal("expected output to contain rss 2.0 root element")
	}
	unhealthyTitle := "[feed-rss] rss-endpoint：服务出现异常 (Unhealthy)"
	healthyTitle := "[feed-rss] healthy-endpoint：服务恢复正常 (Healthy)"
	if !strings.Contains(output, unhealthyTitle) {
		t.Fatal("expected output to contain unhealthy event title, got", output)
	}
	if !strings.Contains(output, healthyTitle) {
		t.Fatal("expected output to contain healthy event title, got", output)
	}
	if indexUnhealthy, indexHealthy := strings.Index(output, unhealthyTitle), strings.Index(output, healthyTitle); indexUnhealthy > indexHealthy {
		t.Fatalf("expected newest event first, got unhealthy at index %d and healthy at index %d", indexUnhealthy, indexHealthy)
	}
	if !strings.Contains(output, `isPermaLink="false"`) {
		t.Fatal("expected guids to be non-permalinks")
	}
}

func TestHistoryFeedHandler_GetAtom(t *testing.T) {
	ep := &endpoint.Endpoint{Group: "feed-atom", Name: "atom-endpoint"}
	if err := store.Get().InsertEndpointResult(ep, &endpoint.Result{Success: true, Duration: time.Second, Timestamp: time.Now().Add(-10 * time.Minute)}); err != nil {
		t.Fatal("expected no error, got", err.Error())
	}
	if err := store.Get().InsertEndpointResult(ep, &endpoint.Result{Success: false, Errors: []string{"timeout"}, Duration: time.Second, Timestamp: time.Now().Add(-4 * time.Minute)}); err != nil {
		t.Fatal("expected no error, got", err.Error())
	}
	app := fiber.New()
	app.Get("/history.atom", HistoryFeedHandler{cfg: &config.Config{UI: ui.GetDefaultConfig()}}.GetAtom)
	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/history.atom", http.NoBody))
	if err != nil {
		t.Fatal("expected no error, got", err.Error())
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
	if contentType := response.Header.Get("Content-Type"); !strings.Contains(contentType, "application/atom+xml") {
		t.Fatalf("expected Content-Type to contain application/atom+xml, got %s", contentType)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal("expected no error, got", err.Error())
	}
	output := string(body)
	if !strings.Contains(output, `<feed xmlns="http://www.w3.org/2005/Atom">`) {
		t.Fatal("expected output to contain atom feed root element, got", output)
	}
	if !strings.Contains(output, "[feed-atom] atom-endpoint：服务出现异常 (Unhealthy)") {
		t.Fatal("expected output to contain unhealthy event title, got", output)
	}
	if !strings.Contains(output, "urn:gatus:event:feed-atom_atom-endpoint:") {
		t.Fatal("expected output to contain urn-based entry id, got", output)
	}
}
