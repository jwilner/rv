package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/bits"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgtype"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/jwilner/rv/models"
)

// electionPage is the main page for managing an election
type electionPage struct {
	Election *election
	Report   *report
	Form     electionForm
	Zones    []string
	Now      time.Time
}

func (e *electionPage) ScheduleValue() time.Time {
	if e.Election.CloseScheduled(e.Now) {
		return e.Election.CloseTime()
	}
	return e.Now.AddDate(0, 0, 1)
}

// electionForm is the election management form
type electionForm struct {
	form

	Op uint

	CloseDate string
	CloseTime string
	CloseTZ   string
}

func (e *electionForm) unmarshal(vals url.Values) {
	for _, s := range []struct {
		field string
		flag  uint
	}{
		{"unset", updateUnset},
		{"scheduleClose", updateCloseByDate},
		{"closeNow", updateCloseNow},
		{"setPrivate", updateSetPrivate},
		{"setPublic", updateSetPublic},
	} {
		if vals.Get(s.field) != "" {
			e.Op |= s.flag
		}
	}

	e.CloseDate = vals.Get("closeDate")
	e.CloseTime = vals.Get("closeTime")
	e.CloseTZ = vals.Get("closeTZ")
}

type electionUpdate struct {
	update uint
	time   time.Time
}

const (
	updateCloseNow = 1 << iota
	updateCloseByDate
	updateUnset
	updateSetPrivate
	updateSetPublic
)

func (e *electionForm) validate(now time.Time) (eu electionUpdate) {
	fmt.Printf("%+v", *e)
	if bits.OnesCount(e.Op) > 1 {
		e.setErrorf("updateCloseByDate", "Cannot close both now and at date")
		return
	}

	switch e.Op {
	case updateCloseByDate:
		t := e.CloseDate + "T" + e.CloseTime
		if loc, err := time.LoadLocation(e.CloseTZ); err != nil {
			e.setErrorf("closeTZ", "Unrecognized time zone: %v", e.CloseTZ)
		} else if eu.time, err = time.ParseInLocation("2006-01-02T15:04", t, loc); err != nil {
			e.setErrorf("closeDate", "Invalid date time: %v", t)
			e.setErrorf("closeTime", "Invalid date time: %v", t)
		}
		if !eu.time.After(now) {
			e.setErrorf("closeDate", "Must be in future: %v", t)
			e.setErrorf("closeTime", "Must be in future: %v", t)
		}
	}

	eu.update = e.Op

	if !e.checkErrors() {
		eu.update = 0
	}
	return
}

func (h *handler) getElection(w http.ResponseWriter, r *http.Request) {
	var (
		el    *models.Election
		votes []*models.Vote
	)
	err := h.txM.inTx(r.Context(), &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		if el, err = models.Elections(models.ElectionWhere.Key.EQ(keyParam(ctx))).One(ctx, tx); err != nil {
			return fmt.Errorf("models.Elections key=%v: %w", keyParam(ctx), err)
		}
		if votes, err = models.Votes(models.VoteWhere.ElectionID.EQ(el.ID)).All(ctx, tx); err != nil {
			return fmt.Errorf("models.Votes electionId=%d: %w", el.ID, err)
		}
		return
	})
	if handleError(w, r, err) {
		return
	}
	h.tmpls.render(r.Context(), w, "election.html", &electionPage{
		Election: newElection(el),
		Report:   calculateReport(votes),
		Zones:    h.tzes,
		Now:      time.Now().UTC(),
	})
}

