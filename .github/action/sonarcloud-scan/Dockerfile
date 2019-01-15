FROM ubuntu:18.04
MAINTAINER Scott Brown

WORKDIR /github/workspace

#
# The version of SonarCloud CLI we want to use.
#
ENV SONAR_VERSION=3.2.0.1227

#
# The full URL to the SonarCloud CLI linux package from their CDN.
#
ENV SONAR_DOWNLOAD_URL=https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-$SONAR_VERSION-linux.zip

RUN apt-get update -qq && apt-get install unzip -y

ADD $SONAR_DOWNLOAD_URL /tmp/scanner.zip

RUN mkdir -p /tmp/sonar && unzip /tmp/scanner.zip -d /tmp/sonar

CMD /tmp/sonar/sonar-scanner-$SONAR_VERSION-linux/bin/sonar-scanner \
  -Dsonar.projectKey=$SONAR_PROJECT_KEY \
  -Dsonar.organization=$SONAR_ORG \
  -Dsonar.sources=/github/workspace \
  -Dsonar.host.url=https://sonarcloud.io \
  -Dsonar.login=$SONAR_LOGIN \
  -Dsonar.branch.name=$GITHUB_REF

