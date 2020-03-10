const { events, Job, Group } = require("brigadier");
const { Check } = require("@brigadecore/brigade-utils");

const projectName = "cnab-spec";

// Event handlers

events.on("exec", (e, p) => {
  return Group.runAll([
    validate(e, p),
    validateURL(e, p)
  ]);
});
events.on("check_suite:requested", runSuite);
events.on("check_suite:rerequested", runSuite);
events.on("check_run:rerequested", runSuite);
events.on("issue_comment:created", (e, p) => Check.handleIssueComment(e, p, runSuite));
events.on("issue_comment:edited", (e, p) => Check.handleIssueComment(e, p, runSuite));
events.on("push", (e, p) => {
  if (e.revision.ref.startsWith("refs/tags/")) {
    return publish(e, p).run();
  }
});

// Functions/Helpers

function validate(e, project) {
  var validator = new Job(`${projectName}-validate`, "node:8-alpine");
  validator.streamLogs = true;

  validator.tasks = [
    "apk add --update make curl jq",
    "cd /src",
    "make validate-local",
  ];

  return validator;
}

function validateURL(e, project) {
  var validator = new Job(`${projectName}-validate-url`, "node:8-alpine");
  validator.streamLogs = true;

  validator.tasks = [
    "apk add --update make curl jq",
    "cd /src",
    "make validate-url-local",
  ];

  return validator;
}

function publish(e, p) {
  var publisher = new Job(`${projectName}-publish`, "node:8-alpine");

  publisher.env.AZURE_STORAGE_CONNECTION_STRING = p.secrets.azureStorageConnectionString;
  publisher.tasks.push(
    "apk add --update make curl",
    // Fetch az cli needed for publishing
    "curl -sLO https://github.com/carolynvs/az-cli/releases/download/v0.3.2/az-linux-amd64 && \
      chmod +x az-linux-amd64 && \
      mv az-linux-amd64 /usr/local/bin/az",
    "cd /src",
    "make publish"
  );

  return publisher;
}

// Here we can add additional Check Runs, which will run in parallel and
// report their results independently to GitHub
function runSuite(e, p) {
  return runValidation(e, p, "validate")
  .then(() => {
    if (e.revision.ref == "master") {
      validateURL(e, p).run();
    }
  })
  .catch(e => {console.error(e.toString())});
}

// runValidation is a Check Run that is ran as part of a Checks Suite
function runValidation(e, p, jobFunc) {
  var check = new Check(e, p, jobFunc(),
    `https://brigadecore.github.io/kashti/builds/${e.buildID}`);
  return check.run();
}
