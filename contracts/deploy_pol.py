#!/usr/bin/env python3
"""
PoL (Proof of Liquidity) Deployment Automation Script
Automates the deployment of PoL contracts on BeaconKit devnet.
"""

import os
import re
import sys
import json
import time
import shutil
import subprocess
from typing import Dict, List, Optional, Tuple, Any
from dataclasses import dataclass, field, asdict
from datetime import datetime
from pathlib import Path
import logging

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler(f'pol_deployment_{datetime.now().strftime("%Y%m%d_%H%M%S")}.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)


@dataclass
class DeploymentConfig:
    """Configuration for PoL deployment"""
    # Environment settings
    foundry_profile: str = "deploy"
    is_testnet: bool = False
    use_software_wallet: bool = True
    
    # Connection settings
    eth_from: str = "0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4"
    eth_from_pk: str = "0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306"
    rpc_url: str = "http://localhost:8545"
    
    # Deployment parameters
    dry_run: bool = False
    backup_files: bool = True
    state_file: str = "deployment_state.json"
    
    # Contract addresses (will be populated during deployment)
    addresses: Dict[str, str] = field(default_factory=dict)
    token_addresses: List[str] = field(default_factory=list)
    vault_addresses: List[str] = field(default_factory=list)


@dataclass
class DeploymentState:
    """Track deployment progress for recovery"""
    step: str = ""
    completed_steps: List[str] = field(default_factory=list)
    addresses: Dict[str, str] = field(default_factory=dict)
    token_addresses: List[str] = field(default_factory=list)
    vault_addresses: List[str] = field(default_factory=list)
    timestamp: str = ""


class PoLDeployer:
    """Main deployment orchestrator for PoL contracts"""
    
    def __init__(self, config: DeploymentConfig):
        self.config = config
        self.state = DeploymentState()
        self.backup_dir = Path(f"backups_{datetime.now().strftime('%Y%m%d_%H%M%S')}")
        
        if self.config.backup_files:
            self.backup_dir.mkdir(exist_ok=True)
            logger.info(f"Backup directory created: {self.backup_dir}")
    
    def run(self):
        """Execute the complete deployment workflow"""
        try:
            logger.info("Starting PoL deployment automation")
            self.setup_environment()
            
            # Load previous state if exists
            self.load_state()
            
            # Execute deployment steps
            steps = [
                ("predict_addresses", self.predict_addresses),
                ("validate_bgt_config", self.validate_bgt_config),
                ("deploy_bgt", self.deploy_bgt),
                ("deploy_pol", self.deploy_pol_contracts),
                ("change_parameters", self.change_pol_parameters),
                ("deploy_tokens", self.deploy_tokens),
                ("deploy_vaults", self.deploy_reward_vaults),
                ("whitelist_vaults", self.whitelist_vaults),
                ("set_allocations", self.set_default_allocations),
                ("verify_deployment", self.verify_deployment)
            ]
            
            for step_name, step_func in steps:
                if step_name not in self.state.completed_steps:
                    logger.info(f"\n{'='*60}")
                    logger.info(f"Executing step: {step_name}")
                    logger.info(f"{'='*60}")
                    
                    self.state.step = step_name
                    self.save_state()
                    
                    step_func()
                    
                    self.state.completed_steps.append(step_name)
                    self.save_state()
                else:
                    logger.info(f"Skipping completed step: {step_name}")
            
            self.print_deployment_summary()
            logger.info("\nPoL deployment completed successfully!")
            
        except Exception as e:
            logger.error(f"Deployment failed: {e}")
            logger.error(f"Current step: {self.state.step}")
            logger.error("Run again to resume from the last successful step")
            raise
    
    def setup_environment(self):
        """Set up environment variables for deployment"""
        logger.info("Setting up environment variables")
        
        env_vars = {
            "FOUNDRY_PROFILE": self.config.foundry_profile,
            "IS_TESTNET": str(self.config.is_testnet).lower(),
            "USE_SOFTWARE_WALLET": str(self.config.use_software_wallet).lower(),
            "ETH_FROM": self.config.eth_from,
            "RPC_URL": self.config.rpc_url,
            "ETH_FROM_PK": self.config.eth_from_pk
        }
        
        for key, value in env_vars.items():
            os.environ[key] = value
            if key != "ETH_FROM_PK":  # Don't log private key
                logger.debug(f"Set {key}={value}")
    
    def run_command(self, cmd: List[str], capture_output: bool = True) -> subprocess.CompletedProcess:
        """Execute a shell command with error handling"""
        logger.debug(f"Running command: {' '.join(cmd)}")
        
        if self.config.dry_run:
            logger.info(f"[DRY RUN] Would execute: {' '.join(cmd)}")
            return subprocess.CompletedProcess(args=cmd, returncode=0, stdout="", stderr="")
        
        try:
            result = subprocess.run(
                cmd,
                capture_output=capture_output,
                text=True,
                check=True
            )
            return result
        except subprocess.CalledProcessError as e:
            logger.error(f"Command failed: {e.cmd}")
            logger.error(f"Return code: {e.returncode}")
            logger.error(f"Output: {e.stdout}")
            logger.error(f"Error: {e.stderr}")
            raise
    
    def run_forge_script(self, script_path: str, sig: Optional[str] = None, 
                        sig_args: Optional[List[str]] = None) -> str:
        """Execute a forge script and return output"""
        cmd = [
            "forge", "script", script_path,
            "--private-key", self.config.eth_from_pk,
            "--sender", self.config.eth_from,
            "--rpc-url", self.config.rpc_url,
            "--broadcast", "-vv"
        ]
        
        if sig:
            cmd.extend(["--sig", sig])
            if sig_args:
                cmd.extend(sig_args)
        
        result = self.run_command(cmd)
        return result.stdout
    
    def run_cast_command(self, *args) -> str:
        """Execute a cast command and return output"""
        cmd = ["cast"] + list(args)
        result = self.run_command(cmd)
        return result.stdout.strip()
    
    def backup_file(self, file_path: str):
        """Create a backup of a file before modification"""
        if not self.config.backup_files:
            return
        
        source = Path(file_path)
        if source.exists():
            dest = self.backup_dir / source.name
            shutil.copy2(source, dest)
            logger.debug(f"Backed up {source} to {dest}")
    
    def update_solidity_file(self, file_path: str, updates: Dict[str, str]):
        """Update constants in a Solidity file"""
        self.backup_file(file_path)
        
        with open(file_path, 'r') as f:
            content = f.read()
        
        for var_name, new_value in updates.items():
            # Pattern to match Solidity constant declarations
            pattern = rf'(address\s+(?:internal\s+)?constant\s+{re.escape(var_name)}\s*=\s*)([^;]+);'
            replacement = rf'\g<1>{new_value};'
            content = re.sub(pattern, replacement, content)
            logger.debug(f"Updated {var_name} to {new_value}")
        
        with open(file_path, 'w') as f:
            f.write(content)
        
        logger.info(f"Updated {file_path}")
    
    def parse_address_from_output(self, output: str, pattern: str) -> Optional[str]:
        """Extract an address from command output using regex"""
        match = re.search(pattern, output)
        if match:
            return match.group(1)
        return None
    
    def predict_addresses(self):
        """Run POLPredictAddresses script and update POLAddresses.sol"""
        logger.info("Predicting contract addresses...")
        
        output = self.run_forge_script("script/pol/POLPredictAddresses.s.sol")
        
        # Parse addresses from output
        address_patterns = {
            "BERACHEF_ADDRESS": r"BeraChef:\s+(0x[a-fA-F0-9]{40})",
            "REWARD_VAULT_FACTORY_ADDRESS": r"RewardVaultFactory:\s+(0x[a-fA-F0-9]{40})",
            "BGT_ADDRESS": r"BGT:\s+(0x[a-fA-F0-9]{40})",
            "BLOCK_REWARD_CONTROLLER_ADDRESS": r"BlockRewardController:\s+(0x[a-fA-F0-9]{40})",
            "DISTRIBUTOR_ADDRESS": r"Distributor:\s+(0x[a-fA-F0-9]{40})"
        }
        
        updates = {}
        for var_name, pattern in address_patterns.items():
            address = self.parse_address_from_output(output, pattern)
            if address:
                updates[var_name] = address
                self.state.addresses[var_name] = address
                logger.info(f"Predicted {var_name}: {address}")
        
        # Update POLAddresses.sol
        if updates:
            self.update_solidity_file("script/pol/POLAddresses.sol", updates)
    
    def validate_bgt_config(self):
        """Validate BGT and distributor addresses match BeaconKit config"""
        logger.info("Validating BGT configuration...")
        
        bgt_address = self.state.addresses.get("BGT_ADDRESS")
        distributor_address = self.state.addresses.get("DISTRIBUTOR_ADDRESS")
        
        if not bgt_address or not distributor_address:
            logger.warning("BGT or Distributor address not found, skipping validation")
            return
        
        logger.info(f"BGT_ADDRESS: {bgt_address}")
        logger.info(f"DISTRIBUTOR_ADDRESS: {distributor_address}")
        
        # Check BGT balance to verify it's receiving inflation
        try:
            balance = self.run_cast_command(
                "balance", bgt_address, 
                "--rpc-url", self.config.rpc_url
            )
            logger.info(f"Current BGT balance: {balance}")
        except Exception as e:
            logger.warning(f"Could not check BGT balance: {e}")
    
    def deploy_bgt(self):
        """Deploy BGT contract"""
        logger.info("Deploying BGT contract...")
        
        output = self.run_forge_script("script/pol/deployment/2_DeployBGT.s.sol")
        
        # Extract deployed address
        pattern = r"BGT deployed at:\s+(0x[a-fA-F0-9]{40})"
        address = self.parse_address_from_output(output, pattern)
        
        if address:
            self.state.addresses["BGT_DEPLOYED"] = address
            logger.info(f"BGT deployed at: {address}")
    
    def deploy_pol_contracts(self):
        """Deploy core PoL contracts"""
        logger.info("Deploying PoL contracts...")
        
        output = self.run_forge_script("script/pol/deployment/3_DeployPoL.s.sol")
        
        # Parse deployed addresses
        patterns = {
            "BeraChef": r"BeraChef deployed at:\s+(0x[a-fA-F0-9]{40})",
            "BlockRewardController": r"BlockRewardController deployed at:\s+(0x[a-fA-F0-9]{40})",
            "Distributor": r"Distributor deployed at:\s+(0x[a-fA-F0-9]{40})",
            "RewardVaultFactory": r"RewardVaultFactory deployed at:\s+(0x[a-fA-F0-9]{40})"
        }
        
        for name, pattern in patterns.items():
            address = self.parse_address_from_output(output, pattern)
            if address:
                self.state.addresses[name] = address
                logger.info(f"{name} deployed at: {address}")
    
    def change_pol_parameters(self):
        """Configure PoL economic parameters"""
        logger.info("Changing PoL parameters...")
        
        output = self.run_forge_script("script/pol/actions/ChangePOLParameters.s.sol")
        logger.info("PoL parameters updated successfully")
    
    def deploy_tokens(self):
        """Deploy 5 BST tokens for reward vaults"""
        logger.info("Deploying BST tokens...")
        
        self.state.token_addresses = []
        
        for i in range(1, 6):
            logger.info(f"Deploying BST token {i}...")
            
            output = self.run_forge_script(
                "script/misc/testnet/DeployToken.s.sol",
                sig="deployBST(uint256)",
                sig_args=[str(i)]
            )
            
            # Extract token address
            pattern = r"BST deployed at:\s+(0x[a-fA-F0-9]{40})"
            address = self.parse_address_from_output(output, pattern)
            
            if address:
                self.state.token_addresses.append(address)
                logger.info(f"BST{i} deployed at: {address}")
        
        # Update DeployRewardVault.s.sol with token addresses
        if len(self.state.token_addresses) == 5:
            updates = {
                "LP_BERA_HONEY": self.state.token_addresses[0],
                "LP_BERA_ETH": self.state.token_addresses[1],
                "LP_BERA_WBTC": self.state.token_addresses[2],
                "LP_USDC_HONEY": self.state.token_addresses[3],
                "LP_BEE_HONEY": self.state.token_addresses[4]
            }
            
            self.update_solidity_file("script/pol/actions/DeployRewardVault.s.sol", updates)
    
    def deploy_reward_vaults(self):
        """Deploy reward vaults for staking tokens"""
        logger.info("Deploying reward vaults...")
        
        output = self.run_forge_script("script/pol/actions/DeployRewardVault.s.sol")
        
        # Extract vault addresses
        pattern = r"RewardVault deployed at\s+(0x[a-fA-F0-9]{40})\s+for staking token\s+(0x[a-fA-F0-9]{40})"
        matches = re.findall(pattern, output)
        
        self.state.vault_addresses = []
        for vault_addr, token_addr in matches:
            self.state.vault_addresses.append(vault_addr)
            logger.info(f"Vault {vault_addr} deployed for token {token_addr}")
        
        # Update WhitelistRewardVault.s.sol
        if len(self.state.vault_addresses) >= 5:
            # Remove USDS references and update addresses
            self.update_whitelist_script()
    
    def update_whitelist_script(self):
        """Update WhitelistRewardVault.s.sol with vault addresses"""
        file_path = "script/pol/actions/WhitelistRewardVault.s.sol"
        self.backup_file(file_path)
        
        with open(file_path, 'r') as f:
            content = f.read()
        
        # Remove USDS_HONEY references
        content = re.sub(r'.*REWARD_VAULT_USDS_HONEY.*\n', '', content)
        
        # Fix trailing commas in arrays after removing USDS entries
        # This regex finds the last element before the closing bracket and removes trailing comma
        content = re.sub(r',(\s*\])', r'\1', content)
        
        # Update vault addresses
        if len(self.state.vault_addresses) >= 5:
            updates = {
                "REWARD_VAULT_BERA_HONEY": self.state.vault_addresses[0],
                "REWARD_VAULT_BERA_ETH": self.state.vault_addresses[1],
                "REWARD_VAULT_BERA_WBTC": self.state.vault_addresses[2],
                "REWARD_VAULT_USDC_HONEY": self.state.vault_addresses[3],
                "REWARD_VAULT_BEE_HONEY": self.state.vault_addresses[4]
            }
            
            for var_name, address in updates.items():
                pattern = rf'(address\s+(?:internal\s+)?constant\s+{re.escape(var_name)}\s*=\s*)([^;]+);'
                replacement = rf'\g<1>{address};'
                content = re.sub(pattern, replacement, content)
        
        with open(file_path, 'w') as f:
            f.write(content)
        
        logger.info(f"Updated {file_path}")
    
    def whitelist_vaults(self):
        """Whitelist reward vaults and set max weight"""
        logger.info("Whitelisting reward vaults...")
        
        output = self.run_forge_script("script/pol/actions/WhitelistRewardVault.s.sol")
        logger.info("Vaults whitelisted successfully")
        
        # Set max weight per vault
        berachef_address = self.state.addresses.get("BERACHEF_ADDRESS", 
                                                   self.state.addresses.get("BeraChef"))
        if berachef_address:
            logger.info("Setting max weight per vault...")
            self.run_cast_command(
                "send", berachef_address,
                "setMaxWeightPerVault(uint96)", "2000",
                "--private-key", self.config.eth_from_pk,
                "--rpc-url", self.config.rpc_url, "-vv"
            )
            logger.info("Max weight per vault set to 2000")
    
    def set_default_allocations(self):
        """Set default reward allocations"""
        logger.info("Setting default reward allocations...")
        
        # Update SetDefaultRewardAllocation.s.sol
        file_path = "script/pol/actions/SetDefaultRewardAllocation.s.sol"
        self.backup_file(file_path)
        
        with open(file_path, 'r') as f:
            content = f.read()
        
        # Update weights to 2000
        content = re.sub(r'(REWARD_VAULT_\w+_WEIGHT\s*=\s*)\d+', r'\g<1>2000', content)
        
        # Update vault addresses
        if len(self.state.vault_addresses) >= 5:
            updates = {
                "REWARD_VAULT_BERA_HONEY": self.state.vault_addresses[0],
                "REWARD_VAULT_BERA_ETH": self.state.vault_addresses[1],
                "REWARD_VAULT_BERA_WBTC": self.state.vault_addresses[2],
                "REWARD_VAULT_USDC_HONEY": self.state.vault_addresses[3],
                "REWARD_VAULT_BEE_HONEY": self.state.vault_addresses[4]
            }
            
            for var_name, address in updates.items():
                pattern = rf'(address\s+(?:internal\s+)?constant\s+{re.escape(var_name)}\s*=\s*)([^;]+);'
                replacement = rf'\g<1>{address};'
                content = re.sub(pattern, replacement, content)
        
        with open(file_path, 'w') as f:
            f.write(content)
        
        # Run the script
        output = self.run_forge_script(
            "script/pol/actions/SetDefaultRewardAllocation.s.sol:WhitelistIncentiveTokenScript"
        )
        logger.info("Default reward allocations set successfully")
    
    def verify_deployment(self):
        """Verify BGT distribution is working"""
        logger.info("Verifying BGT distribution...")
        
        bgt_address = self.state.addresses.get("BGT_ADDRESS", 
                                               self.state.addresses.get("BGT_DEPLOYED"))
        
        if not bgt_address:
            logger.warning("BGT address not found, skipping verification")
            return
        
        # Check BGT contract balance
        logger.info("Monitoring BGT contract balance changes...")
        balances = []
        
        for i in range(3):
            balance = self.run_cast_command(
                "balance", bgt_address,
                "--rpc-url", self.config.rpc_url
            )
            balances.append(int(balance))
            logger.info(f"Check {i+1}: BGT contract balance = {balance}")
            
            if i < 2:
                time.sleep(5)  # Wait between checks
        
        # Verify balance is increasing
        if len(set(balances)) > 1 and balances[-1] > balances[0]:
            logger.info("✓ BGT distribution verified - contract balance is increasing")
        else:
            logger.warning("⚠ BGT contract balance not increasing - distribution may not be working")
        
        # Output instructions for checking operator balance
        logger.info("\n" + "="*60)
        logger.info("IMPORTANT: Verify Validator Operator BGT Distribution")
        logger.info("="*60)
        logger.info("\nTo verify that your validator operator is receiving BGT rewards,")
        logger.info("run the following command with your operator address:\n")
        logger.info(f"cast call {bgt_address} \\")
        logger.info(f"  \"balanceOf(address)(uint256)\" \\")
        logger.info(f"  <OPERATOR_ADDRESS> \\")
        logger.info(f"  --rpc-url {self.config.rpc_url}")
        logger.info("\nReplace <OPERATOR_ADDRESS> with your validator operator address.")
        logger.info("\nThe balance should increase every time the validator produces a block.")
        logger.info("Run the command multiple times to verify the balance is increasing.")
        logger.info("="*60)
    
    def save_state(self):
        """Save deployment state for recovery"""
        self.state.timestamp = datetime.now().isoformat()
        
        with open(self.config.state_file, 'w') as f:
            json.dump(asdict(self.state), f, indent=2)
        
        logger.debug(f"State saved to {self.config.state_file}")
    
    def load_state(self):
        """Load previous deployment state if exists"""
        if os.path.exists(self.config.state_file):
            with open(self.config.state_file, 'r') as f:
                data = json.load(f)
                self.state = DeploymentState(**data)
                logger.info(f"Loaded previous state from {self.config.state_file}")
                logger.info(f"Completed steps: {', '.join(self.state.completed_steps)}")
    
    def print_deployment_summary(self):
        """Print a summary of the deployment"""
        logger.info("\n" + "="*60)
        logger.info("DEPLOYMENT SUMMARY")
        logger.info("="*60)
        
        logger.info("\nCore Contracts:")
        for name, address in self.state.addresses.items():
            logger.info(f"  {name}: {address}")
        
        logger.info("\nToken Addresses:")
        for i, address in enumerate(self.state.token_addresses, 1):
            logger.info(f"  BST{i}: {address}")
        
        logger.info("\nReward Vault Addresses:")
        for i, address in enumerate(self.state.vault_addresses, 1):
            logger.info(f"  Vault{i}: {address}")
        
        # Save summary to file
        summary_file = f"deployment_summary_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
        with open(summary_file, 'w') as f:
            json.dump({
                "timestamp": datetime.now().isoformat(),
                "config": asdict(self.config),
                "contracts": self.state.addresses,
                "tokens": self.state.token_addresses,
                "vaults": self.state.vault_addresses
            }, f, indent=2)
        
        logger.info(f"\nDeployment summary saved to: {summary_file}")


def main():
    """Main entry point"""
    import argparse
    
    parser = argparse.ArgumentParser(description="Deploy PoL contracts on BeaconKit")
    parser.add_argument("--rpc-url", default="http://localhost:8545", help="RPC endpoint")
    parser.add_argument("--eth-from", default="0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4", 
                       help="Deployer address")
    parser.add_argument("--eth-from-pk", 
                       default="0xfffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306",
                       help="Deployer private key")
    parser.add_argument("--dry-run", action="store_true", help="Simulate deployment without executing")
    parser.add_argument("--no-backup", action="store_true", help="Skip file backups")
    parser.add_argument("--reset", action="store_true", help="Reset deployment state")
    
    args = parser.parse_args()
    
    # Reset state if requested
    if args.reset and os.path.exists("deployment_state.json"):
        os.remove("deployment_state.json")
        logger.info("Deployment state reset")
    
    # Configure deployment
    config = DeploymentConfig(
        rpc_url=args.rpc_url,
        eth_from=args.eth_from,
        eth_from_pk=args.eth_from_pk,
        dry_run=args.dry_run,
        backup_files=not args.no_backup
    )
    
    # Run deployment
    deployer = PoLDeployer(config)
    deployer.run()


if __name__ == "__main__":
    main()