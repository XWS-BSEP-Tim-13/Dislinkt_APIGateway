package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	postGw "github.com/XWS-BSEP-Tim-13/Dislinkt_PostService/infrastructure/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"io/ioutil"
	"net/http"
)

type UploadImageHandler struct {
	postsClientAddress string
}

func (handler *UploadImageHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/post/uploadImage", handler.UploadImage)
	if err != nil {
		panic(err)
	}
}

func NewUploadImageHandler(postsClientAddress string) Handler {
	return &UploadImageHandler{
		postsClientAddress: postsClientAddress,
	}
}

func (handler *UploadImageHandler) UploadImage(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	fmt.Println("File Upload Endpoint Hit")
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	postsClient := services.NewPostsClient(handler.postsClientAddress)
	resp, _ := postsClient.UploadImage(context.TODO(), &postGw.ImageRequest{Image: fileBytes})
	response, _ := json.Marshal(resp.ImagePath)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
