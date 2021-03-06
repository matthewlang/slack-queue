package service

import (
	"github.com/golang/glog"
	"github.com/slack-go/slack"

	"fmt"
	"net/http"
)

func (a *RemoveAction) Handle(action *slack.InteractionCallback, s *QueueService, w http.ResponseWriter) {
	user := &action.User
	ok, err := a.perms.IsAdmin(user)
	if err != nil {
		glog.Errorf("Error checking admin status of %v (%v): %v", user.ID, user.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		glog.Errorf("Permission denied to user %v (%v)", user.ID, user.Name)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var pos int
	var token int64
	var found bool
	// Remove is a block action and should be in the actions for this callback.
	for _, act := range action.ActionCallback.BlockActions {
		if ParseAction(act.ActionID) == removeActionName {
			pos, token, err = ParseActionValue(act.Value)
			found = true
			if err != nil {
				glog.Error("Error parsing action value %v: %v", act.Value, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			break
		}
	}

	if !found {
		glog.Errorf("Remove action not found for remove callback!")
	}

	glog.Infof("Removing position %d with token %d", pos, token)

	req := &RemoveRequest{Pos: pos, Token: token}
	resp := &RemoveResponse{}
	err = s.Remove(req, resp)
	if err != nil {
		glog.Errorf("Unexpected error for remove: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	fu, err := a.ul.Lookup(user.ID)
	if err == nil {
		user = fu
	}

	// Post admin message
	var str string
	if resp.Err != nil {
		glog.Infof("Stale remove for token %d, current token %d", req.Token, resp.Token)
		// str = "Remove failed: Queue has been modified since listing."
	} else {
		glog.Infof("Successfully removed pos %d, new sequence %d", req.Pos, resp.Token)
		str = fmt.Sprintf("%s removed position %d\n", userToLink(user), req.Pos+1)
	}
	a.perms.SendAdminMessage(str)

	// Replace list with updated state.
	updateListInUI(action, s, a.api)
	return
}
