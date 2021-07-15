## Introduce
This is a simple backend service for movie theater ticket booking
The application will never return the client the exact error from the server, if you want details please see the console log.

Required from social distancing, each seat/group of seats will have a distance from each other of x (x is configured in the config file)

By default, service will run with port 8080

## Install and Run
### Requirements
1. Docker/Docker Compose

### Run
The fast way to run the service is by executing "make" target from root folder of the repository:
- `make init`
- `make docker_up`
- `make run_directly`

if want to run by docker image:
- `make init`
- `make docker_up`
- `make build_image`
- `make run_container`

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

Storages layer is in internal/api/{each API}/storages which is mysql, no business logic in this layer.

Use case layer is in internal/api/{each API}/usecases which handle business core and use storage layer to reach DB.

Handlers layer is in internal/api/{each API}/handler to handle HTTP routing, validate data before send to usecase layer and make JSON response for client.

Before entering Handlers layer, middleware will print request information to log as well as validate USER for all APIs except Login

# Interact with the API

Please get postman file in `docs/` directory

__Login__

```curl -X POST  -d '{"user_id":"tester","password":"example"}' "http://localhost:8080/login"```

__Create a screen__

```curl -X POST -H -d '{"number_seat_row":5, "number_seat_column":6}' "http://localhost:8080/screen"```

__Check seat available__

```curl -X POST -H -d '{"screen_id": 1,"location": {"row": 1,"column": 1}}' "http://localhost:8080/check"```

__Book random seats__

```curl -X POST -H -d '{"screen_id": 1,"number": 1}' "http://localhost:8080/booking"```

```curl -X POST -H -d '{"screen_id": 1,"number": 3}' "http://localhost:8080/booking"```

__Book exactly seats__

```curl -X POST -H -d '{"screen_id": 1,"locations": [{"row": 2,"column": 3}]}' "http://localhost:8080/booking"```


### DB Schema
```sql
-- users definition

CREATE TABLE IF NOT EXISTS `user` (
                                      `id` varchar(50) NOT NULL,
                                      `password` varchar(200) NOT NULL,
                                      PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `screen` (
                                        `id` int(11) NOT NULL AUTO_INCREMENT,
                                        `number_seat_row` int(11) NOT NULL,
                                        `number_seat_column` int(11) NOT NULL,
                                        `created_date` datetime NOT NULL,
                                        `user_id` varchar(20) NOT NULL,
                                        PRIMARY KEY (`id`),
                                        KEY `user_id_idx_screen_id` (`user_id`),
                                        CONSTRAINT `screen_id_user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS  `seat` (
                                       `id` int(11) NOT NULL AUTO_INCREMENT,
                                       `row_id` int(11) NOT NULL,
                                       `column_id` int(11) NOT NULL,
                                       `user_id` varchar(50) NOT NULL,
                                       `screen_id` int(11) NOT NULL,
                                       `booked_date` datetime NOT NULL,
                                       PRIMARY KEY (`id`,`row_id`,`column_id`),
                                       KEY `user_id_idx` (`user_id`),
                                       KEY `id_idx` (`screen_id`),
                                       CONSTRAINT `screen_id` FOREIGN KEY (`screen_id`) REFERENCES `screen` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
                                       CONSTRAINT `user_id` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;
```

### Sequence diagram
![auth and create bookings request](https://raw.githubusercontent.com/vanxuanty95/tiki/master/docs/sequence.png)
