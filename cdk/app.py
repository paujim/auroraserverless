#!/usr/bin/env python3

from aws_cdk import core

from serverless_aurora.serverless_aurora_stack import ServerlessAuroraStack


app = core.App()
ServerlessAuroraStack(app, "serverless-aurora")

app.synth()
