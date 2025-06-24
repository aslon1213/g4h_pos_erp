package client

import (
	models "aslon1213/magazin_pos/pkg/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (c *Client) ParseJournalOutputSingle(resp *http.Response) models.JournalOutput {
	output := models.JournalOutput{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
	}
	err = json.Unmarshal(body, &output)
	if err != nil {

		log.Error().Err(err).Str("body", string(body)).Msg("Failed to parse journal output")
		return models.JournalOutput{}
	}
	return output
}
func (c *Client) ParseJournalOutputList(resp *http.Response) models.JournalOutputList {
	output := models.JournalOutputList{}
	err := json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse journal output list")
		return models.JournalOutputList{}
	}
	return output
}

func (c *Client) OpenJournal(journal models.NewJournalEntryInput) (*http.Response, models.JournalOutput, error) {

	body, err := json.Marshal(journal)
	if err != nil {
		return nil, models.JournalOutput{}, err
	}
	endpoint := fmt.Sprintf("/journals")

	resp, err := c.MakeRequest(
		"POST",
		endpoint,
		body,
		map[string]string{
			"Content-Type": "application/json",
		},
		false,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open journal")
		return nil, models.JournalOutput{}, err
	}

	// log.Info().Interface("resp", resp).Msg("Open journal")

	return resp, c.ParseJournalOutputSingle(resp), nil
}

func (c *Client) GetJournalByID(journal_id string) (*http.Response, models.JournalOutput, error) {

	endpoint := fmt.Sprintf("/journals/%s", journal_id)

	resp, err := c.MakeRequest(
		"GET",
		endpoint,
		nil,
		map[string]string{},
		false,
	)
	if err != nil {
		return nil, models.JournalOutput{}, err
	}
	return resp, c.ParseJournalOutputSingle(resp), nil
}

func (c *Client) QueryJournalEntries(branch_id string, queryParams models.JournalQueryParams) (*http.Response, models.JournalOutputList, error) {

	endpoint := fmt.Sprintf("/journals/branch/%s?page=%d&page_size=%d", branch_id, queryParams.Page, queryParams.PageSize)

	resp, err := c.MakeRequest(
		"GET",
		endpoint,
		nil,
		map[string]string{
			"Content-Type": "application/json",
		},
		false,
	)
	if err != nil {
		return nil, models.JournalOutputList{}, err
	}
	return resp, c.ParseJournalOutputList(resp), nil
}

func (c *Client) CloseJournal(journal_id string, input models.CloseJournalEntryInput) (*http.Response, models.JournalOutput, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, models.JournalOutput{}, err
	}

	endpoint := fmt.Sprintf("/journals/%s/close", journal_id)

	resp, err := c.MakeRequest(
		"POST",
		endpoint,
		body,
		map[string]string{
			"Content-Type": "application/json",
		},
		false,
	)
	if err != nil {
		return nil, models.JournalOutput{}, err
	}
	return resp, c.ParseJournalOutputSingle(resp), nil
}

func (c *Client) ReOpenJournal(journal_id string) (*http.Response, models.JournalOutput, error) {

	endpoint := fmt.Sprintf("/journals/%s/reopen", journal_id)

	resp, err := c.MakeRequest(
		"POST",
		endpoint,
		nil,
		map[string]string{
			"Content-Type": "application/json",
		},
		false,
	)
	if err != nil {
		return nil, models.JournalOutput{}, err
	}
	return resp, c.ParseJournalOutputSingle(resp), nil
}

func (c *Client) NewJournalOperation(journal_id string, operation models.JournalOperationInput) (*http.Response, models.JournalOutput, error) {

	body, err := json.Marshal(operation)
	if err != nil {
		return nil, models.JournalOutput{}, err
	}

	endpoint := fmt.Sprintf("/journals/%s/operations", journal_id)

	resp, err := c.MakeRequest(
		"POST",
		endpoint,
		body,
		map[string]string{
			"Content-Type": "application/json",
		},
		false,
	)
	if err != nil {
		return nil, models.JournalOutput{}, err
	}
	return resp, c.ParseJournalOutputSingle(resp), nil
}
