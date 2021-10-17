#!/bin/bash

kaml concat \
  channels/packages/coredns/1.7.0/clusterrole.yaml \
  channels/packages/coredns/1.7.0/clusterrolebinding.yaml \
  <(kustomize build config) | \
kaml replace-image operator=justinsb/coredns-operator:latest | \
kaml normalize-labels

