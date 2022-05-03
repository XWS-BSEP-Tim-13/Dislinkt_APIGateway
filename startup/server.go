package startup

import (
	"context"
	//"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	cfg "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/startup/config"
	companyGw "github.com/XWS-BSEP-Tim-13/Dislinkt_CompanyService/infrastructure/grpc/proto"
	//inventoryGw "github.com/tamararankovic/microservices_demo/common/proto/inventory_service"
	//orderingGw "github.com/tamararankovic/microservices_demo/common/proto/ordering_service"
	//shippingGw "github.com/tamararankovic/microservices_demo/common/proto/shipping_service"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux
}

func NewServer(config *cfg.Config) *Server {
	server := &Server{
		config: config,
		mux:    runtime.NewServeMux(),
	}
	server.initHandlers()
	server.initCustomHandlers()
	return server
}

func (server *Server) initHandlers() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	companyEndpoint := fmt.Sprintf("%s:%s", server.config.CompanyHost, server.config.CompanyPort)
	err := companyGw.RegisterCompanyServiceHandlerFromEndpoint(context.TODO(), server.mux, companyEndpoint, opts)
	if err != nil {
		panic(err)
	}
	//orderingEmdpoint := fmt.Sprintf("%s:%s", server.config.OrderingHost, server.config.OrderingPort)
	//err = orderingGw.RegisterOrderingServiceHandlerFromEndpoint(context.TODO(), server.mux, orderingEmdpoint, opts)
	//if err != nil {
	//	panic(err)
	//}
	//shippingEmdpoint := fmt.Sprintf("%s:%s", server.config.CompanyHost, server.config.CompanyPort)
	//err = shippingGw.RegisterShippingServiceHandlerFromEndpoint(context.TODO(), server.mux, shippingEmdpoint, opts)
	//if err != nil {
	//	panic(err)
	//}
	//inventoryEmdpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	//err = inventoryGw.RegisterInventoryServiceHandlerFromEndpoint(context.TODO(), server.mux, inventoryEmdpoint, opts)
	//if err != nil {
	//	panic(err)
	//}
}

func (server *Server) initCustomHandlers() {

}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.mux))
}
