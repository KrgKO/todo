to deploy no heroku ref: https://elements.heroku.com/buildpacks/heroku/heroku-buildpack-go

## add this command to Gopkg.toml
[metadata.heroku]
  root-package = "github.com/KrgKo/todo"
  go-version = "go1.10.4"
  install = [ "./..." ]
  ensure = "true"