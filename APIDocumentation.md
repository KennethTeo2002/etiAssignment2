# Timetable Microservice Reference

## 1. Add credits

This endpoint is used to credit 20 ETI tokens to all student wallet. This microservice should be automatically called at the start of each semester (Monday 00:00).

#### Endpoint URL

```
http://localhost:8071/addCredits
```

#### Example Request

cURL

```
curl -X POST http://localhost:8071/addCredits
```

#### Response

If all the transactions were made successfully, the response will be a status code `200`, else an error code with a corresponding status message would be returned.

## 2. Allocate bids

This endpoint is used to assign students to classes based on highest bids. This microservice should be automatically called when bidding close (Saturday 23:59).

#### Endpoint URL

```
http://localhost:8072/allocateBid
```

#### Example Request

cURL

```
curl -X POST http://localhost:8072/allocateBid
```

#### Response

If all the bids were allocated successfully, the response will be a status code `200`, else an error code with a corresponding status message would be returned.

## 3. Timetable

Base URL: http://localhost:8073/api/timetable

### 3.1 POST Timetable

This endpoint is used to allocate timeslots to classes. This microservice should be automatically called before bidding starts (Friday 23:59).

#### Endpoint URL

```
http://localhost:8073/api/timetable
```

#### Example Request

cURL

```
curl -X POST http://localhost:8073/api/timetable
```

#### Response

If all the class timeslots are allocated successfully, the response will be a status code `200`, else an error code with a corresponding status message would be returned.

### 3.2 GET Timetable

This endpoint is used to retrieve and generate the timetable for a specific userID and semester.

#### Endpoint URL

```
http://localhost:8073/api/timetable
```

#### Query Parameters

| Name      | Type   | Required   | Description                                              | Example     |
| --------- | ------ | ---------- | -------------------------------------------------------- | ----------- |
| semester  | string | Yes        | Any date of a semester to look up in `DD-MM-YYYY` format | `31-1-2022` |
| studentID | string | Either one | The unique student id given to students                  | `S12345678` |
| tutorID   | string | Either one | The unique tutor id given to tutors                      | `T12345678` |

#### Example Request

cURL

```
curl -X GET http://localhost:8073/api/timetable?semester=31-1-2022&studentID=S12345678
```

#### Response

The response would be a status code `200` with a returned data of text that contains the generated timetable html table, else an error code with a corresponding status message would be returned.
