[![Go](https://github.com/bakito/operator-utils/actions/workflows/go.yml/badge.svg)](https://github.com/bakito/operator-utils/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/bakito/operator-utils)](https://goreportcard.com/report/github.com/bakito/operator-utils) [![GitHub Release](https://img.shields.io/github/release/bakito/operator-utils.svg?style=flat)](https://github.com/bakito/operator-utils/releases)

# Operator Utils
A collection of reusable utils when writing operators

## webhook certs controller
A Controller that automatically creates/updates certs for webhooks.
The certs are stored in a secret. The secret is mounted as volume into a pod.
Once the volume is updated in the pod. The ca certs in the webhook configurations are updated.