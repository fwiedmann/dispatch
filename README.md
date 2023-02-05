# dispatch

This is a PoC.

With dispatch you can automatically dispatch MR assignments in large mono repositories.

This can be useful when:

- several units/squads/teams work on the same repository
- units/squads/teams have their main code bases, but others will also adjust/enhance the code based on new requirements
- It's a fast pacing environment, but code reviews by the main code owners is critical

Correct MR assignment for the responsible unit/squad/team can be error prone. An automated process can minimize the risk of unwanted code changes and may help for better communication and code quality.

![mono-repo](./assets/mono-repo.png)

## Goals

- Correct approval assignment in terms of code responsibilities
- With additional meta data, peers are able to identify the source and target aduience for the MRs (smth. like tags which are visible in the MR overview list)
- Do not remove already assigned peers
- With categories it is possible to assign peers within a team to a specific kind of code base change. Like frontend or backend.

### Not functional requirements

- Directory analysis should cost <= 10sek
- Impotency (Desired assignees should always result in the same assignments), Re-run the CI/CD job should always work
- Stdout logging in different levels

## Constraints

- Develop for the GitLab API
- Has to be run as a GitLab CI/CD job
- No other persistence than the MR itself is available

## Flow

![alternative text](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.github.com/fwiedmann/dispatch/main/flow.plantuml)
