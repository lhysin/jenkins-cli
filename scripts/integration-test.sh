#!/bin/bash

set -e

JENKINS_URL="http://localhost:8080"
ADMIN_USER="admin"
JENKINS_CONTAINER="jenkins"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

wait_for_jenkins() {
    log_info "Waiting for Jenkins to be ready..."
    local max_attempts=60
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s -f "${JENKINS_URL}/api/json" --user "admin:${ADMIN_TOKEN}" > /dev/null 2>&1; then
            log_info "Jenkins is ready!"
            return 0
        fi
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    echo ""
    log_error "Jenkins failed to start after ${max_attempts} attempts"
    return 1
}

get_admin_token() {
    docker exec ${JENKINS_CONTAINER} cat /var/jenkins_home/secrets/initialAdminPassword 2>/dev/null
}

setup_jenkins() {
    local token="$1"
    
    log_info "Setting up test job via Jenkins CLI inside container..."
    
    local job_xml='/var/jenkins_home/jobs/test-job/config.xml'
    
    docker exec ${JENKINS_CONTAINER} mkdir -p /var/jenkins_home/jobs/test-job
    
    docker exec ${JENKINS_CONTAINER} /bin/bash -c "cat > ${job_xml}" << 'EOFXML'
<?xml version='1.1' encoding='UTF-8'?>
<project>
  <description>Test job for jenkins-cli integration</description>
  <keepDependencies>false</keepDependencies>
  <properties/>
  <scm class="hudson.scm.NullSCM"/>
  <canRoam>true</canRoam>
  <disabled>false</disabled>
  <blockBuildWhenDownstreamBuilding>false</blockBuildWhenDownstreamBuilding>
  <blockBuildWhenUpstreamBuilding>false</blockBuildWhenUpstreamBuilding>
  <triggers/>
  <concurrentBuild>false</concurrentBuild>
  <builders>
    <hudson.tasks.Shell>
      <command>echo 'Hello from jenkins-cli test job'</command>
    </hudson.tasks.Shell>
  </builders>
  <publishers/>
  <buildWrappers/>
</project>
EOFXML

    docker exec ${JENKINS_CONTAINER} chown jenkins:jenkins ${job_xml}
    
    log_info "Test job created"
}

build_test_job() {
    local token="$1"
    
    log_info "Building test job..."
    
    local crumb=$(curl -s "${JENKINS_URL}/crumbIssuer/api/json" --user "admin:${token}" | grep -o '"crumb":"[^"]*"' | cut -d'"' -f4)
    
    local response
    local http_code
    
    response=$(curl -s -w "\n%{http_code}" -X POST \
        -H "Jenkins-Crumb:${crumb}" \
        -u "admin:${token}" \
        "${JENKINS_URL}/job/test-job/build" 2>&1)
    
    http_code=$(echo "$response" | tail -n1)
    
    if [ "$http_code" = "201" ]; then
        log_info "Build started successfully"
        sleep 15
        return 0
    else
        log_warn "Build start returned HTTP ${http_code}, continuing..."
        sleep 5
        return 0
    fi
}

run_cli_tests() {
    local token="$1"
    log_info "Running jenkins-cli tests..."
    
    local CLI="./jenkins-cli"
    
    if [ ! -f "$CLI" ]; then
        export PATH="/opt/homebrew/opt/go/bin:$PATH"
        go build -o jenkins-cli ./cmd/
    fi
    
    CLI="./jenkins-cli"
    
    log_info ""
    log_info "--- Testing jobs list ---"
    if $CLI --url "${JENKINS_URL}" --user admin --token "${token}" jobs list; then
        log_info "jobs list: PASS"
    else
        log_error "jobs list: FAIL"
        return 1
    fi
    
    log_info ""
    log_info "--- Testing nodes list ---"
    if $CLI --url "${JENKINS_URL}" --user admin --token "${token}" nodes list; then
        log_info "nodes list: PASS"
    else
        log_error "nodes list: FAIL"
        return 1
    fi
    
    log_info ""
    log_info "--- Testing init command ---"
    rm -f ~/.jenkins-cli/config.yaml
    printf "${JENKINS_URL}\nadmin\n${token}\n" | $CLI init 2>/dev/null || true
    
    log_info ""
    log_info "--- Testing status command (using saved config) ---"
    if $CLI status; then
        log_info "status: PASS"
    else
        log_error "status: FAIL"
        return 1
    fi
    
    log_info ""
    log_info "All CLI tests passed!"
}

cleanup() {
    local token="$1"
    
    log_info "Cleaning up..."
    docker exec ${JENKINS_CONTAINER} rm -rf /var/jenkins_home/jobs/test-job 2>/dev/null || true
    log_info "Cleanup complete"
}

main() {
    log_info "========================================="
    log_info "  Jenkins CLI Integration Test"
    log_info "========================================="
    
    if ! docker ps | grep -q ${JENKINS_CONTAINER}; then
        log_error "Jenkins container not running. Start with: docker run -d --name jenkins -p 8080:8080 jenkins/jenkins:lts"
        exit 1
    fi
    
    ADMIN_TOKEN=$(get_admin_token)
    
    if [ -z "$ADMIN_TOKEN" ]; then
        log_error "Could not get admin token"
        return 1
    fi
    
    log_info "========================================="
    log_info "  Jenkins Admin Token: ${ADMIN_TOKEN}"
    log_info "========================================="
    
    wait_for_jenkins
    
    setup_jenkins "${ADMIN_TOKEN}"
    build_test_job "${ADMIN_TOKEN}"
    run_cli_tests "${ADMIN_TOKEN}"
    cleanup "${ADMIN_TOKEN}"
    
    log_info "========================================="
    log_info "  Integration Test Complete!"
    log_info "========================================="
}

main "$@"
