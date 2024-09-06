package stigg_sidecar

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"errors"
	"fmt"
	"github.com/stiggio/api-client-go"
	"github.com/stiggio/sidecar-sdk-go/generated/stigg/sidecar/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/fs"
	"net/http"
)

//go:embed certs/root-ca.pem
var certFile embed.FS

type ApiClientConfig struct {
	apiKey     string
	httpClient *http.Client
	baseUrl    *string
}

type SidecarClient interface {
	sidecarv1.SidecarServiceClient
	Api() stigg.StiggClient
	Close() error
}

type sidecarClient struct {
	sidecarv1.SidecarServiceClient
	conn *grpc.ClientConn
	api  stigg.StiggClient
}

func NewSidecarClient(apiConfig ApiClientConfig, remoteSidecarHost *string, remoteSidecarPort *int) (SidecarClient, error) {
	serverAddress := getServerAddress(remoteSidecarHost, remoteSidecarPort)

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
	transportCredentials := credentials.NewTLS(tlsConfig)

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		fmt.Printf("Error creating Sidecar client: %v\n", err)
		return nil, err
	}

	client := sidecarv1.NewSidecarServiceClient(conn)

	api := stigg.NewStiggClient(apiConfig.apiKey, apiConfig.httpClient, apiConfig.baseUrl)

	sidecarClient := &sidecarClient{
		client,
		conn,
		api,
	}

	return sidecarClient, nil
}

func (c *sidecarClient) Api() stigg.StiggClient {
	return c.api
}

func (c *sidecarClient) Close() error {
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
		port = 8443
	}

	return fmt.Sprintf("%s:%d", host, port)
}
