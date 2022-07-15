package api

import (
	"encoding/json"
	dom "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/domain"
	events "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/create_post"
	pb "github.com/XWS-BSEP-Tim-13/Dislinkt_CompanyService/infrastructure/grpc/proto"
	"io"
	"net/http"
)

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func decodePostDtoBody(r io.Reader) (*dom.PostDto, error) {

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt dom.PostDto
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeCreatePostBody(r io.Reader) (*events.PostFront, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt events.PostFront
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeHomepageFeedBody(r io.Reader) (*dom.HomepageFeedDto, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt dom.HomepageFeedDto
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeJobOfferDtoBody(r io.Reader) (*dom.JobOfferTokenDto, error) {
	dec := json.NewDecoder(r)
	//dec.DisallowUnknownFields()
	var rt dom.JobOfferTokenDto
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeMessageDtoBody(r io.Reader) (*dom.MessageDto, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt dom.MessageDto
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func mapJobDomainToPb(job dom.JobOfferDto) *pb.JobOfferDto {
	dto := &pb.JobOfferDto{
		JobDescription: job.JobDescription,
		Position:       job.Position,
		Prerequisites:  job.Prerequisites,
		EmploymentType: pb.EmploymentType(job.EmploymentType),
		Company: &pb.Company{
			Id:          job.Company.Id.Hex(),
			CompanyName: job.Company.CompanyName,
			Username:    job.Company.Username,
			Description: job.Company.Description,
			Location:    job.Company.Location,
			Website:     job.Company.Website,
			CompanySize: job.Company.CompanySize,
			Industry:    job.Company.Industry,
		},
	}
	return dto
}
