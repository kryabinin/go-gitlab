// Package gitlab - client options
package gitlab

type (
	// ClientOption to use optional parameters in gitlab client
	ClientOption interface {
		apply(c *client)
	}

	clientOptionHttpClient struct {
		httpClient HTTPClient
	}

	clientOptionBaseUrl struct {
		baseUrl string
	}
)

// ClientOptionHttpClient replaces default http client
func ClientOptionHttpClient(httpClient HTTPClient) ClientOption {
	return &clientOptionHttpClient{httpClient: httpClient}
}

func (opt clientOptionHttpClient) apply(c *client) {
	c.httpClient = opt.httpClient
}

// ClientOptionBaseUrl replaces default base url
func ClientOptionBaseUrl(baseUrl string) ClientOption {
	return &clientOptionBaseUrl{baseUrl: baseUrl}
}

func (opt clientOptionBaseUrl) apply(c *client) {
	c.baseUrl = opt.baseUrl
}
