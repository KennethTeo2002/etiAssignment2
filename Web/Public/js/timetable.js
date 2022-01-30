function changeSem() {
  let sem = prompt("Enter semester date in '30-1-2022' format");
  if (sem != null) {
    var url = new URL(window.location.href);
    var search_params = url.searchParams;

    // update sem
    search_params.set("semester", sem);

    url.search = search_params.toString();

    var new_url = url.toString();

    window.location.href = new_url;
  }
}
function save() {
  var htmlContent = [document.documentElement.innerHTML];
  var bl = new Blob(htmlContent, { type: "text/html" });
  var a = document.createElement("a");
  a.href = URL.createObjectURL(bl);
  a.download = "timetable.html";
  a.hidden = true;
  document.body.appendChild(a);
  a.innerHTML =
    "something random - nobody will see this, it doesn't matter what you put here";
  a.click();
}
