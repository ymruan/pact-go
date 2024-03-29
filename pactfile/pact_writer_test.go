package pactfile

import (
	"github.com/bennycao/pact-go/consumer"
	"github.com/bennycao/pact-go/provider"
	"net/http"
	"testing"
)

func Test_InvalidPath_ShouldThrowError(t *testing.T) {
	pact := NewPactFile("consumer", "provider", nil)
	writer := NewPactFileWriter(pact, "//g34/example")

	if err := writer.Write(); err == nil {
		t.Error("expected error as file path is invalid")
	}
}

func Test_ValidPact_ShouldWritePactFile(t *testing.T) {
	var interactions []*consumer.Interaction
	interactions = append(interactions, getFakeInteraction())

	pact := NewPactFile("consumer", "provider", interactions)
	writer := NewPactFileWriter(pact, "./example")

	if err := writer.Write(); err != nil {
		t.Error(err)
	}

}

func getFakeInteraction() *consumer.Interaction {
	header := make(http.Header)
	header.Add("content-type", "application/json")
	i := consumer.NewInteraction("description of the interaction",
		"some state",
		provider.NewProviderRequest("POST", "/", "param=xyzmk", header),
		provider.NewProviderResponse(201, header))

	i.Request.SetBody(`{ "firstName": "John", "lastName": "Doe" }`)
	i.Response.SetBody(`{"result": true}`)
	return i
}
