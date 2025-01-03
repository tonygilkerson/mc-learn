package main

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"

	"github.com/gin-gonic/gin"

	api "basic/pkg/apiv1"


)



//DEVTODO - See this example: https://github.com/GoogleCloudPlatform/metacontroller/blob/v0.4.0/examples/go/main.go

///////////////////////////////////////////////////////////////////////////////
//	main
///////////////////////////////////////////////////////////////////////////////
func main() {
	r := gin.Default()
	r.POST("/sync", syncHandler)
	r.Run(fmt.Sprintf(":%s", "8080"))
}

///////////////////////////////////////////////////////////////////////////////
//	functions
///////////////////////////////////////////////////////////////////////////////

// Handle sync
func syncHandler(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		println("JSON could not be retrieved")
		c.String(http.StatusBadRequest, "JSON could not be retrieved")
		return
	}

	request := &api.SyncRequest{}
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
