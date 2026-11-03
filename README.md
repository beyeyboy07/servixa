# SERVIXA Backend

REST API backend for **SERVIXA**, a fleet and field service management
system. The application manages customer service requests, work orders,
technician assignments, vehicles, schedules, service visits, work
completion, customer confirmations, and notifications.

> This README documents routes visible in `internal/router/router.go`.
> Request fields and response schemas must be confirmed against the
> corresponding handlers before being treated as a complete API
> contract.

## Table of Contents

-   [Features](#features)
-   [Technology Stack](#technology-stack)
-   [Roles](#roles)
-   [Project Structure](#project-structure)
-   [Getting Started](#getting-started)
-   [Running Tests](#running-tests)
-   [API Reference](#api-reference)
    -   [Health](#health)
    -   [Authentication](#authentication)
    -   [Customers](#customers)
    -   [Technicians](#technicians)
    -   [Vehicles](#vehicles)
    -   [Parts](#parts)
    -   [Service Requests](#service-requests)
    -   [Work Orders](#work-orders)
    -   [Service Visits](#service-visits)
    -   [Customer Confirmation](#customer-confirmation)
    -   [Notifications](#notifications)
-   [Authentication](#authentication-1)
-   [Business Workflow](#business-workflow)
-   [Automated Tests](#automated-tests)
-   [Documentation Notes](#documentation-notes)

## Features

-   JWT-based authentication and role-based access control.
-   Customer profile and service request management.
-   Service request review and rejection.
-   Work order management and cancellation.
-   Technician assignment and reassignment.
-   Vehicle and part management.
-   Work order scheduling.
-   Technician service visits: create, start, and complete.
-   Work completion records.
-   Customer confirmation of completed work.
-   User notifications, unread counts, and read-state management.
-   Automated tests for selected authorization, validation, concurrency,
    and notification scenarios.

## Technology Stack

The codebase uses the following technologies, based on the current
implementation:

-   Go
-   Gin
-   GORM
-   MySQL
-   JWT authentication
-   bcrypt password hashing

Check `go.mod` for the exact dependency versions.

## Roles

  -----------------------------------------------------------------------
  Role                                Access shown in the router
  ----------------------------------- -----------------------------------
  `ADMIN`                             Work order management, service
                                      request review/rejection, vehicle
                                      management, and part management

  `DISPATCHER`                        Work order management, service
                                      request review/rejection, vehicle
                                      management, and part management

  `TECHNICIAN`                        Technician profile, service visits,
                                      and work completion routes

  `CUSTOMER`                          Customer profile, service requests,
                                      and customer confirmation
  -----------------------------------------------------------------------

The table summarizes route-level role restrictions. Individual
service-layer rules may impose additional restrictions.

## Project Structure

``` text
backend/
├── internal/
│   ├── config/          # Application configuration
│   ├── database/        # Database connection
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # Authentication and role authorization
│   ├── model/           # Data models
│   ├── repository/      # Database access
│   ├── router/          # Route registration
│   └── service/         # Business logic
├── tests/               # Automated tests
├── go.mod
├── go.sum
└── main.go
```

## Getting Started

### Prerequisites

-   Go version compatible with `go.mod`
-   MySQL
-   Git

### 1. Get the source code

If the repository is hosted on Git, clone it and enter the backend
directory:

``` bash
git clone <repository-url>
cd servixa/backend
```

Replace `<repository-url>` with the actual repository URL.

### 2. Configure the application

Review the configuration implementation in `internal/config` and the
database initialization in `internal/database`.

Configure the required database connection values and JWT secret using
the mechanism expected by the application. Do not commit credentials,
tokens, or production secrets to source control.

### 3. Run the API

From the backend directory:

``` bash
go mod download
go run .
```

The database and required configuration must be available. The API port
depends on the current application configuration.

### 4. Check health

``` http
GET /health
```

Use the host and port configured for your local environment.

## Running Tests

Run the full test suite:

``` bash
go test ./...
```

Run static analysis:

``` bash
go vet ./...
```

Some tests use the configured database. Ensure the test database
configuration is appropriate before running tests.

## API Reference

All API paths below are relative to the `/api/v1` prefix unless noted
otherwise.

**Base path:** `/api/v1`

The tables document the method, route, and role restriction registered
in `internal/router/router.go`. This reference intentionally does not
invent request fields or response bodies that have not been verified
from the handlers.

### Health

  Method   Endpoint    Access   Description
  -------- ----------- -------- --------------------------
  `GET`    `/health`   Public   Check application health

### Authentication

  ------------------------------------------------------------------------
  Method            Endpoint           Access            Description
  ----------------- ------------------ ----------------- -----------------
  `POST`            `/auth/register`   Public route      Register a user

  `POST`            `/auth/login`      Public route      Authenticate a
                                                         user

  `POST`            `/auth/refresh`    Public route      Refresh
                                                         authentication
                                                         credentials

  `POST`            `/auth/logout`     Public route      Log out /
                                                         invalidate
                                                         credentials as
                                                         implemented by
                                                         the handler
  ------------------------------------------------------------------------

The routes above are registered outside the global authentication
middleware. Check `auth_handler.go` for required fields and the exact
token lifecycle behavior.

### Customers

All customer routes require authentication and the `CUSTOMER` role.

  -----------------------------------------------------------------------
  Method            Endpoint          Access            Description
  ----------------- ----------------- ----------------- -----------------
  `POST`            `/customers`      `CUSTOMER`        Create the
                                                        authenticated
                                                        user's customer
                                                        profile

  `GET`             `/customers/me`   `CUSTOMER`        Retrieve the
                                                        authenticated
                                                        user's customer
                                                        profile
  -----------------------------------------------------------------------

### Technicians

All technician routes require authentication and the `TECHNICIAN` role.

  -------------------------------------------------------------------------
  Method            Endpoint            Access            Description
  ----------------- ------------------- ----------------- -----------------
  `POST`            `/technicians`      `TECHNICIAN`      Create the
                                                          authenticated
                                                          user's technician
                                                          profile

  `GET`             `/technicians/me`   `TECHNICIAN`      Retrieve the
                                                          authenticated
                                                          user's technician
                                                          profile
  -------------------------------------------------------------------------

### Vehicles

All vehicle routes require authentication and either `ADMIN` or
`DISPATCHER`.

  Method   Endpoint          Access                  Description
  -------- ----------------- ----------------------- --------------------
  `POST`   `/vehicles`       `ADMIN`, `DISPATCHER`   Create a vehicle
  `GET`    `/vehicles/:id`   `ADMIN`, `DISPATCHER`   Retrieve a vehicle
  `PUT`    `/vehicles/:id`   `ADMIN`, `DISPATCHER`   Update a vehicle

### Parts

All part routes require authentication and either `ADMIN` or
`DISPATCHER`.

  Method   Endpoint       Access                  Description
  -------- -------------- ----------------------- -----------------
  `POST`   `/parts`       `ADMIN`, `DISPATCHER`   Create a part
  `GET`    `/parts/:id`   `ADMIN`, `DISPATCHER`   Retrieve a part
  `PUT`    `/parts/:id`   `ADMIN`, `DISPATCHER`   Update a part

### Service Requests

Customer routes require authentication and the `CUSTOMER` role. Review
and rejection routes require `ADMIN` or `DISPATCHER`.

  --------------------------------------------------------------------------------------
  Method            Endpoint                         Access            Description
  ----------------- -------------------------------- ----------------- -----------------
  `POST`            `/service-requests`              `CUSTOMER`        Submit a service
                                                                       request

  `GET`             `/service-requests`              `CUSTOMER`        List the
                                                                       authenticated
                                                                       customer's
                                                                       requests

  `GET`             `/service-requests/:id`          `CUSTOMER`        Retrieve a
                                                                       service request

  `PATCH`           `/service-requests/:id/review`   `ADMIN`,          Review a service
                                                     `DISPATCHER`      request

  `PATCH`           `/service-requests/:id/reject`   `ADMIN`,          Reject a service
                                                     `DISPATCHER`      request
  --------------------------------------------------------------------------------------

### Work Orders

The routes in this section are registered under a group requiring
authentication and either `ADMIN` or `DISPATCHER`, unless noted
otherwise.

  ----------------------------------------------------------------------------------------------
  Method            Endpoint                                 Access            Description
  ----------------- ---------------------------------------- ----------------- -----------------
  `POST`            `/work-orders`                           `ADMIN`,          Create a work
                                                             `DISPATCHER`      order

  `GET`             `/work-orders`                           `ADMIN`,          List work orders
                                                             `DISPATCHER`      

  `GET`             `/work-orders/:id`                       `ADMIN`,          Retrieve a work
                                                             `DISPATCHER`      order

  `POST`            `/work-orders/:id/assign-technician`     `ADMIN`,          Assign a
                                                             `DISPATCHER`      technician

  `PATCH`           `/work-orders/:id/reassign-technician`   `ADMIN`,          Reassign a
                                                             `DISPATCHER`      technician

  `POST`            `/work-orders/:id/assign-vehicle`        `ADMIN`,          Assign a vehicle
                                                             `DISPATCHER`      

  `POST`            `/work-orders/:id/schedule`              `ADMIN`,          Schedule a work
                                                             `DISPATCHER`      order

  `PATCH`           `/work-orders/:id/cancel`                `ADMIN`,          Cancel a work
                                                             `DISPATCHER`      order

  `POST`            `/work-orders/:id/completion`            See note below    Record work
                                                                               completion
  ----------------------------------------------------------------------------------------------

**Route configuration note:** `/work-orders/:id/completion` is
registered inside a group that requires `ADMIN` or `DISPATCHER`, and the
route also applies `RequireRoles(TECHNICIAN)`. If the middleware checks
the role on every request as expected, these restrictions conflict and
no single role can satisfy both. Verify and correct the route
registration before relying on this endpoint. The separate technician
completion flow should also be reviewed.

### Service Visits

These routes require authentication and the `TECHNICIAN` role.

  --------------------------------------------------------------------------------------
  Method            Endpoint                         Access            Description
  ----------------- -------------------------------- ----------------- -----------------
  `POST`            `/work-orders/:id/visit`         `TECHNICIAN`      Create a service
                                                                       visit for a work
                                                                       order

  `POST`            `/service-visits/:id/start`      `TECHNICIAN`      Start a service
                                                                       visit

  `POST`            `/service-visits/:id/complete`   `TECHNICIAN`      Complete a
                                                                       service visit
  --------------------------------------------------------------------------------------

### Customer Confirmation

This route requires authentication and the `CUSTOMER` role.

  ---------------------------------------------------------------------------------------
  Method            Endpoint                          Access            Description
  ----------------- --------------------------------- ----------------- -----------------
  `POST`            `/work-orders/:id/confirmation`   `CUSTOMER`        Confirm work
                                                                        completion

  ---------------------------------------------------------------------------------------

### Notifications

All notification routes require authentication. The router does not
apply a role restriction to this group; handlers use the authenticated
user context.

  -------------------------------------------------------------------------------------
  Method            Endpoint                        Access            Description
  ----------------- ------------------------------- ----------------- -----------------
  `GET`             `/notifications`                Authenticated     List the user's
                                                    user              notifications

  `GET`             `/notifications/unread`         Authenticated     List unread
                                                    user              notifications

  `GET`             `/notifications/unread/count`   Authenticated     Count unread
                                                    user              notifications

  `PATCH`           `/notifications/read-all`       Authenticated     Mark the user's
                                                    user              notifications as
                                                                      read

  `PATCH`           `/notifications/:id/read`       Authenticated     Mark one of the
                                                    user              user's
                                                                      notifications as
                                                                      read
  -------------------------------------------------------------------------------------

## Authentication

Protected endpoints require an access token in the HTTP `Authorization`
header:

``` http
Authorization: Bearer <access_token>
```

Obtain tokens through the authentication endpoints. The exact request
fields, token response format, expiration settings, and refresh/logout
behavior are defined by the authentication implementation and should be
documented from `internal/handler/auth_handler.go` and
`internal/service/auth_service.go`.

## Business Workflow

The intended high-level workflow is:

1.  A customer submits a service request.
2.  An administrator or dispatcher reviews or rejects the request.
3.  An administrator or dispatcher creates a work order.
4.  A technician is assigned and a vehicle is assigned as needed.
5.  The work order is scheduled.
6.  The technician creates, starts, and completes a service visit.
7.  Work completion is recorded.
8.  The customer confirms the completed work.
9.  Relevant notifications are made available to users.

The exact status transitions and prerequisites are enforced by service
logic; consult the service implementations for those rules.

## Automated Tests

The current test files cover selected scenarios for:

-   Customer confirmation authorization, concurrency, and notification
    behavior.
-   Notification service and HTTP handler behavior.
-   Service visit completion concurrency.
-   Work completion authorization and concurrency.
-   Work order assignment validation and concurrency.
-   Work order, scheduling, and vehicle assignment concurrency.

These tests do not by themselves establish exhaustive coverage of all
API paths or the complete end-to-end business workflow.

## Documentation Notes

-   This API reference is derived from the route registrations in
    `internal/router/router.go`.
-   Request and response examples should be added after checking each
    handler's actual JSON contract.
-   Confirm the configured API port and environment variable names from
    the source before adding them here.
-   Keep this README updated when routes, roles, request fields, or
    response formats change.
