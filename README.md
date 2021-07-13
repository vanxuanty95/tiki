## Introduce
This is a simple backend for a good old booking service, right now this service can handle login/list/create simple bookings.  
The application will never return the client the exact error from the server, if you want details please see the console log.
The Login API will return a JSON token without setting a cookie so that it can be easily used by Fontend or another service as this application can be part of many services.

## Install and Run
### Requirements
1. Docker/Docker Compose

### Run
The fast way to run the service is by executing "make" target from root folder of the repository:
- `make init`
- `make docker_up`
- `make run`

## Guide

###Run unit test:
- `make unit_test`

###Run integration test:
- `make integration_test`

Before running or integration test, bring mysql online:
- `make init`
- `make docker_up`

After everything done, bring mysql offline:

- `make docker_down`

###Configuration

To change the configuration information about the server, the database you can edit it in the file `config/config.{your_state}.yaml` before running
(By default {your_state} is "local")

#Structure
Separate 2 separate API parts for Booking and User

Storages layer is in internal/api/{each API}/storages which have 2 drivers: mysql and sqlite to interact with database, no business logic in this layer.

Use case layer is in internal/api/{each API}/usecases which handle business core and use storage layer to reach DB.

Handlers layer is in internal/api/{each API}/handler to handle HTTP routing, validate data before send to usecase layer and make JSON response for client.

Before entering Handlers layer, middleware will print request information to log as well as validate USER for all APIs except Login

# Interact with the API

__Login__

```curl -X POST  -d '{"user_id":"firstUser","password":"example"}' "http://localhost:5050/login"```

__Get list bookings__

```curl -H "Authorization: Basic <_your_token_>" "http://localhost:5050/bookings?created_date=2020-06-29"```

```curl -H "Authorization: Basic <_your_token_>" "http://localhost:5050/bookings?created_date=2020-06-29?page=1"```

```curl -H "Authorization: Basic <_your_token_>" "http://localhost:5050/bookings?created_date=2020-06-29?page=1?limit=10"```

__Create a booking__

```curl -X POST -H "Authorization: Basic <_your_token_>" -d '{"content":"your content"}' "http://localhost:5050/bookings"```


### DB Schema
```sql
-- users definition

CREATE TABLE users (
	id TEXT NOT NULL,
	password TEXT NOT NULL,
	max_booking INTEGER DEFAULT 5 NOT NULL,
	CONSTRAINT users_PK PRIMARY KEY (id)
);

INSERT INTO users (id, password, max_booking) VALUES('firstUser', '$2a$14$BdgOuNVBU7sdGW9rIDIIv.MWXDdvTVKyTppb3qW03bmvz/6hhA1FO', 5);

-- bookings definition

CREATE TABLE bookings (
	id TEXT NOT NULL,
	content TEXT NOT NULL,
	user_id TEXT NOT NULL,
    created_date TEXT NOT NULL,
	CONSTRAINT bookings_PK PRIMARY KEY (id),
	CONSTRAINT bookings_FK FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### Sequence diagram
![auth and create bookings request](https://tiki/blob/master/docs/sequence.svg)
