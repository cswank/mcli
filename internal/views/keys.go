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
		// {views: []string{s.body.name}, keys: []interface{}{'f', ui.KeyArrowRight}, keybinding: s.locked(s.body.forward), help: keyHelp{key: "f", body: "(or right arrow) forward to next page"}},
		// {views: []string{s.body.name}, keys: []interface{}{'b', ui.KeyArrowLeft}, keybinding: s.locked(s.body.back), help: keyHelp{key: "b", body: "(or left arrow) backward to prev page"}},
		// {views: []string{s.body.name}, keys: []interface{}{ui.KeyEnter}, keybinding: s.locked(s.enter), help: keyHelp{key: "enter", body: "view item at cursor"}},
		// {views: []string{s.body.name}, keys: []interface{}{ui.KeyEsc}, keybinding: s.locked(s.escape), help: keyHelp{key: "esc", body: "back to previous view"}},
		// {views: []string{s.body.name}, keys: []interface{}{ui.KeyCtrlJ}, keybinding: s.locked(s.jump), help: keyHelp{key: "C-j", body: "jump to a kafka offset"}},
		// {views: []string{s.body.name}, keys: []interface{}{ui.KeyCtrlO}, keybinding: s.locked(s.offset), help: keyHelp{key: "C-o", body: "set the offset in all partitions of topic"}},
		// {views: []string{s.body.name}, keys: []interface{}{ui.KeyCtrlS, '/'}, keybinding: s.locked(s.search), help: keyHelp{key: "C-s", body: "(or /) search kafka messages"}},
		{views: []string{"search-type"}, keys: []interface{}{'m'}, keybinding: s.search.album},
		{views: []string{"search-type"}, keys: []interface{}{'t'}, keybinding: s.search.artist},
		{views: []string{"search-type"}, keys: []interface{}{'k'}, keybinding: s.search.track},
		{views: []string{"search"}, keys: []interface{}{ui.KeyEnter}, keybinding: s.search.exit},
		{views: []string{""}, keys: []interface{}{ui.KeyCtrlD, ui.KeyCtrlC}, keybinding: s.quit, help: keyHelp{key: "C-d (or C-c)", body: "quit"}},
		{views: []string{"login"}, keys: []interface{}{ui.KeyEnter}, keybinding: s.login.next},
		// {views: []string{s.footer.name}, keys: []interface{}{ui.KeyEsc}, keybinding: s.footer.bail},
		// {views: []string{s.body.name}, keys: []interface{}{'h'}, keybinding: s.showHelp, help: keyHelp{key: "h", body: "toggle help"}},
		// {views: []string{s.help.name}, keys: []interface{}{'h'}, keybinding: s.hideHelp},
	}
}
