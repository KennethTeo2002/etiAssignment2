const express = require("express");
const axios = require("axios");
const url = require("url");
const schedule = require("node-schedule");
const app = express();

const addCreditsMicroserviceURL = "http://localhost:8071/addCredits";
const timetableAPIURL = "http://localhost:8072/api/timetable";
const allocateBidURL = "http://localhost:8073/allocateBid";

var timetablehtml = "";

app.use(express.static("public"));

app.use("/js", express.static(__dirname + "public/js"));

app.use(express.urlencoded({ extended: true }));
app.use(express.json());

app.set("view engine", "ejs");

app.get("/", (req, res) => {
  res.render("index");
});

const autoAddCredit = schedule.scheduleJob(
  { dayOfWeek: 1, hour: 0, minute: 0 },
  () => {
    console.log("Free 20 credits to all!");
    axios
      .post(addCreditsMicroserviceURL)
      .then(function (response) {
        console.log(response.data);
      })
      .catch(function (error) {
        console.error(error.response.data);
      });
  }
);

app.get("/addToken", (req, res) => {
  console.log("adding tokens");
  axios
    .post(addCreditsMicroserviceURL)
    .then(function (response) {
      console.log(response.data);
    })
    .catch(function (error) {
      console.error(error.response.data);
    });

  res.redirect("/");
});

const autoAllocateClassSchedule = schedule.scheduleJob(
  { dayOfWeek: 6, hour: 9, minute: 13 },
  () => {
    console.log("Class schedule!");
    axios
      .post(timetableAPIURL)
      .then(function (response) {
        console.log(response.data);
      })
      .catch(function (error) {
        console.error(error);
      });
  }
);

app.get("/allocateSchedule", (req, res) => {
  console.log("allocate schedule for classes");

  axios
    .post(timetableAPIURL)
    .then(function (response) {
      console.log(response.data);
    })
    .catch(function (error) {
      console.error(error);
    });

  res.redirect("/");
});

const autoAllocateBids = schedule.scheduleJob(
  { dayOfWeek: 6, hour: 23, minute: 59 },
  () => {
    console.log("allocate students to classes!");
    axios
      .post(timetableAPIURL)
      .then(function (response) {
        console.log(response.data);
      })
      .catch(function (error) {
        console.error(error);
      });
  }
);

app.get("/allocateBids", (req, res) => {
  console.log("allocate bids to classes");
  axios
    .post(allocateBidURL)
    .then(function (response) {
      console.log(response);
    })
    .catch(function (error) {
      console.error(error);
    });

  res.redirect("/");
});

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

app.get("/timetable", (req, res) => {
  var apiurl = timetableAPIURL;
  if (req.query.studentid) {
    apiurl += "?studentID=" + req.query.studentid;
  } else if (req.query.tutorid) {
    apiurl += "?tutorID=" + req.query.tutorid;
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
    });
  res.locals.query = req.query;
  res.render("timetable", {
    timetabledata: timetablehtml,
  });
});

app.post("/changeSem", (req, res) => {
  var urlParams = "";
  if (req.query.studentid) {
    urlParams += "?studentID=" + req.query.studentid;
  } else if (req.query.tutorid) {
    urlParams += "?tutorID=" + req.query.tutorid;
  }
  if (req.query.semester) {
    urlParams += "&semester=" + req.body.sem;
  }
  res.redirect("/timetable" + urlParams);
});

app.get("/saveTT", (req, res) => {});

app.listen(8070);
