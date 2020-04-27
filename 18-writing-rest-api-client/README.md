## packagemain #18: Writing REST API Client in Go

API clients are very helpful when you're shipping your REST APIs to the public. And Go makes it easy, for you as a developer, as well as for your users, thanks to its idiomatic design and type system. But what defines a good API client?

In this tutorial, we're going to review some best practices of writing a good SDK in Go.

We'll be using [Facest.io API](https://docs.facest.io) as an example.

Before we begin to write any code, we should study the API to understand the main aspects of it such as:

- What is the Base URL of the API and can it be changed later?
- Does it support versioning?
- What are the possible errors?
- How clients should authenticate?

Understanding all of this will help you to put a right structure.

Let's start with the basics. Create a repository, pick a correct name, ideally matching the API service name. Initialize go modules. And create our main struct to hold user-specific information. This struct will contain API endpoints as functions later.

This struct should be flexible but also limited so the user can't see internal fields.

We make fields `BaseURL` and `HTTPClient` exportable, so users can use their own HTTP client if necessary.

```go
package facest

import (
	"net/http"
	"time"
)

const (
	BaseURLV1 = "https://api.facest.io/v1"
)

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: BaseURLV1,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}
```

Now let's move on and implement "Get Faces" endpoint, which returns the list of results and supports pagination, which means our function should support pagination options as input.

As I noticed in API, success responses and error responses always follow the same structure, so we can define them separately from data types and don't make them exported since this is not relevant information to the user.

```go
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
```

Make sure you don't write all endpoints in the same .go file, but group them and use separate files. For example you may group by resource type, anything that starts with `/v1/faces` goes into `faces.go` file.

I usually start by defining the types, you can do it manually or by converting JSON to go using [JSON-to-Go tool](https://mholt.github.io/json-to-go/).

```go
package facest

import "time"

type FacesList struct {
	Count      int    `json:"count"`
	PagesCount int    `json:"pages_count"`
	Faces      []Face `json:"faces"`
}

type Face struct {
	FaceToken  string      `json:"face_token"`
	FaceID     string      `json:"face_id"`
	FaceImages []FaceImage `json:"face_images"`
	CreatedAt  time.Time   `json:"created_at"`
}

type FaceImage struct {
	ImageToken string    `json:"image_token"`
	ImageURL   string    `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
}
```

The `GetFaces` function should support pagination and we can do this by adding func arguments, but these arguments are optional, and they may be changed in the future. So it makes sense to group them into a special struct:

```go
type FacesListOptions struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}
```

One more argument our function should support, and it's the context, which will let users control the API call. Users can create a Context, pass it to our func. Simple use case: cancel API call if it takes more than 5 seconds.

Now our function skeleton may look like this:

```go
func (c *Client) GetFaces(ctx context.Context, options *FacesListOptions) (*FacesList, error) {
	return nil, nil
}
```

Now it's time to make API call itself:

```go
func (c *Client) GetFaces(ctx context.Context, options *FacesListOptions) (*FacesList, error) {
	limit := 100
	page := 1
	if options != nil {
		limit = options.Limit
		page = options.Page
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/faces?limit=%d&page=%d", c.BaseURL, limit, page), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := FacesList{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
```

```go
func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	fullResponse := successResponse{
		Data: v,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}
```

Since all API endpoints act in the same manner, helper function `sendRequest` is created to avoid code duplication. It will set common headers (content type, auth header), make request, check for errors, parse response.

Note that we're considering status codes < 200 and >= 400 as errors and parse response into `errorResponse`. It depends on the API design though, your API may handle errors differently.

## Tests

So now we have SDK with single API endpoint covered, which is enough for this example, but is it enough to be able to ship this to users? Probably yes, but let's focus on few more things.

Tests are almost required here, and there can be 2 types of them: unit tests and integration tests. For the second one we'll call real API. Let's write a simple test.


```go
// +build integration

package facest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFaces(t *testing.T) {
	c := NewClient(os.Getenv("FACEST_INTEGRATION_API_KEY"))

	ctx := context.Background()
	res, err := c.GetFaces(nil)

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, res, "expecting non-nil result")

	assert.Equal(t, 1, res.Count, "expecting 1 face found")
	assert.Equal(t, 1, res.PagesCount, "expecting 1 PAGE found")

	assert.Equal(t, "integration_face_id", res.Faces[0].FaceID, "expecting correct face_id")
	assert.NotEmpty(t, res.Faces[0].FaceToken, "expecting non-empty face_token")
	assert.Greater(t, len(res.Faces[0].FaceImages), 0, "expecting non-empty face_images")
}
```

Note that this test uses env. var where API Key is set. By doing this, we're making sure that they are not public. And later we can configure our build system to propagate this env. var using secrets.

Also, these tests are separated from unit tests (because they take longer to execute):

```shell
go test -v -tags=integration
```

### Documentation.

Make your SDK self-explanatory with clear types and abstractions, don't expose too much information. Usually, it's enough to provide `godoc` link as main documentation.

### Compatibility and Versioning.

Version your SDK updates by publishing new semver to your repository. But make sure you're not breaking anything with new minor/patch releases. Usually your SDK library should follow API updates, so if API releases v2, then there should be an SDK v2 release as well.

## Conclusion

That's it.

One question though: what are the best API Go clients have you seen so far? Please share them in the comments.

You can find the full source code [here](https://github.com/facest/facest-go).