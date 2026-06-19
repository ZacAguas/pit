package main

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func press(key string) tea.KeyPressMsg {
	if len(key) == 1 {
		r := rune(key[0])
		return tea.KeyPressMsg(tea.Key{
			Text: key,
			Code: r,
		})
	}
	switch key {
	case "esc":
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyEsc})
	case "enter":
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter})
	default:
		return tea.KeyPressMsg(tea.Key{Text: key})
	}
}

func update(t *testing.T, m model, msg tea.Msg) model {
	t.Helper() // mark function as a test helper - skipped when printing line/file info

	next, _ := m.Update(msg)
	got, ok := next.(model)
	if !ok {
		t.Fatalf("expected model, got %T", next)
	}

	return got
}

func testModel(t *testing.T) model {
	t.Helper()

	dir := t.TempDir()
	return initialModel(dir, config{}, configFilePath(dir), "", "", nil, nil)
}

func TestInitialModelWithExistingEntryDoesNotLoadCommits(t *testing.T) {
	dir := t.TempDir()
	existing := entry{
		Date:     "2026-06-16",
		Did:      "did work",
		Blocked:  "blocked thing",
		Tomorrow: "next thing",
	}

	m := initialModel(dir, config{Repos: []repoConfig{{Path: dir}}}, configFilePath(dir), "", "2026-06-15", &existing, nil)

	if m.loadingCommits {
		t.Fatal("expected loadingCommits false")
	}
	if got := m.did.Value(); got != existing.Did {
		t.Fatalf("expected did %q, got %q", existing.Did, got)
	}
	if cmd := m.Init(); cmd != nil {
		t.Fatal("expected nil Init command")
	}
}

func TestInitialModelWithReposLoadsCommits(t *testing.T) {
	dir := t.TempDir()

	m := initialModel(dir, config{Repos: []repoConfig{{Path: dir}}}, configFilePath(dir), "", "2026-06-15", nil, nil)

	if !m.loadingCommits {
		t.Fatal("expected loadingCommits true")
	}
	if m.commitSinceDate != "2026-06-15" {
		t.Fatalf("expected commitSinceDate %q, got %q", "2026-06-15", m.commitSinceDate)
	}
	if cmd := m.Init(); cmd == nil {
		t.Fatal("expected Init command")
	}
}

func TestInitialModelWithoutReposDoesNotLoadCommits(t *testing.T) {
	m := testModel(t)

	if m.loadingCommits {
		t.Fatal("expected loadingCommits false")
	}
	if m.commitSinceDate != "" {
		t.Fatalf("expected empty commitSinceDate, got %q", m.commitSinceDate)
	}
	if cmd := m.Init(); cmd != nil {
		t.Fatal("expected nil Init command")
	}
}

func TestInitialModelSeedsDidFromPreviousEntry(t *testing.T) {
	dir := t.TempDir()
	previous := entry{
		Date:     "2026-06-15",
		Tomorrow: "finish git prepopulation",
	}

	m := initialModel(dir, config{}, configFilePath(dir), "", "2026-06-15", nil, &previous)

	if got := m.did.Value(); got != previous.Tomorrow {
		t.Fatalf("expected did %q, got %q", previous.Tomorrow, got)
	}
	if m.loadingCommits {
		t.Fatal("expected loadingCommits false")
	}
}

func TestHOpensHistory(t *testing.T) {
	m := testModel(t)

	// send an H keypress, should switch to history view
	m = update(t, m, press("h"))
	if m.view != historyView {
		t.Fatalf("expected %v, got %v", historyView, m.view)
	}
}

func TestQInHistoryReturnsToToday(t *testing.T) {
	m := testModel(t)
	m.view = historyView // move to history view

	// send a Q keypress, should switch to today view
	m = update(t, m, press("q"))
	if m.view != todayView {
		t.Fatalf("expected %v, got %v", todayView, m.view)
	}
}

