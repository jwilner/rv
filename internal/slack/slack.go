package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/slack-go/slack"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/jwilner/rv/pkg/pb/rvapi"
)

// New returns a mux that serve slack interactive commands at `/interactive` and slash commands at `/slashcommand`.
// The mux performs slack signing verification per https://api.slack.com/authentication/verifying-requests-from-slack
func New(token, signingSecret string, conn *grpc.ClientConn) http.Handler {
	h := newHandler(token, conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/interactive", h.ServeInteractive)
	mux.HandleFunc("/slashcommand", h.ServeSlashCommand)

	return slackVerifyMiddleware(signingSecret, mux)
}

// slackVerifyMiddleware verifies that every inbound request on this handler matches the slack validation logic
// - X-Slack-Request-Timestamp must be no more than five minutes old
// - X-Slack-Signature matches the hmac/sha256 of the message
func slackVerifyMiddleware(signingSecret string, next http.Handler) http.Handler {
	swallowPanic := func(f func() error) (err error) {
		defer func() {
			p := recover()
			if p == nil {
				return
			}
			if pErr, ok := p.(error); ok {
				err = pErr
			}
			err = fmt.Errorf("panicked: %v", err)
		}()
		err = f()
		return
	}

	bufPool := sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := bufPool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			bufPool.Put(buf)
		}()

		if err := swallowPanic(func() error {
			_, err := buf.ReadFrom(r.Body)
			return err
		}); handleErr(w, "buf.ReadFrom", err) {
			return
		}

		verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			log.Printf("Rejecting invalid slack request headers: %v", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		_, _ = verifier.Write(buf.Bytes()) // can't fail and even if it did, ensure would be invalid
		if err := verifier.Ensure(); err != nil {
			log.Printf("Rejecting slack request which does not match scret: %v", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// provide buffered body for later reads
		r.Body = ioutil.NopCloser(buf)

		next.ServeHTTP(w, r)
	})
}

func newHandler(token string, conn *grpc.ClientConn) *handler {
	return &handler{slack: slack.New(token), rver: rvapi.NewRVerClient(conn)}
}

type handler struct {
	slack *slack.Client
	rver  rvapi.RVerClient
}

func (h *handler) ServeInteractive(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		wrongMethod(w)
		return
	}

	if handleErr(w, "ParseForm", r.ParseForm()) {
		return
	}

	var ic slack.InteractionCallback
	if handleErr(w, "Unmarshal", json.NewDecoder(
		strings.NewReader(r.PostFormValue("payload")),
	).Decode(&ic)) {
		return
	}

	go h.handleAsync(&ic)
}

func (h *handler) handleAsync(ic *slack.InteractionCallback) {
	ctx := context.Background()

	var err error
	switch ic.Type {
	case slack.InteractionTypeViewSubmission:
		switch ic.View.Title.Text {
		case "Ranked Choice Vote":
			err = h.handleCreateView(ctx, ic)
		case "Rank Your Choices":
			err = h.handleVote(ctx, ic)
		default:
			err = fmt.Errorf("unknown view: %q", ic.View.Title.Text)
		}
	case slack.InteractionTypeBlockActions:
		if len(ic.ActionCallback.BlockActions) != 1 {
			err = fmt.Errorf("expected actions of length 1; got %v", len(ic.ActionCallback.BlockActions))
			break
		}
		switch ic.ActionCallback.BlockActions[0].ActionID {
		case "launch-vote":
			err = h.handleLaunchVote(ctx, ic)
		case "add_option":
			err = h.handleAddOption(ctx, ic)
		default:
			err = fmt.Errorf("unknown block action: %q", ic.ActionCallback.BlockActions[0].ActionID)
		}
	default:
		err = fmt.Errorf("unknown interactivity type: %q", ic.Type)
	}
	if err != nil {
		log.Printf("handling async failed: %v", err)
	}
}

