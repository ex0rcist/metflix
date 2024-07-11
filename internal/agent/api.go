package agent

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/metrics"
)

type API struct {
	address    *entities.Address
	httpClient *http.Client
}

func NewAPI(address *entities.Address, httpTransport http.RoundTripper) *API {
	if httpTransport == nil {
		httpTransport = http.DefaultTransport
	}

	client := &http.Client{
		Timeout:   2 * time.Second,
		Transport: httpTransport,
	}

	return &API{
		address:    address,
		httpClient: client,
	}
}

func (c *API) Report(name string, metric metrics.Metric) *API {
	// todo: another transport?
	url := "http://" + c.address.String() + fmt.Sprintf("/update/%s/%s/%s", metric.Kind(), name, metric)

	req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
	if err != nil {
		logging.LogError(err, "httpRequest error")
	}

	req.Header.Set("Content-Type", "text/plain")

	logging.LogInfo(fmt.Sprintf("sending POST to %v", url))

	resp, err := c.httpClient.Do(req)

	if err != nil {
		logging.LogError(err, "httpClient error")
		return c
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body) // нужно прочитать ответ для keepalive?
	if err != nil {
		logging.LogError(entities.ErrMetricReport, "error reading response body")
	}

	if resp.StatusCode != http.StatusOK {
		logging.LogError(entities.ErrMetricReport, string(respBody))
	}

	return c
}
