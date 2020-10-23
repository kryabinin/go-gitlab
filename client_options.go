// Package gitlab - client options
package gitlab

type (
	// ClientOption to use optional parameters in slack client
	ClientOption interface {
		apply(c *client)
	}

	withHttpClient struct {
		httpClient HTTPClient
	}

	withBaseUrl struct {
		baseUrl string
	}
)

// WithHttpClient replaces default http client
func WithHttpClient(httpClient HTTPClient) ClientOption {
	return &withHttpClient{httpClient: httpClient}
}

func (opt *withHttpClient) apply(c *client) {
	c.httpClient = opt.httpClient
}

// WithBaseUrl replaces default base url
func WithBaseUrl(baseUrl string) ClientOption {
	return &withBaseUrl{baseUrl: baseUrl}
}

func (opt withBaseUrl) apply(c *client) {
	c.baseUrl = opt.baseUrl
}