func TestEscInHistoryReturnsToToday(t *testing.T) {
	m := testModel(t)
	m.view = historyView

	m = update(t, m, press("esc"))
	if m.view != todayView {
		t.Fatalf("expected %v, got %v", todayView, m.view)
	}
}

func TestQInHistoryFilterDoesNotReturnToToday(t *testing.T) {
	m := testModel(t)
	m.view = historyView
	m.history.SetItems(entriesToListItems([]entry{
		{Date: "2026-06-17", Did: "did work"},
	}))

	m = update(t, m, press("/"))
	m = update(t, m, press("q"))

	if m.view != historyView {
		t.Fatalf("expected %v, got %v", historyView, m.view)
	}
	if got := m.history.FilterValue(); got != "q" {
		t.Fatalf("expected filter value %q, got %q", "q", got)
	}
}

func TestEscInHistoryFilterClearsFilter(t *testing.T) {
	m := testModel(t)
	m.view = historyView
	m.history.SetItems(entriesToListItems([]entry{
		{Date: "2026-06-17", Did: "did work"},
	}))

	m = update(t, m, press("/"))
	m = update(t, m, press("d"))
	m = update(t, m, press("esc"))

	if m.view != historyView {
		t.Fatalf("expected %v, got %v", historyView, m.view)
	}
	if got := m.history.FilterValue(); got != "" {
		t.Fatalf("expected empty filter value, got %q", got)
	}
}

func TestQQuitsInNormalMode(t *testing.T) {
	m := testModel(t)

	_, cmd := m.Update(press("q"))
	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}
}

func TestIEntersEditMode(t *testing.T) {
	m := testModel(t)

	m = update(t, m, press("i"))
	if m.mode != editMode {
		t.Fatalf("expected %v, got %v", editMode, m.mode)
	}
}

func TestEscEntersNormalMode(t *testing.T) {
	m := testModel(t)

	// enter edit mode
	m = update(t, m, press("i"))

	m = update(t, m, press("esc"))
	if m.mode != normalMode {
		t.Fatalf("expected %v, got %v", normalMode, m.mode)
	}
}

func TestQDoesNotQuitInEditMode(t *testing.T) {
	m := testModel(t)

	m.mode = editMode
	m.focus = didField
	m.did.Focus()

	next, cmd := m.Update(press("q"))
	gotModel, ok := next.(model)
	if !ok {
		t.Fatalf("expected model, got %T", next)
	}

	// no command - no quit was issued
	if cmd != nil {
		if _, isQuit := cmd().(tea.QuitMsg); isQuit {
			t.Fatal("expected 'q' not to quit but a tea.Quit command was returned")
		}
	}

	if got := gotModel.did.Value(); got != "q" {
		t.Fatalf("expected q to be inserted into 'did' field, got %q", got)
	}
}

func TestNextFieldNavigationCycles(t *testing.T) {
	m := testModel(t)

	m = update(t, m, press("j"))
	if m.focus != blockedField {
		t.Fatalf("expected blockedField focused, got %v", m.focus)
	}

	m = update(t, m, press("j"))
	if m.focus != tomorrowField {
		t.Fatalf("expected tomorrowField focused, got %v", m.focus)
	}

	m = update(t, m, press("j"))
	if m.focus != didField {
		t.Fatalf("expected didField focused, got %v", m.focus)
	}
}

func TestPreviousFieldNavigationCycles(t *testing.T) {
	m := testModel(t)

	m = update(t, m, press("k"))
	if m.focus != tomorrowField {
		t.Fatalf("expected tomorrowField focused, got %v", m.focus)
	}

	m = update(t, m, press("k"))
	if m.focus != blockedField {
		t.Fatalf("expected blockedField focused, got %v", m.focus)
	}

	m = update(t, m, press("k"))
	if m.focus != didField {
		t.Fatalf("expected didField focused, got %v", m.focus)
	}
}

