# authentication
Carry out identity authentication while ensuring user privacy. Implement strict privacy protection mechanisms to safeguard user data and prevent information leaks. Additionally, provide zero-knowledge proofs for identity information.
# Program Description
## Program URLs:
	* 1. https://github.com/scaleomngt/authentication/tree/main/front/identity_based_encryption --Frontend code
	* 2. https://github.com/scaleomngt/authentication --Server code
	* 3. https://github.com/scaleomngt/authentication/tree/main/aleo/card --Leo code
# Installation Instructions
## Frontend Deployment
(1) Step 1: Download the dependencies required for the program
```
npm install
```
(2) Step 2: Verify if all dependencies are successfully downloaded by running the program locally
```
npm run serve  
```
(3) Step 3: Access the configured IP address and port in a browser to use the application

## Backend Deployment
```
go build
```
## leoDeployment
```
- `cd aleo/card && leo build`
- `export PRIVATE_KEY=<PRIVATE_KEY>`
- Publish£º
snarkos developer deploy card.aleo --private-key $PRIVATE_KEY --query "https://vm.aleo.org/api" --path "build/" \
--broadcast "https://vm.aleo.org/api/testnet3/transaction/broadcast" --fee 600000 \
--record <Record_that_you_just_transferred_credits_to>
```
# Data flow diagram
<img src="https://github.com/scaleomngt/authentication/blob/main/t5.png" alt="drawing" width="800"/>
# Program Execution
* 1. Enter the required information such as name, gender, nationality, date of birth, etc., on the webpage, and click submit.
* 2. Invoke the Leo verification program and receive the transaction hash.
* 3. Generate a QR code based on the transaction hash and provide it to the third party for verification. The third party can scan the QR code using a mobile device or browser to view the verification results, without the need to expose the user's identity data.
<img src="https://github.com/scaleomngt/authentication/blob/main/t1.png" alt="drawing" width="800"/>
<img src="https://github.com/scaleomngt/authentication/blob/main/t3.png" alt="drawing" width="800"/>
<img src="https://github.com/scaleomngt/authentication/blob/main/t4.png" alt="drawing" width="600"/>
