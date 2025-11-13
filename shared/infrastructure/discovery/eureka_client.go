package discovery

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/hudl/fargo"
)

type EurekaConfig struct {
	ServiceName     string
	ServerIP        string
	Port            string
	DiscoveryURL    string
	HealthCheckURL  string
	StatusPageURL   string
	HomePageURL     string
	RenewalInterval time.Duration
	DurationInSecs  int
}

type EurekaClient struct {
	connection *fargo.EurekaConnection
	instance   *fargo.Instance
	config     EurekaConfig
}

// NewEurekaClient creates and configures a new Eureka client
func NewEurekaClient(config EurekaConfig) (*EurekaClient, error) {
	// Create Eureka connection
	connection := fargo.NewConn(config.DiscoveryURL)
	connection.PollInterval = 30 * time.Second

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Parse port
	port, err := strconv.Atoi(config.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %v", err)
	}

	// Get IP address
	ipAddr := config.ServerIP
	if ipAddr == "" || ipAddr == "127.0.0.1" {
		ipAddr = getOutboundIP()
	}

	// Create instance
	instance := &fargo.Instance{
		InstanceId:        fmt.Sprintf("%s:%s:%s", ipAddr, config.ServiceName, config.Port),
		HostName:          hostname,
		App:               config.ServiceName,
		IPAddr:            ipAddr,
		VipAddress:        config.ServiceName,
		SecureVipAddress:  config.ServiceName,
		Status:            fargo.UP,
		Port:              port,
		PortEnabled:       true,
		SecurePort:        443,
		SecurePortEnabled: false,
		HomePageUrl:       config.HomePageURL,
		StatusPageUrl:     config.StatusPageURL,
		HealthCheckUrl:    config.HealthCheckURL,
		DataCenterInfo: fargo.DataCenterInfo{
			Name:  fargo.MyOwn,
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
		},
		LeaseInfo: fargo.LeaseInfo{
			RenewalIntervalInSecs: int32(config.RenewalInterval.Seconds()),
			DurationInSecs:        int32(config.DurationInSecs),
		},
		Metadata: fargo.InstanceMetadata{
			Raw: []byte(`{"management.port":"` + strconv.Itoa(port) + `"}`),
		},
	}

	return &EurekaClient{
		connection: &connection,
		instance:   instance,
		config:     config,
	}, nil
}

// Register registers the service instance with Eureka
func (ec *EurekaClient) Register() error {
	log.Printf("Registering service %s with Eureka at %s", ec.config.ServiceName, ec.config.DiscoveryURL)

	err := ec.connection.RegisterInstance(ec.instance)
	if err != nil {
		return fmt.Errorf("failed to register with Eureka: %v", err)
	}

	log.Printf("Successfully registered service %s (IP: %s, Port: %s)",
		ec.config.ServiceName, ec.instance.IPAddr, ec.config.Port)

	return nil
}

// SendHeartbeat sends a heartbeat to Eureka to keep the registration alive
func (ec *EurekaClient) SendHeartbeat() error {
	err := ec.connection.HeartBeatInstance(ec.instance)
	if err != nil {
		log.Printf("Failed to send heartbeat to Eureka: %v", err)
		return err
	}
	return nil
}

// StartHeartbeat starts sending periodic heartbeats to Eureka
func (ec *EurekaClient) StartHeartbeat() {
	ticker := time.NewTicker(ec.config.RenewalInterval)
	go func() {
		for range ticker.C {
			if err := ec.SendHeartbeat(); err != nil {
				log.Printf("Heartbeat error: %v", err)
				// Try to re-register if heartbeat fails
				if err := ec.Register(); err != nil {
					log.Printf("Re-registration failed: %v", err)
				}
			}
		}
	}()
}

// Deregister removes the service instance from Eureka
func (ec *EurekaClient) Deregister() error {
	log.Printf("Deregistering service %s from Eureka", ec.config.ServiceName)

	err := ec.connection.DeregisterInstance(ec.instance)
	if err != nil {
		return fmt.Errorf("failed to deregister from Eureka: %v", err)
	}

	log.Printf("Successfully deregistered service %s", ec.config.ServiceName)
	return nil
}

// UpdateStatus updates the instance status in Eureka
func (ec *EurekaClient) UpdateStatus(status fargo.StatusType) error {
	ec.instance.Status = status
	return ec.connection.UpdateInstanceStatus(ec.instance, status)
}

// getOutboundIP gets the preferred outbound IP address of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
