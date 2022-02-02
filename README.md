# ETI Assignment 2

Hi, I am Kenneth Teo, a Year 3 student studying for a Diploma in Information Technology at Ngee Ann Polytechnic.

This respository contains the source code for my ETI assignment 2 project.
This assignment is about an education and financial application which has been split into several microservices. The section that I would be working on is **3.15. Timetable**

The task breakdown for this assignment is as follows:

- 3.15.1. Auto award 20 ETI tokens for each student at start of semester
- 3.15.2. Auto allocate based on highest amount
- 3.15.3. Generate timetable
- 3.15.4. Auto refund failed bids

## Architecture diagram

<ins>Add credit microservice</ins>

![Add credit microservice design](addCreditDesign.png?raw=true "Add credit design")

<ins>Allocate bids microservice</ins>

![Allocate bids microservice design](allocateBidsDesign.png?raw=true "Allocate bids design")

<ins>Timetable microservice</ins>

![Timetable microservice design](timetableDesign.png?raw=true "Timetable design")

**\*For a more in-depth documentation of the microservices, please refer to the [API Reference](./APIDocumentation.md) file**

Unlike other packages, the task breakdown for 3.15 are isolated and have no correlation with each other. Thus, I decided to split the project into 3 microservices, AddCredit, AllocateBids and Timetable.

The add credit microservice doesn't follow the REST API format, as it only requires one functionality, which is to add 20 ETI credits to all student wallet. Thus, I did not use the different methods when creating this microservice. Similarly, the allocate bids microservice's only usage is to sort out bids into classes, so no methods were created. Whereas for Timetable, since I have 2 requirements, deconflicting the class schedules by allocating a unique timeslot for each class, as well as retrieving the timetable for a specific user, Timetable is a REST API which uses the POST and GET methods.

### Scalability

Since add credit and allocate bids are backend server functions that should only be called by a time scheduled job running once a week, scaling is not an issue.

Whereas, for timetable, since the GET method is public for users, scaling is required, in case of a surge in students looking up their timetable. So, by isolating this requirement from the previous 2, the timetable API can be scaled up without affecting the add credit and allocate bids microservices.

### Security

For add credit and allocate bids, since they are backend server functions, the endpoints are not exposed to public, these microservices could just be deployed on a local server which is called by an internal clock system, so not much security is required.

Currently for timetable, since other students needs to implement my services into their pages, and it might be hard for them to implement the authentication on their end, I decided not to create any authentication, instead they can run the api call internally on their side with the student's session token, and the api would return the raw html code for the timetable.

However, if I were to create an authenication system, instead of inputting their userID raw, I would require the users to send their session id and sesson key to my api, so that I can refer to the login redux cache and retrieve their userID from there.

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
