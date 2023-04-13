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

<details>

