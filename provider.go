package windns

import (
	"context"
	"fmt"
	"strings"

	"github.com/libdns/libdns"
	"golang.org/x/crypto/ssh"
)

// Provider implements the libdns interfaces for Windows DNS Server
type Provider struct {
	Host     string `json:"host,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	// Implementation for getting records
	// This is a placeholder, you'll need to implement the actual retrieval logic
	return nil, fmt.Errorf("get records not implemented")
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var appendedRecords []libdns.Record

	for _, record := range records {
		err := p.addRecord(zone, record)
		if err != nil {
			return appendedRecords, err
		}
		appendedRecords = append(appendedRecords, record)
	}

	return appendedRecords, nil
}

// DeleteRecords deletes the records from the zone. It returns the records that were deleted.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	// Implementation for deleting records
	// This is a placeholder, you'll need to implement the actual deletion logic
	return nil, fmt.Errorf("delete records not implemented")
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
// It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	// Implementation for setting records
	// This is a placeholder, you'll need to implement the actual setting logic
	return nil, fmt.Errorf("set records not implemented")
}

func (p *Provider) addRecord(zone string, record libdns.Record) error {
	config := &ssh.ClientConfig{
		User: p.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(p.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", p.Host+":22", config)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	var cmd string
	switch record.Type {
	case "CNAME":
		cmd = fmt.Sprintf("Add-DnsServerResourceRecordCName -ZoneName %s -Name %s -HostNameAlias %s", zone, record.Name, record.Value)
	// Add cases for other record types as needed
	default:
		return fmt.Errorf("unsupported record type: %s", record.Type)
	}

	fullCmd := fmt.Sprintf("powershell -Command \"%s\"", cmd)

	output, err := session.CombinedOutput(fullCmd)
	if err != nil {
		return fmt.Errorf("failed to run command: %v, output: %s", err, string(output))
	}

	if strings.Contains(string(output), "Error") {
		return fmt.Errorf("DNS record addition failed: %s", string(output))
	}

	return nil
}
