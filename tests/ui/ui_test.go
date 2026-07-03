package ui_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/marcuwynu23/git-share/internal/ui"
)

func TestDashboard(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ui.Dashboard(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html; charset=utf-8", ct)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "git share") {
		t.Error("response should contain 'git share'")
	}
	if !strings.Contains(body, "Repository") {
		t.Error("response should contain 'Repository'")
	}
	if !strings.Contains(body, "Branch") {
		t.Error("response should contain 'Branch'")
	}
	if !strings.Contains(body, "Clone") {
		t.Error("response should contain 'Clone'")
	}
}

func TestDashboardHTML(t *testing.T) {
	if !strings.Contains(ui.DashboardHTML, "<!DOCTYPE html>") {
		t.Error("DashboardHTML should be valid HTML")
	}
	if !strings.Contains(ui.DashboardHTML, "loadInfo") {
		t.Error("DashboardHTML should contain the info loader script")
	}
	if !strings.Contains(ui.DashboardHTML, "git clone") {
		t.Error("DashboardHTML should contain clone instructions")
	}
}
