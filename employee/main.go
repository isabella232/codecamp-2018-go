package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	company_pb "github.com/namely/codecamp-2018-go/employee/gen/protos/company"
	employee_pb "github.com/namely/codecamp-2018-go/employee/gen/protos/employee"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	port        = flag.Int("port", 50051, "The server port")
	companyAddr = flag.String("company_addr", "company:50051", "The host:port of company service")
)

// EmployeeServer is a gRPC server for employee data.
type EmployeeServer struct {
	companies     map[string]*EmployeeCollection
	mutex         *sync.RWMutex
	conn          *grpc.ClientConn
	companyClient company_pb.CompanyServiceClient
}

func newServer() *EmployeeServer {
	s := &EmployeeServer{}
	s.companies = make(map[string]*EmployeeCollection)
	s.mutex = &sync.RWMutex{}

	// For our example, we disable TLS.
	var err error
	s.conn, err = grpc.Dial(*companyAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	s.companyClient = company_pb.NewCompanyServiceClient(s.conn)

	return s
}

// EmployeeCollection is a collection of employees at a company.
type EmployeeCollection struct {
	employees map[string]*employee_pb.Employee
	mutex     *sync.RWMutex
}

func newEmployeeCollection() *EmployeeCollection {
	c := &EmployeeCollection{}
	c.employees = make(map[string]*employee_pb.Employee)
	c.mutex = &sync.RWMutex{}
	return c
}

// AddEmployee assigns a uuid to the employee and adds it to the collection.
func (ec *EmployeeCollection) AddEmployee(employee *employee_pb.Employee) *employee_pb.Employee {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()
	// Copy employee
	ee := &employee_pb.Employee{}
	*ee = *employee
	ee.EmployeeUuid = uuid.Must(uuid.NewV4()).String()
	ec.employees[ee.EmployeeUuid] = ee
	return ee
}

// GetEmployeeCollection gets the employee collection for the given company uuid. It returns false
// if none is found
func (s *EmployeeServer) GetEmployeeCollection(companyUUID string) (*EmployeeCollection, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	c, ok := s.companies[companyUUID]
	return c, ok
}

// GetOrCreateEmployeeCollection looks up the employee collection, or creates one if it doesn't exist
func (s *EmployeeServer) GetOrCreateEmployeeCollection(companyUUID string) *EmployeeCollection {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	ec, ok := s.companies[companyUUID]
	if !ok {
		ec = newEmployeeCollection()
		s.companies[companyUUID] = ec
	}
	return ec
}

// SaveEmployee saves the employee
func (s *EmployeeServer) SaveEmployee(employee *employee_pb.Employee) *employee_pb.Employee {
	ec := s.GetOrCreateEmployeeCollection(employee.CompanyUuid)
	return ec.AddEmployee(employee)
}

// CreateEmployee creates a new employee in the given company.
func (s *EmployeeServer) CreateEmployee(ctx context.Context, req *employee_pb.CreateEmployeeRequest) (*employee_pb.Employee, error) {
	if req.Employee == nil {
		req.Employee = &employee_pb.Employee{}
	}
	// The server should assign the UUID
	if req.Employee.EmployeeUuid != "" {
		return nil, status.Error(codes.InvalidArgument, "don't set EmployeeUuid for create, it will be assigned")
	}
	// The employee must have a name.
	if req.Employee.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "employee must have a name assigned")
	}
	// Check that the Employee's company exists by calling company service.
	_, err := s.companyClient.GetCompany(ctx, &company_pb.GetCompanyRequest{
		CompanyUuid: req.Employee.CompanyUuid,
	})
	if err != nil {
		log.Printf("company %s does not exist: %v", req.Employee.CompanyUuid, err)
		return nil, status.Error(codes.InvalidArgument, "company does not exist")
	}

	// If we're here, we can create the employee
	return s.SaveEmployee(req.Employee), nil
}

// ListEmployees returns all emplyees in the given company.
func (s *EmployeeServer) ListEmployees(ctx context.Context, req *employee_pb.ListEmployeesRequest) (*employee_pb.ListEmployeesResponse, error) {
	ec, ok := s.GetEmployeeCollection(req.CompanyUuid)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "could not find company")
	}

	// Create our response
	resp := &employee_pb.ListEmployeesResponse{}
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	// Copy the employees in the collection to the response
	resp.Employees = make([]*employee_pb.Employee, 0)
	for _, v := range ec.employees {
		resp.Employees = append(resp.Employees, v)
	}

	return resp, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("error listening: %v", err)
	}
	server := grpc.NewServer()
	employeeServer := newServer()
	defer employeeServer.conn.Close()
	employee_pb.RegisterEmployeeServiceServer(server, employeeServer)
	server.Serve(lis)
}
