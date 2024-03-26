---
title: Synchronize Pull Request Previews in Automation
description: Learn how to synchronize pull request previews with Terramate CLI to Terramate Cloud in automation.
---

# Synchronize Pull Request Previews in Automation

::: warning
Pull Request Previews are currently only supported in GitHub Actions.
More integrations are on the roadmap.
:::

## Automation Blueprints

Terramate CLI can run in any CI/CD and we provide Blueprints for various CI/CD platforms:

- [GitHub Actions Blueprints](../../cli/automation/github-actions/preview-workflow.md)
- GitLab CI Blueprints 🚧
- Bitbucket Pipelines Blueprints 🚧
- Azure DevOps Blueprints 🚧

## Required Permission

To gather metadata from GitHub about the pull request associated with the preview, a `GITHUB_TOKEN` needs to be exposed or a valid GitHub CLI configuration needs to be available.

## Best Practices

- Restrict elevated access to your cloud providers (AWS, Google Cloud, or Azure) and access to Terraform State to automation flows.
- Use OpenID Connect (OIDC) to authenticate to your Cloud Provider to use short-lived credentials - Terramate CLI uses OIDC by default.
- Ensure that all drift checks run on all stacks even if some errors are detected using the `--continue-on-error` command line option
- Run a Drift Check right after the deployment and synchronize the result to Terramate Cloud to get an immediate health check and ensure the deployment is stable.
- Run a Drift Check at least every 24 hours to get a detailed history of when drifts were introduced.
- Set up notifications to get informed about newly detected drifts in stacks
