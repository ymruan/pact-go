package consumer

import (
	"bytes"
	"github.com/bennycao/pact-go/provider"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func Test_MatchingInteractionFound_ReturnsCorrectResponse(t *testing.T) {
	mockHttpServer := NewHttpMockService()
	interaction := getFakeInteraction()
	mockHttpServer.ClearInteractions()
	mockHttpServer.RegisterInteraction(interaction)
	url := mockHttpServer.Start()
	defer mockHttpServer.Stop()

	client := &http.Client{}

	req, err := interaction.ToHttpRequest(url)
	resp, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != interaction.Response.Status {
		t.Errorf("The response status is %s, expected %v", resp.Status, interaction.Response.Status)
		t.FailNow()
	}

	defer resp.Body.Close()

	if expectedBody, err := interaction.Response.GetBody(); err != nil {
		t.Error(err)
	} else {
		if actualBody, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Error(err)
		} else {
			if bytes.Compare(expectedBody, actualBody) != 0 {
				t.Error("The response body does not match")
			}
		}
	}

}

func Test_MatchingInteractionNotFound_Returns404(t *testing.T) {
	mockHttpServer := NewHttpMockService()
	interaction := getFakeInteraction()

	url := mockHttpServer.Start()
	defer mockHttpServer.Stop()

	client := &http.Client{}

	req, err := interaction.ToHttpRequest(url)
	resp, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("The response status is %s, expected %v", resp.Status, 404)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bodyText := strings.TrimSpace(string(body)) //had trim, not sure why there are trailing spaces

	if err != nil {
		t.Error(err)
	}

	if bodyText != notFoundError.Error() {
		t.Errorf("The expected response was '%s' but recieved '%s'", notFoundError.Error(), bodyText)
	}
}

func getFakeInteraction() *Interaction {
	header := make(http.Header)
	header.Add("content-type", "application/json")
	i := NewInteraction("description of the interaction",
		"some state",
		provider.NewProviderRequest("GET", "/", "param=xyzmk", header),
		provider.NewProviderResponse(201, header))
	i.Request.SetBody(`{ "firstName": "John", "lastName": "Doe" }`)

	return i
}
