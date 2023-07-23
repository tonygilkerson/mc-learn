package main

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"

	"github.com/gin-gonic/gin"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)


type Controller struct {
	ApiVersion      string `json:"apiVersion"`
	Kind            string `json:"kind"`
	meta.ObjectMeta `json:"metadata"`
}

type SyncRequest struct {
	Parent Controller `json:"parent"`
}


///////////////////////////////////////////////////////////////////////////////
//	main
///////////////////////////////////////////////////////////////////////////////
func main() {
	r := gin.Default()
	r.POST("/basic", basicHandler)
	r.Run(fmt.Sprintf(":%s", "8080"))
}

///////////////////////////////////////////////////////////////////////////////
//	functions
///////////////////////////////////////////////////////////////////////////////

// Handle basic
func basicHandler(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		println("JSON could not be retrieved")
		c.String(http.StatusBadRequest, "JSON could not be retrieved")
		return
	}

	request := &SyncRequest{}
	// Unmarshal the JSON into the struct.
	err = json.Unmarshal(body, request)
	if err != nil {
		fmt.Printf("JSON could not be unmarshalled: %s\n%s\n",err,string(body[:]))
		c.String(http.StatusBadRequest, "JSON could not be unmarshalled.")
		return
	}

	managedField := request.Parent.ObjectMeta.ManagedFields[0]
	fmt.Printf("API Version: [%s], Kind: [%s], Operation: [%s]\n", request.Parent.ApiVersion, request.Parent.Kind,managedField.Operation)


	// Return success
	c.String(http.StatusOK, "Success")
}
