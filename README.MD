to deploy no heroku ref: https://elements.heroku.com/buildpacks/heroku/heroku-buildpack-go

## add this command to Gopkg.toml
[metadata.heroku]
  root-package = "github.com/KrgKo/todo"
  go-version = "go1.10.4"
  install = [ "./..." ]
  ensure = "true"

## heroku command -> native
- heroku login
- heroku create
- heroku config:set ENV=sss

## heroku deploy
- git push heroku master

## heroku command -> docker
- heroku container:login
- heroku container:push web
- heroku container:release web

Note: `do not set port on heroku -> it will set by itself`