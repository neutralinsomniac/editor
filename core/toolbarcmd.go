package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmigpin/editor/core/cmdutil"
	"github.com/jmigpin/editor/core/toolbardata"
	"github.com/jmigpin/editor/ui"
	"github.com/jmigpin/editor/ui/tautil"
)

func ToolbarCmdFromLayout(ed *Editor, layout *ui.Layout) {
	ta := layout.Toolbar.TextArea
	tbsd := toolbardata.NewStringData(ta.Str())
	part, ok := tbsd.GetPartAtIndex(ta.CursorIndex())
	if !ok {
		return
	}
	p0 := part.Args[0].Unquote()
	switch p0 {
	case "Exit":
		ed.Close()
	case "SaveSession":
		cmdutil.SaveSession(ed, part)
	case "OpenSession":
		cmdutil.OpenSession(ed, part)
	case "DeleteSession":
		cmdutil.DeleteSession(ed, part)
	case "ListSessions":
		cmdutil.ListSessions(ed)
	case "NewColumn":
		_ = ed.ui.Layout.Cols.NewColumn()
	case "SaveAllFiles":
		cmdutil.SaveRowsFiles(ed)
	case "ReloadAll":
		cmdutil.ReloadRows(ed)
	case "ReloadAllFiles":
		cmdutil.ReloadRowsFiles(ed)
	case "NewRow":
		cmdutil.NewRow(ed)
	case "ReopenRow":
		ed.reopenRow.Reopen()
	case "XdgOpenDir":
		cmdutil.XdgOpenDirShortcut(ed)
	case "RowDirectory":
		cmdutil.OpenRowDirectory(ed)
	case "FontRunes":
		var u string
		for i := 0; i < 10000; {
			start := i
			var w string
			for j := 0; j < 25; j++ {
				w += string(rune(i))
				i++
			}
			u += fmt.Sprintf("%d: %s\n", start, w)
		}
		ed.Errorf("%s", u)
	default:
		// try running row command
		erow, ok := ed.ActiveERow()
		if ok {
			ok := rowToolbarCmd(erow.(*ERow), part)
			if ok {
				return
			}
		}
		// TODO: consider running external command in new row
		err := fmt.Errorf("unknown layout command (no row is selected or it's also not a row command): %v", part.Str)
		ed.Error(err)
	}
}

func ToolbarCmdFromRow(erow *ERow) {
	err := toolbarCmdFromRow2(erow)
	if err != nil {
		erow.Ed().Error(err)
	}
}
func toolbarCmdFromRow2(erow *ERow) error {
	tbsd := erow.ToolbarSD()
	ta := erow.Row().Toolbar
	part, ok := tbsd.GetPartAtIndex(ta.CursorIndex())
	if !ok {
		return errors.New("missing part at index")
	}

	// don't allow commands on row first part
	if part == tbsd.Parts[0] {
		return errors.New("running a command on first part")
	}

	if len(part.Args) == 0 {
		return errors.New("empty part")
	}

	ok = rowToolbarCmd(erow, part)
	if ok {
		return nil
	}

	// external command
	cmd := strings.TrimSpace(part.Str)
	cmdutil.ExternalCmd(erow, cmd)
	return nil
}

// Returns true if cmd was handled.
func rowToolbarCmd(erow *ERow, part *toolbardata.Part) bool {
	row := erow.Row()
	p0 := part.Args[0].Str
	switch p0 {
	case "Save":
		cmdutil.SaveRowFile(erow)
	case "Reload":
		cmdutil.ReloadRow(erow)
	case "Close":
		row.Close()
	case "CloseColumn":
		row.Col.Cols.CloseColumnEnsureOne(row.Col)
	case "Find":
		s := part.JoinArgsFromIndex(1).Unquote()
		tautil.Find(row.TextArea, s)
	case "GotoLine":
		s := part.JoinArgsFromIndex(1).Str
		tautil.GotoLine(row.TextArea, s)
	case "Replace":
		cmdutil.Replace(erow, part)
	case "Stop":
		cmdutil.RowCtxCancel(row)
	case "ListDir":
		tree, hidden := false, false
		cmdutil.ListDirEd(erow, tree, hidden)
	case "ListDirSub":
		tree, hidden := true, false
		cmdutil.ListDirEd(erow, tree, hidden)
	case "ListDirHidden":
		tree, hidden := false, true
		cmdutil.ListDirEd(erow, tree, hidden)
	case "FWStatus":
		FWStatus(erow)
	default:
		return false
	}
	return true
}