func (h *handler) checkIn(ctx context.Context) (context.Context, error) {
	var md metadata.MD
	_, err := h.rver.CheckIn(ctx, &rvapi.CheckInRequest{}, grpc.Header(&md))
	if err != nil {
		return nil, err
	}
	return metadata.AppendToOutgoingContext(ctx, "rv-token", md["rv-token"][0]), nil
}

func (h *handler) handleCreateView(ctx context.Context, ic *slack.InteractionCallback) error {
	req := viewToCreateRequest(ic)
	if req == nil {
		return errors.New("invalid view")
	}

	var err error
	if ctx, err = h.checkIn(ctx); err != nil {
		return fmt.Errorf("unable to check in: %w", err)
	}

	resp, err := h.rver.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create: %w", err)
	}

	channelID := ic.View.PrivateMetadata

	_, _, err = h.slack.PostMessageContext(
		ctx,
		channelID,
		slack.MsgOptionBlocks(renderElection(resp.Election, nil)...),
	)
	if err != nil {
		return fmt.Errorf("unable to post message content: %w", err)
	}
	return nil
}

func (h *handler) handleLaunchVote(ctx context.Context, ic *slack.InteractionCallback) error {
	ballotKey := ic.ActionCallback.BlockActions[0].Value

	resp, err := h.rver.GetView(ctx, &rvapi.GetViewRequest{BallotKey: ballotKey})
	if err != nil {
		return fmt.Errorf("rver.GetView: %w", err)
	}

	mvr := slack.ModalViewRequest{
		PrivateMetadata: ic.Channel.ID + "|" + ic.Message.Timestamp + "|" + ballotKey,
		Type:            slack.VTModal,
		Title: &slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "Rank Your Choices",
		},
		Close: &slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "Cancel",
		},
		Submit: &slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "Submit",
		},
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				&slack.HeaderBlock{
					Type: slack.MBTHeader,
					Text: &slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: resp.Election.Question,
					},
				},
			},
		},
	}
	for i := range resp.Election.Choices {
		options := make([]*slack.OptionBlockObject, 0, len(resp.Election.Choices))
		for i, c := range resp.Election.Choices {
			options = append(options, &slack.OptionBlockObject{
				Text: &slack.TextBlockObject{Type: slack.PlainTextType, Text: c},
				// value field has a length limit, so use ints for values and handle the indirection later
				Value: strconv.Itoa(i),
			})
		}

		mvr.Blocks.BlockSet = append(
			mvr.Blocks.BlockSet,
			&slack.InputBlock{
				Type:    slack.MBTInput,
				BlockID: fmt.Sprintf("choice_block_%d", i),
				Label:   &slack.TextBlockObject{Type: slack.PlainTextType, Text: fmt.Sprintf("Choice %d", i+1)},
				Element: &slack.SelectBlockElement{
					Type: slack.OptTypeStatic,
					Placeholder: &slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: fmt.Sprintf("Choice %d...", i+1),
					},
					ActionID: fmt.Sprintf("choice_input_%d", i),
					Options:  options,
				},
				Optional: true,
			},
		)
	}

	_, err = h.slack.OpenViewContext(ctx, ic.TriggerID, mvr)
	if err != nil {
		return fmt.Errorf("slack.OpenViewContext: %w", err)
	}
	return nil
}

func (h *handler) handleVote(ctx context.Context, ic *slack.InteractionCallback) error {
	parts := strings.Split(ic.View.PrivateMetadata, "|")
	if len(parts) != 3 {
		return errors.New("invalid view")
	}
	channelID, msgTS, ballotKey := parts[0], parts[1], parts[2]

	req := viewToVoteRequest(ballotKey, ic)
	if req == nil {
		return errors.New("invalid view")
	}
	var err error
	if ctx, err = h.checkIn(ctx); err != nil {
		return fmt.Errorf("checkIn: %w", err)
	}
	_, err = h.rver.Vote(ctx, req)
	if err != nil {
		return fmt.Errorf("vote: %w", err)
	}

	el, err := h.rver.GetView(ctx, &rvapi.GetViewRequest{BallotKey: ballotKey})
	if err != nil {
		return fmt.Errorf("rver.GetView: %w", err)
	}

	rep, err := h.rver.Report(ctx, &rvapi.ReportRequest{BallotKey: ballotKey})
	if err != nil {
		return fmt.Errorf("rver.Report: %w", err)
	}

	_, _, _, err = h.slack.UpdateMessageContext(
		ctx,
		channelID,
		msgTS,
		slack.MsgOptionBlocks(renderElection(el.Election, rep.Report)...),
	)
	if err != nil {
		return fmt.Errorf("slack.UpdateMessageContext: %w", err)
	}
	return nil
}

