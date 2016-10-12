package edit

import (
	"jmigpin/editor/ui"
	"os/exec"
	"path"
)

func filemanagerShortcut(ed *Editor, row *ui.Row) {
	dir := ""
	tsd := ed.rowToolbarStringData(row)
	d, ok := tsd.FirstPartDirectory()
	if ok {
		dir = d
	} else {
		f, ok := tsd.FirstPartFilename()
		if ok {
			dir = path.Dir(f)
		}
	}
	c := exec.Command("filemanager.sh", dir)
	go c.Run()
}
