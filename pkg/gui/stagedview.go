package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/pkg/git"
	"github.com/jroimartin/gocui"
)

// staged view
func (gui *Gui) openStageView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView(stageViewFeature.Name, 6, 5, maxX/2-1, int(0.75*float32(maxY))-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = stageViewFeature.Title
	}
	entity := gui.getSelectedRepository()
	if err := refreshStagedView(g, entity); err != nil {
		return err
	}
	return gui.focusToView(stageViewFeature.Name)
}

func (gui *Gui) resetChanges(g *gocui.Gui, v *gocui.View) error {
	entity := gui.getSelectedRepository()

	_, cy := v.Cursor()
	_, oy := v.Origin()
	if len(stagedFiles) <= 0 || len(stagedFiles) <= cy+oy {
		return nil
	}
	if err := git.Reset(entity, stagedFiles[cy+oy], git.ResetOptions{}); err != nil {
		return err
	}
	return refreshAllStatusView(g, entity, true)
}

func (gui *Gui) resetAllChanges(g *gocui.Gui, v *gocui.View) error {
	entity := gui.getSelectedRepository()
	ref, err := entity.Repository.Head()
	if err != nil {
		return err
	}
	if err := git.ResetAll(entity, git.ResetOptions{
		Hash:  ref.Hash().String(),
		Rtype: git.ResetMixed,
	}); err != nil {
		return err
	}
	return refreshAllStatusView(g, entity, true)
}

// refresh the main view and re-render the repository representations
func refreshStagedView(g *gocui.Gui, entity *git.RepoEntity) error {
	stageView, err := g.View(stageViewFeature.Name)
	if err != nil {
		return err
	}
	stageView.Clear()
	_, cy := stageView.Cursor()
	_, oy := stageView.Origin()
	for i, file := range stagedFiles {
		var prefix string
		if i == cy+oy {
			prefix = prefix + selectionIndicator
		}
		fmt.Fprintf(stageView, "%s%s%s %s\n", prefix, green.Sprint(string(file.X)), red.Sprint(string(file.Y)), file.Name)
	}
	return nil
}
