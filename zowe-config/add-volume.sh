#!/bin/bash

oc -n kzed-system set volume deployment/kzed-controller-manager --add --type=secret --secret-name=zowe-config --mount-path=/zowe-config