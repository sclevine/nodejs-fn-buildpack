# Node.js Riff Function Buildpack

[![CF Slack](https://www.google.com/s2/favicons?domain=www.slack.com) Join us on Slack](https://cloudfoundry.slack.com/messages/buildpacks/)

A Cloud Foundry [buildpack](http://docs.cloudfoundry.org/buildpacks/) for Node based apps.

### Buildpack User Documentation

Official buildpack documentation can be found at [node buildpack docs](http://docs.cloudfoundry.org/buildpacks/node/index.html).

### Building the Buildpack

To build this buildpack, run the following commands from the buildpack's directory:

1. Source the .envrc file in the buildpack directory.

   ```bash
   source .envrc
   ```
   To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Install buildpack-packager

    ```bash
    (cd src/fn/vendor/github.com/cloudfoundry/libbuildpack/packager/buildpack-packager && go install)
    ```

1. Build the buildpack

    ```bash
    buildpack-packager build [ --cached=(true|false) ]
    ```

1. Use in Cloud Foundry

   Upload the buildpack to your Cloud Foundry and optionally specify it by name

    ```bash
    cf create-buildpack [BUILDPACK_NAME] [BUILDPACK_ZIP_FILE_PATH] 1
    cf push my_app [-b BUILDPACK_NAME]
    ```

### Help and Support

Join the #buildpacks channel in our [Slack community](http://slack.cloudfoundry.org/).

### Reporting Issues

Open an issue on this project.

### Active Development

The project backlog is on [Pivotal Tracker](https://www.pivotaltracker.com/projects/1042066).
