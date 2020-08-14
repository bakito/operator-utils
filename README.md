[![Build Status](https://travis-ci.com/bakito/operator-utils.svg?branch=master)](https://travis-ci.com/bakito/operator-utils) [![Go Report Card](https://goreportcard.com/badge/github.com/bakito/operator-utils)](https://goreportcard.com/report/github.com/bakito/operator-utils) [![GitHub Release](https://img.shields.io/github/release/bakito/operator-utils.svg?style=flat)](https://github.com/bakito/operator-utils/releases)   <a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-70%25-brightgreen.svg?longCache=true&style=flat)</a>

# Operator Utils
A collection of reusable utils when writing operators

## wehook certs controller
A Controller that automatically creates/updates certs for webhooks.
The certs are stored in a secret. The secret is mounted as volumne into a pod.
Once the volumne is updated in the pod. The ca certs in the webhook configurations are updated.