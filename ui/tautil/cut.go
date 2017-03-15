package tautil

func Cut(ta Texta) {
	if !ta.SelectionOn() {
		return
	}
	a, b := SelectionStringIndexes(ta)
	ta.SetClipboardString(ta.Str()[a:b])
	ta.EditOpen()
	ta.EditDelete(a, b)
	ta.EditClose()
	ta.SetSelectionOn(false)
	ta.SetCursorIndex(a)
}
