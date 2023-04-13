# Welcome to Buffalo Test Examples

This repository aims to show some examples of how I've done different tests in
Buffalo apps. I'm not saying that these are best-practice, just that I've found
these techniques and patterns useful for validating what's happening.

Please feel free to raise issues and discussions where you have a comment or
suggestion. And I'm always happy for someone to raise a PR to show off what
they've discovered too.

# Examples

Click on each of these named examples to see more details...

<details>
  <summary>Testing in Github Actions</summary>

Here you can see an example of what's necessary to test your app using Github
Actions.

## Important files

* [Compose file](docker-compose.yml) - this puts your tests together with the database
* [database.yml](database.yml) - defines your database details
* [Tests script](scripts/docker-test.sh) - This does the testing for you.
* [PR workflow](.github/workflows/pr.yml) - this runs the tests using docker compose
* [Linting workfow](.github/workflows/golangci-lint.yml) - this enforces common linting rules

## What's happening

Let's start with the Docker Compose file. This sets up the DB and makes your app depend on
it, setting the database details using environment variables. These variables are consumed
in `database.yml`. The compose file also gives access to all of the files in the project and
runs the tests script, that sets up the DB by running the migrations and then executes your
tests.

If you would like to run these tests from inside docker locally, just use the following command:
```shell
docker-compose down && docker-compose up --remove-orphans --abort-on-container-exit test
```

The PR workflow is the bit that triggers your tests when you open a PR and push new commits
to it.

Lastly, we have the linting workflow - this is just a good idea really.

</details>

<details>

  <summary>Little user testing tweaks</summary>

This is a really simple one, I've created a couple of helpers that have been useful when
defining user interactions.

## Important files

* [actions_test.go](actions/actions_test.go) - test helpers sit here
* [fixture data](fixtures/base-data.toml) - test data to be inserted in to the database
* [example test](actions/home_test.go) - `Test_HomeHandler_LoggedIn` shows a helper being used

## What's happening

These are a couple of very simple helpers to:

1. Set the current session using the email address as the key
2. Get the a user details by ID

In both of these helpers, we just fail immediately if there's a problem. Don't bother
propagating the error back up because this is a fundamental failure in the test environment.

Also, don't forget to load the fixtures at the beginning of your tests!

</details>

<details>

  <summary>Workflow tests</summary>

Workflow tests allow you to put together several actions so that you can test more complex
behaviours. This work by taking advantage of the session and other internal details to
retain user state between actions.

The idea is that you can define behaviour that works across different endpoints in your app
so that you can think about a feature from beginning to end. It's been particularly powerful
when used with tests around user auth.

## Important files

* [Workflow tests](actions/workflow_test.go) - this is where the individual workflow tests live

</details>

