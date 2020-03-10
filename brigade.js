const { events, Job, Group } = require("brigadier");
const { Check } = require("@brigadecore/brigade-utils");

const projectName = "cnab-spec";
// Currently a very lenient regex.
// Could be made more strict.  Some examples: cnab-core-v1.0.0, cnab-claim-v1.0.0-DRAFT+abc1234
const releaseTagRegex = /^refs\/tags\/(.*)/;

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
  let matchStr = e.revision.ref.match(releaseTagRegex);
  if (matchStr) {
    let matchTokens = Array.from(matchStr);
    let version = matchTokens[1];
    return publish(p, version).run();
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

function publish(p, version) {
  var publisher = new Job(`${projectName}-publish`, "node:8-alpine");

  publisher.env.AZURE_STORAGE_CONNECTION_STRING = p.secrets.azureStorageConnectionString;
  publisher.tasks.push(
    "apk add --update make curl",
    // Fetch az cli needed for publishing
    "curl -sLO https://github.com/carolynvs/az-cli/releases/download/v0.3.2/az-linux-amd64 && \
      chmod +x az-linux-amd64 && \
      mv az-linux-amd64 /usr/local/bin/az",
    "cd /src",
    `VERSION=${version} make publish`
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
