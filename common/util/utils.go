package util

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

// processBody processBody
func ProcessBody(req *gin.Context) ([]byte, error) {
	reqBody, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("Internal JSON error when reading request body:%s", err)
	}
	defer req.Request.Body.Close()

	// Incoming request body may not be empty, client must supply request payload
	if string(reqBody) == "" {
		return nil, errors.New("Client must supply a payload for order requests")
	}
	//logger.Infof("Req body: %s", string(reqBody))

	return reqBody, nil
}
