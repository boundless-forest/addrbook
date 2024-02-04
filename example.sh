go build .

# Create workspace

./addrbook workspace new --name MyProject1
./addrbook workspace new --name MyProject2

# Save contract address

./addrbook workspace save --workspace MyProject1 --contract Contract1 --address 0x000000000000000000000000000000000000000x --note "Ethereum side"
./addrbook workspace save --workspace MyProject1 --contract Contract2 --address 0x000000000000000000000000000000000000000x --note "Sepolia side"
./addrbook workspace save --workspace MyProject2 --contract Demo --address 0x000000000000000000000000000000000000000x --note "In this section, you have the opportunity to document details about this contract, including the specific network where it was deployed or the identity of the individual or team responsible for its deployment. This information can serve as a valuable reference for future maintenance, updates, or audits."

# Update contract address

./addrbook workspace update --workspace MyProject1 --contract Contract2 --address 0x000000000000000000000000000000000000000y --note "Updated note version"
