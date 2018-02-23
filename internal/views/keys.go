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
		{views: []string{"body"}, keys: []interface{}{'n', ui.KeyArrowDown}, keybinding: s.body.next, help: keyHelp{key: "n", body: "(or down arrow) move cursor down"}},
		{views: []string{"body"}, keys: []interface{}{'p', ui.KeyArrowUp}, keybinding: s.body.prev, help: keyHelp{key: "p", body: "(or up arrow) move cursor up"}},
		{views: []string{"body"}, keys: []interface{}{ui.KeyEnter}, keybinding: s.body.enter, help: keyHelp{key: "enter", body: "select item at cursor"}},
		{views: []string{"body"}, keys: []interface{}{ui.KeyEsc}, keybinding: s.escape, help: keyHelp{key: "escape", body: "go back to the previous view"}},
		{views: []string{"search-type"}, keys: []interface{}{'m'}, keybinding: s.search.album},
		{views: []string{"search-type"}, keys: []interface{}{'t'}, keybinding: s.search.artist},
		{views: []string{"search-type"}, keys: []interface{}{'k'}, keybinding: s.search.track},
		{views: []string{"search"}, keys: []interface{}{ui.KeyEnter}, keybinding: s.search.exit},
		{views: []string{"search"}, keys: []interface{}{ui.KeyEsc}, keybinding: s.search.escape},
		{views: []string{""}, keys: []interface{}{ui.KeyCtrlD, ui.KeyCtrlC}, keybinding: s.quit, help: keyHelp{key: "C-d (or C-c)", body: "quit"}},
		{views: []string{"login"}, keys: []interface{}{ui.KeyEnter}, keybinding: s.login.next},
	}
}
