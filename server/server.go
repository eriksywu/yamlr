package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gitlab.com/erikwu09/yamlr/app"
	"gitlab.com/erikwu09/yamlr/models"
	"gopkg.in/yaml.v2"
)

type yamlMetadataService struct {
	router  *mux.Router
	port    int
	manager *app.MetadataManager
	logger  *log.Logger
}

func (s *yamlMetadataService) Run() {
	p := strconv.Itoa(s.port)
	s.logger.Println("Starting service...")
	s.logger.Fatal(http.ListenAndServe(":"+p, s.router))
}

// singleton instance of service
var service *yamlMetadataService

// Instantiates the YAML Metadata Service
func BuildYamlApp(port int, manager *app.MetadataManager, logger *log.Logger) *yamlMetadataService {
	if service == nil {
		service = &yamlMetadataService{port: port}
		service.router = buildRouter(service, logger)
		service.manager = manager
		service.logger = logger
	}
	return service
}

func buildRouter(app *yamlMetadataService, logger *log.Logger) *mux.Router {
	logger.Println("registering handlers routes")
	router := mux.NewRouter()
	router.HandleFunc("/api/metadata", app.createMetadataHandler).Methods("POST")
	router.HandleFunc("/api/search", app.searchMetadataHandler).Methods("POST")
	router.HandleFunc("/api/metadata/{guid}", app.updateMetadataHandler).Methods("PUT")
	router.HandleFunc("/api/metadata/{guid}", app.getMetadataHandler).Methods("GET")
	return router
}

func (s *yamlMetadataService) createMetadataHandler(res http.ResponseWriter, req *http.Request) {
	s.logger.Println("received request for Create Metadata")
	metadata, err := s.getPayload(req)
	if err != nil {
		s.writeErrorResponse(res, err)
		return
	}
	id, err := s.manager.CreateMetadata(*metadata)
	if err != nil {
		s.writeErrorResponse(res, err)
		return
	}
	s.logger.Printf("new metadata created with id %s \n ", id.String())
	res.Write([]byte(id.String()))
}

func (s *yamlMetadataService) searchMetadataHandler(res http.ResponseWriter, req *http.Request) {
	s.logger.Println("received request for Search Metadata")
	metadata, err := s.getPayload(req)
	if err != nil {
		s.writeErrorResponse(res, err)
	}
	results, err := s.manager.SearchMetadata(*metadata)
	if err != nil {
		s.writeErrorResponse(res, err)
		return
	}
	responseBody, err := yaml.Marshal(&results)
	if err != nil {
		s.writeErrorResponse(res, err)
		return
	}
	res.Write(responseBody)
}

func (s *yamlMetadataService) updateMetadataHandler(res http.ResponseWriter, req *http.Request) {
	s.logger.Println("received request for Update Metadata")
	params := mux.Vars(req)
	id, err := uuid.Parse(params["guid"])
	if err != nil {
		s.writeErrorResponse(res, err)
		return
	}
	metadata, err := s.getPayload(req)
	if err != nil {
		s.writeErrorResponse(res, err)
		return
	}
	err = s.manager.UpdateMetadata(*metadata, id)
	if err != nil {
		s.writeErrorResponse(res, err)
		return
	}
}

func (s *yamlMetadataService) getMetadataHandler(res http.ResponseWriter, req *http.Request) {
	s.logger.Println("received request for Get Metadata")
	params := mux.Vars(req)
	id, err := uuid.Parse(params["guid"])
	if err != nil {
		s.writeErrorResponse(res, err)
		return
	}
	result, err := s.manager.GetMetadata(id)
	if err != nil {
		s.writeErrorResponse(res, err)
		return
	}
	responseBody, err := yaml.Marshal(result)
	if err != nil {
		s.writeErrorResponse(res, err)
		return
	}
	res.Write(responseBody)

}
func (s *yamlMetadataService) writeErrorResponse(res http.ResponseWriter, err error) {
	s.logger.Printf("encountered error: %s", err.Error())
	res.Header().Set("Content-Type", "application/json")
	var errorBody []byte
	switch err.(type) {
	case app.ValidationError:
		res.WriteHeader(http.StatusBadRequest)
		errorBody = formatValidationError(err.(app.ValidationError))
	case app.AggregatedValidationError:
		res.WriteHeader(http.StatusBadRequest)
		errorBody = formatAggregatedValidationError(err.(app.AggregatedValidationError))
	default:
		res.WriteHeader(http.StatusInternalServerError)
		body := struct {
			ErrorMessage string
		}{
			err.Error(),
		}
		errorBody, _ = json.Marshal(body)
	}
	res.Write([]byte(errorBody))
}

func formatAggregatedValidationError(err app.AggregatedValidationError) []byte {
	errorMessage := make([]interface{}, 0, len(err.Errors()))
	for _, e := range err.Errors() {
		body := struct {
			Path    string
			Reason  string
			Message string
		}{
			Path:    e.Path,
			Reason:  e.Reason,
			Message: e.Error(),
		}
		errorMessage = append(errorMessage, body)
	}
	errorBody, _ := json.Marshal(errorMessage)
	return errorBody
}

func formatValidationError(err app.ValidationError) []byte {
	body := struct {
		Path    string
		Reason  string
		Message string
	}{
		Path:    err.Path,
		Reason:  err.Reason,
		Message: err.Error(),
	}
	errorBody, _ := json.Marshal(body)
	return errorBody
}

func (s *yamlMetadataService) getPayload(req *http.Request) (*models.Metadata, error) {
	reqBody, err := ioutil.ReadAll(req.Body)
	s.logger.Printf("received payload: %s \n", reqBody)
	if err != nil {
		return nil, err
	}
	metadata := models.Metadata{}
	err = yaml.Unmarshal(reqBody, &metadata)
	if err != nil {
		return nil, app.ValidationError{Path: ".", Reason: err.Error()}
	}
	return &metadata, nil
}
