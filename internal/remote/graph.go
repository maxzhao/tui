package remote

import (
	"context"
	"net/url"
	"strconv"
)

type Reference struct {
	Source     string   `json:"source"`
	Line       *int     `json:"line"`
	NodeID     *string  `json:"nodeId"`
	Kind       string   `json:"kind"`
	Raw        string   `json:"raw"`
	Anchor     string   `json:"anchor"`
	Status     string   `json:"status"`
	Path       *string  `json:"path"`
	Candidates []string `json:"candidates"`
}

func (c *Client) References(ctx context.Context, kind, path string) ([]Reference, error) {
	out := []Reference{}
	err := c.Get(ctx, "/api/links/"+kind+"?path="+url.QueryEscape(path), &out)
	return out, err
}
func (c *Client) Orphans(ctx context.Context) ([]map[string]string, error) {
	out := []map[string]string{}
	err := c.Get(ctx, "/api/links/orphans", &out)
	return out, err
}

type PendingThread struct {
	Path        string `json:"path"`
	ThreadID    string `json:"threadId"`
	AnchorText  string `json:"anchorText"`
	AnchorLine  int    `json:"anchorLine"`
	AnchorFound bool   `json:"anchorFound"`
	Messages    []struct {
		ID        string `json:"id"`
		Author    string `json:"author"`
		Body      string `json:"body"`
		CreatedAt int64  `json:"createdAt"`
	} `json:"messages"`
}

func (c *Client) PendingComments(ctx context.Context, author string, wait bool, timeout int) ([]PendingThread, error) {
	out := []PendingThread{}
	client := *c
	httpClient := *c.http
	if wait {
		httpClient.Timeout = 0
	}
	client.http = &httpClient
	params := url.Values{"author": {author}, "wait": {strconv.FormatBool(wait)}, "timeout": {strconv.Itoa(timeout)}}
	err := client.Get(ctx, "/api/comments/pending?"+params.Encode(), &out)
	return out, err
}
