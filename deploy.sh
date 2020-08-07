#!/usr/bin/env bash
openssl aes-256-cbc -K $encrypted_49353df1995f_key -iv $encrypted_49353df1995f_iv -in key.gpg.enc -out key.gpg -d
echo "Deploying staging artifact"
./gradlew publishToSonatype -x test -x sign -PossrhUsername=${SONATYPE_USERNAME} -PossrhPassword=${SONATYPE_PASSWORD} -Psigning.keyId=${GPG_KEY_ID} -Psigning.password=${GPG_KEY_PASSPHRASE} -Psigning.secretKeyRingFile=key.gpg
