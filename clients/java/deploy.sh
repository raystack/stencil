#!/usr/bin/env bash
echo "Deploying staging artifact"
./gradlew clean  publishToSonatype closeAndReleaseStagingRepository -PossrhUsername=${SONATYPE_USERNAME} -PossrhPassword=${SONATYPE_PASSWORD} -Psigning.keyId=${GPG_KEY_ID} -Psigning.password=${GPG_KEY_PASSPHRASE} -Psigning.secretKeyRingFile=private.key --console=verbose
