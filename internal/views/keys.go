package views

import ui "github.com/jroimartin/gocui"

type keyHelp struct {
	key  string
	body string
}

type key struct {
	views      []string
	keys       []interface{}
	keybinding func(*ui.Gui, *ui.View) error

	help struct {
		key  string
		body string
	}
}

func (s *screen) getKeys() []key {
	return []key{
		{views: []string{"body"}, keys: []interface{}{'n', ui.KeyCtrlN, ui.KeyArrowDown}, keybinding: s.body.next, help: keyHelp{key: "n", body: "(or down arrow) move cursor down"}},
		{views: []string{"body"}, keys: []interface{}{'p', ui.KeyCtrlP, ui.KeyArrowUp}, keybinding: s.body.prev, help: keyHelp{key: "p", body: "(or up arrow) move cursor up"}},
		{views: []string{"body"}, keys: []interface{}{'s'}, keybinding: s.showSearch, help: keyHelp{key: "s", body: "search"}},
		{views: []string{"body"}, keys: []interface{}{'y'}, keybinding: s.showHistory, help: keyHelp{key: "y", body: "history"}},
		{views: []string{"body"}, keys: []interface{}{'v'}, keybinding: s.volumeDown, help: keyHelp{key: "v", body: "volume down"}},
		{views: []string{"body"}, keys: []interface{}{'V'}, keybinding: s.volumeUp, help: keyHelp{key: "V", body: "volume up"}},
		{views: []string{"body"}, keys: []interface{}{'m'}, keybinding: s.goToAlbum, help: keyHelp{key: "m", body: "go to album at cursor"}},
		{views: []string{"body"}, keys: []interface{}{'t'}, keybinding: s.goToArtist, help: keyHelp{key: "t", body: "go to artist at cursor"}},
		{views: []string{"body"}, keys: []interface{}{'a'}, keybinding: s.playAlbum, help: keyHelp{key: "a", body: "play entire album"}},
		{views: []string{"body"}, keys: []interface{}{'q'}, keybinding: s.queue, help: keyHelp{key: "q", body: "view play queue"}},
		{views: []string{"body"}, keys: []interface{}{'P'}, keybinding: s.playlists, help: keyHelp{key: "P", body: "get saved playlists"}},
		{views: []string{"body"}, keys: []interface{}{'d'}, keybinding: s.removeFromQueue, help: keyHelp{key: "d", body: "remove track from queue"}},
		{views: []string{"body"}, keys: []interface{}{'l'}, keybinding: s.body.albumLink, help: keyHelp{key: "l", body: "copy a link to the current album to clipboard"}},
		{views: []string{"body"}, keys: []interface{}{'f'}, keybinding: s.play.next, help: keyHelp{key: "f", body: "fast forward to next song in queue"}},
		{views: []string{"body"}, keys: []interface{}{ui.KeyEnter}, keybinding: s.enter, help: keyHelp{key: "enter", body: "select item at cursor"}},
		{views: []string{"body"}, keys: []interface{}{ui.KeyEsc}, keybinding: s.escape, help: keyHelp{key: "escape", body: "go back to the previous view"}},
		{views: []string{"body"}, keys: []interface{}{ui.KeySpace}, keybinding: s.pause, help: keyHelp{key: "space", body: "pause/unpause"}},
		{views: []string{"search-type"}, keys: []interface{}{'m'}, keybinding: s.search.album},
		{views: []string{"search-type"}, keys: []interface{}{'t'}, keybinding: s.search.artist},
		{views: []string{"search-type"}, keys: []interface{}{'k'}, keybinding: s.search.track},
		{views: []string{"search-type"}, keys: []interface{}{'y'}, keybinding: s.showHistory},
		{views: []string{"search-type"}, keys: []interface{}{ui.KeyEsc}, keybinding: s.escapeSearch},
		{views: []string{"search"}, keys: []interface{}{ui.KeyEnter}, keybinding: s.search.exit},
		{views: []string{"search"}, keys: []interface{}{ui.KeyEsc}, keybinding: s.escapeSearch},
		{views: []string{""}, keys: []interface{}{ui.KeyCtrlD, ui.KeyCtrlC}, keybinding: s.quit, help: keyHelp{key: "C-d (or C-c)", body: "quit"}},
		{views: []string{"login"}, keys: []interface{}{ui.KeyEnter}, keybinding: s.login.next},
		{views: []string{"body"}, keys: []interface{}{'h'}, keybinding: s.showHelp, help: keyHelp{key: "h", body: "toggle help menu"}},
		{views: []string{"help"}, keys: []interface{}{'h'}, keybinding: s.hideHelp},
	}
}
