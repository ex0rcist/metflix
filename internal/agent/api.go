package agent

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/rs/zerolog/log"
)

type API struct {
	baseURL    string
	httpClient *http.Client
	err        error
}

func NewAPI(serverAddr string, httpTransport http.RoundTripper) *API {
	if httpTransport == nil {
		httpTransport = http.DefaultTransport
	}

	client := &http.Client{
		Timeout:   2 * time.Second,
		Transport: httpTransport,
	}

	return &API{
		baseURL:    serverAddr,
		httpClient: client,
		err:        nil,
	}
}

func (c *API) Report(name string, metric metrics.Metric) *API {
	log.Info().Msg("reporting stats ... ")

	url := c.baseURL + fmt.Sprintf("/update/%s/%s/%s", metric.Kind(), name, metric)

	req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "text/plain")

	log.Info().Msg(fmt.Sprintf("sending POST to %v", url))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.err = err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body) // нужно прочитать ответ для keepalive?
	if err != nil {
		c.err = err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(respBody) // todo: log
	}

	return c
}
