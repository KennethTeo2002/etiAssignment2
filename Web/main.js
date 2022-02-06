const express = require("express");
const axios = require("axios");
const url = require("url");
const schedule = require("node-schedule");
const app = express();

const addCreditsMicroserviceURL = "http://10.31.11.11:8071/addCredits";
const allocateBidURL = "http://10.31.11.11:8072/allocateBid";
const timetableAPIURL = "http://10.31.11.11:8073/api/timetable";

var timetablehtml = "";

app.use(express.static("public"));
app.use("/js", express.static(__dirname + "Public/js"));

app.use(express.urlencoded({ extended: true }));
app.use(express.json());

app.set("view engine", "ejs");

app.get("/", (req, res) => {
  res.render("index");
});

// 3.15.1 Add credits at start of semester
// automated job on monday 00:00
const autoAddCredit = schedule.scheduleJob(
  { dayOfWeek: 1, hour: 0, minute: 0 },
  () => {
    console.log("Free 20 credits for all!");
    axios
      .post(addCreditsMicroserviceURL)
      .then(function (response) {
        console.log(response.data);
      })
      .catch(function (error) {
        if (error.response) {
          console.error(error.response.data);
        } else {
          console.error("Failed to connect to microservice");
        }
      });
  }
);

// manual api call on testing panel
app.get("/addToken", (req, res) => {
  console.log("adding tokens");
  axios
    .post(addCreditsMicroserviceURL)
    .then(function (response) {
      console.log(response.data);
    })
    .catch(function (error) {
      if (error.response) {
        console.error(error.response.data);
      } else {
        console.error("Failed to connect to microservice");
      }
    });

  res.redirect("/");
});

// 3.15.2 allocate bids for classes
// automated job on saturday 23:59
const autoAllocateBids = schedule.scheduleJob(
  { dayOfWeek: 6, hour: 23, minute: 59 },
  () => {
    console.log("Sorting and allocating algo :)");
    axios
      .post(timetableAPIURL)
      .then(function (response) {
        console.log(response.data);
      })
      .catch(function (error) {
        if (error.response) {
          console.error(error.response.data);
        } else {
          console.error("Failed to connect to microservice");
        }
      });
  }
);

// manual api call on testing panel
app.get("/allocateBids", (req, res) => {
  console.log("allocate bids to classes");
  axios
    .post(allocateBidURL)
    .then(function (response) {
      console.log(response);
    })
    .catch(function (error) {
      if (error.response) {
        console.error(error.response.data);
      } else {
        console.error("Failed to connect to microservice");
      }
    });

  res.redirect("/");
});

// 3.15.3 Set time schedule for classes
// automated job on friday 23:59
const autoAllocateClassSchedule = schedule.scheduleJob(
  { dayOfWeek: 5, hour: 23, minute: 59 },
  () => {
    console.log("Populate class schedules!");
    axios
      .post(timetableAPIURL)
      .then(function (response) {
        console.log(response.data);
      })
      .catch(function (error) {
        if (error.response) {
          console.error(error.response.data);
        } else {
          console.error("Failed to connect to microservice");
        }
      });
  }
);

// manual api call on testing panel
app.get("/allocateSchedule", (req, res) => {
  console.log("allocate schedule for classes");

  axios
    .post(timetableAPIURL)
    .then(function (response) {
      console.log(response.data);
    })
    .catch(function (error) {
      if (error.response) {
        console.error(error.response.data);
      } else {
        console.error("Failed to connect to microservice");
      }
    });

  res.redirect("/");
});

//redirect to timetable with url query
app.post("/timetable", (req, res) => {
  var urlParams = "";
  if (req.body.studentid) {
    urlParams += "?studentID=" + req.body.studentid;
  } else if (req.body.tutorid) {
    urlParams += "?tutorID=" + req.body.tutorid;
  }
  if (req.body.semester) {
    urlParams += "&semester=" + req.body.semester;
  }
  res.redirect("/timetable" + urlParams);
});

// 3.15.3 retrieve class timetable for user
app.get("/timetable", (req, res) => {
  //break down url to query api
  var apiurl = timetableAPIURL;
  if (req.query.studentID) {
    apiurl += "?studentID=" + req.query.studentID;
  } else if (req.query.tutorID) {
    apiurl += "?tutorID=" + req.query.tutorID;
  }
  if (req.query.semester) {
    apiurl += "&semester=" + req.query.semester;
  }
  console.log(apiurl);
  axios
    .get(apiurl)
    .then(function (response) {
      timetablehtml = response.data;
    })
    .catch(function (error) {
      timetablehtml = "<h1>Error getting timetable</h1>";
      timetablehtml = apiurl + error;
    });
  res.locals.query = req.query;
  res.render("timetable", {
    timetabledata: timetablehtml,
  });
});

app.listen(8070);
