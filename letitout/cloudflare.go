package letitout

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"net"
	"os"
	"strings"
)

type CloudFlare struct {
	Token string `yaml:"token"`
}

func backupRecords(api *cloudflare.API, zone string, hostname string) {
	backup := fmt.Sprintf("_letitout.%s", hostname)

	backupRecords, err := api.DNSRecords(zone, cloudflare.DNSRecord{
		Name: backup,
		Type: "TXT",
	})

	if err != nil {
		fmt.Println("Failed to fetch backup records:", err)
		os.Exit(1)
	}

	if backupRecords != nil {
		// Already backed up
		return
	}

	targetRecords, err := api.DNSRecords(zone, cloudflare.DNSRecord{
		Name: hostname,
	})

	if err != nil {
		fmt.Println("Failed to fetch target records:", err)
		os.Exit(1)
	}

	if targetRecords == nil {
		// Nothing to backup
		return
	}

	if len(targetRecords) > 1 {
		fmt.Printf("Multiple DNS entries found for %s, don't know how to back it up.", hostname)
		os.Exit(1)
	}

	// Backup entry
	_, err = api.CreateDNSRecord(zone, cloudflare.DNSRecord{
		Name: backup,
		Type: "TXT",
		Content: fmt.Sprintf("%s|%s|%d|%t", targetRecords[0].Type, targetRecords[0].Content, targetRecords[0].TTL, targetRecords[0].Proxied),
	})

	if err != nil {
		fmt.Println("Failed to backup DNS record:", err)
		os.Exit(1)
	}
}

func deleteRecord(api *cloudflare.API, record cloudflare.DNSRecord) {
	err := api.DeleteDNSRecord(record.ZoneID, record.ID)
	if err != nil {
		fmt.Printf("Failed to delete DNS record %s -> %s (%s).", record.Name, record.Content, record.Type)
		os.Exit(1)
	}
}

func updateRecord(api *cloudflare.API, record cloudflare.DNSRecord) {
	err := api.UpdateDNSRecord(record.ZoneID, record.ID, record)
	if err != nil {
		fmt.Printf("Failed to update DNS record %s -> %s (%s).", record.Name, record.Content, record.Type)
		os.Exit(1)
	}
}

func createRecord(api *cloudflare.API,  record cloudflare.DNSRecord) {
	_, err := api.CreateDNSRecord(record.ZoneID, record)
	if err != nil {
		fmt.Printf("Failed to create DNS record %s -> %s (%s).", record.Name, record.Content, record.Type)
		os.Exit(1)
	}
}

func updateRecords(api *cloudflare.API, zone string, hostname string, server *Server) {
	serverAddress := server.Address
	if strings.Index(server.Address, ":") >= 0 {
		parts := strings.Split(server.Address, ":")
		serverAddress = parts[0]
	}

	ip := net.ParseIP(serverAddress)
	if ip == nil {
		fmt.Println("Invalid server address found:", serverAddress)
		os.Exit(1)
	}

	recordType := "A"
	if ip.To4() == nil {
		recordType = "AAAA"
	}

	// Backup entry first
	backupRecords(api, zone, hostname)

	targetRecords, err := api.DNSRecords(zone, cloudflare.DNSRecord{
		Name: hostname,
	})

	if err != nil {
		fmt.Println("Failed to fetch target records:", err)
		os.Exit(1)
	}

	// Target doesn't exist, create it
	if targetRecords == nil {
		createRecord(api, cloudflare.DNSRecord{
			ZoneID:  zone,
			Type:    recordType,
			Name:    hostname,
			Content: serverAddress,
			Proxied: true,
		})
		return
	}

	if len(targetRecords) > 1 {
		fmt.Printf("Multiple DNS entries found for %s, don't know how to update it.", hostname)
		os.Exit(1)
	}

	targetRecord := targetRecords[0]

	if targetRecord.Type != recordType {
		deleteRecord(api, targetRecord)
		targetRecord.Type = recordType
		targetRecord.Content = serverAddress
		targetRecord.Proxied = true
		createRecord(api, targetRecord)
	} else {
		targetRecord.Content = serverAddress
		targetRecord.Proxied = true
		updateRecord(api, targetRecord)
	}
}

func updateCloudFlare(zone string, hostname string, credentials CloudFlare, server *Server) {
	api, err := cloudflare.NewWithAPIToken(credentials.Token)
	if err != nil {
		fmt.Println("Failed to instantiate CloudFlare API:", err)
		os.Exit(1)
	}

	id, err := api.ZoneIDByName(zone)
	if err != nil {
		fmt.Println("Failed to fetch zone id:", err)
		os.Exit(1)
	}

	updateRecords(api, id, hostname, server)
}

func UpdateDns(hostname string, server *Server) {
	parts := strings.Split(hostname, ".")
	for i := len(parts) - 1; i >= 0; i -= 1 {
		part := parts[i:len(parts)]
		zone := strings.Join(part, ".")

		credentials, ok := config.CloudFlare[zone]
		if ok == false {
			continue
		}

		updateCloudFlare(zone, hostname, credentials, server)
		return
	}

	fmt.Printf("No cloudflare entries found for hostname %s, skipping DNS update.\n", hostname)
}
