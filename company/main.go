package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/namely/codecamp-2018-go/company/gen/protos/company"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// CompanyServer is a gRPC server for company data.
type CompanyServer struct {
	companies map[string]*pb.Company
	mutex     *sync.RWMutex
}

func newServer() *CompanyServer {
	s := &CompanyServer{}
	s.companies = make(map[string]*pb.Company)
	s.mutex = &sync.RWMutex{}
	return s
}

// CreateCompany creates a new company if none exists.
func (s *CompanyServer) CreateCompany(ctx context.Context, req *pb.CreateCompanyRequest) (*pb.Company, error) {
	// To make the code easier to write, set an empty company if none exists.
	if req.Company == nil {
		req.Company = &pb.Company{}
	}
	// For now, only validation is that caller doesn't set company uuid
	if req.Company.CompanyUuid != "" {
		return nil, status.Error(codes.InvalidArgument, "should not set company uuid, the server will assign one")
	}
	// Make our own copy of the Company struct.
	c := &pb.Company{
		CompanyUuid:     uuid.Must(uuid.NewV4()).String(),
		CeoEmployeeUuid: req.Company.CeoEmployeeUuid,
	}

	// Office location is optional and might be nil
	if req.Company.OfficeLocation != nil {
		c.OfficeLocation = &pb.Address{
			Address1: req.Company.OfficeLocation.Address1,
			Address2: req.Company.OfficeLocation.Address2,
			Zip:      req.Company.OfficeLocation.Zip,
			State:    req.Company.OfficeLocation.State,
		}
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.companies[c.CompanyUuid] = c
	return c, nil
}

// GetCompany returns the company if it exists
func (s *CompanyServer) GetCompany(ctx context.Context, req *pb.GetCompanyRequest) (*pb.Company, error) {
	// Look up the company in our map of companies.
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	c, ok := s.companies[req.CompanyUuid]
	if !ok {
		// If it doesn't exist, we return a not found error.
		return nil, status.Error(codes.NotFound, "could not find company")
	}
	// Otherwise, return the company to the client.
	return c, nil
}

// SetCompanyCeo sets the company's CEO.
func (s *CompanyServer) SetCompanyCeo(ctx context.Context, req *pb.SetCompanyCeoRequest) (*empty.Empty, error) {
	ceoUUID, err := uuid.FromString(req.CeoEmployeeUuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "could not parse ceo uuid")
	}
	// Look up the company in our map of companies.
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	c, ok := s.companies[req.CompanyUuid]
	if !ok {
		// If it doesn't exist, we return a not found error.
		return nil, status.Error(codes.NotFound, "could not find company")
	}
	c.CeoEmployeeUuid = ceoUUID.String()
	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("error listening: %v", err)
	}
	server := grpc.NewServer()

	pb.RegisterCompanyServiceServer(server, newServer())
	server.Serve(lis)
}
