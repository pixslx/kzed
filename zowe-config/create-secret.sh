#!/bin/bash

oc -n kzed-system create secret generic zowe-config --from-file=zowe.config.json --from-file=zowe.schema.json --from-literal SYSUID=<SYSUID>