func viewToVoteRequest(ballotKey string, ic *slack.InteractionCallback) *rvapi.VoteRequest {
	vote := rvapi.VoteRequest{
		BallotKey: ballotKey,
		Name:      ic.User.Name,
	}
	for _, r := range ic.View.Blocks.BlockSet {
		i, ok := r.(*slack.InputBlock)
		if !ok {
			continue
		}

		switch e := i.Element.(type) {
		case *slack.SelectBlockElement:
			opt := ic.View.State.Values[i.BlockID][e.ActionID].SelectedOption.Text
			if opt == nil || opt.Text == "" {
				break
			}
			vote.Choices = append(vote.Choices, opt.Text)
		}
	}
	return &vote
}

func (h *handler) handleAddOption(ctx context.Context, ic *slack.InteractionCallback) error {
	m := slack.ModalViewRequest{
		Type:            slack.VTModal,
		PrivateMetadata: ic.View.PrivateMetadata,
		Title:           ic.View.Title,
		Close:           ic.View.Close,
		Submit:          ic.View.Submit,
		Blocks:          ic.View.Blocks,
		CallbackID:      ic.CallbackID,
	}

	m.Blocks.BlockSet = append(
		m.Blocks.BlockSet[:len(m.Blocks.BlockSet)-1],
		slack.NewInputBlock(
			fmt.Sprintf("choice_block_%d", len(m.Blocks.BlockSet)-2),
			&slack.TextBlockObject{Type: slack.PlainTextType, Text: fmt.Sprintf("Choice %v", len(m.Blocks.BlockSet)-1)},
			slack.NewPlainTextInputBlockElement(
				&slack.TextBlockObject{
					Type: slack.PlainTextType,
					Text: fmt.Sprintf("choice %v...", len(m.Blocks.BlockSet)-1),
				},
				fmt.Sprintf("choice_input_%d", len(m.Blocks.BlockSet)-2),
			),
		),
		m.Blocks.BlockSet[len(m.Blocks.BlockSet)-1],
	)

	_, err := h.slack.UpdateViewContext(ctx, m, "", ic.Hash, ic.View.ID)
	if err != nil {
		return fmt.Errorf("slack.UpdateViewContext: %w", err)
	}
	return nil
}

func viewToCreateRequest(cb *slack.InteractionCallback) *rvapi.CreateRequest {
	l := 0
	for k := range cb.View.State.Values {
		if strings.HasPrefix(k, "choice_block_") {
			idx, err := strconv.Atoi(k[len("choice_block_"):])
			if err != nil {
				log.Printf("parsing %q: %v", k, err)
				return nil
			}
			if newL := idx + 1; newL > l {
				l = newL
			}
		}
	}
	req := rvapi.CreateRequest{
		Choices: make([]string, l),
	}
	for k, v := range cb.View.State.Values {
		if k == "question_block" {
			req.Question = v["question_input"].Value
			continue
		}
		if !strings.HasPrefix(k, "choice_block_") {
			continue
		}
		idx, _ := strconv.Atoi(k[len("choice_block_"):]) // already parsed, err impossible
		req.Choices[idx] = v[fmt.Sprintf("choice_input_%d", idx)].Value
	}
	return &req
}

