package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/infrastructure/services"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/jwt"
	authGw "github.com/XWS-BSEP-Tim-13/Dislinkt_AuthenticationService/infrastructure/grpc/proto"
	companyGw "github.com/XWS-BSEP-Tim-13/Dislinkt_CompanyService/infrastructure/grpc/proto"
	"github.com/go-stomp/stomp"
	"log"
)

type ReceiveJobOffer struct {
	authClientAddress    string
	companyClientAddress string
}

func NewReceiveJobOfferHandler(authClientAddress, companyClientAddress string) *ReceiveJobOffer {
	return &ReceiveJobOffer{
		authClientAddress:    authClientAddress,
		companyClientAddress: companyClientAddress,
	}
}

func (handler *ReceiveJobOffer) Connect() (*stomp.Conn, error) {
	return stomp.Dial("tcp", "activemq:61613")
}

func (handler *ReceiveJobOffer) HandleReceive() {
	conn, err := handler.Connect()
	fmt.Printf("Connecting to queue\n")
	if err != nil {
		fmt.Printf("Error while connecting %s\n", err)
		panic(err)
	}
	sub, err := conn.Subscribe("jobOffer.queue", stomp.AckAuto)
	if err != nil {
		fmt.Printf("Error while subscribing %s\n", err)
	}
	defer conn.Disconnect()
	defer sub.Unsubscribe()
	for {
		fmt.Printf("#")
		m := <-sub.C
		log.Println("Message : ", m.Body)
		//handler.DecodeBody(m.Err, string(m.Body))
	}
}

func (handler *ReceiveJobOffer) DecodeBody(err error, body string) error {
	fmt.Printf("Decoding body %s\n", body)
	if err != nil {
		return err
	}
	var jobTokenDto domain.JobOfferTokenDto
	json.Unmarshal([]byte(body), &jobTokenDto)
	fmt.Printf("Species: %s, Description: %s", jobTokenDto.JobOffer.Position, jobTokenDto.Token)
	_, claims, err := jwt.ParseJwtWithEmail(jobTokenDto.Token)
	if err != nil {
		return err
	}
	authClient := services.NewAuthClient(handler.authClientAddress)
	resp, err := authClient.CheckIfUserExist(context.TODO(), &authGw.CheckIfUserExistsRequest{Username: claims.Email})
	if !resp.Exists {
		return err
	}
	fmt.Printf("Received first response \n")
	companyClient := services.NewCompanyClient(handler.companyClientAddress)
	pb := mapJobDomainToPb(jobTokenDto.JobOffer)
	_, err = companyClient.CreateJobOffer(context.TODO(), &companyGw.JobOfferRequest{Dto: pb})
	if err != nil {
		return err
	}
	return nil
}
