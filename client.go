package stigg_sidecar

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"errors"
	"fmt"
	"github.com/stiggio/api-client-go/v3"
	"github.com/stiggio/sidecar-sdk-go/v3/generated/stigg/sidecar/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io/fs"
	"net/http"
)

//go:embed certs/root-ca.pem
var certFile embed.FS

type ApiClientConfig struct {
	ApiKey     string
	HttpClient *http.Client
	BaseUrl    *string
}

type RemoteSidecarConfig struct {
	UseLegacyTls bool
}

type SidecarClient struct {
	sidecarv1.SidecarServiceClient
	conn *grpc.ClientConn
	Api  stigg.StiggClient
}

func NewSidecarClient(apiConfig ApiClientConfig, remoteSidecarHost *string, remoteSidecarPort *int, remoteSidecarOptions RemoteSidecarConfig) (*SidecarClient, error) {
	serverAddress := getServerAddress(remoteSidecarHost, remoteSidecarPort)

	var transportCredentials credentials.TransportCredentials
	if remoteSidecarOptions.UseLegacyTls {
		rootPem, err := fs.ReadFile(certFile, "certs/root-ca.pem")
		if err != nil {
			return nil, err
		}

		root := x509.NewCertPool()
		if !root.AppendCertsFromPEM(rootPem) {
			return nil, errors.New("failed to parse root certificate")
		}

		tlsConfig := &tls.Config{
			RootCAs: root,
		}
		transportCredentials = credentials.NewTLS(tlsConfig)
	} else {
		transportCredentials = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		fmt.Printf("Error creating Sidecar client: %v\n", err)
		return nil, err
	}

	client := sidecarv1.NewSidecarServiceClient(conn)

	api := stigg.NewStiggClient(apiConfig.ApiKey, apiConfig.HttpClient, apiConfig.BaseUrl)

	sidecarClient := &SidecarClient{
		client,
		conn,
		api,
	}

	return sidecarClient, nil
}

func (c *SidecarClient) Close() error {
	return c.conn.Close()
}

func getServerAddress(remoteSidecarHost *string, remoteSidecarPort *int) string {
	var host string
	if remoteSidecarHost != nil {
		host = *remoteSidecarHost
	} else {
		host = "localhost"
	}

	var port int
	if remoteSidecarPort != nil {
		port = *remoteSidecarPort
	} else {
		port = 80
	}

	return fmt.Sprintf("%s:%d", host, port)
}
