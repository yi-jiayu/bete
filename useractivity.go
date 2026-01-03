package bete

import (
	"context"
	"encoding/json"
)

type EtaRequestMessageActivity struct {
	StopID string   `json:"stop_id"`
	Filter []string `json:"filter"`
}

func (a EtaRequestMessageActivity) name() string {
	return "eta_request_message"
}

type NonEtaMessageActivity struct{}

func (a NonEtaMessageActivity) name() string {
	return "non_eta_message"
}

type EtaRefreshActivity struct {
	StopID          string   `json:"stop_id"`
	Filter          []string `json:"filter"`
	Format          string   `json:"format"`
	IsInlineMessage bool     `json:"is_inline_message"`
}

func (a EtaRefreshActivity) name() string {
	return "eta_refresh"
}

type EtaResendActivity struct {
	StopID string   `json:"stop_id"`
	Filter []string `json:"filter"`
	Format string   `json:"format"`
}

func (a EtaResendActivity) name() string {
	return "eta_resend"
}

type InlineQueryActivity struct {
	Query string `json:"query"`
}

func (a InlineQueryActivity) name() string {
	return "inline_query"
}

type InlineQueryResultSelectedActivity struct {
	ResultID string `json:"result_id"`
}

func (a InlineQueryResultSelectedActivity) name() string {
	return "inline_query_result_selected"
}

type FavouriteAddedActivity struct{}

type FavouriteRemovedActivity struct{}

type FavouriteEtaRequestedActivity struct {
	StopID string `json:"stop_id"`
}

type LocationQueryActivity struct{}

func (a LocationQueryActivity) name() string {
	return "location_query"
}

type LocationQueryResultSelectedActivity struct {
	StopID string `json:"stop_id"`
}

func (a LocationQueryResultSelectedActivity) name() string {
	return "location_query_result_selected"
}

type Activity interface {
	name() string
}
type UserActivityRepository struct {
	DB Conn
}

func (r *UserActivityRepository) Record(ctx context.Context, userID int64, activity Activity) error {
	if r == nil {
		return nil
	}
	data, err := json.Marshal(activity)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(`insert into user_activity (user_id, type, data) values ($1, $2, $3)`, userID, activity.name(), data)
	return err
}
