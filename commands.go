package main

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

type saveEntryMsg struct {
	err error
}

func saveEntryCmd(path string, e entry) tea.Cmd {
	return func() tea.Msg {
		return saveEntryMsg{
			err: saveEntry(path, e),
		}
	}
}

type clearMessageMsg struct{}

func clearMessageAfter(seconds int) tea.Cmd {
	return tea.Tick(time.Duration(seconds)*time.Second, func(time.Time) tea.Msg {
		return clearMessageMsg{}
	})
}

type loadEntriesMsg struct {
	entries []entry
	err     error
}

func loadEntriesCmd(dataDir string) tea.Cmd {
	return func() tea.Msg {
		entries, err := getAllEntries(dataDir)
		return loadEntriesMsg{entries: entries, err: err}
	}
}

type trackRepoMsg struct {
	cfg      config
	repoPath string
	err      error
}

func trackRepoCmd(configPath string, cfg config, repoPath string) tea.Cmd {
	return func() tea.Msg {
		nextConfig, normalizedRepoPath, err := configWithRepo(cfg, repoPath)
		if err != nil {
			return trackRepoMsg{err: err}
		}
		if err := saveConfig(configPath, nextConfig); err != nil {
			return trackRepoMsg{err: err}
		}
		return trackRepoMsg{
			cfg:      nextConfig,
			repoPath: normalizedRepoPath,
		}
	}
}

type queryReposCommitsMsg struct {
	commits  string
	warnings []string
}

func queryReposCommitsCmd(repos []repoConfig, commitSinceDate string, fallbackEmail string) tea.Cmd {
	return func() tea.Msg {
		var allCommits []string
		var warnings []string

		for _, repo := range repos {
			commits, err := queryRepoCommits(repo, commitSinceDate, fallbackEmail)
			if err != nil {
				warnings = append(warnings, "Could not load commits for "+repo.Path)
				continue
			}
			if commits != "" {
				allCommits = append(allCommits, commits)
			}
		}
		return queryReposCommitsMsg{
			commits:  strings.Join(allCommits, "\n"),
			warnings: warnings,
		}
	}
}
