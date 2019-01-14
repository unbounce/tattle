workflow "New workflow" {
  on = "push"
  resolves = ["Analyze with SonarCloud Scanner"]
}

action "Analyze with SonarCloud Scanner" {
  uses = "actions/docker/cli@c08a5fc9e0286844156fefff2c141072048141f6"
  secrets = ["SONAR_LOGIN"]
  env = {
    SONAR_DOWNLOAD_URL = "https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-3.2.0.1227-linux.zip"
    SONAR_ORG = "scottbrown-github"
  }
  runs = "./sonar-scanner-3.2.0.1227-linux/bin/sonar-scanner -Dsonar.projectKey=tattle -Dsonar.organization=$SONAR_ORG -Dsonar.sources=. -Dsonar.host.url=https://sonarcloud.io -Dsonar.login=${SONAR_LOGIN} -Dsonar.branch.name=$GITHUB_BRANCH -X"
}
