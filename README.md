# stigg-sidecar-sdk

Stigg Go SDK makes it easier to interact with Stigg Sidecar

## Documentation

See https://docs.stigg.io/docs/sidecar-sdk

## Installation

```shell
    go get github.com/stiggio/sidecar-sdk-go/v5
```

## Usage

Initialize the client:

```go

package xyz

import (
  "github.com/stiggio/sidecar-sdk-go/v5"
)

host := "localhost"
port := 80
client, err := NewSidecarClient(ApiClientConfig{ApiKey: "<SERVER-API-KEY>"}, &host, &port, RemoteSidecarConfig{})

// With legacy TLS enabled
// client, err := NewSidecarClient(ApiClientConfig{ApiKey: "<SERVER-API-KEY>"}, &host, &port, RemoteSidecarConfig{UseLegacyTls: true})

```

Get single entitlement of a customer

```go

package xyz

import (
  "context"
  "fmt"
  "github.com/stiggio/sidecar-sdk-go/v5"
  sidecarv1 "github.com/stiggio/sidecar-sdk-go/v5/generated/stigg/sidecar/v1"
)

func main() {
  host := "localhost"
  port := 80
  client, err := NewSidecarClient(ApiClientConfig{apiKey: "<SERVER-API-KEY>"}, &host, &port)

  req := sidecarv1.GetBooleanEntitlementRequest{
    CustomerId: "customer-demo-01",
    ResourceId: nil,
  }
  
  resp, err := client.GetBooleanEntitlement(context.Background(), &req)
  if err != nil {
    fmt.Printf("Failed to get boolean entitlement: %v", err)
    return
  }

  fmt.Printf("Has access: %v", resp.HasAccess)
}

```

Accessing the `api` client:

```go

package xyz

import (
  "context"
  "fmt"
  "github.com/stiggio/sidecar-sdk-go/v5"
  sidecarv1 "github.com/stiggio/sidecar-sdk-go/v5/generated/stigg/sidecar/v1"
)

func main() {
  host := "localhost"
  port := 80
  client, err := NewSidecarClient(ApiClientConfig{apiKey: "<SERVER-API-KEY>"}, &host, &port)
  
  customerId := "test-customer-6923842"

  result, err := client.Api().ProvisionCustomer(context.Background(), stigg.ProvisionCustomerInput{
    AdditionalMetaData:       nil,
    AwsMarketplaceCustomerID: nil,
    BillingID:                nil,
    BillingInformation:       nil,
    CouponRefID:              nil,
    CreatedAt:                nil,
    CrmID:                    nil,
    CustomerID:               nil,
    Email:                    nil,
    EnvironmentID:            nil,
    ExcludeFromExperiment:    nil,
    Name:                     nil,
    RefID:                    &customerId,
    SalesforceID:             nil,
    ShouldSyncFree:           nil,
    SubscriptionParams:       nil,
  })

  if err != nil {
    fmt.Printf("Provision failed: %v", err)
    panic(err)
  }

  fmt.Printf("Customer provisioned: %v", result)
}

```

### License

See the [LICENSE](LICENSE) file for details
