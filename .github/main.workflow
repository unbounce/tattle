workflow "Scan for Vulnerabilities with SonarCloud" {
  on = "push"
  resolves = ["sonarcloud-scan"]
}

action "sonarcloud-scan" {
  uses = "./.github/action/sonarcloud-scan"
  secrets = ["SONAR_LOGIN"]
}

