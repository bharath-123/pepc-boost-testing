# PEPC-Boost Custom Devnet Testing

This repo contains code to send ToB txs to a PEPC-Boost relay deployed in a custom kurtosis devnet.

## Usage

Add the I.P address of your PEPC-Boost relay in `constants/addresses.go` by substituting it in MevRelayerUrl variable and also add the I.P address of your
execution layer client in `constants/addresses.go` by substituting it in the EcUrl:

Then run the following command to send ToB txs to the PEPC-Boost relay

```bash
go run main.go
```

To inspect the block post submitting a ToB tx, you can uncomment the block inspection code and comment out the ToB tx submission code. Once that is run, you can run the code similar 
to submitting a ToB tx.