func (h *handler) openView(channelID, triggerID string) {
	_, err := h.slack.OpenView(triggerID, slack.ModalViewRequest{
		Type:            slack.VTModal,
		PrivateMetadata: channelID,
		Title: &slack.TextBlockObject{
			Type:  slack.PlainTextType,
			Text:  "Ranked Choice Vote",
			Emoji: true,
		},
		Submit: &slack.TextBlockObject{
			Type:  slack.PlainTextType,
			Text:  "Submit",
			Emoji: true,
		},
		Close: &slack.TextBlockObject{
			Type:  slack.PlainTextType,
			Text:  "Cancel",
			Emoji: true,
		},
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				&slack.InputBlock{
					Type:    slack.MBTInput,
					BlockID: "question_block",
					Label:   &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Question"},
					Element: &slack.PlainTextInputBlockElement{
						Type:        slack.METPlainTextInput,
						ActionID:    "question_input",
						Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "question..."},
					},
				},
				&slack.InputBlock{
					Type:    slack.MBTInput,
					BlockID: "choice_block_0",
					Label:   &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Choice 1"},
					Element: &slack.PlainTextInputBlockElement{
						Type:        slack.METPlainTextInput,
						ActionID:    "choice_input_0",
						Placeholder: &slack.TextBlockObject{Type: slack.PlainTextType, Text: "choice 1..."},
					},
				},
				&slack.ActionBlock{
					Type:    slack.MBTAction,
					BlockID: "action_block",
					Elements: slack.BlockElements{
						ElementSet: []slack.BlockElement{
							&slack.ButtonBlockElement{
								Type:     slack.METButton,
								ActionID: "add_option",
								Text: &slack.TextBlockObject{
									Type:  slack.PlainTextType,
									Text:  "Add another",
									Emoji: true,
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("OpenView: %v", err)
	}
}

func (h *handler) ServeSlashCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "POST" {
		wrongMethod(w)
		return
	}

	if handleErr(w, "ParseForm", r.ParseForm()) {
		return
	}

	triggerID := r.FormValue("trigger_id")
	channelID := r.FormValue("channel_id")
	go h.openView(channelID, triggerID)
}

func handleErr(w http.ResponseWriter, tag string, err error) bool {
	if err == nil {
		return false
	}
	log.Printf("Received error %v: %v\n", tag, err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	return true
}

func wrongMethod(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

type election interface {
	GetBallotKey() string
	GetQuestion() string
	GetChoices() []string
}

func renderElection(el election, rep *rvapi.Report) []slack.Block {
	var b strings.Builder
	for _, c := range el.GetChoices() {
		_, _ = fmt.Fprintf(&b, "- %v\n", c)
	}

	blocks := []slack.Block{
		&slack.HeaderBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{
				Type: slack.PlainTextType,
				Text: el.GetQuestion(),
			},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: b.String(),
			},
		},
	}

	if rep != nil && rep.Winner != "" {
		rounds := rep.Rounds
		var b strings.Builder
		printf := func(format string, args ...interface{}) {
			_, _ = fmt.Fprintf(&b, format, args...)
		}
		if len(rounds) == 1 {
			printf(":tada::tada::tada: Winner after 1 round: %v :tada::tada::tada:\n\n", rep.Winner)
		} else {
			printf(":tada::tada::tada: Winner after %d rounds: %v :tada::tada::tada:\n\n", len(rounds), rep.Winner)
		}

		for i := len(rounds) - 1; i >= 0; i-- {
			printf("*Round %d*:\n", i+1)
			for _, t := range rounds[i].Tallies {
				printf("- %v: %v\n", t.Choice, t.Count)
			}
		}

		blocks = append(
			blocks,
			&slack.SectionBlock{
				Type: slack.MBTSection,
				Text: &slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: b.String(),
				},
			},
		)
	}

	return append(
		blocks,
		&slack.ActionBlock{
			Type: slack.MBTAction,
			Elements: slack.BlockElements{
				ElementSet: []slack.BlockElement{
					&slack.ButtonBlockElement{
						Type: slack.METButton,
						Text: &slack.TextBlockObject{
							Type: slack.PlainTextType,
							Text: "Click here to vote",
						},
						ActionID: "launch-vote",
						Value:    el.GetBallotKey(),
					},
				},
			},
		},
	)
}
