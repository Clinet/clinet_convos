package convos

import (
	"time"
)

var ConvoServices []ConvoService

func init() {
	ConvoServices = make([]ConvoService, 0)
}

type ConvoService interface {
	Login()                                                                     error                                 //Login to the conversation service
	Query(query *ConversationQuery, lastState *ConversationState) (*ConversationResponse, error) //Query the service with a conversation
}

type Conversation struct {
	History []*ConversationState //Conversation state history in order
}

//NewConversation returns an empty conversation
func NewConversation() (*Conversation) {
	return &Conversation{
		History: make([]*ConversationState, 0),
	}
}

//QueryText returns a new conversation state for the given query text and appends it to the convo history
func (convo *Conversation) QueryText(queryText string) (*ConversationState) {
	newState := &ConversationState{
		Query: &ConversationQuery{
			Time: time.Now(),
			Text: queryText,
		},
		Errors: make([]error, 0),
	}

	//Query all available convo services in registered order
	for i := 0; i < len(ConvoServices); i++ {
		resp, err := ConvoServices[i].Query(newState.Query, convo.LastState())
		if err != nil {
			newState.Errors = append(newState.Errors, err)
		} else {
			newState.Response = resp
			break
		}
	}

	if newState.Response != nil {
		convo.History = append(convo.History, newState) //Only add successful responses to the history
	}

	return newState
}

//LastResponse returns the most recent conversation state
func (convo *Conversation) LastState() (*ConversationState) {
	if len(convo.History) == 0 {
		return nil
	}
	return convo.History[len(convo.History)-1]
}