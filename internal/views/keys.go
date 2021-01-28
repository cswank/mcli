package views

import ui "github.com/awesome-gocui/gocui"

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
		{views: []string{"body"}, keys: []interface{}{'N'}, keybinding: s.body.nextPage, help: keyHelp{key: "N", body: "next page of results"}},
		{views: []string{"body"}, keys: []interface{}{'p', ui.KeyCtrlP, ui.KeyArrowUp}, keybinding: s.body.prev, help: keyHelp{key: "p", body: "(or up arrow) move cursor up"}},
		{views: []string{"body"}, keys: []interface{}{'P'}, keybinding: s.body.prevPage, help: keyHelp{key: "P", body: "previous page of results"}},
		{views: []string{"body"}, keys: []interface{}{'s'}, keybinding: s.showSearchDialog, help: keyHelp{key: "s", body: "search"}},
		{views: []string{"body"}, keys: []interface{}{'y'}, keybinding: s.showHistoryDialog, help: keyHelp{key: "y", body: "history"}},
		{views: []string{"body"}, keys: []interface{}{'S'}, keybinding: s.showSeekDialog, help: keyHelp{key: "S", body: "seek to a song time (ex. 2:31)"}},
		{views: []string{"body", "volume"}, keys: []interface{}{'v'}, keybinding: s.volumeDown, help: keyHelp{key: "v", body: "volume down"}},
		{views: []string{"body", "volume"}, keys: []interface{}{'V'}, keybinding: s.volumeUp, help: keyHelp{key: "V", body: "volume up"}},
		{views: []string{"body"}, keys: []interface{}{'m'}, keybinding: s.goToAlbum, help: keyHelp{key: "m", body: "go to album at cursor"}},
		{views: []string{"body"}, keys: []interface{}{'t'}, keybinding: s.goToArtist, help: keyHelp{key: "t", body: "go to artist at cursor"}},
		{views: []string{"body"}, keys: []interface{}{'k'}, keybinding: s.goToArtistTracks, help: keyHelp{key: "k", body: "view the tracks of the artist at cursor"}},
		{views: []string{"body"}, keys: []interface{}{'a'}, keybinding: s.playAlbum, help: keyHelp{key: "a", body: "play entire album"}},
		{views: []string{"body"}, keys: []interface{}{'q'}, keybinding: s.queue, help: keyHelp{key: "q", body: "view play queue"}},
		{views: []string{"body"}, keys: []interface{}{'T'}, keybinding: s.playlists, help: keyHelp{key: "T", body: "get saved playlists"}},
		{views: []string{"body"}, keys: []interface{}{'d'}, keybinding: s.removeFromQueue, help: keyHelp{key: "d", body: "remove track from queue"}},
		{views: []string{"body"}, keys: []interface{}{'f'}, keybinding: s.next, help: keyHelp{key: "f", body: "fast forward to next song in queue"}},
		{views: []string{"body"}, keys: []interface{}{'r'}, keybinding: s.rewind, help: keyHelp{key: "r", body: "rewind the current song to the beginning"}},
		{views: []string{"body"}, keys: []interface{}{ui.KeyEnter}, keybinding: s.enter, help: keyHelp{key: "enter", body: "select item at cursor"}},
		{views: []string{"body"}, keys: []interface{}{ui.KeyEsc}, keybinding: s.escape, help: keyHelp{key: "escape", body: "go back to the previous view"}},
		{views: []string{"body"}, keys: []interface{}{ui.KeySpace}, keybinding: s.pause, help: keyHelp{key: "space", body: "pause/unpause"}},
		{views: []string{"search-type"}, keys: []interface{}{'m'}, keybinding: s.search.album},
		{views: []string{"search-type"}, keys: []interface{}{'t'}, keybinding: s.search.artist},
		{views: []string{"search-type"}, keys: []interface{}{'k'}, keybinding: s.search.track},
		{views: []string{"search-type"}, keys: []interface{}{ui.KeyEsc}, keybinding: s.escapeSearch},
		{views: []string{"artist-dialog"}, keys: []interface{}{'m'}, keybinding: s.artistDialog.albums},
		{views: []string{"artist-dialog"}, keys: []interface{}{'k'}, keybinding: s.artistDialog.tracks},
		{views: []string{"history-type"}, keys: []interface{}{'r'}, keybinding: s.history.recent},
		{views: []string{"history-type"}, keys: []interface{}{'p'}, keybinding: s.history.played},
		{views: []string{"search"}, keys: []interface{}{ui.KeyEnter}, keybinding: s.search.exit},
		{views: []string{"search"}, keys: []interface{}{ui.KeyEsc}, keybinding: s.escapeSearch},
		{views: []string{"seek"}, keys: []interface{}{ui.KeyEnter}, keybinding: s.seeker.exit},
		{views: []string{"seek"}, keys: []interface{}{ui.KeyEsc}, keybinding: s.escapeSeek},
		{views: []string{"body"}, keys: []interface{}{'i'}, keybinding: s.importMusic, help: keyHelp{key: "i", body: "import new music into the database"}},
		{views: []string{"body"}, keys: []interface{}{'h'}, keybinding: s.showHelp, help: keyHelp{key: "h", body: "toggle help menu"}},
		{views: []string{"help"}, keys: []interface{}{'h'}, keybinding: s.hideHelp},
		//{views: []string{"body"}, keys: []interface{}{'H'}, keybinding: s.showManual, help: keyHelp{key: "H", body: "toggle manual"}},
		//{views: []string{"manual"}, keys: []interface{}{'H'}, keybinding: s.hideManual},
		{views: []string{""}, keys: []interface{}{'Q', ui.KeyCtrlD, ui.KeyCtrlC}, keybinding: s.quit, help: keyHelp{key: "C-d (or C-c or Q)", body: "quit"}},
	}
}
