package accounts

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// The qualification test (DQ-3, US-24): pick the word for each English
// meaning, from four verified words of the language.
const (
	testQuestions = 10
	testOptions   = 4
	testPassMark  = 8
	retakeAfter   = 24 * time.Hour
)

// Question is one meaning and the words to choose from; one of them is right.
type Question struct {
	GlossEn string
	Options []Option
}

// Option is a word to choose, identified by its entry.
type Option struct {
	ID       uuid.UUID
	Headword string
}

// Test is a started qualification test.
type Test struct {
	ID        uuid.UUID
	Language  string
	Questions []Question
}

// StartTest draws a test in a language for the signed-in user. A user takes
// it once a day until they pass, and not while a pass awaits confirmation.
func (s *Service) StartTest(ctx context.Context, language string) (Test, error) {
	me, err := Require(ctx, nil)
	if err != nil {
		return Test{}, err
	}
	if slices.ContainsFunc(me.Approvals, func(a Approval) bool { return a.Language == language }) {
		return Test{}, httpx.Conflict("You already review " + language)
	}
	t := Test{Language: language}
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		q := dbgen.New(tx)
		if err := mayStart(ctx, q, me.ID, language); err != nil {
			return err
		}
		pool, err := q.QualificationPool(ctx, language)
		if err != nil {
			return fmt.Errorf("qualification pool: %w", err)
		}
		questions, expected, ok := drawTest(pool)
		if !ok {
			return httpx.Conflict("This language does not have enough verified words for a test yet")
		}
		t.Questions = questions
		t.ID, err = q.CreateQualification(ctx, dbgen.CreateQualificationParams{UserID: me.ID, LanguageCode: language, Expected: expected})
		return err
	})
	return t, err
}

func mayStart(ctx context.Context, q *dbgen.Queries, userID, language string) error {
	last, err := q.LatestQualification(ctx, dbgen.LatestQualificationParams{UserID: userID, LanguageCode: language})
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("latest qualification: %w", err)
	case last.Passed.Bool && !last.ConfirmedAt.Valid:
		return httpx.Conflict("Your pass is waiting for a senior reviewer to confirm it")
	case time.Since(last.CreatedAt) < retakeAfter:
		return httpx.Conflict("You can take the test once a day; try again tomorrow")
	}
	return nil
}

// drawTest makes questions from a shuffled pool: each question's options are
// its word and words with other meanings, shuffled. expected holds the right
// option for each question.
func drawTest(pool []dbgen.QualificationPoolRow) (questions []Question, expected []uuid.UUID, ok bool) {
	if len(pool) < testQuestions+testOptions-1 {
		return nil, nil, false
	}
	for _, q := range pool[:testQuestions] {
		options := []Option{{ID: q.EntryID, Headword: q.Headword}}
		for _, i := range rand.Perm(len(pool)) {
			o := pool[i]
			if len(options) == testOptions {
				break
			}
			if o.GlossEn != q.GlossEn && !slices.ContainsFunc(options, func(x Option) bool { return x.Headword == o.Headword }) {
				options = append(options, Option{ID: o.EntryID, Headword: o.Headword})
			}
		}
		if len(options) < testOptions {
			return nil, nil, false
		}
		rand.Shuffle(len(options), func(i, j int) { options[i], options[j] = options[j], options[i] })
		questions = append(questions, Question{GlossEn: q.GlossEn, Options: options})
		expected = append(expected, q.EntryID)
	}
	return questions, expected, true
}

// Result is a handed-in test's score.
type Result struct {
	Score, Total int
	Passed       bool
}

// AnswerTest scores the signed-in user's test, once.
func (s *Service) AnswerTest(ctx context.Context, id uuid.UUID, answers []uuid.UUID) (Result, error) {
	me, err := Require(ctx, nil)
	if err != nil {
		return Result{}, err
	}
	var r Result
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		q := dbgen.New(tx)
		t, err := lockTest(ctx, q, id)
		switch {
		case err != nil:
			return err
		case t.UserID != me.ID:
			return httpx.ErrNotFound
		case t.AnsweredAt.Valid:
			return httpx.Conflict("This test was already handed in")
		case len(answers) != len(t.Expected):
			return httpx.BadRequest(fmt.Sprintf("Answer all %d questions", len(t.Expected)))
		}
		r.Total = len(t.Expected)
		for i, want := range t.Expected {
			if answers[i] == want {
				r.Score++
			}
		}
		r.Passed = r.Score >= testPassMark
		return q.ScoreQualification(ctx, dbgen.ScoreQualificationParams{ID: id, Score: int32(r.Score), Passed: r.Passed})
	})
	return r, err
}

// ConfirmTest approves the taker of a passed test as a reviewer of its
// language. Senior reviewers of the language, or admins, confirm others' passes.
func (s *Service) ConfirmTest(ctx context.Context, id uuid.UUID) error {
	me, err := Require(ctx, nil)
	if err != nil {
		return err
	}
	return database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		q := dbgen.New(tx)
		t, err := lockTest(ctx, q, id)
		switch {
		case err != nil:
			return err
		case !me.SeniorOrAdmin(t.LanguageCode):
			return httpx.ErrForbidden
		case t.UserID == me.ID:
			return httpx.Forbidden("Another senior reviewer confirms your test")
		case !t.Passed.Bool:
			return httpx.Conflict("Only a passed test can be confirmed")
		case t.ConfirmedAt.Valid:
			return httpx.Conflict("This test is already confirmed")
		}
		if err := q.GrantReviewer(ctx, dbgen.GrantReviewerParams{UserID: t.UserID, LanguageCode: t.LanguageCode, ApprovedBy: database.Text(me.ID)}); err != nil {
			return fmt.Errorf("grant reviewer: %w", err)
		}
		return q.ConfirmQualification(ctx, dbgen.ConfirmQualificationParams{ID: id, ConfirmedBy: database.Text(me.ID)})
	})
}

func lockTest(ctx context.Context, q *dbgen.Queries, id uuid.UUID) (dbgen.LockQualificationRow, error) {
	t, err := q.LockQualification(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return t, httpx.ErrNotFound
	}
	return t, err
}

// PendingTests lists passed tests waiting for the signed-in senior reviewer.
func (s *Service) PendingTests(ctx context.Context) ([]dbgen.PendingQualificationsRow, error) {
	me, err := Require(ctx, nil)
	if err != nil {
		return nil, err
	}
	languages := me.SeniorLanguages()
	if len(languages) == 0 {
		return nil, httpx.ErrForbidden
	}
	var rows []dbgen.PendingQualificationsRow
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		rows, err = dbgen.New(tx).PendingQualifications(ctx, languages)
		return err
	})
	return rows, err
}
