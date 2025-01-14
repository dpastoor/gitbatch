package git

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

// RemoteBranch is the wrapper of go-git's Reference struct. In addition to
// that, it also holds name of the remote branch
type RemoteBranch struct {
	Name      string
	Reference *plumbing.Reference
	Deleted   bool
}

// NextRemoteBranch iterates to the next remote branch
func (remote *Remote) NextRemoteBranch() error {
	currentRemoteIndex := remote.findCurrentRemoteBranchIndex()
	if currentRemoteIndex == len(remote.Branches)-1 {
		remote.Branch = remote.Branches[0]
	} else {
		remote.Branch = remote.Branches[currentRemoteIndex+1]
	}
	return nil
}

// PreviousRemoteBranch iterates to the previous remote branch
func (remote *Remote) PreviousRemoteBranch() error {
	currentRemoteIndex := remote.findCurrentRemoteBranchIndex()
	if currentRemoteIndex == 0 {
		remote.Branch = remote.Branches[len(remote.Branches)-1]
	} else {
		remote.Branch = remote.Branches[currentRemoteIndex-1]
	}
	return nil
}

// returns the active remote branch index
func (remote *Remote) findCurrentRemoteBranchIndex() int {
	currentRemoteIndex := 0
	for i, rb := range remote.Branches {
		if rb.Reference.Hash() == remote.Branch.Reference.Hash() {
			currentRemoteIndex = i
		}
	}
	return currentRemoteIndex
}

// search for the remote branches of the remote. It takes the go-git's repo
// pointer in order to get storer struct
func (remote *Remote) loadRemoteBranches(entity *RepoEntity) error {
	remote.Branches = make([]*RemoteBranch, 0)
	bs, err := remoteBranchesIter(entity.Repository.Storer)
	if err != nil {
		log.Warn("Cannot initiate iterator " + err.Error())
		return err
	}
	defer bs.Close()
	err = bs.ForEach(func(b *plumbing.Reference) error {
		deleted := false
		if strings.Split(b.Name().Short(), "/")[0] == remote.Name {
			remote.Branches = append(remote.Branches, &RemoteBranch{
				Name:      b.Name().Short(),
				Reference: b,
				Deleted:   deleted,
			})
		}
		return nil
	})
	if err != nil {
		return err
	}
	return err
}

// create an iterator for the references. it checks if the reference is a hash
// reference
func remoteBranchesIter(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		log.Warn("Cannot find references " + err.Error())
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		if ref.Type() == plumbing.HashReference {
			return ref.Name().IsRemote()
		}
		return false
	}, refs), nil
}

// switches to the given remote branch
func (remote *Remote) switchRemoteBranch(remoteBranchName string) error {
	for _, rb := range remote.Branches {
		if rb.Name == remoteBranchName {
			remote.Branch = rb
			return nil
		}
	}
	return errors.New("Remote branch not found.")
}
