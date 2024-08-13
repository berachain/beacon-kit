# Smart Contract Deployment Instructions

## General guidelines about `forge-config.yaml`

- The format of the `repository` should be `github.com/<OrgName/Username>/<Repositoryname>`. Refer to the Kurtosis documentation on [locators](https://docs.kurtosis.com/advanced-concepts/locators) for more information.

- The `script_path` should be relative to the `repository` and `contracts_path` (if present).

- For the wallet, currently only `private_key` is supported.

- For `rpc_url`, 
    - URL for externally running network.

    - If you want to spin up a devnet locally, you could use Kurtosis `make start-devnet`.
        `rpc_url` would be `http://HOST_IP_ADDRESS:8547` , Do not change the port as 8547 is the public port for the Erigon node.



## There could be different cases -

- **Contract directory does not use git submodules**

    If the contract directory does not use git submodules, provide the path up to the contracts in the `repository` and leave `contracts_path` as an empty string.

- **Contract directory use git submodules**

    If the contract directory uses git submodules, then we need to clone the whole repository to get the submodules. In that case, we need to provide the `repository` at the root level and `contracts_path` where the contracts are present.


## Notes: 

If there's a contract present locally, the only two options supported by kurtosis are :

- The directory should be present inside the same kurtosis package as the current one.

- Use github URL.

    Edge scenario:
    This would not be supported if the contracts are part of this repository - beacon-kit, in that case, fork the repository into your user profile and use that as a repository.

    I know this is a kind of workaround, until Kurtosis supports local URLs, we don't really have a choice.

    **Example for running contracts in beacon-kit**

    ```bash
    repository: "github.com/nidhi-singh02/beacon-kit"
    contracts_path: "contracts"
    script_path: "script/DeployAndCallERC20.s.sol"
    contract_name: "DeployAndCallERC20"
    ```

## Example for GitHub hosted repository

```bash

repository: "github.com/nidhi-singh02/solidity-scripting/"
contracts_path: ""
script_path: "script/NFT.s.sol"
contract_name: "MyScript"
```
