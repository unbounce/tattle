workflow "Scan for Vulnerabilities with SonarCloud" {
  on = "push"
  resolves = ["sonarcloud-scan"]
}

action "sonarcloud-scan" {
  uses = "docker://unbounce/sonarcloud-github"
  secrets = ["SONAR_LOGIN"]
  env = {
    "SONAR_PROJECT_KEY" = "unbounce_tattle"
    "SONAR_ORG" = "unbounceinc"
  }
}

