package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testApplication/models"
	"testApplication/repositories/postgres"
	"testApplication/utils"
	"testing"
)

func TestClientHandler_GetClientById(t *testing.T) {

	utils.LoadConf()
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(recorder)

	pg := postgres.InitConnectionNoMigration()
	handler, err := NewClientHandler(pg)
	if err != nil {
		return
	}

	ctx.Request = &http.Request{
		URL:    &url.URL{},
		Method: http.MethodGet,
	}

	want, err := json.MarshalIndent(models.Client{
		Id:   1,
		Name: "New name",
	}, "", "    ")
	if err != nil {
		return
	}

	tests := []struct {
		params     gin.Params
		wantStatus int
		want       string
	}{
		{
			[]gin.Param{
				{
					Key:   "id",
					Value: "1",
				},
			},
			http.StatusOK,
			string(want),
		},
	}

	for _, tt := range tests {
		ctx.Params = tt.params

		handler.GetClientById(ctx)

		assert.Equal(t, recorder.Code, tt.wantStatus)
		got := recorder.Body.String()

		assert.Equal(t, got, tt.want)
		ctx.Params = nil
	}
}
