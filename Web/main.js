const express = require("express");
const axios = require("axios");
const url = require("url");
const addCreditsMicroserviceURL = "http://localhost:8071/addCredits";
const timetableAPIURL = "http://localhost:8072/api/timetable";
const allocateBidURL = "http://localhost:8073/allocateBid";

const app = express();

app.use(express.urlencoded({ extended: true }));
app.use(express.json());

app.set("view engine", "ejs");

app.get("/", (req, res) => {
  res.render("index");
});

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
  console.log("get timetable");

  res.redirect(url.format({ pathname: "/timetable", query: req.query }));
});

app.get("/timetable", (req, res) => {
  console.log("print timetable");
  var url = allocateBidURL;
  if (req.body.studentid) {
    url += "?studentID=" + req.body.studentid;
  } else if (req.body.tutorid) {
    url += "?tutorID=" + req.body.tutorid;
  }

  url += "&semester=" + req.body.semester;
  console.log(url);
  axios
    .get(url)
    .then(function (response) {
      timetablehtml = response.data;
      console.log(response);
    })
    .catch(function (error) {
      console.error(error);
    });
  res.render("timetable", timetablehtml);
});

app.listen(8070);
