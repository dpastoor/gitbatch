package git

import (
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var stashCommand = "stash"

// StashedItem holds the required fields for a stashed change
type StashedItem struct {
	StashID     int
	BranchName  string
	Hash        string
	Description string
	EntityPath  string
}

func stashGet(entity *RepoEntity, option string) string {
	args := make([]string, 0)
	args = append(args, stashCommand)
	args = append(args, option)
	out, err := GenericGitCommandWithOutput(entity.AbsPath, args)
	if err != nil {
		log.Warn("Error while stash command")
		return "?"
	}
	return out
}

func (entity *RepoEntity) loadStashedItems() error {
	entity.Stasheds = make([]*StashedItem, 0)
	output := stashGet(entity, "list")
	stashIDRegex := regexp.MustCompile(`stash@{[\d]+}:`)
	stashIDRegexInt := regexp.MustCompile(`[\d]+`)
	stashBranchRegex := regexp.MustCompile(`[\w]+: `)
	stashHashRegex := regexp.MustCompile(`[\w]{7}`)

	stashlist := strings.Split(output, "\n")
	for _, stashitem := range stashlist {
		// find id
		id := stashIDRegexInt.FindString(stashIDRegex.FindString(stashitem))
		i, err := strconv.Atoi(id)
		if err != nil {
			// probably something isn't right let's continue over this iteration
			log.Trace("cannot initiate stashed item")
			continue
		}
		// trim id section
		trimmed := stashIDRegex.Split(stashitem, 2)[1]

		// find branch
		stashBranchRegexMatch := stashBranchRegex.FindString(trimmed)
		branchName := stashBranchRegexMatch[:len(stashBranchRegexMatch)-2]

		// trim branch section
		trimmed = stashBranchRegex.Split(trimmed, 2)[1]
		hash := stashHashRegex.FindString(trimmed)

		// trim hash
		desc := stashHashRegex.Split(trimmed, 2)[1][1:]

		entity.Stasheds = append(entity.Stasheds, &StashedItem{
			StashID:     i,
			BranchName:  branchName,
			Hash:        hash,
			Description: desc,
			EntityPath:  entity.AbsPath,
		})
	}
	return nil
}

// Stash is the wrapper of convetional "git stash" command
func (entity *RepoEntity) Stash() (output string, err error) {
	args := make([]string, 0)
	args = append(args, stashCommand)

	output, err = GenericGitCommandWithErrorOutput(entity.AbsPath, args)
	entity.Refresh()
	return output, err
}

// Pop is the wrapper of "git stash pop" command that used for a file
func (stashedItem *StashedItem) Pop() (output string, err error) {
	args := make([]string, 0)
	args = append(args, stashCommand)
	args = append(args, "pop")
	args = append(args, "stash@{"+strconv.Itoa(stashedItem.StashID)+"}")
	output, err = GenericGitCommandWithErrorOutput(stashedItem.EntityPath, args)
	return output, err
}

// Show is the wrapper of "git stash show -p " command
func (stashedItem *StashedItem) Show() (output string, err error) {
	args := make([]string, 0)
	args = append(args, stashCommand)
	args = append(args, "show")
	args = append(args, "-p")
	args = append(args, "stash@{"+strconv.Itoa(stashedItem.StashID)+"}")
	output, err = GenericGitCommandWithErrorOutput(stashedItem.EntityPath, args)
	return output, err
}