func TestNumberKeysJumpToFields(t *testing.T) {
	m := testModel(t)

	m = update(t, m, press("2"))
	if m.focus != blockedField {
		t.Fatalf("expected blockedField focused, got %v", m.focus)
	}

	m = update(t, m, press("3"))
	if m.focus != tomorrowField {
		t.Fatalf("expected tomorrowField focused, got %v", m.focus)
	}

	m = update(t, m, press("1"))
	if m.focus != didField {
		t.Fatalf("expected didField focused, got %v", m.focus)
	}
}

func TestBulletizeCurrentField(t *testing.T) {
	m := testModel(t)
	m.focus = didField
	m.did.SetValue("first\n\n- already\nsecond")

	m = update(t, m, press("b"))

	want := "- first\n\n- already\n- second"
	if got := m.did.Value(); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestBInEditModeTypesIntoField(t *testing.T) {
	m := testModel(t)
	m.mode = editMode
	m.focus = didField
	m.did.Focus()

	m = update(t, m, press("b"))

	if got := m.did.Value(); got != "b" {
		t.Fatalf("expected %q, got %q", "b", got)
	}
}

func TestCurrentEntryUsesTextFields(t *testing.T) {
	m := testModel(t)

	const did = "did work"
	const blocked = "blocked thing"
	const tomorrow = "next thing"
	m.did.SetValue(did)
	m.blocked.SetValue(blocked)
	m.tomorrow.SetValue(tomorrow)

	got := m.currentEntry()
	if got.Did != did {
		t.Fatalf("expected %q got %q", did, got.Did)
	}
	if got.Blocked != blocked {
		t.Fatalf("expected %q got %q", blocked, got.Blocked)
	}
	if got.Tomorrow != tomorrow {
		t.Fatalf("expected %q got %q", tomorrow, got.Tomorrow)
	}
}

func TestCCopiesCurrentEntryShowsMessage(t *testing.T) {
	m := testModel(t)

	next, cmd := m.Update(press("c"))
	got, ok := next.(model)
	if !ok {
		t.Fatalf("expected model, got %T", next)
	}

	if got.message != "Copied to clipboard" {
		t.Fatalf("expected copy message, got %q", got.message)
	}

	if cmd == nil {
		t.Fatal("expected command, got nil")
	}
}

func TestClearMessageClearsMessage(t *testing.T) {
	m := testModel(t)
	m.message = "Copied to clipboard"

	next, _ := m.Update(clearMessageMsg{})
	got, ok := next.(model)
	if !ok {
		t.Fatalf("expected model, got %T", next)
	}

	if got.message != "" {
		t.Fatalf("expected message to be cleared, got %q", got.message)
	}
}

func TestQueryReposCommitsMessageSetsDidAndClearsLoading(t *testing.T) {
	m := testModel(t)
	m.loadingCommits = true

	next, _ := m.Update(queryReposCommitsMsg{commits: "- Add config"})
	got := next.(model)

	if got.loadingCommits {
		t.Fatal("expected loadingCommits false")
	}
	if got.did.Value() != "- Add config" {
		t.Fatalf("expected commits in did field, got %q", got.did.Value())
	}
	if got.message != "" {
		t.Fatalf("expected empty message, got %q", got.message)
	}
}

func TestQueryReposCommitsMessageAppendsToExistingDid(t *testing.T) {
	m := testModel(t)
	m.loadingCommits = true
	m.did.SetValue("finish git prepopulation")

	next, _ := m.Update(queryReposCommitsMsg{commits: "- Add commits"})
	got := next.(model)

	want := "finish git prepopulation\n\n- Add commits"
	if got.did.Value() != want {
		t.Fatalf("expected did %q, got %q", want, got.did.Value())
	}
}

func TestQueryReposCommitsMessageStoresWarnings(t *testing.T) {
	m := testModel(t)
	m.loadingCommits = true

	next, _ := m.Update(queryReposCommitsMsg{
		warnings: []string{"Could not load commits for /tmp/project"},
	})
	got := next.(model)

	if got.loadingCommits {
		t.Fatal("expected loadingCommits false")
	}
	if len(got.commitWarnings) != 1 {
		t.Fatalf("expected one warning, got %#v", got.commitWarnings)
	}
	want := "Could not load commits for 1 repo"
	if got.message != want {
		t.Fatalf("expected %q, got %q", want, got.message)
	}
}

func TestCommitWarningMessagePluralizesRepos(t *testing.T) {
	got := commitWarningMessage(2)
	want := "Could not load commits for 2 repos"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSaveEntrySuccessShowsSavedMessage(t *testing.T) {
	m := testModel(t)

	next, _ := m.Update(saveEntryMsg{err: nil})
	got := next.(model)

	if got.message != "Saved" {
		t.Fatalf("expected Saved message, got %q", got.message)
	}
}

func TestSaveEntryErrorShowsFailureMessage(t *testing.T) {
	m := testModel(t)

	next, _ := m.Update(saveEntryMsg{err: errors.New("permission denied")})
	got := next.(model)

	want := "Save failed: permission denied"
	if got.message != want {
		t.Fatalf("expected %q, got %q", want, got.message)
	}
}

func TestAWithUntrackedRepoReturnsTrackRepoCommand(t *testing.T) {
	m := testModel(t)
	m.untrackedRepoPath = filepath.Join(t.TempDir(), "project")

	_, cmd := m.Update(press("a"))
	if cmd == nil {
		t.Fatal("expected command, got nil")
	}
}

func TestAWithoutUntrackedRepoDoesNothing(t *testing.T) {
	m := testModel(t)

	_, cmd := m.Update(press("a"))
	if cmd != nil {
		t.Fatal("expected nil command")
	}
}

func TestTrackRepoSuccessUpdatesConfigAndClearsUntrackedRepo(t *testing.T) {
	m := testModel(t)
	repoPath := filepath.Join(t.TempDir(), "project")
	normalizedRepoPath, err := normalizeRepoPath(repoPath)
	if err != nil {
		t.Fatal(err)
	}
	m.untrackedRepoPath = repoPath

	next, _ := m.Update(trackRepoMsg{
		cfg:      config{Repos: []repoConfig{{Path: normalizedRepoPath}}},
		repoPath: normalizedRepoPath,
	})
	got := next.(model)

	if got.untrackedRepoPath != "" {
		t.Fatalf("expected untracked repo path to be cleared, got %q", got.untrackedRepoPath)
	}
	if !configHasRepo(got.config, normalizedRepoPath) {
		t.Fatalf("expected config to contain %q, got %#v", normalizedRepoPath, got.config.Repos)
	}
	if got.message != "Tracking repo: "+normalizedRepoPath {
		t.Fatalf("expected tracking message, got %q", got.message)
	}
}

func TestTrackRepoErrorKeepsUntrackedRepo(t *testing.T) {
	m := testModel(t)
	m.untrackedRepoPath = "/tmp/project"

	next, _ := m.Update(trackRepoMsg{err: errors.New("permission denied")})
	got := next.(model)

	if got.untrackedRepoPath != m.untrackedRepoPath {
		t.Fatalf("expected untracked repo path %q, got %q", m.untrackedRepoPath, got.untrackedRepoPath)
	}
	want := "Could not track repo: permission denied"
	if got.message != want {
		t.Fatalf("expected %q, got %q", want, got.message)
	}
}

func TestVInTodayPreviewsCurrentEntry(t *testing.T) {
	m := testModel(t)
	m = update(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	m.did.SetValue("did work")
	m.blocked.SetValue("blocked thing")
	m.tomorrow.SetValue("next thing")

	m = update(t, m, press("v"))

	if m.view != detailView {
		t.Fatalf("expected %v, got %v", detailView, m.view)
	}
	if !m.previewingCurrentEntry {
		t.Fatal("expected previewingCurrentEntry true")
	}
	got := m.detail.View()
	for _, want := range []string{"Standup", "did", "work", "blocked", "thing", "next"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected preview to contain %q, got %q", want, got)
		}
	}
}

func TestQInTodayPreviewReturnsToToday(t *testing.T) {
	m := testModel(t)
	m.view = detailView
	m.previewingCurrentEntry = true

	m = update(t, m, press("q"))

	if m.view != todayView {
		t.Fatalf("expected %v, got %v", todayView, m.view)
	}
	if m.previewingCurrentEntry {
		t.Fatal("expected previewingCurrentEntry false")
	}
}

func TestCInTodayPreviewCopiesCurrentEntryShowsMessage(t *testing.T) {
	m := testModel(t)
	m.view = detailView
	m.previewingCurrentEntry = true
	m.did.SetValue("did work")

	next, cmd := m.Update(press("c"))
	got, ok := next.(model)
	if !ok {
		t.Fatalf("expected model, got %T", next)
	}

	if got.message != "Copied to clipboard" {
		t.Fatalf("expected copy message, got %q", got.message)
	}
	if cmd == nil {
		t.Fatal("expected command, got nil")
	}
}

func TestEnterInHistoryWithSelectedItemOpensDetail(t *testing.T) {
	m := testModel(t)
	m.view = historyView
	m = update(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	m.history.SetItems(entriesToListItems([]entry{
		{Date: "2026-06-17", Did: "did work"},
	}))

	m = update(t, m, press("enter"))
	if m.view != detailView {
		t.Fatalf("expected %v, got %v", detailView, m.view)
	}
	if got := m.detail.View(); !strings.Contains(got, "Standup") || !strings.Contains(got, "did") || !strings.Contains(got, "work") {
		t.Fatalf("expected detail viewport to contain rendered entry, got %q", got)
	}
}

func TestEnterInEmptyHistoryDoesNotOpenDetail(t *testing.T) {
	m := testModel(t)
	m.view = historyView

	m = update(t, m, press("enter"))
	if m.view != historyView {
		t.Fatalf("expected %v, got %v", historyView, m.view)
	}
}

func TestQInDetailReturnsToHistory(t *testing.T) {
	m := testModel(t)
	m.view = detailView

	m = update(t, m, press("q"))
	if m.view != historyView {
		t.Fatalf("expected %v, got %v", historyView, m.view)
	}
}

func TestEscInDetailReturnsToHistory(t *testing.T) {
	m := testModel(t)
	m.view = detailView

	m = update(t, m, press("esc"))
	if m.view != historyView {
		t.Fatalf("expected %v, got %v", historyView, m.view)
	}
}

func TestCInDetailCopiesSelectedEntryShowsMessage(t *testing.T) {
	m := testModel(t)
	m.view = detailView
	m.history.SetItems(entriesToListItems([]entry{
		{Date: "2026-06-17", Did: "did work"},
	}))

	next, cmd := m.Update(press("c"))
	got, ok := next.(model)
	if !ok {
		t.Fatalf("expected model, got %T", next)
	}

	if got.message != "Copied to clipboard" {
		t.Fatalf("expected copy message, got %q", got.message)
	}
	if cmd == nil {
		t.Fatal("expected command, got nil")
	}
}

func TestResizeInDetailKeepsRenderedContent(t *testing.T) {
	m := testModel(t)
	m.view = historyView
	m = update(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	m.history.SetItems(entriesToListItems([]entry{
		{Date: "2026-06-17", Did: "did work"},
	}))
	m = update(t, m, press("enter"))

	m = update(t, m, tea.WindowSizeMsg{Width: 40, Height: 12})
	if m.view != detailView {
		t.Fatalf("expected %v, got %v", detailView, m.view)
	}
	if got := m.detail.View(); !strings.Contains(got, "Standup") || !strings.Contains(got, "did") || !strings.Contains(got, "work") {
		t.Fatalf("expected detail viewport to keep rendered entry after resize, got %q", got)
	}
}