func (h *handler) postElection(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()

	if handleError(w, r, r.ParseForm()) {
		return
	}

	var ef electionForm
	ef.unmarshal(r.PostForm)

	parsed := ef.validate(now)

	var (
		el    *election
		votes []*models.Vote
	)
	err := h.txM.inTx(r.Context(), nil, func(ctx context.Context, tx *sql.Tx) (err error) {
		elM, err := models.Elections(models.ElectionWhere.Key.EQ(keyParam(r.Context()))).One(ctx, tx)
		if err != nil {
			return fmt.Errorf("models.Elections key=%v: %w", keyParam(ctx), err)
		}
		el = newElection(elM)
		if votes, err = models.Votes(models.VoteWhere.ElectionID.EQ(el.ID)).All(ctx, tx); err != nil {
			return fmt.Errorf("models.Votes electionId=%d: %w", el.ID, err)
		}
		var modified []string
		switch parsed.update {
		case updateCloseNow:
			if !el.CanCloseNow(now) {
				ef.setErrorf("closeNow", "cannot close now -- already closed")
				return
			}
			modified = el.SetClose(time.Now().UTC())
		case updateCloseByDate:
			if !el.CanSchedule(now) {
				ef.setErrorf("scheduleClose", "cannot schedule close -- already closed")
				return
			}
			modified = el.SetClose(parsed.time)
		case updateUnset:
			if !el.CanUnset(now) {
				ef.setErrorf("unset", "cannot reopen -- already open")
				return
			}
			modified = el.UnsetClose()
		case updateSetPublic:
			if el.Public() {
				ef.setErrorf("setPublic", "cannot set public -- already public")
				return
			}
			modified = el.SetPublic()
		case updateSetPrivate:
			if !el.Public() {
				ef.setErrorf("setPrivate", "cannot set private -- already private")
				return
			}
			modified = el.UnsetPublic()
		}
		if len(modified) == 0 {
			log.Println("no update necessary")
			return
		}
		if _, err := el.Update(ctx, tx, boil.Whitelist(modified...)); err != nil {
			return fmt.Errorf("models.Election.Update %v failed: %w", modified, err)
		}
		return
	})
	if handleError(w, r, err) {
		return
	}
	h.tmpls.render(r.Context(), w, "election.html", &electionPage{
		Form:     ef,
		Election: el,
		Report:   calculateReport(votes),
		Zones:    h.tzes,
		Now:      now,
	})
}

// Elections are either:
// - open (close NULL)
// - scheduled (close a time in the future)
// - closed (close a time not in the future)
// Updates:
// - schedule close: open or scheduled -> scheduled (set a close time in future)
// - close now: open -> closed (set close to now)
// - unset: closed -> open (set close to null)
// - unset: scheduled -> open (set close to null)

func newElection(e *models.Election) *election {
	el := &election{e, nil}
	if e.Close.Status == pgtype.Undefined {
		e.Close.Status = pgtype.Null
		e.CloseTZ = null.String{}
	}
	if e.CreatedAt.Status == pgtype.Undefined {
		_ = e.CreatedAt.Set(time.Now().UTC())
	}
	return el
}

type election struct {
	*models.Election

	loc *time.Location
}

func (e *election) CloseScheduled(now time.Time) bool {
	return e.Close.Status == pgtype.Present && now.Before(e.CloseTime())
}

func (e *election) Closed(now time.Time) bool {
	return e.Close.Status == pgtype.Present && !now.Before(e.CloseTime())
}

func (e *election) CloseTime() time.Time {
	if e.loc == nil && e.CloseTZ.Valid {
		var err error
		if e.loc, err = time.LoadLocation(e.CloseTZ.String); err != nil {
			panic(err) // this should never happen
		}
		_ = e.Close.Set(e.Close.Time.In(e.loc))
	}
	return e.Close.Time
}

func (e *election) CanCloseNow(now time.Time) bool {
	return !e.Closed(now)
}

func (e *election) CanSchedule(now time.Time) bool {
	return !e.Closed(now)
}

func (e *election) CanUnset(now time.Time) bool {
	return e.Closed(now) || e.CloseScheduled(now)
}

func (e *election) SetClose(t time.Time) []string {
	e.loc = t.Location()
	e.CloseTZ = null.NewString(e.loc.String(), true)
	_ = e.Close.Set(t)

	return []string{models.ElectionColumns.Close, models.ElectionColumns.CloseTZ}
}

func (e *election) UnsetClose() []string {
	e.loc = nil
	e.CloseTZ = null.String{}
	_ = e.Close.Set(nil)

	return []string{models.ElectionColumns.Close, models.ElectionColumns.CloseTZ}
}

const (
	electionFlagPublic = "public"
)

func (e *election) Public() bool {
	return contains(e.Flags, electionFlagPublic)
}

func (e *election) SetPublic() []string {
	e.Flags = add(e.Flags, electionFlagPublic)
	return []string{models.ElectionColumns.Flags}
}

func (e *election) UnsetPublic() []string {
	l := remove(e.Flags, electionFlagPublic)
	e.Flags = e.Flags[:l]
	return []string{models.ElectionColumns.Flags}
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func remove(haystack []string, needle string) int {
	cur := 0
	for i, v := range haystack {
		if haystack[i] == needle {
			continue
		}
		haystack[cur] = v
		cur++
	}
	return cur
}

func add(haystack []string, needle string) []string {
	if contains(haystack, needle) {
		return haystack
	}
	return append(haystack, needle)
}
