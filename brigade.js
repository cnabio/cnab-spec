const { events, Job, Group } = require("brigadier");

const projectName = "cnab-spec";

// Event handlers

events.on("exec", (e, p) => {
  validate(e, p).run();
});
events.on("check_suite:requested", runSuite);
events.on("check_suite:rerequested", runSuite);
events.on("check_run:rerequested", runSuite);

// Functions/Helpers

function validate(e, project) {
  var validator = new Job(`${projectName}-validate`, "node:8-alpine");
  validator.streamLogs = true;

  validator.tasks = [
    "apk add --update make",
    "cd /src",
    "make validate-local",
  ];

  return validator
}

// Here we can add additional Check Runs, which will run in parallel and
// report their results independently to GitHub
function runSuite(e, p) {
  runValidation(e, p).catch(e => {console.error(e.toString())});
}

// runValidation is a Check Run that is ran as part of a Checks Suite
function runValidation(e, p) {
  // Create Notification object (which is just a Job to update GH using the Checks API)
  var note = new Notification(`validation`, e, p);
  note.conclusion = "";
  note.title = "Run Validation";
  note.summary = "Running the schema validation for " + e.revision.commit;
  note.text = "Ensuring all bundle.json files adhere to json schema specs"

  // Send notification, then run, then send pass/fail notification
  return notificationWrap(validate(e, p), note)
}

// A GitHub Check Suite notification
class Notification {
  constructor(name, e, p) {
    this.proj = p;
    this.payload = e.payload;
    this.name = name;
    this.externalID = e.buildID;
    this.detailsURL = `https://brigadecore.github.io/kashti/builds/${ e.buildID }`;
    this.title = "running check";
    this.text = "";
    this.summary = "";

    // count allows us to send the notification multiple times, with a distinct pod name
    // each time.
    this.count = 0;

    // One of: "success", "failure", "neutral", "cancelled", or "timed_out".
    this.conclusion = "neutral";
  }

  // Send a new notification, and return a Promise<result>.
  run() {
    this.count++
    var j = new Job(`${ this.name }-${ this.count }`, "deis/brigade-github-check-run:latest");
    j.env = {
      CHECK_CONCLUSION: this.conclusion,
      CHECK_NAME: this.name,
      CHECK_TITLE: this.title,
      CHECK_PAYLOAD: this.payload,
      CHECK_SUMMARY: this.summary,
      CHECK_TEXT: this.text,
      CHECK_DETAILS_URL: this.detailsURL,
      CHECK_EXTERNAL_ID: this.externalID
    }
    return j.run();
  }
}

// Helper to wrap a job execution between two notifications.
async function notificationWrap(job, note, conclusion) {
  if (conclusion == null) {
    conclusion = "success"
  }
  await note.run();
  try {
    let res = await job.run()
    const logs = await job.logs();

    note.conclusion = conclusion;
    note.summary = `Task "${ job.name }" passed`;
    note.text = note.text = "```" + res.toString() + "```\nTest Complete";
    return await note.run();
  } catch (e) {
    const logs = await job.logs();
    note.conclusion = "failure";
    note.summary = `Task "${ job.name }" failed for ${ e.buildID }`;
    note.text = "```" + logs + "```\nFailed with error: " + e.toString();
    try {
      return await note.run();
    } catch (e2) {
      console.error("failed to send notification: " + e2.toString());
      console.error("original error: " + e.toString());
      return e2;
    }
  }
}