# ETI Assignment 2

Hi, I am Kenneth Teo, a Year 3 student studying for a Diploma in Information Technology at Ngee Ann Polytechnic.

This respository contains the source code for my ETI assignment 2 project.
This assignment is about an education and financial application which has been split into several microservices. The section that I would be working on is **3.15. Timetable**

The task breakdown for this assignment is as follows:

- 3.15.1. Auto award 20 ETI tokens for each student at start of semester
- 3.15.2. Auto allocate based on highest amount
- 3.15.3. Generate timetable
- 3.15.4. Auto refund failed bids

## Design consideration of microservices

## Architecture diagram

images of diagram

**\*For a more in-depth documentation of the microservices, please refer to the [API Reference](./APIDocumentation.md) file**

## Link to your container image

Each service has been publicly published onto Docker Hub.

Front-end web view: https://hub.docker.com/r/kennethteo2002/frontend-testingtimetable

Add credit microservice: https://hub.docker.com/r/kennethteo2002/microservice-addcredits

Allocate bids microservice: https://hub.docker.com/r/kennethteo2002/microservice-bidallocation

Timetable API: https://hub.docker.com/r/kennethteo2002/microservice-timetable

## Instructions for setting up and running your microservices

After setting up the services, the applications would be hosted from http://localhost:8070-8073

**Automatic deployment**

This deployment uses the `docker-compose.yml` file to automatically build and run the application containers.

Prerequisite:

- Downloaded git repository to local storage
- Docker Destop installed

Steps:

1. Open a command terminal and navigate to project ROOT directory
2. Run command `docker-compose up --build`

**Manual deployment**

In order to run the project without docker, run each of the following commands in seperate command prompts.

```
# run front-end application
npm start --prefix Web

# run add credits microservice
go run AddCredits/main.go

# run allocate bids microservice
go Allocation/run main.go

# run timetable microservice
go run Timetable/main.go
```
