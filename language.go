package main

import(
  "code.google.com/p/goconf/conf"
)

type Language struct {
  name string
  config *conf.ConfigFile
  filename string
  docker_image string
  apparmor_profile string
}

func (lang *Language) readConfig() (err error) {
  lang.name, err = lang.config.GetString("general", "name")

  if err == nil {
    lang.filename, err = lang.config.GetString("general", "filename")
  }

  if err == nil {
    lang.docker_image, err = lang.config.GetString("general", "docker_image")
  }

  if err == nil {
    lang.apparmor_profile, err = lang.config.GetString("general", "apparmor_profile")
  }

  return
}

func LoadLanguage(configName string) (lang Language, err error) {
  lang.config, err = conf.ReadConfigFile(configName)
  if err == nil {
    err = lang.readConfig()
  }
  return
}
