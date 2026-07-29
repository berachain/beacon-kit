# Smart Contract Deployment Instructions

## General guidelines about `config.yaml`

- The format of the `repository` should be `github.com/<OrgName/Username>/<Repositoryname>`. Refer to the Kurtosis documentation on [locators](https://docs.kurtosis.com/advanced-concepts/locators) for more information.

- The `script_path` should be relative to the `repository` and `contracts_path` (if present).

- For the wallet, currently only `private_key` is supported.

- For `network`, network name as defined in hardhat.config.ts of the contract repository.

- If the smart contract has prerequisites or dependencies, there are two options:
    1) If the contracts repository has its dependency built-in, provide the full path to the dependency script.
    2) Else, to set it locally via kurtosis package. Ensure that the `dependency.sh` file inside the `dependency` folder is completed.
    
  Set the `dependency` `type` to "git" or "local". This is necessary when additional setup is required before deployment.
  In case of no dependency, set it to "none".

## Example for GitHub hosted repository having dependency in the local kurtosis package

```yaml
deployment:
  repository: "github.com/nidhi-singh02/CrocSwap-protocol"  # give repo name if there are submodules, else give the folder till contracts
  contracts_path: ""  # give the path till contracts, if the repository is the contract folder itself, then leave it empty
  script_path: "misc/scripts/deploy.ts"  # this must be relative to the repository path + contracts_path(if applicable)
  network: "bartio"
  wallet:
    type: "private_key"  # currently only private_key wallet is supported. Do not change the type.
    value: "0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306"
  dependency:
    type: local  # type can be "local" or "git".
    path: "dependency.sh"  # Full path for local dependency.
```

## Example for GitHub hosted repository having dependency on Git

```yaml
deployment:
  repository: "github.com/nidhi-singh02/CrocSwap-protocol"
  contracts_path: ""
  script_path: "misc/scripts/deploy.ts"
  network: "bartio"
  wallet:
    type: "private_key"
    value: "0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306"
  dependency:
    type: git
    path: "misc/scripts/dependency/dependency.sh"
```
