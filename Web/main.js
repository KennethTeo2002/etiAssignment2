const express = require("express");
const http = require("http");

const app = express();

app.use(express.urlencoded({ extended: true }));
app.use(express.json());

app.set("view engine", "ejs");

app.get("/", (req, res) => {
  console.log("helloworld");
  res.render("index");
});

app.get("/addToken", (req, res) => {
  console.log("adding tokens");

  const requestDetails = {
    hostname: "docker-addcredits",
    port: 8071,
    path: "/addCredits",
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  };

  const request = https.request(requestDetails, (response) => {
    console.log(`statusCode: ${response.statusCode}`);
  });
  request.on("error", (error) => {
    console.error(error);
  });

  request.end();

  res.redirect("/");
});

app.listen(8070);
